package hindley_milner

import (
	"fmt"

	"github.com/pkg/errors"
)

// Cloner is any type that can clone
type Cloner interface {
	Clone() interface{}
}

// Fresher keeps track of all the TypeVariables that has been generated so far. It has one method - Fresh(), which is to create a new TypeVariable
type Fresher interface {
	Fresh() TypeVariable
}

func saveExprContext(t Type, source *Expression) Type {
	return t.WithContext(CreateCodeContext(*source))
}

type IntrospectionConstraint struct {
	tv TypeVariable
	argTV Type
	context CodeContext
}

type OverloadConstraint struct {
	tv TypeVariable
	name string
	alternatives []*Scheme
	context CodeContext
}

type inferer struct {
	env Env
	retEnv Env
	cs  Constraints
	t   Type
	ocs []OverloadConstraint
	ics []IntrospectionConstraint
	returns []Type

	count int
	config *InferConfiguration
	csflag int
}

func newInferer(env Env, config *InferConfiguration) *inferer {
	return &inferer{
		env: env,
		config: config,
	}
}

func (infer *inferer) Fresh() TypeVariable {
	retVal := infer.count
	infer.count++
	return TypeVariable{
		value: int16(retVal),
	}
}

func (infer *inferer) lookup(isLiteral bool, name string, source Expression) error {
	if infer.env.IsBuiltin(name) {
		//fmt.Printf("%s IS BUILTIN\n", name)
		scheme, _ := infer.env.SchemeOf(name)
		//fmt.Printf("SO TYPE IS %v\n", scheme)
		infer.t = Instantiate(infer, scheme.Clone())
		return nil
	}
	s, ok := infer.env.SchemeOf(name)
	if !ok {
		return UndefinedSymbol{
			Name: name,
			Source: source,
			IsLiteral: isLiteral,
			IsVariable: !isLiteral,
		}
	}
	infer.t = Instantiate(infer, s)
	return nil
}

func (infer *inferer) resolveProxy(expr Expression, exprType ExpressionType) (Expression, ExpressionType) {
	// Resolve proxies
	for {
		if expr == nil {
			return expr, exprType
		}
		if exprType != E_PROXY {
			break
		}
		expr = expr.Body()
		exprType = expr.ExpressionType()
	}
	return expr, exprType
}

func (infer *inferer) consGen(expr Expression, forceType ExpressionType, isTop bool, isOpaqueTop bool) (err error) {

	defer infer.cleanupConstraints()

	if expr == nil {
		tv := infer.Fresh()
		infer.t = tv
		return nil
	}

	// fallbacks

	exprType := expr.ExpressionType()
	if forceType != E_NONE {
		exprType = forceType
	}

	// Resolve unions
	expr, exprType = infer.resolveProxy(expr, exprType)

	// Optionaly get ident dependencies
	if exprWithDeps, ok := expr.(ExpressionWithIdentifiersDeps); ok {
		idents := exprWithDeps.GetIdentifierDeps()
		for _, name := range idents {
			// Declare identifiers
			infer.env.Add(name,
				NewScheme(TypeVarSet{TVar('a')}, TVar('a')))
		}
	}

	// explicit types/inferers - can fail
	switch et := expr.(type) {
	case Typer:
		if infer.t = et.Type(); infer.t != nil {
			infer.t = saveExprContext(infer.t, &expr)
			return nil
		}
	case Inferer:
		if infer.t, err = et.Infer(infer.env, infer); err == nil && infer.t != nil {
			infer.t = saveExprContext(infer.t, &expr)
			return nil
		}

		err = nil // reset errors
	}

	switch exprType {
	case E_CUSTOM:
		et := expr.(CustomExpression)

		context := CustomExpressionEnv{
			Env:                 infer.env,
			InferencedType:      infer.t,
			LookupEnv:           func(isLiteral bool, name string) error {
				return infer.lookup(isLiteral, name, et)
			},
			GenerateConstraints: func(expr Expression) (e error, env Env, i Type, constraints Constraints) {
				backupT := infer.t
				backupEnv := infer.env
				backupCS := infer.cs
				if err = infer.consGen(et.Body(), E_NONE, false, false); err != nil {
					return err, nil, nil, nil
				}
				retT := infer.t
				retEnv := infer.env
				retCS := infer.cs
				infer.t = backupT
				infer.env = backupEnv
				infer.cs = backupCS
				return nil, retEnv, retT, retCS
			},
		}
		err, newEnv, newType, newCS := et.GenerateConstraints(context)
		if err != nil {
			return err
		}
		infer.env = newEnv
		infer.t = newType
		infer.cs = newCS

	case E_INTROSPECTION:
		et := expr.(IntrospectionExpression)

		if err = infer.consGen(et.Body(), E_NONE, false, false); err != nil {
			return err
		}
		bodyType := infer.t

		tvArg := infer.Fresh()
		tvRet := infer.Fresh()

		infer.t = tvRet
		infer.ics = append(infer.ics, IntrospectionConstraint{
			tv:      tvArg,
			argTV: bodyType,
			context: CreateCodeContext(expr),
		})
		//infer.cs = append(infer.cs, Constraint{
		//	a:       tvArg,
		//	b:       bodyType,
		//	context: CodeContext{},
		//})

		return nil

	case E_LITERAL:
		et := expr.(Literal)
		if len(et.Name().GetNames()) != 1 {
			return fmt.Errorf("Literal entity cannot conntain other value than one variable name. You cannot use Names batch here.")
		}
		name := et.Name().GetNames()[0]
		//fmt.Printf("GOT %s\n", name)
		if infer.env.IsOverloaded(name) {
			tv := infer.Fresh()
			infer.t = tv
			infer.ocs = append(infer.ocs, OverloadConstraint{
				name: name,
				tv:           tv,
				alternatives: infer.env.OverloadedAlternatives(name),
				context: CreateCodeContext(expr),
			})
			return nil
		}
		return infer.lookup(true, name, et)

	case E_VAR:
		et := expr.(Var)
		if len(et.Name().GetNames()) != 1 {
			return fmt.Errorf("Var entity cannot conntain other value than one variable name. You cannot use Names batch here.")
		}
		name := et.Name().GetNames()[0]
		if err = infer.lookup(false, name, et); err != nil {
			infer.env.Add(name, &Scheme{t: et.Type()})
			err = nil
		}

	case E_TYPE:
		et := expr.(EmbeddedType)
		scheme := et.EmbeddedType()
		//tempName := fmt.Sprintf("__embt_%d", int(rand.Int63()))
		//infer.env.Add(tempName, scheme)
		//err = infer.lookup(false, tempName, et)
		//infer.env.Remove(tempName)
		infer.t = Instantiate(infer, scheme)

	case E_RETURN:
		et := expr.(Return)
		if err = infer.consGen(et.Body(), E_NONE, false, false); err != nil {
			return err
		}
		infer.returns = append(infer.returns, infer.t)
		tv := infer.Fresh()
		infer.t = tv
		if defaultTyper, ok := et.(DefaultTyper); ok {
			infer.cs = append(infer.cs, Constraint{
				a:       tv,
				b:       Instantiate(infer, defaultTyper.DefaultType()),
				context: CreateCodeContext(expr),
			})
		} else if infer.config.CreateDefaultEmptyType() != nil {
			infer.cs = append(infer.cs, Constraint{
				a:       tv,
				b:       Instantiate(infer, infer.config.CreateDefaultEmptyType()),
				context: CreateCodeContext(expr),
			})
		}
		infer.t = tv

	case E_LAMBDA, E_FUNCTION:
		et := expr.(Lambda)

		env := []Env{}
		tv := []TypeVariable{}

		// Clear returns
		rets := infer.returns // backup
		infer.returns = []Type{}

		names := et.Args().GetNames()
		for _, name := range names {
			varType := et.Args().GetTypeOf(name)

			tv = append(tv, infer.Fresh())
			env = append(env, infer.env) // backup

			infer.env = infer.env.Clone()
			infer.env.Remove(name)
			sc := new(Scheme)
			sc.t = tv[len(tv)-1].WithContext(CreateCodeContext(expr))
			infer.env.Add(name, sc)

			if varType != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       sc.t,
					b:       Instantiate(infer, varType),
					context: CreateCodeContext(expr),
				})
			}
		}

		if exprType == E_FUNCTION {
			if err = infer.consGen(et.Body(), E_NONE, false, false); err != nil {
				return err
			}
			bodyType := infer.Fresh()
			infer.t = bodyType
		} else {
			if err = infer.consGen(et.Body(), E_NONE, false, false); err != nil {
				return err
			}
		}

		// here we have an error on zero parameters
		for i:=0; i<len(names); i++ {
			infer.t = NewFnType(tv[len(tv)-1], infer.t).WithContext(CreateCodeContext(expr))

			infer.env = env[len(env)-1] // restore backup
			env = env[:len(env)-1]
			tv = tv[:len(tv)-1]
		}
		if len(names) == 0 {
			infer.t = NewFnType(infer.t).WithContext(CreateCodeContext(expr))
		}

		for _, ret := range infer.returns {
			infer.cs = append(infer.cs, Constraint{
				a:       ret,
				b:       infer.t.(*FunctionType).Ret(true),
				context: CreateCodeContext(expr),
			})
		}
		if true || len(infer.returns) == 0 {
			r := infer.t.(*FunctionType).Ret(true)
			if defaultTyper, ok := et.(DefaultTyper); ok {
				infer.cs = append(infer.cs, Constraint{
					a:       r,
					b:       Instantiate(infer, defaultTyper.DefaultType()),
					context: CreateCodeContext(expr),
				})
			} else if infer.config.CreateDefaultEmptyType() != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       r,
					b:       Instantiate(infer, infer.config.CreateDefaultEmptyType()),
					context: CreateCodeContext(expr),
				})
			}
		}

		infer.returns = rets

	case E_APPLICATION:
		et := expr.(Apply)
		firstExec := true
		batchErr := ApplyBatch(et.Body(), func(body Expression) error {
			if firstExec {
				if err = infer.consGen(et.Fn(), E_NONE, false, false); err != nil {
					return err
				}
				firstExec = false
			}
			fnType, fnCs := infer.t, infer.cs

			if err = infer.consGen(body, E_NONE, false, false); err != nil {
				return err
			}
			bodyType, bodyCs := infer.t, infer.cs

			tv := infer.Fresh()
			cs := append(fnCs, bodyCs...)
			cs = append(cs, Constraint{fnType, saveExprContext(NewFnType(bodyType, tv), &expr), CreateCodeContext(expr)})

			infer.t = tv
			infer.t = saveExprContext(infer.t, &expr)
			infer.cs = cs

			return nil
		})
		if batchErr != nil {
			return batchErr
		}

	case E_BLOCK, E_OPAQUE_BLOCK:
		et := expr.(Block)
		env := infer.env // backup
		//t := infer.t

		for _, statement := range et.GetContents().Expressions() {
			//if true { break }
			//infer.t = t
			//infer.env = env

			if err = infer.consGen(statement, E_NONE, false, exprType == E_OPAQUE_BLOCK && isOpaqueTop); err != nil {
				return err
			}

			//infer.t = t
			//infer.env = env
		}

		tv := infer.Fresh()
		if defaultTyper, ok := et.(DefaultTyper); ok {
			infer.cs = append(infer.cs, Constraint{
				a:       tv,
				b:       Instantiate(infer, defaultTyper.DefaultType()),
				context: CreateCodeContext(expr),
			})
		} else if infer.config.CreateDefaultEmptyType() != nil {
			infer.cs = append(infer.cs, Constraint{
				a:       tv,
				b:       Instantiate(infer, infer.config.CreateDefaultEmptyType()),
				context: CreateCodeContext(expr),
			})
		}

		if isTop || isOpaqueTop {
			infer.retEnv = infer.env.Clone()
		}

		infer.t = tv
		if exprType != E_OPAQUE_BLOCK {
			infer.env = env
		}


	case E_LET_RECURSIVE, E_DECLARATION, E_FUNCTION_DECLARATION:

		// env := infer.env // backup

		et := expr.(LetBase)
		names := et.Var().GetNames()
		types := []*Scheme{}

		definitions := []Expression{}
		if len(names) == 1 {
			if exprType == E_FUNCTION_DECLARATION {
				definitions = append(definitions, et.(Lambda))
			} else {
				if batch, ok := et.(Let).Def().(Batch); ok {
					definitions = append(definitions, batch.Expressions()[0])
				} else {
					definitions = append(definitions, et.(Let).Def())
				}
				if et.Var().HasTypes() {
					types = append(types, et.Var().GetTypeOf(names[0]))
				}
			}
		} else if len(names) > 1 {
			if exprType == E_FUNCTION_DECLARATION {
				for range names {
					definitions = append(definitions, et.(Lambda))
				}
			} else {
				for i, expr := range et.(Let).Def().(Batch).Expressions() {
					definitions = append(definitions, expr)
					if et.Var().HasTypes() {
						types = append(types, et.Var().GetTypeOf(names[i]))
					}
				}
			}
		} else {
			panic("Invalid number of identifiers returned by Var() of the declaration/let: zero.")
		}

		if len(types) != len(names) {
			types = []*Scheme{}
		}

		for i, _ := range names {
			name := names[i]
			def := definitions[i]
			body := expr.Body()
			tv := infer.Fresh()

			var defExpectedType *Scheme = nil
			if len(types) > 0 {
				defExpectedType = types[i]
			}

			//if exprType == E_FUNCTION_DECLARATION {
			//	def = def.(Lambda)
			//} else {
			//	def = def.(Let).Def()
			//}

			//
			//if len(et.Var().GetNames()) != 1 {
			//	return fmt.Errorf("LetRec entity cannot conntain other value than one variable name. You cannot use Names batch here.")
			//}
			//name := et.Var().GetNames()[0]

			infer.env = infer.env.Clone()
			infer.env.Remove(name)
			infer.env.Add(name, &Scheme{tvs: TypeVarSet{tv}, t: tv})

			nonVal := false
			defResolved, _ := infer.resolveProxy(def, def.ExpressionType())
			if block, ok := defResolved.(Block); ok && exprType == E_DECLARATION {
				if len(block.GetContents().Expressions()) == 0 {
					nonVal = true
				}
			}

			if nonVal {
				if defExpectedType != nil {
					infer.t = Instantiate(infer, defExpectedType)
				} else {
					tv := infer.Fresh()
					infer.t = tv
				}
			} else if exprType == E_FUNCTION_DECLARATION {
				//fn := expr.(Lambda)
				if err = infer.consGen(def, E_FUNCTION, false, false); err != nil {
					return err
				}
			} else {
				//def := expr.(Let).Def()
				if err = infer.consGen(def, E_NONE, false, false); err != nil {
					return err
				}
			}
			defType, defCs := infer.t, infer.cs

			s := newSolver()
			s.solve(defCs)
			if s.err != nil {
				return err
			}

			sc := Generalize(infer.env.Apply(s.sub).(Env), saveExprContext(defType.Apply(s.sub).(Type), &expr))

			infer.env.Remove(name)
			infer.env.Add(name, sc)

			if exprType == E_DECLARATION || exprType == E_FUNCTION_DECLARATION {
				retType := infer.Fresh()
				if defaultTyper, ok := expr.(DefaultTyper); ok {
					infer.cs = append(infer.cs, Constraint{
						a:       retType,
						b:       Instantiate(infer, defaultTyper.DefaultType()),
						context: CreateCodeContext(expr),
					})
				} else if infer.config.CreateDefaultEmptyType() != nil {
					infer.cs = append(infer.cs, Constraint{
						a:       retType,
						b:       Instantiate(infer, infer.config.CreateDefaultEmptyType()),
						context: CreateCodeContext(expr),
					})
				}
				infer.t = retType
			} else {
				if err = infer.consGen(body, E_NONE, false, false); err != nil {
					return err
				}
				infer.t = infer.t.Apply(s.sub).(Type)
				infer.t = saveExprContext(infer.t, &expr)
				infer.cs = infer.cs.Apply(s.sub).(Constraints)
			}

			infer.cs = append(infer.cs, defCs...)
			// Add expected type
			if defExpectedType != nil {
				actualType := defType
				if exprType == E_LET_RECURSIVE {
					actualType = sc.t
				}
				//fmt.Printf("Expect %v to be %v\n", sc.t, Instantiate(infer, defExpectedType))
				infer.cs = append(infer.cs, Constraint{
					a:       actualType,
					b:       Instantiate(infer, defExpectedType),
					context: CreateCodeContext(expr),
				})
			}
		}

	case E_LET:
		et := expr.(Let)
		if len(et.Var().GetNames()) != 1 {
			return fmt.Errorf("Let entity cannot conntain other value than one variable name. You cannot use Names batch here.")
		}
		name := et.Var().GetNames()[0]

		env := infer.env

		if err = infer.consGen(et.Def(), E_NONE, false, false); err != nil {
			return err
		}
		defType, defCs := infer.t, infer.cs

		s := newSolver()
		s.solve(defCs)
		if s.err != nil {
			return err
		}

		sc := Generalize(env.Apply(s.sub).(Env), saveExprContext(defType.Apply(s.sub).(Type), &expr))
		infer.env = infer.env.Clone()
		infer.env.Remove(name)
		infer.env.Add(name, sc)

		if err = infer.consGen(et.Body(), E_NONE, false, false); err != nil {
			return err
		}

		infer.t = infer.t.Apply(s.sub).(Type)
		infer.t = saveExprContext(infer.t, &expr)
		infer.cs = infer.cs.Apply(s.sub).(Constraints)
		infer.cs = append(infer.cs, defCs...)

	default:
		return errors.Errorf("Expression of %T is unhandled", expr)
	}

	return nil
}

func (inferer *inferer) cleanupConstraintsRemoveDuplicates() {
	hashtable := map[string]Constraint{}
	for _, cs := range inferer.cs {
		hashtable[fmt.Sprintf("%v", cs)] = cs
	}
	inferer.cs = []Constraint{}
	for _, v := range hashtable {
		inferer.cs = append(inferer.cs, v)
	}
}

func (infer *inferer) cleanupConstraints() {
	infer.cleanupConstraintsRemoveDuplicates()

	if infer.csflag >= 1 {
		infer.csflag = 0
	} else {
		infer.csflag = infer.csflag + 1
		return
	}
	//prevLen := len(infer.cs)

	cs := Constraints{}
	freeVars := map[int16]map[Type]interface{}{}
	contexts := map[int16]CodeContext{}
	for _, cons := range infer.cs {
		if !cons.a.Eq(cons.b) {
			if tv, ok := cons.a.(TypeVariable); ok {
				if _, has := freeVars[tv.value]; !has {
					freeVars[tv.value] = map[Type]interface{}{}
					contexts[tv.value] = tv.context
				}
				has := false
				for q, _ := range freeVars[tv.value] {
					if q.Eq(cons.b) {
						has = true
						break
					}
				}
				if !has {
					freeVars[tv.value][cons.b] = true
				}
			} else {
				cs = append(cs, cons)
			}
		}
	}
	for id, context := range contexts {
		cons := freeVars[id]
		for b, _ := range cons {
			cs = append(cs, Constraint{
				a:       TypeVariable{
					value:   id,
					context: context,
				},
				b:       b,
				context: context,
			})
		}
	}

	infer.cs = cs

	//fmt.Printf("Cleanup %d => %d\n", prevLen, len(infer.cs))
}

// Instantiate takes a fresh name generator, an a polytype and makes a concrete type out of it.
//
// If ...
// 		  Γ ⊢ e: T1  T1 ⊑ T
//		----------------------
//		       Γ ⊢ e: T
//
func Instantiate(f Fresher, s *Scheme) Type {
	l := len(s.tvs)
	tvs := make(TypeVarSet, l)

	var sub Subs
	if l > 30 {
		sub = make(mSubs)
	} else {
		sub = newSliceSubs(l)
	}

	for i, tv := range s.tvs {
		fr := f.Fresh()
		tvs[i] = fr
		sub = sub.Add(tv, fr)
	}

	return s.t.Apply(sub).(Type)
}

// Generalize takes an env and a type and creates the most general possible type - which is a polytype
//
// Generalization
//
// If ...
//		  Γ ⊢ e: T1  T1 ∉ free(Γ)
//		---------------------------
//		   Γ ⊢ e: ∀ α.T1
func Generalize(env Env, t Type) *Scheme {
	logf("generalizing %v over %v", t, env)
	enterLoggingContext()
	defer leaveLoggingContext()
	var envFree, tFree, diff TypeVarSet

	if env != nil {
		envFree = env.FreeTypeVar()
	}

	tFree = t.FreeTypeVar()

	switch {
	case envFree == nil && tFree == nil:
		goto ret
	case len(envFree) > 0 && len(tFree) > 0:
		defer ReturnTypeVarSet(envFree)
		defer ReturnTypeVarSet(tFree)
	case len(envFree) > 0 && len(tFree) == 0:
		// cannot return envFree because envFree will just be sorted and set
	case len(envFree) == 0 && len(tFree) > 0:
		// return ?
	}

	diff = tFree.Difference(envFree)

ret:
	return &Scheme{
		tvs: diff,
		t:   t,
	}
}

type InferConfiguration struct {
	CreateDefaultEmptyType func() *Scheme
	OnConstrintGenerationStarted *func()
	OnConstrintGenerationFinished *func()
	OnSolvingStarted *func()
	OnSolvingFinished *func()
	OnPostprocessingStarted *func()
	OnPostprocessingFinished *func()
}

func NewInferConfiguration() *InferConfiguration {
	return &InferConfiguration{
		CreateDefaultEmptyType: func() *Scheme { return nil },
	}
}

// Infer takes an env, and an expression, and returns a scheme.
func Infer(env Env, expr Expression, config *InferConfiguration) (*Scheme, Env, error) {
	if expr == nil {
		return nil, nil, errors.Errorf("Cannot infer a nil expression")
	}

	if env == nil {
		env = CreateSimpleEnv(map[string][]*Scheme{})
	}

	infer := newInferer(env, config)
	if config.OnConstrintGenerationStarted != nil {
		(*config.OnConstrintGenerationStarted)()
	}
	err := infer.consGen(expr, E_NONE, true, true)
	if config.OnConstrintGenerationFinished != nil {
		(*config.OnConstrintGenerationFinished)()
	}
	if err != nil {
		return nil, nil, err
	}

	s := newSolver()
	if config.OnSolvingStarted != nil {
		(*config.OnSolvingStarted)()
	}
	s.solve(infer.cs)
	if config.OnSolvingFinished != nil {
		(*config.OnSolvingFinished)()
	}
	if s.err != nil {
		return nil, nil, s.err
	}

	//fmt.Printf("SOLVED NOW OCS ARE:\n%#v\n%#v", infer.ocs, infer.cs)
	if config.OnPostprocessingStarted != nil {
		(*config.OnPostprocessingStarted)()
	}
	cleanCS := infer.cs
	for _, ocs := range infer.ocs {
		hasCleanRun := false
		cs := Constraints{}
		for _, alt := range ocs.alternatives {
			cs = Constraints{}
			//copy(cs, cleanCS)
			for _, c := range cleanCS {
				cs = append(cs, c)
			}
			cs = append(cs, Constraint{
				a:       ocs.tv,
				b:       Instantiate(infer, alt),
				context: ocs.tv.context,
			})
			s2 := newSolver()
			s2.solve(cs)
			if s2.err == nil {
				//fmt.Printf("\n\nOLDS CS: %#v\nNEW CS: %#v\n\n", cleanCS, cs)
				hasCleanRun = true
				break
			}
		}
		if hasCleanRun {
			cleanCS = cs
		} else {
			if config.OnPostprocessingFinished != nil {
				(*config.OnPostprocessingFinished)()
			}
			return nil, nil, InvalidOverloadCandidatesError{
				Name:       ocs.name,
				Candidates: ocs.alternatives,
				Context: ocs.context,
			}
		}
	}
	if config.OnPostprocessingFinished != nil {
		defer (*config.OnPostprocessingFinished)()
	}
	infer.cs = cleanCS

	//fmt.Printf("CONSTRAINTS: %v\n", infer.cs)
	for _, ics := range infer.ics {
		//fmt.Printf("GOT INTRO: %v => %v\n", ics.tv, ics.argTV)
		introExpr := (*ics.context.Source).(IntrospectionExpression)
		introExpr.OnTypeReturned(ics.argTV)
	}

	if s.err != nil {
		return nil, nil, s.err
	}

	if infer.t == nil {
		return nil, nil, errors.Errorf("infer.t is nil")
	}

	t := infer.t.Apply(s.sub).(Type)
	ret, err := closeOver(t)
	if err != nil {
		return nil, nil, err
	}

	var retEnv Env = nil
	if infer.retEnv != nil {
		retEnv = infer.retEnv
	} else {
		infer.retEnv = infer.env
	}

	return ret, retEnv, nil
}

// Unify unifies the two types and returns a list of substitutions.
func Unify(a, b Type, context Constraint) (sub Subs, err error) {
	//logf("%v ~ %v", a, b)
	//enterLoggingContext()
	//defer leaveLoggingContext()

	switch at := a.(type) {
	case TypeVariable:
		return bind(at, b, context, a)
	default:
		if a.Eq(b) {
			return nil, nil
		}

		if btv, ok := b.(TypeVariable); ok {
			return bind(btv, a, context, b)
		}
		atypes := a.Types()
		btypes := b.Types()
		defer ReturnTypes(atypes)
		defer ReturnTypes(btypes)

		if len(atypes) == 0 && len(btypes) == 0 {
			goto e
		}

		return unifyMany(atypes, btypes, a, b, context)

	e:
	}
	err = UnificationWrongTypeError{
		TypeA:      a,
		TypeB:      b,
		Constraint: context,
	}
	return
}

func unifyMany(a, b Types, contextA, contextB Type, context Constraint) (sub Subs, err error) {
	//logf("UnifyMany %v %v", a, b)
	//enterLoggingContext()
	//defer leaveLoggingContext()

	if len(a) != len(b) {
		return nil, UnificationLengthError{
			TypeA:      contextA,
			TypeB:      contextB,
			Constraint: context,
		}
	}

	for i, at := range a {
		bt := b[i]

		if sub != nil {
			at = at.Apply(sub).(Type)
			bt = bt.Apply(sub).(Type)
		}

		var s2 Subs
		if s2, err = Unify(at, bt, context); err != nil {
			return nil, err
		}

		if sub == nil {
			sub = s2
		} else {
			sub2 := compose(sub, s2)
			defer ReturnSubs(s2)
			if sub2 != sub {
				defer ReturnSubs(sub)
			}
			sub = sub2
		}
	}
	return
}

func bind(tv TypeVariable, t Type, context Constraint, tvt Type) (sub Subs, err error) {
	//logf("Binding %v to %v", tv, t)
	switch {
	// case tv == t:
	case occurs(tv, t):
		err = UnificationRecurrentTypeError{
			Type:       t,
			Variable:   tv,
			VariableTypeSource: tvt,
			Constraint: context,
		}
	default:
		ssub := BorrowSSubs(1)
		ssub.s[0] = Substitution{tv, t}
		sub = ssub
	}
	//logf("Sub %v", sub)
	return
}

func occurs(tv TypeVariable, s Substitutable) bool {
	ftv := s.FreeTypeVar()
	defer ReturnTypeVarSet(ftv)

	return ftv.Contains(tv)
}

func closeOver(t Type) (sch *Scheme, err error) {
	sch = Generalize(nil, t)
	err = sch.Normalize()
	//logf("closeoversch: %v", sch)
	return
}

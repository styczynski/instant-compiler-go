package hindley_milner

import (
	"fmt"
	"reflect"

	"github.com/alecthomas/repr"
	"github.com/pkg/errors"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type Cloner interface {
	Clone() interface{}
}

type Fresher interface {
	Fresh() TypeVariable
}

func saveExprContext(t Type, source *generic_ast.Expression) Type {
	return t.WithContext(CreateCodeContext(*source))
}

func wrapEnvError(err error, source *generic_ast.Expression) error {
	if e, ok := err.(*envError); ok {
		if e.builtin {
			return BuiltinRedefinedError{
				Name:    e.variableName,
				Context: CreateCodeContext(*source),
			}
		}
		return VariableRedefinedError{
			Name:               e.variableName,
			PreviousDefinition: e.oldDef,
			Context:            CreateCodeContext(*source),
		}
	}
	return err
}

type IntrospectionConstraint struct {
	tv      TypeVariable
	argTV   Type
	context CodeContext
}

type OverloadConstraint struct {
	tv           TypeVariable
	name         string
	alternatives []*Scheme
	context      CodeContext
}

type inferer struct {
	env             Env
	retEnv          Env
	cs              Constraints
	t               Type
	ocs             []OverloadConstraint
	ics             []IntrospectionConstraint
	returns         []Type
	blockScopeLevel int

	count  int
	config *InferConfiguration
	csflag int
}

func newInferer(env Env, config *InferConfiguration) *inferer {
	return &inferer{
		env:    env,
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

func (infer *inferer) lookup(isLiteral bool, name string, source generic_ast.Expression) error {
	t, err, isInstance := infer.env.Lookup(infer, name)
	if err != nil {
		return UndefinedSymbol{
			Name:       name,
			Source:     source,
			IsLiteral:  isLiteral,
			IsVariable: !isLiteral,
		}
	}
	if isInstance {
		t = t.WithContext(CreateCodeContext(source))
	}
	infer.t = t
	return nil
}

func (infer *inferer) resolveProxy(expr generic_ast.Expression, exprType ExpressionType) (generic_ast.Expression, ExpressionType) {

	for {
		if expr == nil {
			return expr, exprType
		}
		if exprType != E_PROXY {
			break
		}
		expr = expr.Body()
		exprType = expr.(HMExpression).ExpressionType()
	}
	return expr, exprType
}

func (infer *inferer) TypeOf(et generic_ast.Expression) (Type, error) {
	defer infer.cleanupConstraints()

	if err := infer.consGen(et, E_NONE, false, false); err != nil {
		return nil, err
	}
	actType := infer.t
	tv := infer.Fresh()
	infer.t = tv
	infer.cs = append(infer.cs, Constraint{
		a:       tv.WithContext(CreateCodeContext(et)),
		b:       actType,
		context: CreateCodeContext(et),
	})

	logf("TYPEOF [%s]: {%v}\n", et, tv)
	return tv, nil
}

func (infer *inferer) consGen(expr generic_ast.Expression, forceType ExpressionType, isTop bool, isOpaqueTop bool) (err error) {

	defer infer.cleanupConstraints()

	if expr == nil {
		tv := infer.Fresh()
		infer.t = tv
		return nil
	}

	exprType := expr.(HMExpression).ExpressionType()
	if forceType != E_NONE {
		exprType = forceType
	}

	expr, exprType = infer.resolveProxy(expr, exprType)

	if exprWithDeps, ok := expr.(ExpressionWithIdentifiersDeps); ok {
		idents := exprWithDeps.GetIdentifierDeps(nil)
		for _, name := range idents.GetNames() {
			if objType := idents.GetTypeOf(name); objType != nil {

				/*_, osx1, osx2 :=*/
				infer.env.AddPrototype(infer, name,
					objType,
					infer.blockScopeLevel)

			} else {
				tv := infer.Fresh().WithContext(objType.t.GetContext())

				_, osx1, osx2, _, err := infer.env.Add(infer, name,
					NewScheme(nil, tv),
					infer.blockScopeLevel, false)
				if err != nil {
					return wrapEnvError(err, &expr)
				}
				if osx1 != nil && osx2 != nil {
					infer.cs = append(infer.cs, Constraint{
						a:       Instantiate(infer, osx1),
						b:       Instantiate(infer, osx2),
						context: CreateCodeContext(expr),
					})
				}
			}
		}
	}

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

		err = nil
	}

	switch exprType {
	case E_CUSTOM:
		et := expr.(CustomExpression)

		context := CustomExpressionEnv{
			Env:            infer.env,
			InferencedType: infer.t,
			LookupEnv: func(isLiteral bool, name string) error {
				return infer.lookup(isLiteral, name, et)
			},
			GenerateConstraints: func(expr generic_ast.Expression) (e error, env Env, i Type, constraints Constraints) {
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
			argTV:   bodyType,
			context: CreateCodeContext(expr),
		})
		infer.cs = append(infer.cs, Constraint{
			a:       tvArg,
			b:       bodyType,
			context: CreateCodeContext(expr),
		})

		return nil

	case E_LITERAL:
		et := expr.(Literal)
		if len(et.Name().GetNames()) != 1 {
			return fmt.Errorf("Literal entity cannot conntain other value than one variable name. You cannot use Names batch here.")
		}
		name := et.Name().GetNames()[0]

		if infer.env.IsOverloaded(name) {
			tv := infer.Fresh().WithContext(CreateCodeContext(expr))
			infer.t = tv
			infer.ocs = append(infer.ocs, OverloadConstraint{
				name:         name,
				tv:           tv.(TypeVariable),
				alternatives: infer.env.OverloadedAlternatives(name),
				context:      CreateCodeContext(expr),
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
			_, s1, s2, _, err := infer.env.Add(infer, name, &Scheme{t: et.Type()}, infer.blockScopeLevel, false)
			if err != nil {
				return wrapEnvError(err, &expr)
			}
			if s1 != nil && s2 != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       Instantiate(infer, s1),
					b:       Instantiate(infer, s2),
					context: CreateCodeContext(expr),
				})
			}
			err = nil
		}

	case E_TYPE:
		et := expr.(EmbeddedType)
		scheme := et.EmbeddedType(infer)

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
				a:       tv.WithContext(CreateCodeContext(expr)),
				b:       Instantiate(infer, defaultTyper.DefaultType(infer)),
				context: CreateCodeContext(expr),
			})
		} else if infer.config.CreateDefaultEmptyType() != nil {
			infer.cs = append(infer.cs, Constraint{
				a:       tv.WithContext(CreateCodeContext(expr)),
				b:       Instantiate(infer, infer.config.CreateDefaultEmptyType()),
				context: CreateCodeContext(expr),
			})
		}
		infer.t = tv

	case E_LAMBDA, E_FUNCTION:
		et := expr.(Lambda)

		env := []Env{}
		tv := []TypeVariable{}

		rets := infer.returns
		infer.returns = []Type{}

		args := et.Args(infer)
		names := args.GetNames()
		for _, name := range names {
			varType := args.GetTypeOf(name)

			newVar := infer.Fresh()
			if varType != nil {
				newVar = newVar.WithContext(varType.t.GetContext()).(TypeVariable)
			}
			tv = append(tv, newVar)

			env = append(env, infer.env)

			infer.env = infer.env.Clone()

			sc := new(Scheme)
			sc.t = tv[len(tv)-1].WithContext(CreateCodeContext(expr))
			_, s1, s2, _, err := infer.env.Add(infer, name, sc, infer.blockScopeLevel, false)
			if err != nil {
				return wrapEnvError(err, &expr)
			}
			if s1 != nil && s2 != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       Instantiate(infer, s1),
					b:       Instantiate(infer, s2),
					context: CreateCodeContext(expr),
				})
			}

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
			bodyType := infer.Fresh().WithContext(CreateCodeContext(et.Body())).(TypeVariable)
			infer.t = bodyType
		} else {
			if err = infer.consGen(et.Body(), E_NONE, false, false); err != nil {
				return err
			}
		}

		for i := 0; i < len(names); i++ {
			infer.t = NewFnType(tv[len(tv)-1], infer.t).WithContext(CreateCodeContext(expr))

			infer.env = env[len(env)-1]
			env = env[:len(env)-1]
			tv = tv[:len(tv)-1]
		}
		if len(names) == 0 {
			infer.t = NewFnType(infer.t).WithContext(CreateCodeContext(expr))
		}

		for _, ret := range infer.returns {
			r := infer.t.(*FunctionType).Ret(true)
			infer.cs = append(infer.cs, Constraint{
				a:       r.WithContext(ret.GetContext()),
				b:       ret.WithContext(CreateCodeContext(expr)),
				context: ret.GetContext(),
			})
		}

		r := infer.t.(*FunctionType).Ret(true)
		if defaultTyper, ok := et.(DefaultTyper); ok {
			infer.cs = append(infer.cs, Constraint{
				a:       r.WithContext(CreateCodeContext(expr)),
				b:       Instantiate(infer, defaultTyper.DefaultType(infer)),
				context: CreateCodeContext(expr),
			})
		}

		if len(infer.returns) == 0 {
			if infer.config.CreateDefaultEmptyType() != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       r.WithContext(CreateCodeContext(expr)),
					b:       Instantiate(infer, infer.config.CreateDefaultEmptyType()),
					context: CreateCodeContext(expr),
				})
			}
		}

		infer.returns = rets

	case E_TYPE_EQUALITY:
		et := expr.Body().(Batch).Expressions()
		if err = infer.consGen(et[0], E_NONE, false, false); err != nil {
			return err
		}
		aType, aCs := infer.t, infer.cs
		if err = infer.consGen(et[1], E_NONE, false, false); err != nil {
			return err
		}
		bType, bCs := infer.t, infer.cs
		cs := append(aCs, bCs...)
		cs = append(cs, Constraint{
			a:       aType,
			b:       bType,
			context: CreateCodeContext(expr),
		})
		infer.cs = cs
		infer.t = aType
		return nil

	case E_APPLICATION:
		et := expr.(Apply)
		firstExec := true
		logf("\n\nAPPLICATION START %v\n", expr)
		batchErr := ApplyBatch(et.Body(), func(body generic_ast.Expression) error {
			if firstExec {
				if err = infer.consGen(et.Fn(infer), E_NONE, false, false); err != nil {
					return err
				}
				firstExec = false
			}
			fnType, fnCs := infer.t, infer.cs

			if err = infer.consGen(body, E_NONE, false, false); err != nil {
				return err
			}
			bodyType, bodyCs := infer.t, infer.cs

			tv := infer.Fresh().WithContext(CreateCodeContext(body)).(TypeVariable)
			applyCs := Constraint{fnType, saveExprContext(NewFnType(bodyType, tv), &expr), CreateCodeContext(expr)}
			cs := append(fnCs, bodyCs...)
			cs = append(cs, applyCs)

			infer.t = tv
			infer.t = saveExprContext(infer.t, &expr)
			infer.cs = cs

			logf("  -> [%v] (%v) bodyType is (%v) and fn type is (%v) --> %v\n", tv, infer.t, bodyType, fnType, applyCs)

			return nil
		})
		logf("\n\nAPPLICATION END %v\n\n", expr)
		if batchErr != nil {
			return batchErr
		}

	case E_BLOCK, E_OPAQUE_BLOCK:
		et := expr.(Block)
		env := infer.env
		if exprType != E_OPAQUE_BLOCK {
			logf("BLOCK_SCOPE level++ (old value: %d)\n", infer.blockScopeLevel)
			infer.blockScopeLevel++
		}

		for _, statement := range et.GetContents().Expressions() {

			if err = infer.consGen(statement, E_NONE, false, exprType == E_OPAQUE_BLOCK && isOpaqueTop); err != nil {
				return err
			}

		}

		tv := infer.Fresh().WithContext(CreateCodeContext(et)).(TypeVariable)
		if defaultTyper, ok := et.(DefaultTyper); ok {
			infer.cs = append(infer.cs, Constraint{
				a:       tv,
				b:       Instantiate(infer, defaultTyper.DefaultType(infer)),
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

		if exprType != E_OPAQUE_BLOCK {
			logf("BLOCK_SCOPE level-- (old value: %d)\n", infer.blockScopeLevel)
			infer.blockScopeLevel--
		}
		infer.t = tv
		if exprType != E_OPAQUE_BLOCK {
			infer.env = env
		}

	case E_LET_RECURSIVE, E_DECLARATION, E_FUNCTION_DECLARATION:

		et := expr.(LetBase)
		vars := et.Var(infer)
		names := vars.GetNames()
		types := []*Scheme{}

		definitions := []generic_ast.Expression{}
		if len(names) == 1 {
			if exprType == E_FUNCTION_DECLARATION {
				definitions = append(definitions, et.(Lambda))
			} else {
				def := et.(Let).Def(infer)
				if batch, ok := def.(Batch); ok {
					definitions = append(definitions, batch.Expressions()[0])
				} else {
					definitions = append(definitions, def)
				}
				if vars.HasTypes() {
					types = append(types, vars.GetTypeOf(names[0]))
				}
			}
		} else if len(names) > 1 {
			if exprType == E_FUNCTION_DECLARATION {
				for range names {
					definitions = append(definitions, et.(Lambda))
				}
			} else {
				def := et.(Let).Def(infer)
				for i, expr := range def.(Batch).Expressions() {
					definitions = append(definitions, expr)
					if vars.HasTypes() {
						types = append(types, vars.GetTypeOf(names[i]))
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
			tv := infer.Fresh().WithContext(CreateCodeContext(expr)).(TypeVariable)

			var defExpectedType *Scheme = nil
			if len(types) > 0 {
				defExpectedType = types[i]
			}

			infer.env = infer.env.Clone()
			has := infer.env.Has(name)
			if !has {
				infer.env.Remove(name)
			}
			_, s1, s2, varEnvDef, err := infer.env.Add(infer, name, &Scheme{tvs: TypeVarSet{tv}, t: tv}, infer.blockScopeLevel, false)
			if err != nil {
				return wrapEnvError(err, &expr)
			}
			if s1 != nil && s2 != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       Instantiate(infer, s1),
					b:       Instantiate(infer, s2),
					context: CreateCodeContext(expr),
				})
			}

			nonVal := false
			defResolved, _ := infer.resolveProxy(def, def.(HMExpression).ExpressionType())
			if block, ok := defResolved.(Block); ok && exprType == E_DECLARATION {
				if len(block.GetContents().Expressions()) == 0 {
					nonVal = true
				}
			}

			if nonVal {
				if defExpectedType != nil {
					infer.t = Instantiate(infer, defExpectedType)
				} else {
					tv := infer.Fresh().WithContext(CreateCodeContext(expr)).(TypeVariable)
					infer.t = tv
				}
			} else if exprType == E_FUNCTION_DECLARATION {

				if err = infer.consGen(def, E_FUNCTION, false, false); err != nil {
					return err
				}
			} else {

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
			logf("\nDefinition type [%s]: %v\n", name, saveExprContext(defType.Apply(s.sub).(Type), &expr))
			//Instantiate(infer, defExpectedType)
			logf("\n |-> Expected type [%s]: %v\n", name, defExpectedType)

			sc := Generalize(infer.env.Apply(s.sub).(Env), saveExprContext(defType.Apply(s.sub).(Type), &expr))

			if !has {
				infer.env.Remove(name)
			}
			_, os1, os2, varEnvDef2, err := infer.env.Add(infer, name, sc, infer.blockScopeLevel, false)
			if err != nil {
				if _, isDup := err.(*envError); isDup && varEnvDef.GetUID() == varEnvDef2.GetUID() {

				} else {
					return wrapEnvError(err, &expr)
				}
			}
			if os1 != nil && os2 != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       Instantiate(infer, os1),
					b:       Instantiate(infer, os2),
					context: CreateCodeContext(expr),
				})
			}

			if exprType == E_DECLARATION || exprType == E_FUNCTION_DECLARATION {
				retType := infer.Fresh().WithContext(CreateCodeContext(expr)).(TypeVariable)
				if defaultTyper, ok := expr.(DefaultTyper); ok {
					infer.cs = append(infer.cs, Constraint{
						a:       retType,
						b:       Instantiate(infer, defaultTyper.DefaultType(infer)),
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

			if defExpectedType != nil {
				actualType := defType
				if exprType == E_LET_RECURSIVE {
					actualType = sc.t
				}

				infer.cs = append(infer.cs, Constraint{
					a:       actualType,
					b:       Instantiate(infer, defExpectedType),
					context: CreateCodeContext(expr),
				})
			}
		}

	case E_LET, E_REDEFINABLE_LET:
		et := expr.(Let)
		vars := et.Var(infer)

		if len(vars.GetNames()) != 1 {
			return fmt.Errorf("Let entity cannot conntain other value than one variable name. You cannot use Names batch here.")
		}
		name := vars.GetNames()[0]

		env := infer.env

		if err = infer.consGen(et.Def(infer), E_NONE, false, false); err != nil {
			return err
		}
		defType, defCs := infer.t, infer.cs

		s := newSolver()
		s.solve(defCs)
		if s.err != nil {
			return err
		}

		logf("PATRZ CWELU: %v\n", defType.Apply(s.sub).(Type))
		sc := Generalize(env.Apply(s.sub).(Env), saveExprContext(defType.Apply(s.sub).(Type), &expr))
		infer.env = infer.env.Clone()

		_, s1, s2, _, err := infer.env.Add(infer, name, sc, infer.blockScopeLevel, exprType == E_REDEFINABLE_LET)
		if err != nil {
			return wrapEnvError(err, &expr)
		}
		if s1 != nil && s2 != nil {
			infer.cs = append(infer.cs, Constraint{
				a:       Instantiate(infer, s1),
				b:       Instantiate(infer, s2),
				context: CreateCodeContext(expr),
			})
		}

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
	order := []string{}
	for _, cs := range inferer.cs {
		key := fmt.Sprintf("%v", cs)
		if entry, ok := hashtable[key]; ok {
			if entry.context.Source == nil && cs.context.Source != nil {
				hashtable[key] = cs
			}
		} else {
			hashtable[key] = cs
			order = append(order, key)
		}
	}
	inferer.cs = []Constraint{}
	for _, key := range order {
		inferer.cs = append(inferer.cs, hashtable[key])
	}
}

func (infer *inferer) cleanupConstraints() {
	infer.cleanupConstraintsRemoveDuplicates()
	return

	if infer.csflag >= 1 {
		infer.csflag = 0
	} else {
		infer.csflag = infer.csflag + 1
		return
	}

	cs := Constraints{}
	freeVars := map[int16]map[Type]interface{}{}
	contexts := map[int16]CodeContext{}
	for _, cons := range infer.cs {
		if !cons.a.Eq(cons.b) {
			if tv, ok := cons.a.(TypeVariable); ok {
				if _, has := freeVars[tv.value]; !has {
					freeVars[tv.value] = map[Type]interface{}{}
					contexts[tv.value] = tv.context
				} else {
					if contexts[tv.value].Source == nil && tv.context.Source != nil {
						contexts[tv.value] = tv.context
					}
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
				a: TypeVariable{
					value:   id,
					context: context,
				},
				b:       b,
				context: context,
			})
		}
	}

	infer.cs = cs

}

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

	case len(envFree) == 0 && len(tFree) > 0:

	}

	diff = tFree.Difference(envFree)

ret:
	return &Scheme{
		tvs: diff,
		t:   t,
	}
}

type InferConfiguration struct {
	CreateDefaultEmptyType        func() *Scheme
	OnConstrintGenerationStarted  *func()
	OnConstrintGenerationFinished *func()
	OnSolvingStarted              *func()
	OnSolvingFinished             *func()
	OnPostprocessingStarted       *func()
	OnPostprocessingFinished      *func()
}

func NewInferConfiguration() *InferConfiguration {
	return &InferConfiguration{
		CreateDefaultEmptyType: func() *Scheme { return nil },
	}
}

func Infer(env Env, expr generic_ast.Expression, config *InferConfiguration) (*Scheme, Env, error) {
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

	if config.OnPostprocessingStarted != nil {
		(*config.OnPostprocessingStarted)()
	}
	cleanCS := infer.cs
	for _, ocs := range infer.ocs {
		hasCleanRun := false
		cs := Constraints{}
		for _, alt := range ocs.alternatives {
			cs = Constraints{}

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
				Context:    ocs.context,
			}
		}
	}
	if config.OnPostprocessingFinished != nil {
		defer (*config.OnPostprocessingFinished)()
	}
	infer.cs = cleanCS

	for _, ics := range infer.ics {

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

func Unify(a, b Type, context Constraint) (sub Subs, err error) {

	if sa, ok := a.(UnionTypeCheck); ok {
		if reflect.TypeOf(a) == reflect.TypeOf(b) {
			err := sa.CheckIfCanUnionTypes(b)
			if err != nil {
				return nil, UnificationWrongTypeError{
					TypeA:      a,
					TypeB:      b,
					Constraint: context,
					Details:    err.Error(),
				}
			}
		}
	}

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
	logf("Binding %v to %v", tv, t)
	switch {

	case occurs(tv, t):
		if tv.Eq(t) {
			ssub := BorrowSSubs(1)
			ssub.s[0] = Substitution{tv, t}
			sub = ssub
			return
		}
		logf("KURWA 1\n%v\nKURWA 2\n%v\nKURWA 3\n%v\nKURWA 4\n%v\n",
			repr.String(*t.GetContext().Source),
			repr.String(*tv.GetContext().Source),
			repr.String(*tvt.GetContext().Source),
			repr.String(*context.context.Source))

		err = UnificationRecurrentTypeError{
			Type:               t,
			Variable:           tv,
			VariableTypeSource: tvt,
			Constraint:         context,
		}
	default:
		ssub := BorrowSSubs(1)
		ssub.s[0] = Substitution{tv, t}
		sub = ssub
	}

	return
}

func occurs(tv TypeVariable, s Substitutable) bool {
	ftv := s.FreeTypeVar()
	defer ReturnTypeVarSet(ftv)

	return ftv.Contains(tv)
}

func closeOver(t Type) (sch *Scheme, err error) {
	logf("PATRZ SPERMA: %v\n", t)
	sch = Generalize(nil, t)
	err = sch.Normalize()

	return
}

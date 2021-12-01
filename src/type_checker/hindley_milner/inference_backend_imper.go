package hindley_milner

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type ImperInferenceBackend struct {
	env             Env
	retEnv          Env
	t               Type
	cs              Constraints
	returns         []Type
	callQueue       []generic_ast.Expression
	blockScopeLevel int
	config          *InferConfiguration
	count           int
}

func CreateImperInferenceBackend(env Env, config *InferConfiguration) *ImperInferenceBackend {
	return &ImperInferenceBackend{
		env:             env,
		config:          config,
		returns:         []Type{},
		cs:              Constraints{},
		callQueue:       []generic_ast.Expression{},
		t:               nil,
		retEnv:          env,
		blockScopeLevel: 0,
		count:           0,
	}
}

func (infer *ImperInferenceBackend) Fresh() TypeVariable {
	retVal := infer.count
	infer.count++
	return TypeVariable{
		value: int16(retVal),
	}
}

func (infer *ImperInferenceBackend) lookup(isLiteral bool, name string, source generic_ast.Expression) error {
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

func (infer *ImperInferenceBackend) resolveProxy(expr generic_ast.Expression, exprType ExpressionType) (generic_ast.Expression, ExpressionType) {
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

func (infer *ImperInferenceBackend) ProgramType() Type {
	return infer.t
}

func (infer *ImperInferenceBackend) GetEnv() Env {
	return infer.env
}

func (infer *ImperInferenceBackend) GetReturnEnv() Env {
	if infer.retEnv != nil {
		return infer.retEnv
	} else {
		infer.retEnv = infer.env
		return infer.env
	}
}

func (infer *ImperInferenceBackend) Constraints() Constraints {
	return infer.cs
}

func (infer *ImperInferenceBackend) OverrideConstraints(cs Constraints) {
	infer.cs = cs
}

func (infer *ImperInferenceBackend) GetOverloadConstraints() []OverloadConstraint {
	return []OverloadConstraint{}
}

func (infer *ImperInferenceBackend) GetIntrospectionConstraints() []IntrospectionConstraint {
	return []IntrospectionConstraint{}
}

func (infer *ImperInferenceBackend) TypeOf(et generic_ast.Expression, contextExpressions ...generic_ast.Expression) (Type, error) {
	return nil, nil
}

func (infer *ImperInferenceBackend) GenerateConstraints(expr generic_ast.Expression, forceType ExpressionType, isTop bool, isOpaqueTop bool) (err error) {

	defer func() {
		if err != nil {
			return
		}
		newCS := Constraints{}
		allSubs, _ := SubsDisjointConcat()
		for _, cs := range infer.cs {

			// Problem?
			if cs.FreeTypeVar().Len() > 0 {
				subs, err0 := Unify(cs.a, cs.b, cs, infer.env.GetIntrospecionListener())
				if err0 != nil {
					err = err0
					return
				}
				subsOk := true
				if allSubs, subsOk = SubsDisjointConcat(allSubs, subs); !subsOk {
					err = UnificationWrongTypeError{
						TypeA:      cs.a,
						TypeB:      cs.b,
						Constraint: cs,
						Details:    fmt.Sprintf("Polymorphic function cannot be bound to concrete types"),
					}
					return
				}
				//infer.t = infer.t.Apply(subs).(Type)
			} else if !TypeEq(cs.a, cs.b) {
				newCS = append(newCS, cs)
			}
		}
		infer.t = infer.t.Apply(allSubs).(Type)
		infer.cs = newCS
		for _, cs := range infer.cs {
			if cs.FreeTypeVar().Len() > 0 {
				err = fmt.Errorf("Output constraints contains type variables")
				return
			}
		}
		if infer.t.GetContext().Source == nil {
			for i := len(infer.callQueue) - 1; i >= 0; i-- {
				queueItem := infer.callQueue[i]
				ok := false
				for s := 0; s < 10; s++ {
					if b, ok := queueItem.(Batch); ok {
						if len(b.Exp) == 0 {
							break
						}
						queueItem = b.Exp[0]
					} else {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				infer.t = infer.t.WithContext(CreateCodeContext(queueItem))
				break
			}
		}
		infer.callQueue = infer.callQueue[:len(infer.callQueue)-1]
		logf(">> [END] Generate constrints: %v %v (%v)\n", reflect.TypeOf(expr), infer.cs, infer.t)
	}()

	infer.callQueue = append(infer.callQueue, expr)

	if expr == nil {
		return nil
	}

	exprType := expr.(HMExpression).ExpressionType()
	if forceType != E_NONE {
		exprType = forceType
	}

	expr, exprType = infer.resolveProxy(expr, exprType)

	if exprWithDeps, ok := expr.(ExpressionWithIdentifiersDeps); ok {
		err, idents := exprWithDeps.GetIdentifierDeps(infer, false)
		if err != nil {
			return err
		}
		for _, name := range idents.GetNames() {
			if objType := idents.GetTypeOf(name); objType != nil {
				infer.env.AddPrototype(infer, name,
					objType,
					infer.blockScopeLevel)
			} else {
				return fmt.Errorf("GetIdentifierDeps returned variable with no type.")
			}
		}
	}

	// fmt.Printf("THE ENV:")
	// PrintEnv(infer.env)
	// fmt.Printf("END")

	logf(">> Generate constrints %v (%v)\n", reflect.TypeOf(expr), exprType)

	switch et := expr.(type) {
	case Typer:
		if infer.t = et.Type(); infer.t != nil {
			infer.t = saveExprContext(infer.t, &expr)
			logf("TYPED = %v\n", infer.t)
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
				if err = infer.GenerateConstraints(et.Body(), E_NONE, false, false); err != nil {
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
		//et := expr.(IntrospectionExpression)

		// TODO: Implement

		return nil

	case E_LITERAL:
		et := expr.(Literal)
		if len(et.Name().GetNames()) != 1 {
			return fmt.Errorf("Literal entity cannot conntain other value than one variable name. You cannot use Names batch here.")
		}
		name := et.Name().GetNames()[0]
		// if infer.env.IsOverloaded(name) {
		// 	alternatives := infer.env.OverloadedAlternatives(name)
		// 	types := []Type{}
		// 	for _, alt := range alternatives {
		// 		types = append(types, alt.Concrete())
		// 	}
		// 	infer.t = NewUnionType(types)
		// 	return nil
		// }
		return infer.lookup(true, name, et)

	case E_VAR:
		et := expr.(Var)
		if len(et.Name().GetNames()) != 1 {
			return fmt.Errorf("Var entity cannot conntain other value than one variable name. You cannot use Names batch here.")
		}
		name := et.Name().GetNames()[0]
		if err = infer.lookup(false, name, et); err != nil {
			_, _, _, _, err := infer.env.Add(infer, name, Concreate(et.Type()), infer.blockScopeLevel, false)
			if err != nil {
				return wrapEnvError(err, &expr)
			}
			err = nil
		}

	case E_TYPE:
		et := expr.(EmbeddedType)
		scheme := et.EmbeddedType(infer)
		infer.t = Instantiate(infer, scheme)

	case E_RETURN:
		et := expr.(Return)
		if err = infer.GenerateConstraints(et.Body(), E_NONE, false, false); err != nil {
			return err
		}
		infer.t = infer.t.WithContext(CreateCodeContext(et))
		infer.returns = append(infer.returns, infer.t)

		if defaultTyper, ok := et.(DefaultTyper); ok {
			infer.t = Instantiate(infer, defaultTyper.DefaultType(infer))
		} else if infer.config.CreateDefaultEmptyType() != nil {
			infer.t = Instantiate(infer, infer.config.CreateDefaultEmptyType())
		}

	case E_LAMBDA, E_FUNCTION:
		et := expr.(Lambda)

		env := []Env{}
		rets := infer.returns
		types := []*Scheme{}
		infer.returns = []Type{}

		args := et.Args(infer)
		names := args.GetNames()
		logf("FUC A: %v\n", args)
		for _, name := range names {
			logf("FUC B: %v\n", name)
			varType := args.GetTypeOf(name)
			logf("FUC B2: %v\n", varType)

			sc := varType
			types = append(types, sc)
			//newVar = newVar.WithContext(newVar.GetContext()).(TypeVariable)
			env = append(env, infer.env)

			infer.env = infer.env.Clone()
			logf("FUC C\n")
			_, s1, s2, _, err := infer.env.Add(infer, name, sc, infer.blockScopeLevel, false)
			if err != nil {
				return wrapEnvError(err, &expr)
			}
			logf("FUC D\n")
			if s1 != nil && s2 != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       Instantiate(infer, s1),
					b:       Instantiate(infer, s2),
					context: CreateCodeContext(expr),
				})
			}

			logf("FUC E\n")
			if varType != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       sc.t,
					b:       Instantiate(infer, varType),
					context: CreateCodeContext(expr),
				})
			}
		}

		if exprType == E_FUNCTION {
			logf("FUC F\n")
			if err = infer.GenerateConstraints(et.Body(), E_NONE, false, false); err != nil {
				return err
			}
		} else {
			if err = infer.GenerateConstraints(et.Body(), E_NONE, false, false); err != nil {
				return err
			}
		}
		logf("FUC G\n")

		// Add return type
		if true {
			if defaultTyper, ok := et.(DefaultTyper); ok {
				rT := defaultTyper.DefaultType(infer)
				infer.t = rT.Concrete()
			} else if rT := infer.config.CreateDefaultEmptyType(); rT != nil {
				infer.t = rT.Concrete()
			}
		}

		for i := 0; i < len(names); i++ {
			infer.t = NewFnType(types[len(types)-1].Concrete(), infer.t).WithContext(CreateCodeContext(expr))

			infer.env = env[len(env)-1]
			env = env[:len(env)-1]
			types = types[:len(types)-1]
		}

		logf("FUC H => %v\n", infer.t)
		if len(names) == 0 {
			infer.t = NewFnType(infer.t).WithContext(CreateCodeContext(expr))
		}
		logf("FUC H0 => <%v> %v\n", et.Args(nil), infer.t)

		for _, ret := range infer.returns {
			r := infer.t.(*FunctionType).Ret(true)
			logf("RETURN CONSTRAINT HERE: %v ~ %v\n", r, ret)
			infer.cs = append(infer.cs, Constraint{
				a:       r.WithContext(ret.GetContext()),
				b:       ret.WithContext(CreateCodeContext(expr)),
				context: ret.GetContext(),
			})
		}
		logf("FUC I\n")

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
		logf("FUC J\n")

		infer.returns = rets

	case E_TYPE_EQUALITY:
		et := expr.Body().(Batch).Expressions()
		if err = infer.GenerateConstraints(et[0], E_NONE, false, false); err != nil {
			return err
		}
		aType, aCs := infer.t, infer.cs
		if err = infer.GenerateConstraints(et[1], E_NONE, false, false); err != nil {
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
		if defaultTyper, ok := expr.(DefaultTyper); ok {
			infer.t = Instantiate(infer, defaultTyper.DefaultType(infer))
		} else if infer.config.CreateDefaultEmptyType() != nil {
			infer.t = Instantiate(infer, infer.config.CreateDefaultEmptyType())
		}
		return nil

	case E_APPLICATION:
		et := expr.(Apply)
		firstExec := true

		// collect all argument types
		argTypes := []Type{}
		batchErr := ApplyBatch(et.Body(), func(body generic_ast.Expression) error {
			if err = infer.GenerateConstraints(body, E_NONE, false, false); err != nil {
				return err
			}
			argTypes = append(argTypes, infer.t)
			return nil
		})
		if batchErr != nil {
			return batchErr
		}

		logf("\n\nAPPLICATION START %v\n", expr)
		var originalFnType *FunctionType
		index := 0
		//argCS := infer.cs
		batchErr = ApplyBatch(et.Body(), func(body generic_ast.Expression) error {
			if firstExec {
				if err = infer.GenerateConstraints(et.Fn(infer), E_NONE, false, false); err != nil {
					return err
				}
				if unionType, ok := infer.t.(*Union); ok {
					allArgs := []Type{}
					allArgs = append(allArgs, argTypes...)
					allArgs = append(allArgs, TVar(0))
					err, resolvedType := unionType.FindMatchingFunction(et, argTypes)
					fmt.Printf("UNIONIZED %v FOR %v %v\n", err, resolvedType, argTypes)
					if err != nil {
						return err
					}
					if resolvedType == nil {
						return UnificationWrongTypeError{
							TypeA: infer.t,
							TypeB: NewFnType(allArgs...),
							Constraint: Constraint{
								a:       infer.t,
								b:       NewFnType(allArgs...).WithContext(CreateCodeContext(et)),
								context: CreateCodeContext(et),
							},
							Details: fmt.Sprintf("Function overload not found"),
						}
					}
					infer.t = resolvedType
				}
				if _, ok := infer.t.(*FunctionType); !ok {
					return UnificationWrongTypeError{
						TypeA: infer.t,
						TypeB: NewFnType(TVar(0), TVar(1)).WithContext(CreateCodeContext(et)),
						Constraint: Constraint{
							a:       infer.t,
							b:       NewFnType(TVar(0), TVar(1)).WithContext(CreateCodeContext(et)),
							context: CreateCodeContext(et),
						},
						Details: fmt.Sprintf("Value is not a function. You can call only functions. Got type: %v", infer.t),
					}
				}
				originalFnType = infer.t.(*FunctionType)
				firstExec = false
			}
			fnType := infer.t
			bodyType := infer.t
			if false {
				if err = infer.GenerateConstraints(body, E_NONE, false, false); err != nil {
					return err
				}
			} else {
				bodyType = argTypes[index]
			}

			if _, ok := fnType.(*FunctionType); !ok {
				return UnificationWrongTypeError{
					TypeA: fnType,
					TypeB: NewFnType(TVar(0), TVar(1)),
					Constraint: Constraint{
						a:       fnType,
						b:       NewFnType(TVar(0), TVar(1)),
						context: CreateCodeContext(et),
					},
					Details: fmt.Sprintf("Function is applied to wrong number of arguments. Expected: %d", originalFnType.CountArgs()),
				}
			}

			expectedType := NewFnType(
				bodyType,
				fnType.(*FunctionType).Ret(false),
			)
			applyCs := Constraint{fnType, saveExprContext(expectedType, &expr), CreateCodeContext(expr)}
			infer.cs = append(infer.cs, applyCs)

			logf("FUNCTION RETURN DEBUG: %v from %v\n", fnType, expectedType)

			infer.t = fnType.(*FunctionType).Ret(false)
			infer.t = saveExprContext(infer.t, &expr)

			logf(" KURWAAA -> [%v] (%v) bodyType is (%v) and fn type is (%v) --> %v\n", nil, infer.t, bodyType, fnType, applyCs)

			index++
			return nil
		})
		logf("\n\nAPPLICATION END %v WITH %v AND %v\n\n", expr, infer.cs, infer.t)
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
			if err = infer.GenerateConstraints(statement, E_NONE, false, exprType == E_OPAQUE_BLOCK && isOpaqueTop); err != nil {
				return err
			}
		}

		if defaultTyper, ok := et.(DefaultTyper); ok {
			infer.t = Instantiate(infer, defaultTyper.DefaultType(infer))
		} else if infer.config.CreateDefaultEmptyType() != nil {
			infer.t = Instantiate(infer, infer.config.CreateDefaultEmptyType())
		}

		if isTop || isOpaqueTop {
			infer.retEnv = infer.env.Clone()
		}

		if exprType != E_OPAQUE_BLOCK {
			logf("BLOCK_SCOPE level-- (old value: %d)\n", infer.blockScopeLevel)
			infer.blockScopeLevel--
		}
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
				if vars.HasTypes() {
					types = append(types, vars.GetTypeOf(names[0]))
				}
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
				for i, _ := range names {
					definitions = append(definitions, et.(Lambda))
					if vars.HasTypes() {
						types = append(types, vars.GetTypeOf(names[i]))
					}
				}
			} else {
				def := et.(Let).Def(infer)
				for i, expr := range def.(Batch).Expressions() {
					definitions = append(definitions, expr)
					if vars.HasTypes() {
						types = append(types, vars.GetTypeOf(names[i]))
					}
				}
				fmt.Printf("DEFSGEJ: %v\n", vars)
			}
		} else {
			panic("Invalid number of identifiers returned by Var() of the declaration/let: zero.")
		}

		logf("HERE A\n")

		if len(types) != len(names) {
			types = []*Scheme{}
		}

		for i, _ := range names {
			logf("HERE B\n")

			redef := false

			name := names[i]
			def := definitions[i]
			body := expr.Body()

			var defExpectedType *Scheme = nil
			if len(types) > 0 {
				defExpectedType = types[i]
			}

			if defExpectedType == nil {
				return fmt.Errorf("Missing expected vars type for function parameter.")
			}

			logf("HERE B2: %v\n", defExpectedType.t)
			tv := defExpectedType.t.WithContext(CreateCodeContext(expr))
			logf("HERE B2 NEXT\n")

			infer.env = infer.env.Clone()
			logf("HERE B3\n")
			has := infer.env.Has(name)
			if !has {
				logf("HERE B4 %v\n", infer.cs)
				infer.env.Remove(name)
			}
			logf("HERE C %v\n", infer.cs)
			_, s1, s2, varEnvDef, err := infer.env.Add(infer, name, Concreate(tv), infer.blockScopeLevel, redef)
			if err != nil {
				return wrapEnvError(err, &expr)
			}
			logf("HERE C2 %v\n", infer.cs)
			if s1 != nil && s2 != nil {
				infer.cs = append(infer.cs, Constraint{
					a:       Instantiate(infer, s1),
					b:       Instantiate(infer, s2),
					context: CreateCodeContext(expr),
				})
			}

			logf("HERE D %v\n", infer.cs)

			nonVal := false
			defResolved, _ := infer.resolveProxy(def, def.(HMExpression).ExpressionType())
			if block, ok := defResolved.(Block); ok && exprType == E_DECLARATION {
				if len(block.GetContents().Expressions()) == 0 {
					nonVal = true
				}
			}

			logf("HERE E %v\n", infer.cs)
			if nonVal {
				if defExpectedType != nil {
					fmt.Printf("SKURWYSYSYN %v\n", defExpectedType)
					infer.t = Instantiate(infer, defExpectedType)
				} else {
					return fmt.Errorf("Expected concrete type in function node")
				}
			} else if exprType == E_FUNCTION_DECLARATION {
				if err = infer.GenerateConstraints(def, E_FUNCTION, false, false); err != nil {
					return err
				}
			} else {
				logf("HELLO F! ==> %v <%v>\n", def, reflect.TypeOf(et))
				if err = infer.GenerateConstraints(def, E_NONE, false, false); err != nil {
					return err
				}
			}
			logf("HERE F %v => %v <%s?> context {\n%v\n}\n", infer.cs, infer.t, name, et)
			defType, defCs := infer.t, infer.cs

			s := newSolver()
			s.solve(defCs, infer.env.GetIntrospecionListener())
			fmt.Printf("HERE F (SOLVED)\n")
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
			logf("HERE G %v <%s?> with new def: %v\n", infer.cs, name, sc)
			_, os1, os2, varEnvDef2, err := infer.env.Add(infer, name, sc, infer.blockScopeLevel, redef)
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
				if defaultTyper, ok := expr.(DefaultTyper); ok {
					infer.t = Instantiate(infer, defaultTyper.DefaultType(infer))
				} else if infer.config.CreateDefaultEmptyType() != nil {
					infer.t = Instantiate(infer, infer.config.CreateDefaultEmptyType())
				}
			} else {
				if err = infer.GenerateConstraints(body, E_NONE, false, false); err != nil {
					return err
				}
				infer.t = infer.t.Apply(s.sub).(Type)
				infer.t = saveExprContext(infer.t, &expr)
				infer.cs = infer.cs.Apply(s.sub).(Constraints)
			}
			logf("HERE H %v\n", infer.cs)

			infer.cs = append(infer.cs, defCs...)

			if defExpectedType != nil {
				actualType := defType
				if exprType == E_LET_RECURSIVE {
					actualType = sc.t
				}

				logf("CHECK RETURN: %v ~ %v\n", actualType, defExpectedType)
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

		if err = infer.GenerateConstraints(et.Def(infer), E_NONE, false, false); err != nil {
			return err
		}
		defType, defCs := infer.t, infer.cs

		s := newSolver()
		s.solve(defCs, infer.env.GetIntrospecionListener())
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

		if err = infer.GenerateConstraints(et.Body(), E_NONE, false, false); err != nil {
			return err
		}

		infer.t = infer.t.Apply(s.sub).(Type)
		infer.t = saveExprContext(infer.t, &expr)
		infer.cs = infer.cs.Apply(s.sub).(Constraints)
		infer.cs = append(infer.cs, defCs...)

	default:
		return errors.Errorf("Expression of %T is unhandled", expr)
	}

	//infer.t = defaultEmptyType
	return nil
}

package hindley_milner

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type ImperInferenceBackend struct {
	env             Env
	retEnv          Env
	t               Type
	blockScopeLevel int
	config          *InferConfiguration
}

func CreateImperInferenceBackend(env Env, config *InferConfiguration) *ImperInferenceBackend {
	return &ImperInferenceBackend{
		env:    env,
		config: config,
	}
}

func (infer *ImperInferenceBackend) Fresh() TypeVariable {
	return TVar(0)
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
	return Constraints{}
}

func (infer *ImperInferenceBackend) OverrideConstraints(cs Constraints) {
	// Nothing
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
	defaultEmptyType := TypeConst{
		value:   "void",
		context: CreateCodeContext(expr),
	}

	if expr == nil {
		infer.t = defaultEmptyType
		return nil
	}

	exprType := expr.(HMExpression).ExpressionType()
	if forceType != E_NONE {
		exprType = forceType
	}

	expr, exprType = infer.resolveProxy(expr, exprType)

	if exprWithDeps, ok := expr.(ExpressionWithIdentifiersDeps); ok {
		idents := exprWithDeps.GetIdentifierDeps(infer, false)
		for _, name := range idents.GetNames() {
			if objType := idents.GetTypeOf(name); objType != nil {
				infer.env.Add(infer, name,
					objType,
					infer.blockScopeLevel, false)
			} else {
				return fmt.Errorf("GetIdentifierDeps returned variable with no type.")
			}
		}
	}

	fmt.Printf("THE ENV:")
	PrintEnv(infer.env)
	fmt.Printf("END")

	infer.t = defaultEmptyType
	return nil
}

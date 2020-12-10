package type_checker

import (
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LatteTypeChecker struct {}

func CreateLatteTypeChecker() *LatteTypeChecker {
	return &LatteTypeChecker{}
}

func (tc *LatteTypeChecker) Test(c *context.ParsingContext) {
	Example_greenspun()
}

func (tc *LatteTypeChecker) GetEnv() hindley_milner.SimpleEnv {
	return hindley_milner.CreateSimpleEnv(map[string]*hindley_milner.Scheme{
		"||":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"&&":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"+":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"-":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"/":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"*":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"!":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"<=":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		">=":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"==":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"!=":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"<":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		">":      hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
	})
}

func (tc *LatteTypeChecker) Check(program *ast.LatteProgram, c *context.ParsingContext) {
	var scheme *hindley_milner.Scheme
	var err error
	var retEnv hindley_milner.Env

	config := hindley_milner.NewInferConfiguration()
	config.CreateDefaultEmptyType = func() *hindley_milner.Scheme { return hindley_milner.NewScheme(nil, ast.CreatePrimitive(ast.T_VOID)) }

	scheme, retEnv, err = hindley_milner.Infer(tc.GetEnv(), program, config)
	if err != nil {
		log.Printf("%+v", errors.Cause(err))
	}
	simpleType, ok := scheme.Type()
	fmt.Printf("simple Type: %v | isMonoType: %v | err: %v\n", simpleType, ok, err)
	hindley_milner.PrintEnv(retEnv)
}


/*
var scheme *hindley_milner.Scheme
	var err error
	config := hindley_milner.NewInferConfiguration()
	config.CreateDefaultEmptyType = func() *hindley_milner.Scheme { return hindley_milner.NewScheme(nil, Prim(Void)) }

	scheme, err = hindley_milner.Infer(env, fac, config)
	if err != nil {
		log.Printf("%+v", errors.Cause(err))
	}
	simpleType, ok := scheme.Type()
	fmt.Printf("simple Type: %v | isMonoType: %v | err: %v\n", simpleType, ok, err)
 */
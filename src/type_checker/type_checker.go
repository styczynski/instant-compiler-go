package type_checker

import (
	"fmt"
	"runtime/debug"

	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LatteTypeChecker struct {}

func CreateLatteTypeChecker() *LatteTypeChecker {
	return &LatteTypeChecker{}
}

func (tc *LatteTypeChecker) Test(c *context.ParsingContext) {
	// Nothing
}

func (tc *LatteTypeChecker) GetEnv() hindley_milner.SimpleEnv {
	return hindley_milner.CreateSimpleEnv(map[string][]*hindley_milner.Scheme{
		"+":  []*hindley_milner.Scheme{
			hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
				ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
			)),
			hindley_milner.NewScheme(nil, hindley_milner.NewFnType(
				ast.CreatePrimitive(ast.T_STRING), ast.CreatePrimitive(ast.T_STRING), ast.CreatePrimitive(ast.T_STRING),
			)),
		},
		"true":      hindley_milner.SingleDef(nil, ast.CreatePrimitive(ast.T_BOOL)),
		"false":      hindley_milner.SingleDef(nil, ast.CreatePrimitive(ast.T_BOOL)),
		"||":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"&&":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		//"+":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
		//	ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		//)),
		"-":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"/":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"*":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"!":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"--":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"++":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"<=":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		">=":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"==":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"!=":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"<":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		">":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"printInt":      hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_VOID),
		)),
		"=":  hindley_milner.SingleDef(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(hindley_milner.TVar('a'), hindley_milner.TVar('a'), hindley_milner.TVar('a'))),
		"if":  hindley_milner.SingleDef(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a'), hindley_milner.TVar('b')},
			hindley_milner.NewFnType(ast.CreatePrimitive(ast.T_BOOL), hindley_milner.TVar('a'), hindley_milner.TVar('b'), ast.CreatePrimitive(ast.T_VOID))),
		"while":  hindley_milner.SingleDef(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(ast.CreatePrimitive(ast.T_BOOL), hindley_milner.TVar('a'), ast.CreatePrimitive(ast.T_VOID))),
		"[]":  hindley_milner.SingleDef(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(
				hindley_milner.NewSignedTupleType("array", hindley_milner.TVar('a')),
				ast.CreatePrimitive(ast.T_INT),
				hindley_milner.TVar('a'),
			)),
		"[_]": hindley_milner.SingleDef(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(
				hindley_milner.NewSignedTupleType("array", hindley_milner.TVar('a')),
				hindley_milner.TVar('a'),
			)),
	})
}

type TypeCheckingError struct {
	message string
	textMessage string
}

func (e TypeCheckingError) Error() string {
    return e.message
}

func (e TypeCheckingError) CliMessage() string {
	return e.textMessage
}

func wrapTypeCheckingError(err error, c *context.ParsingContext) error {
	if undef, ok := err.(hindley_milner.UndefinedSymbol); ok {
		src := undef.Source.(interface{}).(ast.NodeWithPosition)
		message, textMessage := c.FormatParsingError(
			"Unknown Symbol",
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			undef.Error(),
			)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
		}
	} else if wrongType, ok := err.(hindley_milner.UnificationWrongTypeError); ok {
		src := wrongType.Source().(interface{}).(ast.NodeWithPosition)
		causeInfo := ""

		if wrongType.IsCausedByBuiltin() {
			causeInfo = fmt.Sprintf("Caused by internal definition: %s", wrongType.GetCauseName())
		} else {
			//sourceA := (*wrongType.TypeA.GetContext().Source).(interface{}).(ast.PrintableNode)
			//sourceB := (*wrongType.TypeB.GetContext().Source).(interface{}).(ast.PrintableNode)
			//causeInfo = fmt.Sprintf("First type comes from: %s and the second one from: N/A.", sourceA.Print(c))
		}

		message, textMessage := c.FormatParsingError(
			"Type Mismatch",
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", wrongType.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
		}
	} else if wrongTypeLen, ok := err.(hindley_milner.UnificationLengthError); ok {
		src := wrongTypeLen.Source().(interface{}).(ast.NodeWithPosition)

		causeInfo := ""

		if wrongType.IsCausedByBuiltin() {
		causeInfo = fmt.Sprintf("Caused by internal definition: %s", wrongType.GetCauseName())
		} else {
		//sourceA := (*wrongType.TypeA.GetContext().Source).(interface{}).(ast.PrintableNode)
		//sourceB := (*wrongType.TypeB.GetContext().Source).(interface{}).(ast.PrintableNode)
		//causeInfo = fmt.Sprintf("First type comes from: %s and the second one from: N/A.", sourceA.Print(c))
		}

		message, textMessage := c.FormatParsingError(
			"Type Mismatch",
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", wrongTypeLen.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
		}
	} else if noOverloadCandidates, ok := err.(hindley_milner.InvalidOverloadCandidatesError); ok {
		src := noOverloadCandidates.Source().(interface{}).(ast.NodeWithPosition)

		causeInfo := ""

		message, textMessage := c.FormatParsingError(
			"No overload candidates",
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", noOverloadCandidates.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
		}
	}
	panic(fmt.Sprintf("Unknown error: [%v]\n", err))
	return TypeCheckingError{
		message: "Unknown error\n",
		textMessage: "Unknown error\n",
	}
}

func (tc *LatteTypeChecker) Check(program *ast.LatteProgram, c *context.ParsingContext) error {
	debug.SetGCPercent(-1)
	//var scheme *hindley_milner.Scheme
	var err error
	//var retEnv hindley_milner.Env

	config := hindley_milner.NewInferConfiguration()
	config.CreateDefaultEmptyType = func() *hindley_milner.Scheme { return hindley_milner.NewScheme(nil, ast.CreatePrimitive(ast.T_VOID)) }

	_, _, err = hindley_milner.Infer(tc.GetEnv(), program, config)
	if err != nil {
		return wrapTypeCheckingError(err, c)
	}
	return nil
	//simpleType, ok := scheme.Type()
	//fmt.Printf("simple Type: %v | isMonoType: %v | err: %v\n", simpleType, ok, err)
	//hindley_milner.PrintEnv(retEnv)
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
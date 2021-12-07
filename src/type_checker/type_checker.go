package type_checker

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser"
	"github.com/styczynski/latte-compiler/src/parser/ast"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LatteTypeChecker struct{}

func CreateLatteTypeChecker() *LatteTypeChecker {
	return &LatteTypeChecker{}
}

func (tc *LatteTypeChecker) Test(c *context.ParsingContext) {
	// Nothing
}

func (tc *LatteTypeChecker) GetEnv() *hindley_milner.SimpleEnv {
	return hindley_milner.CreateSimpleEnv(map[string][]*hindley_milner.Scheme{
		"+": hindley_milner.SingleDef(nil, hindley_milner.NewUnionType([]hindley_milner.Type{
			hindley_milner.NewFnType(
				ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
			),
			hindley_milner.NewFnType(
				ast.CreatePrimitive(ast.T_STRING), ast.CreatePrimitive(ast.T_STRING), ast.CreatePrimitive(ast.T_STRING),
			),
		})),
		"null": hindley_milner.SingleDef(nil,
			hindley_milner.NewSignedStructType("", nil),
		),
		"true":  hindley_milner.SingleDef(nil, ast.CreatePrimitive(ast.T_BOOL)),
		"false": hindley_milner.SingleDef(nil, ast.CreatePrimitive(ast.T_BOOL)),
		"||": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"&&": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"-": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"/": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"*": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"%": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"unary_!": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_BOOL), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"unary_-": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"--": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"++": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT),
		)),
		"<=": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		">=": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"==": hindley_milner.SingleDef(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(
				hindley_milner.TVar('a'), hindley_milner.TVar('a'), ast.CreatePrimitive(ast.T_BOOL),
			)),
		"!=": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"<": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		">": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_BOOL),
		)),
		"readInt": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_VOID_ARG), ast.CreatePrimitive(ast.T_INT),
		)),
		"readString": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_VOID_ARG), ast.CreatePrimitive(ast.T_STRING),
		)),
		"printInt": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_INT), ast.CreatePrimitive(ast.T_VOID),
		)),
		"printString": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(
			ast.CreatePrimitive(ast.T_STRING), ast.CreatePrimitive(ast.T_VOID),
		)),
		"[]": hindley_milner.SingleDef(
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
	message     string
	textMessage string
	errorName   string
}

func (e TypeCheckingError) ErrorName() string {
	return e.errorName
}

func (e TypeCheckingError) Error() string {
	return e.message
}

func (e TypeCheckingError) CliMessage() string {
	return e.textMessage
}

func wrapTypeCheckingError(err error, c *context.ParsingContext) *TypeCheckingError {
	if undef, ok := err.(hindley_milner.UndefinedSymbol); ok {
		src := undef.Source.(interface{}).(generic_ast.NodeWithPosition)
		errorName := "Unknown symbol"
		message, textMessage := c.FormatParsingError(
			errorName,
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
			errorName:   errorName,
		}
	} else if wrongType, ok := err.(hindley_milner.UnificationWrongTypeError); ok {
		srcNode := wrongType.Source().(interface{})
		src := srcNode.(generic_ast.NodeWithPosition)
		causeInfo := []string{}

		if wrongType.IsCausedByBuiltin() {
			causeInfo = []string{fmt.Sprintf("Caused by internal definition: %s", wrongType.GetCauseName())}
		} else {
			var sourceA *generic_ast.NodeWithPosition = nil
			var sourceB *generic_ast.NodeWithPosition = nil

			if wrongType.TypeA.GetContext().Source == nil || wrongType.TypeB.GetContext().Source == nil {
				causeInfo = append(causeInfo, fmt.Sprintf("The type mismatch occured in %s", (*wrongType.Constraint.Context().Source).(interface{}).(generic_ast.NodeWithPosition).Begin().String()))
			}
			if wrongType.TypeA.GetContext().Source != nil {
				v := (*wrongType.TypeA.GetContext().Source).(interface{}).(generic_ast.NodeWithPosition)
				sourceA = &v
				causeInfo = append(causeInfo, fmt.Sprintf("First type comes from: %s.", v.Begin().String()))
			}
			if wrongType.TypeB.GetContext().Source != nil {
				v := (*wrongType.TypeB.GetContext().Source).(interface{}).(generic_ast.NodeWithPosition)
				sourceB = &v
				causeInfo = append(causeInfo, fmt.Sprintf("Second type from: %s.", v.Begin().String()))
			}

			if sourceA != nil && sourceB != nil {
				if custA, ok := (*sourceA).(hindley_milner.HMExpressionWithCustomMismatchErrorDescription); ok {
					causeInfo = append(causeInfo, custA.OnTypeMismatch(*sourceA, *sourceB)...)
				} else if custB, ok := (*sourceB).(hindley_milner.HMExpressionWithCustomMismatchErrorDescription); ok {
					causeInfo = append(causeInfo, custB.OnTypeMismatch(*sourceB, *sourceA)...)
				}
			}
		}

		errorName := "Type Mismatch"
		message, textMessage := c.FormatParsingError(
			errorName,
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s %s", wrongType.Error(), strings.Join(causeInfo, "\n                  ")),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	} else if wrongTypeLen, ok := err.(hindley_milner.UnificationLengthError); ok {
		src := wrongTypeLen.Source().(interface{}).(generic_ast.NodeWithPosition)

		causeInfo := ""

		if wrongType.IsCausedByBuiltin() {
			causeInfo = fmt.Sprintf("Caused by internal definition: %s", wrongType.GetCauseName())
		} else {
			//sourceA := (*wrongType.TypeA.GetContext().Source).(interface{}).(ast.PrintableNode)
			//sourceB := (*wrongType.TypeB.GetContext().Source).(interface{}).(ast.PrintableNode)
			//causeInfo = fmt.Sprintf("First type comes from: %s and the second one from: N/A.", sourceA.Print(c))
		}

		errorName := "Type Mismatch"
		message, textMessage := c.FormatParsingError(
			errorName,
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
			errorName:   errorName,
		}
	} else if noOverloadCandidates, ok := err.(hindley_milner.InvalidOverloadCandidatesError); ok {
		src := noOverloadCandidates.Source().(interface{}).(generic_ast.NodeWithPosition)

		causeInfo := ""

		errorName := "No overload candidates"
		message, textMessage := c.FormatParsingError(
			errorName,
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
			errorName:   errorName,
		}
	} else if reccurentTypeError, ok := err.(hindley_milner.UnificationRecurrentTypeError); ok {
		//src := reccurentTypeError.Source().(interface{}).(generic_ast.NodeWithPosition)

		causeInfo := ""
		errorName := "Recurrent type"
		message, textMessage := c.FormatParsingError(
			errorName,
			undef.Error(),
			0,  //src.Begin().Line,
			0,  //src.Begin().Column,
			"", //src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", reccurentTypeError.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	} else if varRedef, ok := err.(hindley_milner.VariableRedefinedError); ok {
		src := varRedef.Source().(interface{}).(generic_ast.NodeWithPosition)
		causeInfo := ""

		previousDef := varRedef.PreviousDefinition.Source

		if previousDef != nil {
			pos := (*previousDef).(interface{}).(generic_ast.NodeWithPosition).Begin()
			causeInfo = fmt.Sprintf("Previous definition can be found at %s.", pos.String())
		}

		errorName := "Variable redefined"
		message, textMessage := c.FormatParsingError(
			errorName,
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", varRedef.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	} else if builtinRedef, ok := err.(hindley_milner.BuiltinRedefinedError); ok {
		src := builtinRedef.Source().(interface{}).(generic_ast.NodeWithPosition)
		causeInfo := ""

		errorName := "Builtin redefined"
		message, textMessage := c.FormatParsingError(
			errorName,
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", builtinRedef.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	} else if genericASTError, ok := err.(hindley_milner.ASTError); ok {
		src := genericASTError.Source().(interface{}).(generic_ast.NodeWithPosition)
		causeInfo := ""

		errorName := genericASTError.Name
		message, textMessage := c.FormatParsingError(
			errorName,
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", genericASTError.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	} else if invalidReturnType, ok := err.(hindley_milner.InvalidReturnTypeError); ok {
		src := invalidReturnType.Source().(interface{}).(generic_ast.NodeWithPosition)
		causeInfo := ""

		errorName := "Invalid return"
		message, textMessage := c.FormatParsingError(
			errorName,
			undef.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			fmt.Sprintf("%s%s", invalidReturnType.Error(), causeInfo),
		)
		return &TypeCheckingError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	}
	// panic(fmt.Sprintf("Unknown error: [%v]\n", err))
	// TODO: Format error
	return &TypeCheckingError{
		message:     err.Error(),
		textMessage: err.Error(),
		errorName:   "Unknown error",
	}
}

type LatteTypecheckedProgram struct {
	Program           parser.LatteParsedProgram
	TypeCheckingError *TypeCheckingError
	filename          string
}

func (p LatteTypecheckedProgram) Filename() string {
	return p.filename
}

func (p LatteTypecheckedProgram) Resolve() LatteTypecheckedProgram {
	return p
}

type LatteTypecheckedProgramPromise interface {
	Resolve() LatteTypecheckedProgram
}

type LatteTypecheckedProgramPromiseChan <-chan LatteTypecheckedProgram

func (p LatteTypecheckedProgramPromiseChan) Resolve() LatteTypecheckedProgram {
	return <-p
}

func (tc *LatteTypeChecker) checkAsync(programPromise parser.LatteParsedProgramPromise, c *context.ParsingContext) LatteTypecheckedProgramPromise {
	r := make(chan LatteTypecheckedProgram)
	ctx := c.Copy()
	go func() {
		program := programPromise.Resolve()
		defer errors.GeneralRecovery(ctx, "Typechecking", program.Filename(), func(message string, textMessage string) {
			r <- LatteTypecheckedProgram{
				Program: program,
				TypeCheckingError: &TypeCheckingError{
					message:     message,
					textMessage: textMessage,
					errorName:   "PANIC (Typechecking)",
				},
				filename: program.Filename(),
			}
		}, func() {
			close(r)
		})
		//fmt.Printf("%#v\n", program)
		if program.ParsingError() != nil {
			r <- LatteTypecheckedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}
		if program.Context() != nil {
			ctx = program.Context()
		}
		c.EventsCollectorStream.Start("Typechecking", c, program)
		defer c.EventsCollectorStream.End("Typechecking", c, program)

		OnConstrintGenerationStarted := func() {
			c.EventsCollectorStream.Start("Generating constraints", c, program)
		}
		OnConstrintGenerationFinished := func() {
			c.EventsCollectorStream.End("Generating constraints", c, program)
		}
		OnSolvingStarted := func() {
			c.EventsCollectorStream.Start("Solving constraints", c, program)
		}
		OnSolvingFinished := func() {
			c.EventsCollectorStream.End("Solving constraints", c, program)
		}
		OnPostprocessingStarted := func() {
			c.EventsCollectorStream.Start("Postprocessing", c, program)
		}
		OnPostprocessingFinished := func() {
			c.EventsCollectorStream.End("Postprocessing", c, program)
		}

		config := hindley_milner.NewInferConfiguration()
		config.CreateDefaultEmptyType = func() *hindley_milner.Scheme { return hindley_milner.NewScheme(nil, ast.CreatePrimitive(ast.T_VOID)) }

		config.OnConstrintGenerationStarted = &OnConstrintGenerationStarted
		config.OnConstrintGenerationFinished = &OnConstrintGenerationFinished
		config.OnSolvingStarted = &OnSolvingStarted
		config.OnSolvingFinished = &OnSolvingFinished
		config.OnPostprocessingStarted = &OnPostprocessingStarted
		config.OnPostprocessingFinished = &OnPostprocessingFinished

		env := tc.GetEnv()
		infer := hindley_milner.CreateImperInferenceBackend(env, config)
		_, _, err := hindley_milner.Infer(env, program.AST(), config, infer)
		if err != nil {
			r <- LatteTypecheckedProgram{
				Program:           program,
				TypeCheckingError: wrapTypeCheckingError(err, ctx),
				filename:          program.Filename(),
			}
			return
		}
		r <- LatteTypecheckedProgram{
			Program:  program,
			filename: program.Filename(),
		}
	}()

	return LatteTypecheckedProgramPromiseChan(r)
}

func (tc *LatteTypeChecker) Check(programs []parser.LatteParsedProgramPromise, c *context.ParsingContext) []LatteTypecheckedProgramPromise {
	ret := []LatteTypecheckedProgramPromise{}
	for _, programPromise := range programs {
		ret = append(ret, tc.checkAsync(programPromise, c))
	}
	return ret
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

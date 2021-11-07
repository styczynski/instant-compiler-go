package jvm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/compiler/jvm/jasmine"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_COMPILER_BACKEND, CompilerJVMBackendFactory{})
}

type CompilerJVMBackendFactory struct{}

func (CompilerJVMBackendFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateCompilerJVMBackend()
}

func (CompilerJVMBackendFactory) Params(argSpec *config.EntityArgSpec) {
	argSpec.AddString("test", "huj", "elo")
}

func (CompilerJVMBackendFactory) EntityName() string {
	return "jvm"
}

type CompilerJVMBackend struct {
	state *compiler.CompilerState
}

func (backend CompilerJVMBackend) optimizeStackBiAlloc(left generic_ast.Expression, right generic_ast.Expression, supportsSwap bool) ([]jasmine.JasmineInstruction, int64) {
	cl, dl := backend.compileExpression(left)
	cr, dr := backend.compileExpression(right)

	d1 := dl + 1 // d1 = max(dl+1, dr) - swapped
	if d1 < dr {
		d1 = dr
	}

	d2 := dr + 1 // d2 = max(dl, dr+1) - normal
	if d2 < dl {
		d2 = dl
	}

	ret := []jasmine.JasmineInstruction{}
	if d1 > d2 {
		ret = append(ret, cl...)
		ret = append(ret, cr...)
		return ret, d2
	} else {
		ret = append(ret, cr...)
		ret = append(ret, cl...)
		if supportsSwap {
			ret = append(ret, &jasmine.JasmineSwap{})
		}
		return ret, d1
	}
}

func (backend CompilerJVMBackend) compileExpression(expr generic_ast.Expression) ([]jasmine.JasmineInstruction, int64) {
	if _, ok := (expr.(*ast.Empty)); ok {
		return []jasmine.JasmineInstruction{}, 0
	}
	if expr, ok := (expr.(*ast.LatteProgram)); ok {
		ret := []jasmine.JasmineInstruction{}
		maxDepth := int64(0)
		for _, stmt := range expr.Definitions {
			compiledValue, s := backend.compileExpression(stmt)
			if s > maxDepth {
				maxDepth = s
			}
			ret = append(ret, compiledValue...)
		}
		return ret, int64(maxDepth)
	}
	if expr, ok := (expr.(*ast.Statement)); ok {
		ret := []jasmine.JasmineInstruction{}
		if expr.IsAssignment() {
			compiledValue, s := backend.compileExpression(expr.Assignment.Value)
			loc := backend.state.Define(expr.Assignment.TargetName)

			ret = append(ret, compiledValue...)
			ret = append(ret, &jasmine.JasmineStoreInt{
				Index: loc,
			})
			return ret, s
		} else if expr.IsExpression() {
			ret = append(ret, &jasmine.JasmineGetStatic{
				Source: "java/lang/System/out",
				Object: "Ljava/io/PrintStream;",
			})
			compiledValue, s := backend.compileExpression(expr.Expression)
			ret = append(ret, compiledValue...)
			ret = append(ret, &jasmine.JasmineInvokeStatic{
				Target:  "java/io/PrintStream/println",
				Args:    []string{"I"},
				Return:  "V",
				Virtual: true,
			})
			return ret, s + 1
		} else if expr.IsEmpty() {
			return ret, 0
		}
	}
	if expr, ok := (expr.(*ast.Addition)); ok {
		if !expr.HasNext() {
			return backend.compileExpression(expr.Multiplication)
		}
		ret, s := backend.optimizeStackBiAlloc(expr.Multiplication, expr.Next, expr.Op == "-")
		ret = append(ret, jasmine.CreateJasmineIntOp(expr.Op))
		return ret, s
	}
	if expr, ok := (expr.(*ast.Expression)); ok {
		return backend.compileExpression(&expr.Addition)
	}
	if expr, ok := (expr.(*ast.Multiplication)); ok {
		if !expr.HasNext() {
			return backend.compileExpression(expr.Primary)
		}
		ret, s := backend.optimizeStackBiAlloc(expr.Primary, expr.Next, expr.Op == "/")
		ret = append(ret, jasmine.CreateJasmineIntOp(expr.Op))
		return ret, s
	}
	if expr, ok := (expr.(*ast.Primary)); ok {
		if expr.IsVariable() {
			loc := backend.state.GetVariableFromScope(*expr.Variable)
			return []jasmine.JasmineInstruction{
				&jasmine.JasmineLoadInt{
					Index: loc,
				},
			}, 1
		} else if expr.IsInt() {
			return []jasmine.JasmineInstruction{
				&jasmine.JasmineConstInt{
					Index: *expr.Int,
				},
			}, 1
		} else if expr.IsSubexpression() {
			return backend.compileExpression(expr.SubExpression)
		}
	}
	panic(fmt.Sprintf("Invalid instruction given to compileExpression(): %s", reflect.TypeOf(expr)))
}

func (backend CompilerJVMBackend) RunCompiledCode(runContext compiler.CompiledCodeRunContext, c *context.ParsingContext) ([]string, *compiler.RunError) {
	callOut, err := runContext.Call("java", "rror", "-cp", "$OUTPUT_DIR", "$CLASS_NAME")
	if err != nil {
		return nil, err
	}

	out := []string{}
	for _, line := range callOut {
		if len(line) > 0 && !strings.Contains(line, "_JAVA_OPTIONS") {
			out = append(out, line)
		}
	}

	return out, nil
}

func (backend CompilerJVMBackend) Compile(program type_checker.LatteTypecheckedProgram, c *context.ParsingContext, b *compiler.BuildContext) compiler.LatteCompiledProgramPromiseChan {
	ret := make(chan compiler.LatteCompiledProgram)
	go func() {

		ast := program.Program.AST()
		outputCode, maxStack := backend.compileExpression(ast)

		// output := jasmine.JasmineProgram{
		// 	StackLimit:   maxStack,
		// 	LocalsLimit:  int64(backend.state.ScopeSize()),
		// 	Instructions: outputCode,
		// }

		className := b.GetVariable("INPUT_FILE_BASE")
		b.SetCompilerMeta("CLASS_NAME", className)

		output := jasmine.JasmineProgram{
			Instructions: []jasmine.JasmineInstruction{
				&jasmine.JasmineClass{
					Name:  fmt.Sprintf("public %s", className),
					Super: "java/lang/Object",
					Methods: []*jasmine.JasmineMethod{
						&jasmine.JasmineMethod{
							Name:        "<init>",
							Returns:     "V",
							StackLimit:  1,
							LocalsLimit: 1,
							Body: []jasmine.JasmineInstruction{
								&jasmine.JasmineReferenceLoad{
									Index: 0,
								},
								&jasmine.JasmineInvokeStatic{
									Target:  "java/lang/Object/<init>",
									Special: true,
									Return:  "V",
								},
								&jasmine.JasmineReturn{},
							},
						},
						&jasmine.JasmineMethod{
							Name:        "public static main",
							Args:        []string{"[Ljava/lang/String;"},
							Returns:     "V",
							StackLimit:  maxStack,
							LocalsLimit: int64(backend.state.ScopeSize()) + 1,
							Body:        append(outputCode, &jasmine.JasmineReturn{}),
						},
					},
				},
			},
		}
		output.Normalize()

		validationErr := output.Validate()

		if validationErr != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  &output,
				CompilationError: validationErr,
			}
			return
		}

		b.WriteBuildFile("code.jasmine", []byte(output.ProgramToText()))

		validationErr = b.Call("java", "rror", "-jar", "$ROOT/lib/jasmin.jar", "-d", "$BUILD_DIR/out", "$BUILD_DIR/code.jasmine")
		if validationErr != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  &output,
				CompilationError: validationErr,
			}
			return
		}

		outputJVMBytecode := b.ReadBuildFile("out/%s.class", className)
		b.WriteOutput("JVM bytecode file", "class", outputJVMBytecode)
		b.WriteOutput("Jasmine source", "j", []byte(output.ProgramToText()))

		ret <- compiler.LatteCompiledProgram{
			Program:          program,
			CompiledProgram:  &output,
			CompilationError: validationErr,
		}
	}()

	return ret
}

func (CompilerJVMBackend) BackendName() string {
	return "JVM Jasmine backend"
}

func CreateCompilerJVMBackend() compiler.CompilerBackend {
	return CompilerJVMBackend{
		state: compiler.CreateCompilerState(),
	}
}

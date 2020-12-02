package type_checker

import (
	"fmt"
	"log"
	"strings"

	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
	"github.com/pkg/errors"
)

const digits = "0123456789"

type TyperExpression interface {
	hindley_milner.Expression
	hindley_milner.Typer
}

type λ struct {
	args []string
	body hindley_milner.Expression
}

func (n λ) Name() hindley_milner.NameGroup     { return hindley_milner.Names(n.args) }
func (n λ) Body() hindley_milner.Expression { return n.body }
func (n λ) IsLambda() bool   { return true }

type lit string

func (n lit) Name() hindley_milner.NameGroup     { return hindley_milner.Name(string(n)) }
func (n lit) Body() hindley_milner.Expression { return n }
func (n lit) Type() hindley_milner.Type {
	switch {
	case strings.ContainsAny(digits, string(n)) && strings.ContainsAny(digits, string(n[0])):
		return Prim(Float)
	case string(n) == "true" || string(n) == "false":
		return Prim(Bool)
	default:
		return nil
	}
}
func (n lit) IsLit() bool    { return true }
func (n lit) IsLambda() bool { return true }

type app struct {
	f   hindley_milner.Expression
	args []hindley_milner.Expression
}

func (n app) Fn() hindley_milner.Expression   { return n.f }
func (n app) Body() hindley_milner.Expression {
	return hindley_milner.Batch{
		n.args,
	}
}
//func (n app) Arg() hindley_milner.Expression  { return n.arg }

type let struct {
	name string
	def  hindley_milner.Expression
	in   hindley_milner.Expression
}

func (n let) Name() hindley_milner.NameGroup     { return hindley_milner.Name(n.name) }
func (n let) Def() hindley_milner.Expression  { return n.def }
func (n let) Body() hindley_milner.Expression { return n.in }

type letrec struct {
	name string
	def  hindley_milner.Expression
	in   hindley_milner.Expression
}

func (n letrec) Name() hindley_milner.NameGroup           { return hindley_milner.Name(n.name) }
func (n letrec) Def() hindley_milner.Expression        { return n.def }
func (n letrec) Body() hindley_milner.Expression       { return n.in }
func (n letrec) Children() []hindley_milner.Expression { return []hindley_milner.Expression{n.def, n.in} }
func (n letrec) IsRecursive() bool      { return true }

type prim struct {
	val primid
	context hindley_milner.CodeContext
}

type primid byte

const (
	Float primid = iota
	Bool
)

func Prim(val primid) prim {
	return prim{
		val: val,
	}
}

// implement Type
func (t prim) Name() string                                   { return t.String() }
func (t prim) Apply(hindley_milner.Subs) hindley_milner.Substitutable                       { return t }
func (t prim) FreeTypeVar() hindley_milner.TypeVarSet                        { return nil }
func (t prim) Normalize(hindley_milner.TypeVarSet, hindley_milner.TypeVarSet) (hindley_milner.Type, error) { return t, nil }
func (t prim) Types() hindley_milner.Types                                   { return nil }
func (t prim) Eq(other hindley_milner.Type) bool {
	if ot, ok := other.(prim); ok {
		return ot.val == t.val
	}
	return false
}

//func (t prim) SetCodeContext(c *hindley_milner.EnvCodeContext) hindley_milner.Type {
//	return t
//}
//
//func (t prim) GetCodeContext() *hindley_milner.EnvCodeContext {
//	return nil
//}

func (t prim) Format(s fmt.State, c rune) { fmt.Fprintf(s, t.String()) }
func (t prim) String() string {
	name := "?"
	switch t.val {
	case Float:
		name = "Float"
	case Bool:
		name = "Bool"
	}
	return fmt.Sprintf("%s%s", hindley_milner.TypeStringPrefix(t), name)
}

func (t prim) MapTypes(mapper hindley_milner.TypeMapper) hindley_milner.Type {
	return mapper(t)
}

func (t prim) WithContext(c hindley_milner.CodeContext) hindley_milner.Type {
	return prim{
		val:     t.val,
		context: c,
	}
}

func (t prim) GetContext() hindley_milner.CodeContext {
	return t.context
}

//Phillip Greenspun's tenth law says:
//		"Any sufficiently complicated C or Fortran program contains an ad hoc, informally-specified, bug-ridden, slow implementation of half of Common Lisp."
//
// So let's implement a half-arsed lisp (Or rather, an AST that can optionally be executed upon if you write the correct interpreter)!
func Example_greenspun() {
	// haskell envy in a greenspun's tenth law example function!
	//
	// We'll assume the following is the "input" code
	// 		let fac n = if n == 0 then 1 else n * fac (n - 1) in fac 5
	// and what we have is the AST

	//fac := letrec{
	//	"fac",
	//	λ{
	//		"n",
	//		app{
	//			app{
	//				app{
	//					lit("if"),
	//					app{
	//						lit("isZero"),
	//						lit("n"),
	//					},
	//				},
	//				lit("1"),
	//			},
	//			app{
	//				app{lit("mul"), lit("n")},
	//				app{lit("fac"), app{lit("--"), lit("n")}},
	//			},
	//		},
	//	},
	//	app{lit("fac"), lit("5")},
	//}

	// but first, let's start with something simple:
	// let x = 3 in x+5
	//fac := let{
	//	"x",
	//	lit("3"),
	//	app{
	//		app{
	//			lit("+"),
	//			lit("5"),
	//		},
	//		lit("x"),
	//	},
	//}

	fac := let{
			name: "test",
			def:  λ{
				args: []string{ "x", "y" },
				body: app{
					lit("+"),
					[]hindley_milner.Expression{
						lit("x"),
						lit("y"),
					},
				},
			},
			in:   app{
				lit("test"),
				[]hindley_milner.Expression{
					lit("2"),
					lit("5"),
				},
			},
	}

	env := hindley_milner.CreateSimpleEnv(map[string]*hindley_milner.Scheme{
		"--":     hindley_milner.NewScheme(hindley_milner.TypeVarSet{hindley_milner.TVar('a')}, hindley_milner.NewFnType(hindley_milner.TVar('a'), hindley_milner.TVar('a'))),
		"if":     hindley_milner.NewScheme(hindley_milner.TypeVarSet{hindley_milner.TVar('a')}, hindley_milner.NewFnType(Prim(Bool), hindley_milner.TVar('a'), hindley_milner.TVar('a'), hindley_milner.TVar('a'))),
		"isZero": hindley_milner.NewScheme(nil, hindley_milner.NewFnType(Prim(Float), Prim(Bool))),
		"mul":    hindley_milner.NewScheme(nil, hindley_milner.NewFnType(Prim(Float), Prim(Float), Prim(Float))),
		"+":      hindley_milner.NewScheme(hindley_milner.TypeVarSet{hindley_milner.TVar('a')}, hindley_milner.NewFnType(hindley_milner.TVar('a'), hindley_milner.TVar('a'), hindley_milner.TVar('a'))),
	})

	var scheme *hindley_milner.Scheme
	var err error
	scheme, err = hindley_milner.Infer(env, fac)
	if err != nil {
		log.Printf("%+v", errors.Cause(err))
	}
	simpleType, ok := scheme.Type()
	fmt.Printf("simple Type: %v | isMonoType: %v | err: %v\n", simpleType, ok, err)

	//scheme, err = hindley_milner.Infer(env, fac)
	//if err != nil {
	//	panic(err)
	//	log.Printf("%+v", errors.Cause(err))
	//}
	//
	//facType, ok := scheme.Type()
	//fmt.Printf("fac Type: %v | isMonoType: %v | err: %v", facType, ok, err)

	// Output:
	// simple Type: Float | isMonoType: true | err: <nil>
	// fac Type: Float | isMonoType: true | err: <nil>

}

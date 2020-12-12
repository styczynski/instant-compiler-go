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
	args map[string]*hindley_milner.Scheme
	body hindley_milner.Expression
}

func (n λ) Args() hindley_milner.NameGroup     { return hindley_milner.NamesWithTypesFromMap(n.args) }
func (n λ) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(n.body)
}
func (n λ) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n.body)
}
func (n λ) Body() hindley_milner.Expression { return n.body }
func (n λ) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_FUNCTION }


type ret struct {
	expr hindley_milner.Expression
}
func (n ret) Body() hindley_milner.Expression {
	return n.expr
}
func (n ret) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(ret{
		expr: mapper(n.expr),
	})
}
func (n ret) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n.expr)
	mapper(n)
}
func (n ret) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_RETURN }


type em struct {}
func (n em) Body() hindley_milner.Expression {
	return em{}
}
func (n em) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(n)
}
func (n em) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n)
}
func (n em) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_TYPE }
func (n em) Type() *hindley_milner.Scheme {
	return hindley_milner.NewScheme(nil, hindley_milner.NewFnType(Prim(Float), Prim(Float), Prim(Float)))
}

type lit string

func (n lit) Name() hindley_milner.NameGroup     { return hindley_milner.Name(string(n)) }
func (n lit) Body() hindley_milner.Expression { return n }
func (n lit) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(n)
}
func (n lit) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n)
}
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
// TODO: Lit/lambda needed?
func (n lit) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }

type app struct {
	f   hindley_milner.Expression
	args []hindley_milner.Expression
}

func (n app) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	mappedExp := []hindley_milner.Expression{}
	for _, arg := range n.args {
		mappedExp = append(mappedExp, mapper(arg))
	}
	return mapper(app{
		f: mapper(n.f),
		args: mappedExp,
	})
}
func (n app) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n.f)
	for _, arg := range n.args {
		mapper(arg)
	}
	mapper(n)
}
func (n app) Fn() hindley_milner.Expression   { return n.f }
func (n app) Body() hindley_milner.Expression {
	return hindley_milner.Batch{
		n.args,
	}
}
func (n app) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_APPLICATION }


//func (n app) Arg() hindley_milner.Expression  { return n.arg }

type let struct {
	name string
	def  hindley_milner.Expression
	in   hindley_milner.Expression
}

func (n let) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(letrec{
		def: mapper(n.def),
		in: mapper(n.in),
	})
}
func (n let) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n.def)
	mapper(n.in)
	mapper(n)
}
func (n let) Var() hindley_milner.NameGroup     { return hindley_milner.Name(n.name) }
func (n let) Def() hindley_milner.Expression  { return n.def }
func (n let) Body() hindley_milner.Expression { return n.in }
func (n let) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LET }

type decl struct {
	name string
	def  hindley_milner.Expression
}

func (n decl) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(decl{
		def: mapper(n.def),
	})
}
func (n decl) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n.def)
	mapper(n)
}
func (n decl) Var() hindley_milner.NameGroup           { return hindley_milner.Name(n.name) }
func (n decl) Def() hindley_milner.Expression        { return n.def }
func (n decl) Body() hindley_milner.Expression       { return n }
func (n decl) Children() []hindley_milner.Expression { return []hindley_milner.Expression{n.def} }
func (n decl) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_DECLARATION }


type letrec struct {
	name string
	def  hindley_milner.Expression
	in   hindley_milner.Expression
}

func (n letrec) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(letrec{
		def: mapper(n.def),
		in: mapper(n.in),
	})
}
func (n letrec) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(n.def)
	mapper(n.in)
	mapper(n)
}
func (n letrec) Var() hindley_milner.NameGroup           { return hindley_milner.Name(n.name) }
func (n letrec) Def() hindley_milner.Expression        { return n.def }
func (n letrec) Body() hindley_milner.Expression       { return n.in }
func (n letrec) Children() []hindley_milner.Expression { return []hindley_milner.Expression{n.def, n.in} }
func (n letrec) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LET_RECURSIVE }


type prim struct {
	val primid
	context hindley_milner.CodeContext
}

type primid byte

const (
	Float primid = iota
	Bool
	Void
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
	case Void:
		name = "Void"
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

	//fac := hindley_milner.Batch{[]hindley_milner.Expression{
	//	decl{
	//		name: "test",
	//		def: λ{
	//			args: map[string]*hindley_milner.Scheme{
	//				"x": hindley_milner.NewScheme(nil, Prim(Float)),
	//				"y": hindley_milner.NewScheme(nil, Prim(Float)),
	//			},
	//			body: hindley_milner.Batch{Exp: []hindley_milner.Expression{
	//				app{
	//					em{},
	//					[]hindley_milner.Expression{
	//						lit("y"),
	//						lit("x"),
	//					},
	//				},
	//				app{
	//					em{},
	//					[]hindley_milner.Expression{
	//						lit("x"),
	//						lit("y"),
	//					},
	//				},
	//				ret{
	//					lit("2"),
	//				},
	//			}},
	//		},
	//	},
	//	app{
	//		lit("test"),
	//		[]hindley_milner.Expression{
	//			lit("2"),
	//			lit("5"),
	//		},
	//	},
	//}}

	fac := hindley_milner.Batch{[]hindley_milner.Expression{
		decl{
			name: "test",
			def: λ{
				args: map[string]*hindley_milner.Scheme{
					"x": hindley_milner.NewScheme(nil, Prim(Float)),
				},
				body: hindley_milner.Batch{Exp: []hindley_milner.Expression{
					ret{
						lit("2"),
					},
				}},
			},
		},
		app{
			lit("test"),
			[]hindley_milner.Expression{
				lit("5"),
			},
		},
	}}

	env := hindley_milner.CreateSimpleEnv(map[string][]*hindley_milner.Scheme{
		"--":     hindley_milner.SingleDef(hindley_milner.TypeVarSet{hindley_milner.TVar('a')}, hindley_milner.NewFnType(hindley_milner.TVar('a'), hindley_milner.TVar('a'))),
		"if":     hindley_milner.SingleDef(hindley_milner.TypeVarSet{hindley_milner.TVar('a')}, hindley_milner.NewFnType(Prim(Bool), hindley_milner.TVar('a'), hindley_milner.TVar('a'), hindley_milner.TVar('a'))),
		"isZero": hindley_milner.SingleDef(nil, hindley_milner.NewFnType(Prim(Float), Prim(Bool))),
		"mul":    hindley_milner.SingleDef(nil, hindley_milner.NewFnType(Prim(Float), Prim(Float), Prim(Float))),
		"+":      hindley_milner.SingleDef(hindley_milner.TypeVarSet{hindley_milner.TVar('a')}, hindley_milner.NewFnType(hindley_milner.TVar('a'), hindley_milner.TVar('a'), hindley_milner.TVar('a'))),
	})

	var scheme *hindley_milner.Scheme
	var err error
	config := hindley_milner.NewInferConfiguration()
	config.CreateDefaultEmptyType = func() *hindley_milner.Scheme { return hindley_milner.NewScheme(nil, Prim(Void)) }

	scheme, _, err = hindley_milner.Infer(env, fac, config)
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

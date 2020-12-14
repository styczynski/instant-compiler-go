package hindley_milner

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

// SignedTuple is a basic record/tuple type. It takes an optional name.
type SignedTuple struct {
	ts   []Type
	name string
	context CodeContext
}

// NewSignedTupleType creates a new SignedTuple Type
func NewSignedTupleType(name string, ts ...Type) *SignedTuple {
	return &SignedTuple{
		ts:   ts,
		name: name,
	}
}

func (t *SignedTuple) Apply(subs Subs) Substitutable {
	ts := make([]Type, len(t.ts))
	for i, v := range t.ts {
		ts[i] = v.Apply(subs).(Type)
	}
	return NewSignedTupleType(t.name, ts...)
}

func (t *SignedTuple) FreeTypeVar() TypeVarSet {
	var tvs TypeVarSet
	for _, v := range t.ts {
		tvs = v.FreeTypeVar().Union(tvs)
	}
	return tvs
}

func (t *SignedTuple) Name() string {
	return t.name
}

func (t *SignedTuple) Normalize(k, v TypeVarSet) (Type, error) {
	ts := make([]Type, len(t.ts))
	var err error
	for i, tt := range t.ts {
		if ts[i], err = tt.Normalize(k, v); err != nil {
			return nil, err
		}
	}
	return NewSignedTupleType(t.name, ts...), nil
}

func (t *SignedTuple) Types() Types {
	ts := BorrowTypes(len(t.ts))
	copy(ts, t.ts)
	return ts
}

func (t *SignedTuple) Eq(other Type) bool {
	if ot, ok := other.(*SignedTuple); ok {
		if len(ot.ts) != len(t.ts) {
			return false
		}
		if ot.name != t.name {
			return false
		}
		for i, v := range t.ts {
			if !v.Eq(ot.ts[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (t *SignedTuple) Format(f fmt.State, c rune) {
	f.Write([]byte(fmt.Sprintf("%s<", t.name)))
	f.Write([]byte(TypeStringPrefix(t)))
	for i, v := range t.ts {
		if i < len(t.ts)-1 {
			fmt.Fprintf(f, "%v, ", v)
		} else {
			fmt.Fprintf(f, "%v>", v)
		}
	}

}

func (t *SignedTuple) MapTypes(mapper TypeMapper) Type {
	newSignedTuple := &SignedTuple{
		ts:   []Type{},
		name: t.name,
		context: t.context,
	}
	for _, v := range t.ts {
		newSignedTuple.ts = append(newSignedTuple.ts, v.MapTypes(mapper))
	}
	return mapper(newSignedTuple)
}

func (t *SignedTuple) WithContext(c CodeContext) Type {
	return &SignedTuple{
		ts:   t.ts,
		name: t.name,
		context: c,
	}
}

func (t *SignedTuple) GetContext() CodeContext {
	return t.context
}

func (t *SignedTuple) String() string { return fmt.Sprintf("%s%v", TypeStringPrefix(t), t) }

// Clone implements Cloner
func (t *SignedTuple) Clone() interface{} {
	retVal := new(SignedTuple)
	ts := BorrowTypes(len(t.ts))
	for i, tt := range t.ts {
		if c, ok := tt.(Cloner); ok {
			ts[i] = c.Clone().(Type)
		} else {
			ts[i] = tt
		}
	}
	retVal.ts = ts
	retVal.name = t.name

	return retVal
}


//

type SignedTupleUnwrapExpr struct {
	name string
	len int16
	index int16
	expr  generic_ast.Expression
}

func ExpressionSignedTupleGet(name string, len int, index int, expr generic_ast.Expression) *SignedTupleUnwrapExpr {
	return &SignedTupleUnwrapExpr{
		name:  name,
		len:   int16(len),
		index: int16(index),
		expr:  expr,
	}
}

func (ast *SignedTupleUnwrapExpr) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	// TODO
	return mapper(parent, ast, context)
}

func (ast *SignedTupleUnwrapExpr) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *SignedTupleUnwrapExpr) ExpressionType() ExpressionType {
	return E_APPLICATION
}

func (ast *SignedTupleUnwrapExpr) Fn() generic_ast.Expression {
	args := []TypeVariable{}
	argsTypes := []Type{}
	for i := int16(0); i<ast.len; i++ {
		args = append(args, TVar(i))
		argsTypes = append(argsTypes, TVar(i))
	}
	return EmbeddedTypeExpr{
		GetType: func() *Scheme {
			return NewScheme(args, NewFnType(
				NewSignedTupleType(ast.name, argsTypes...),
				TVar(ast.index)))
		},
	}
}

func (ast *SignedTupleUnwrapExpr) Body() generic_ast.Expression {
	return ast.expr
}
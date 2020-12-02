package hindley_milner

import (
	"fmt"

	"github.com/pkg/errors"
)

// Cloner is any type that can clone
type Cloner interface {
	Clone() interface{}
}

// Fresher keeps track of all the TypeVariables that has been generated so far. It has one method - Fresh(), which is to create a new TypeVariable
type Fresher interface {
	Fresh() TypeVariable
}

func saveExprContext(t Type, source *Expression) Type {
	return t.WithContext(CreateCodeContext(*source))
}

type inferer struct {
	env Env
	cs  Constraints
	t   Type

	count int
}

func newInferer(env Env) *inferer {
	return &inferer{
		env: env,
	}
}

func (infer *inferer) Fresh() TypeVariable {
	retVal := letters[infer.count]
	infer.count++
	return TypeVariable{
		value: rune(retVal),
	}
}

func (infer *inferer) lookup(isLiteral bool, name string) error {
	s, ok := infer.env.SchemeOf(name)
	if !ok {
		return UndefinedSymbol{
			Name: name,
			IsLiteral: isLiteral,
			IsVariable: !isLiteral,
		}
	}
	infer.t = Instantiate(infer, s)
	return nil
}

func (infer *inferer) consGen(expr Expression) (err error) {

	// explicit types/inferers - can fail
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

		err = nil // reset errors
	}

	// fallbacks

	switch et := expr.(type) {
	case Batch:
		panic(fmt.Errorf("Batch cannot be used directly inside the expression."))
	case Literal:
		return infer.lookup(true, et.Name())

	case Var:
		if err = infer.lookup(false, et.Name()); err != nil {
			infer.env.Add(et.Name(), &Scheme{t: et.Type()})
			err = nil
		}

	case Lambda:
		tv := infer.Fresh()
		env := infer.env // backup

		infer.env = infer.env.Clone()
		infer.env.Remove(et.Name())
		sc := new(Scheme)
		sc.t = tv
		infer.env.Add(et.Name(), sc)

		if err = infer.consGen(et.Body()); err != nil {
			return errors.Wrapf(err, "Unable to infer body of %v. Body: %v", et, et.Body())
		}

		infer.t = NewFnType(tv, infer.t)
		infer.t = saveExprContext(infer.t, &expr)
		infer.env = env // restore backup

	case Apply:

		firstExec := true
		batchErr := ApplyBatch(et.Body(), func(body Expression) error {
			if firstExec {
				if err = infer.consGen(et.Fn()); err != nil {
					return errors.Wrapf(err, "Unable to infer Fn of Apply: %v. Fn: %v", et, et.Fn())
				}
				firstExec = false
			}
			fnType, fnCs := infer.t, infer.cs

			if err = infer.consGen(body); err != nil {
				return errors.Wrapf(err, "Unable to infer body of Apply: %v. Body: %v", et, body)
			}
			bodyType, bodyCs := infer.t, infer.cs

			tv := infer.Fresh()
			cs := append(fnCs, bodyCs...)
			cs = append(cs, Constraint{fnType, saveExprContext(NewFnType(bodyType, tv), &expr), CreateCodeContext(expr)})

			infer.t = tv
			infer.t = saveExprContext(infer.t, &expr)
			infer.cs = cs

			return nil
		})
		if batchErr != nil {
			return batchErr
		}

	case LetRec:
		tv := infer.Fresh()
		// env := infer.env // backup

		infer.env = infer.env.Clone()
		infer.env.Remove(et.Name())
		infer.env.Add(et.Name(), &Scheme{tvs: TypeVarSet{tv}, t: tv})

		if err = infer.consGen(et.Def()); err != nil {
			return errors.Wrapf(err, "Unable to infer the definition of a letRec %v. Def: %v", et, et.Def())
		}
		defType, defCs := infer.t, infer.cs

		s := newSolver()
		s.solve(defCs)
		if s.err != nil {
			return errors.Wrapf(s.err, "Unable to solve constraints of def: %v", defCs)
		}

		sc := Generalize(infer.env.Apply(s.sub).(Env), saveExprContext(defType.Apply(s.sub).(Type), &expr))

		infer.env.Remove(et.Name())
		infer.env.Add(et.Name(), sc)

		if err = infer.consGen(et.Body()); err != nil {
			return errors.Wrapf(err, "Unable to infer body of letRec %v. Body: %v", et, et.Body())
		}

		infer.t = infer.t.Apply(s.sub).(Type)
		infer.t = saveExprContext(infer.t, &expr)
		infer.cs = infer.cs.Apply(s.sub).(Constraints)
		infer.cs = append(infer.cs, defCs...)

	case Let:
		env := infer.env

		if err = infer.consGen(et.Def()); err != nil {
			return errors.Wrapf(err, "Unable to infer the definition of a let %v. Def: %v", et, et.Def())
		}
		defType, defCs := infer.t, infer.cs

		s := newSolver()
		s.solve(defCs)
		if s.err != nil {
			return errors.Wrapf(s.err, "Unable to solve for the constraints of a def %v", defCs)
		}

		sc := Generalize(env.Apply(s.sub).(Env), saveExprContext(defType.Apply(s.sub).(Type), &expr))
		infer.env = infer.env.Clone()
		infer.env.Remove(et.Name())
		infer.env.Add(et.Name(), sc)

		if err = infer.consGen(et.Body()); err != nil {
			return errors.Wrapf(err, "Unable to infer body of let %v. Body: %v", et, et.Body())
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

// Instantiate takes a fresh name generator, an a polytype and makes a concrete type out of it.
//
// If ...
// 		  Γ ⊢ e: T1  T1 ⊑ T
//		----------------------
//		       Γ ⊢ e: T
//
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

// Generalize takes an env and a type and creates the most general possible type - which is a polytype
//
// Generalization
//
// If ...
//		  Γ ⊢ e: T1  T1 ∉ free(Γ)
//		---------------------------
//		   Γ ⊢ e: ∀ α.T1
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
		// cannot return envFree because envFree will just be sorted and set
	case len(envFree) == 0 && len(tFree) > 0:
		// return ?
	}

	diff = tFree.Difference(envFree)

ret:
	return &Scheme{
		tvs: diff,
		t:   t,
	}
}

// Infer takes an env, and an expression, and returns a scheme.
//
// The Infer function is the core of the HM type inference system. This is a reference implementation and is completely servicable, but not quite performant.
// You should use this as a reference and write your own infer function.
//
// Very briefly, these rules are implemented:
//
// Var
//
// If x is of type T, in a collection of statements Γ, then we can infer that x has type T when we come to a new instance of x
//		 x: T ∈ Γ
//		-----------
//		 Γ ⊢ x: T
//
// Apply
//
// If f is a function that takes T1 and returns T2; and if x is of type T1;
// then we can infer that the result of applying f on x will yield a result has type T2
//		 Γ ⊢ f: T1→T2  Γ ⊢ x: T1
//		-------------------------
//		     Γ ⊢ f(x): T2
//
//
// Lambda Abstraction
//
// If we assume x has type T1, and because of that we were able to infer e has type T2
// then we can infer that the lambda abstraction of e with respect to the variable x,  λx.e,
// will be a function with type T1→T2
//		  Γ, x: T1 ⊢ e: T2
//		-------------------
//		  Γ ⊢ λx.e: T1→T2
//
// Let
//
// If we can infer that e1 has type T1 and if we take x to have type T1 such that we could infer that e2 has type T2,
// then we can infer that the result of letting x = e1 and substituting it into e2 has type T2
//		  Γ, e1: T1  Γ, x: T1 ⊢ e2: T2
//		--------------------------------
//		     Γ ⊢ let x = e1 in e2: T2
//
func Infer(env Env, expr Expression) (*Scheme, error) {
	if expr == nil {
		return nil, errors.Errorf("Cannot infer a nil expression")
	}

	if env == nil {
		env = make(SimpleEnv)
	}

	infer := newInferer(env)
	if err := infer.consGen(expr); err != nil {
		return nil, err
	}

	s := newSolver()
	s.solve(infer.cs)

	if s.err != nil {
		return nil, s.err
	}

	if infer.t == nil {
		return nil, errors.Errorf("infer.t is nil")
	}

	t := infer.t.Apply(s.sub).(Type)
	return closeOver(t)
}

// Unify unifies the two types and returns a list of substitutions.
// These are the rules:
//
// Type Constants and Type Constants
//
// Type constants (atomic types) have no substitution
//		c ~ c : []
//
// Type Variables and Type Variables
//
// Type variables have no substitutions if there are no instances:
// 		a ~ a : []
//
// Default Unification
//
// if type variable 'a' is not in 'T', then unification is simple: replace all instances of 'a' with 'T'
// 		     a ∉ T
//		---------------
//		 a ~ T : [a/T]
//
func Unify(a, b Type, context Constraint) (sub Subs, err error) {
	logf("%v ~ %v", a, b)
	enterLoggingContext()
	defer leaveLoggingContext()

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
	err = errors.Errorf("Unification Fail: %v [%s] ~ %v [%s] cannot be unified", a, a.GetContext().String(), b, b.GetContext().String())
	return
}

func unifyMany(a, b Types, contextA, contextB Type, context Constraint) (sub Subs, err error) {
	logf("UnifyMany %v %v", a, b)
	enterLoggingContext()
	defer leaveLoggingContext()

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
	// case tv == t:
	case occurs(tv, t):
		err = UnificationRecurrentTypeError{
			Type:       t,
			Variable:   tv,
			VariableTypeSource: tvt,
			Constraint: context,
		}
	default:
		ssub := BorrowSSubs(1)
		ssub.s[0] = Substitution{tv, t}
		sub = ssub
	}
	logf("Sub %v", sub)
	return
}

func occurs(tv TypeVariable, s Substitutable) bool {
	ftv := s.FreeTypeVar()
	defer ReturnTypeVarSet(ftv)

	return ftv.Contains(tv)
}

func closeOver(t Type) (sch *Scheme, err error) {
	sch = Generalize(nil, t)
	err = sch.Normalize()
	logf("closeoversch: %v", sch)
	return
}

package hindley_milner

import (
	"fmt"
	"reflect"

	"github.com/alecthomas/repr"
	"github.com/pkg/errors"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type Cloner interface {
	Clone() interface{}
}

type Fresher interface {
	Fresh() TypeVariable
}

func saveExprContext(t Type, source *generic_ast.Expression) Type {
	return t.WithContext(CreateCodeContext(*source))
}

func wrapEnvError(err error, source *generic_ast.Expression) error {
	if e, ok := err.(*envError); ok {
		if e.builtin {
			return BuiltinRedefinedError{
				Name:    e.variableName,
				Context: CreateCodeContext(*source),
			}
		}
		return VariableRedefinedError{
			Name:               e.variableName,
			PreviousDefinition: e.oldDef,
			Context:            CreateCodeContext(*source),
		}
	}
	return err
}

type IntrospectionConstraint struct {
	tv      TypeVariable
	argTV   Type
	context CodeContext
}

type OverloadConstraint struct {
	tv           TypeVariable
	name         string
	alternatives []*Scheme
	context      CodeContext
}

type InferenceBackend interface {
	Fresher
	GenerateConstraints(expr generic_ast.Expression, forceType ExpressionType, isTop bool, isOpaqueTop bool) (err error)
	GetEnv() Env
	GetReturnEnv() Env
	ProgramType() Type
	Constraints() Constraints
	OverrideConstraints(c Constraints)
	GetOverloadConstraints() []OverloadConstraint
	GetIntrospectionConstraints() []IntrospectionConstraint
}

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

	case len(envFree) == 0 && len(tFree) > 0:

	}

	diff = tFree.Difference(envFree)

ret:
	return &Scheme{
		tvs: diff,
		t:   t,
	}
}

type InferConfiguration struct {
	CreateDefaultEmptyType        func() *Scheme
	OnConstrintGenerationStarted  *func()
	OnConstrintGenerationFinished *func()
	OnSolvingStarted              *func()
	OnSolvingFinished             *func()
	OnPostprocessingStarted       *func()
	OnPostprocessingFinished      *func()
}

func NewInferConfiguration() *InferConfiguration {
	return &InferConfiguration{
		CreateDefaultEmptyType: func() *Scheme { return nil },
	}
}

func Infer(env Env, expr generic_ast.Expression, config *InferConfiguration, infer InferenceBackend) (*Scheme, Env, error) {
	env.RegisterIntrospectionListener(NewIntrospecionSimpleListener())

	if expr == nil {
		return nil, nil, errors.Errorf("Cannot infer a nil expression")
	}

	if env == nil {
		env = CreateSimpleEnv(map[string][]*Scheme{})
	}

	if config.OnConstrintGenerationStarted != nil {
		(*config.OnConstrintGenerationStarted)()
	}
	err := infer.GenerateConstraints(expr, E_NONE, true, true)
	if config.OnConstrintGenerationFinished != nil {
		(*config.OnConstrintGenerationFinished)()
	}
	if err != nil {
		return nil, nil, err
	}

	s := newSolver()
	if config.OnSolvingStarted != nil {
		(*config.OnSolvingStarted)()
	}
	cs := infer.Constraints()
	inferEnv := infer.GetEnv()
	s.solve(cs, inferEnv.GetIntrospecionListener())
	if config.OnSolvingFinished != nil {
		(*config.OnSolvingFinished)()
	}
	if s.err != nil {
		return nil, nil, s.err
	}

	if config.OnPostprocessingStarted != nil {
		(*config.OnPostprocessingStarted)()
	}
	cleanCS := infer.Constraints()
	for _, ocs := range infer.GetOverloadConstraints() {
		hasCleanRun := false
		cs := Constraints{}
		for _, alt := range ocs.alternatives {
			cs = Constraints{}

			for _, c := range cleanCS {
				cs = append(cs, c)
			}
			cs = append(cs, Constraint{
				a:       ocs.tv,
				b:       Instantiate(infer, alt),
				context: ocs.tv.context,
			})
			s2 := newSolver()
			s2.solve(cs, infer.GetEnv().GetIntrospecionListener())
			if s2.err == nil {

				hasCleanRun = true
				break
			}
		}
		if hasCleanRun {
			cleanCS = cs
		} else {
			if config.OnPostprocessingFinished != nil {
				(*config.OnPostprocessingFinished)()
			}
			return nil, nil, InvalidOverloadCandidatesError{
				Name:       ocs.name,
				Candidates: ocs.alternatives,
				Context:    ocs.context,
			}
		}
	}
	if config.OnPostprocessingFinished != nil {
		defer (*config.OnPostprocessingFinished)()
	}
	infer.OverrideConstraints(cleanCS)

	if s.err != nil {
		return nil, nil, s.err
	}

	if infer.ProgramType() == nil {
		return nil, nil, errors.Errorf("infer.t is nil")
	}

	t := infer.ProgramType().Apply(s.sub).(Type)
	ret, err := closeOver(t)
	if err != nil {
		return nil, nil, err
	}

	//fmt.Printf("PIZDO TRY TO GET ALL INTROSPECTION\n")
	//fmt.Printf("%v\n", infer.cs)

	for _, ics := range infer.GetIntrospectionConstraints() {
		is := newSolver()
		is.solve(infer.Constraints(), infer.GetEnv().GetIntrospecionListener())
		h := infer.GetReturnEnv().GetIntrospecionListener().GetIntrospectionVariable(ics.tv).Apply(is.sub).(Type)
		introExpr := (*ics.context.Source).(IntrospectionExpression)
		introExpr.OnTypeReturned(h)
	}

	//fmt.Printf("EEEEEEND PIZDO TRY TO GET ALL INTROSPECTION!\n")

	return ret, infer.GetReturnEnv(), nil
}

func Unify(a, b Type, context Constraint, listener IntrospecionListener) (sub Subs, err error) {
	//fmt.Printf("UNION %v %v\n", a, b)

	// aTV, bTV := false, false

	// _, aTV = a.(*TypeVariable)
	// _, bTV = b.(*TypeVariable)

	//a, b = b, a

	if unionA, ok := a.(*Union); ok {
		allSubs := []Subs{}
		for _, v := range unionA.types {
			if subs, err := Unify(b, v.(Type), context, listener); err != nil {
				return nil, fmt.Errorf("Cannot match union type: %w", err)
			} else {
				allSubs = append(allSubs, subs)
			}
		}
		return SubsConcat(allSubs...), nil
	}
	if unionB, ok := b.(*Union); ok {
		allSubs := []Subs{}
		var lastErr error = nil
		for _, v := range unionB.types {
			if subs, err := Unify(a, v.(Type), context, listener); err != nil {
				lastErr = fmt.Errorf("Cannot match union type: %w", err)
			} else {
				allSubs = append(allSubs, subs)
				lastErr = nil
				break
			}
		}
		if lastErr != nil {
			return nil, lastErr
		}
		return SubsConcat(allSubs...), nil
	}
	//return SubsConcat(allSubs...), nil

	if sa, ok := a.(UnionableType); ok {
		if reflect.TypeOf(a) == reflect.TypeOf(b) {
			if TypeEq(a, b) {
				return nil, nil
			}

			defer ReturnTypes(a.Types())
			defer ReturnTypes(b.Types())

			subs, err := sa.Union(b, context, listener)
			if err != nil {
				return nil, UnificationWrongTypeError{
					TypeA:      a,
					TypeB:      b,
					Constraint: context,
					Details:    err.Error(),
				}
			}
			return subs, nil
		}
	}

	switch at := a.(type) {
	case TypeVariable:
		return bind(at, b, context, a, listener)
	default:
		if TypeEq(a, b) {
			return nil, nil
		}

		if btv, ok := b.(TypeVariable); ok {
			return bind(btv, a, context, b, listener)
		}
		atypes := a.Types()
		btypes := b.Types()
		defer ReturnTypes(atypes)
		defer ReturnTypes(btypes)

		if len(atypes) == 0 && len(btypes) == 0 {
			goto e
		}

		return unifyMany(atypes, btypes, a, b, context, listener)

	e:
	}
	err = UnificationWrongTypeError{
		TypeA:      a,
		TypeB:      b,
		Constraint: context,
	}
	return
}

func unifyMany(a, b Types, contextA, contextB Type, context Constraint, listener IntrospecionListener) (sub Subs, err error) {

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
		if s2, err = Unify(at, bt, context, listener); err != nil {
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

func bind(tv TypeVariable, t Type, context Constraint, tvt Type, listener IntrospecionListener) (sub Subs, err error) {
	logf("Binding %v to %v", tv, t)
	listener.OnApplySingle(tv, t)

	switch {

	case occurs(tv, t):
		if TypeEq(tv, t) {
			ssub := BorrowSSubs(1)
			ssub.s[0] = Substitution{tv, t}
			sub = ssub
			return
		}
		logf("KURWA 1\n%v\nKURWA 2\n%v\nKURWA 3\n%v\nKURWA 4\n%v\n",
			repr.String(*t.GetContext().Source),
			repr.String(*tv.GetContext().Source),
			repr.String(*tvt.GetContext().Source),
			repr.String(*context.context.Source))

		err = UnificationRecurrentTypeError{
			Type:               t,
			Variable:           tv,
			VariableTypeSource: tvt,
			Constraint:         context,
		}
	default:
		ssub := BorrowSSubs(1)
		ssub.s[0] = Substitution{tv, t}
		sub = ssub
	}

	return
}

func occurs(tv TypeVariable, s Substitutable) bool {
	ftv := s.FreeTypeVar()
	defer ReturnTypeVarSet(ftv)

	return ftv.Contains(tv)
}

func closeOver(t Type) (sch *Scheme, err error) {
	logf("PATRZ SPERMA: %v\n", t)
	sch = Generalize(nil, t)
	err = sch.Normalize()

	return
}

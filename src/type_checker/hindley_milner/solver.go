package hindley_milner

type solver struct {
	sub Subs
	err error
}

func newSolver() *solver {
	return new(solver)
}

func (s *solver) solve(cs Constraints, listener IntrospecionListener) {

	//logf("SOLVE CALL\n")

	if s.err != nil {
		return
	}

	switch len(cs) {
	case 0:
		return
	default:
		var sub Subs
		c := cs[0]
		sub, s.err = Unify(c.a, c.b, c, listener)
		defer ReturnSubs(s.sub)

		s.sub = compose(sub, s.sub)
		cs = cs[1:].Apply(s.sub).(Constraints)
		s.solve(cs, listener)

	}

	return
}

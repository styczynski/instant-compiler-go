package hindley_milner

import "fmt"

type solver struct {
	sub Subs
	err error
}

func newSolver() *solver {
	return new(solver)
}

func (s *solver) solve(cs Constraints, listener IntrospecionListener) {

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		err := recover()
		if err != nil {
			fmt.Printf("PIZDA NA RYJ! %v\n", err)
			panic(err)
		} else {
			fmt.Printf("SEEMS TO BE OK?\n")
		}
	}()

	logf("SOLVE CALL: %v\n", cs)

	if s.err != nil {
		return
	}

	switch len(cs) {
	case 0:
		return
	default:
		var sub Subs
		fmt.Printf("SOLVE A\n")
		c := cs[0]
		fmt.Printf("SOLVE B\n")
		sub, s.err = Unify(c.a, c.b, c, listener)
		fmt.Printf("SOLVE C\n")
		defer ReturnSubs(s.sub)

		fmt.Printf("SOLVE D\n")
		s.sub = compose(sub, s.sub)
		fmt.Printf("SOLVE E\n")
		cs = cs[1:].Apply(s.sub).(Constraints)
		fmt.Printf("SOLVE F\n")
		s.solve(cs, listener)

	}

	return
}

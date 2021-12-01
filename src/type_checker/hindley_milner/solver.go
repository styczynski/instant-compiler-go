package hindley_milner

import (
	"github.com/styczynski/latte-compiler/src/logs"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type solver struct {
	sub Subs
	err error
}

func newSolver() *solver {
	return new(solver)
}

func (s *solver) LogContext(c *context.ParsingContext) map[string]interface{} {
	return map[string]interface{}{}
}

func (s *solver) solve(infer InferenceBackend, cs Constraints, listener IntrospecionListener) {
	logs.Debug(s, "Solve constraints: %v", cs)

	if s.err != nil {
		return
	}

	switch len(cs) {
	case 0:
		return
	default:
		var sub Subs
		logs.Debug(s, "Obtain constraint")
		c := cs[0]
		logs.Debug(s, "Perform unification")
		sub, s.err = Unify(c.a, c.b, c, infer, listener)
		defer ReturnSubs(s.sub)

		logs.Debug(s, "Compose substitutions")
		s.sub = compose(sub, s.sub)
		logs.Debug(s, "Apply substitutions")
		cs = cs[1:].Apply(s.sub).(Constraints)
		logs.Debug(s, "Launch recursive solve")
		s.solve(infer, cs, listener)
	}

	return
}

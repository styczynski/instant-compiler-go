package core

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/sat_solver"
)

type ColourableGraph interface {
}

func GraphColouring() {
	ast := &sat_solver.Entry{
		Formula: &sat_solver.Formula{
			And: &sat_solver.And{
				Arg1: &sat_solver.Formula{Variable: &sat_solver.Variable{Name: "a"}},
				Arg2: &sat_solver.Formula{Variable: &sat_solver.Variable{Name: "b"}},
			},
		},
	}
	context := sat_solver.NewSATContext(sat_solver.SATConfiguration{
		InputFile:              "",
		EnableSelfVerification: false,
		EnableEventCollector:   true,
		EnableSolverTracing:    true,
		EnableCNFConversion:    true,
		EnableASTOptimization:  false,
		EnableCNFOptimizations: false,
		SolverName:             "cdcl",
		LoaderName:             "haskell",
	})
	err, result := RunSATSolverOnLoadedFormula(ast, context)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("SOLVER RESULT = %s", result.Brief())
}

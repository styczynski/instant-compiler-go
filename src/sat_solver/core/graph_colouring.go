package core

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/sat_solver"
	"github.com/yourbasic/graph"
)

func convertGraphColouringAssignment(g graph.Iterator, assignment map[string]bool) (error, map[int]int) {
	colouring := map[int]int{}
	n := g.Order()
	// j is graph vertex
	for j := 0; j < n; j++ {
		// i is index of colour
		for i := 0; i < n; i++ {
			if val, ok := assignment[fmt.Sprintf("x[%d,%d]", i, j)]; ok && val {
				if previousi, ok := colouring[j]; ok {
					return fmt.Errorf("Node %d has two colours assigned: %d and %d", j, previousi, i), nil
				}
				colouring[j] = i
			}
		}
	}
	return nil, colouring
}

func buildGraphColouringFormula(g graph.Iterator) *sat_solver.SATFormula {
	n := g.Order() // Order returns the number of vertices.
	cnf := &sat_solver.CNFFormula{
		Variables: []sat_solver.CNFClause{},
	}

	// Assumptions:
	// j is graph vertex (0 to n-1)
	// i is index of colour (0 to n-1)
	// x[i,j]      =   vars.Get(fmt.Sprintf("x[%d,%d]", i, j))
	// not(x[i,j]) =  -vars.Get(fmt.Sprintf("x[%d,%d]", i, j))

	vars := sat_solver.NewSATVariableMapping()

	// j is graph vertex
	for j := 0; j < n; j++ {
		vertexJHasAssignedColour := sat_solver.CNFClause{}
		// i is index of colour
		for i := 0; i < n; i++ {
			// We generate: x[0,j] or x[1,j] or x[2, j]... or x[n-1,j]
			vertexJHasAssignedColour = append(vertexJHasAssignedColour, vars.Get(fmt.Sprintf("x[%d,%d]", i, j)))
			for i2 := 0; i2 < n; i2++ {
				if i != i2 {
					// not(x[i,j]) or not(x[i2,j])
					cnf.Variables = append(cnf.Variables, sat_solver.CNFClause{
						-vars.Get(fmt.Sprintf("x[%d,%d]", i, j)),
						-vars.Get(fmt.Sprintf("x[%d,%d]", i2, j)),
					})
				}
			}
		}
		cnf.Variables = append(cnf.Variables, vertexJHasAssignedColour)
	}

	// i is index of colour
	for i := 0; i < n; i++ {
		// Iterate all pairs of connected indices
		// j is graph vertex
		for j := 0; j < n; j++ {
			g.Visit(j, func(j2 int, c int64) (skip bool) {
				// not(x[i,j]) or not(x[i,j2])
				cnf.Variables = append(cnf.Variables, sat_solver.CNFClause{
					-vars.Get(fmt.Sprintf("x[%d,%d]", i, j)),
					-vars.Get(fmt.Sprintf("x[%d,%d]", i, j2)),
				})
				return
			})
		}
	}

	return sat_solver.NewSATFormula(cnf, vars, nil)
}

func TestComputeGraphColouring() {
	g := graph.New(4)
	g.AddBoth(0, 1) //  0 -- 1
	g.AddBoth(0, 2) //  |    |
	g.AddBoth(2, 3) //  2 -- 3
	g.AddBoth(1, 3)

	err, colouring := ComputeGraphColouring(g)
	if err != nil {
		panic(err)
	}

	for node, colour := range colouring {
		//fmt.Printf("\nNode %d has colour %d", node, colour)
	}
}

func ComputeGraphColouring(g graph.Iterator) (error, map[int]int) {
	context := sat_solver.NewSATContext(sat_solver.SATConfiguration{
		InputFile:              "",
		EnableSelfVerification: false,
		EnableEventCollector:   false,
		EnableSolverTracing:    false,
		EnableCNFConversion:    true,
		EnableASTOptimization:  false,
		EnableCNFOptimizations: false,
		SolverName:             "cdcl",
		LoaderName:             "cnf",
	})
	formula := buildGraphColouringFormula(g)
	err, result := RunSATSolverOnLoadedFormula(formula, context)
	if err != nil {
		return err, nil
	}
	if result.IsSAT() {
		err, colouring := convertGraphColouringAssignment(g, result.GetSatisfyingAssignment())
		if err != nil {
			return err, nil
		}
		return nil, colouring
	}

	return fmt.Errorf("Generated formula that is unsat"), nil
}

package main

import (
	"github.com/skiesel/expsys/rdb"
)

func main() {
	filters := []map[string]string{
		map[string]string{
			"algorithm": "astar",
		},
		map[string]string{
			"algorithm": "speedy",
		},
	}

	names := []string{"A*", "Speedy",}

	dss := rdb.GetDatasets("data", filters, names)

	plotSolutionCosts("Solution Costs", dss)
	plotSolutionCostsFactorOfBest("Solution Costs Factor Best", dss)
}


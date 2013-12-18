package main

import (
	"github.com/skiesel/expsys/rdb"
	"github.com/skiesel/expsys/plots/standardplots"
	"github.com/skiesel/expsys/plots/realtimeplots"
	"github.com/skiesel/expsys/tables/standardtables"
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

	dss := rdb.GetDatasetsWithPathKeys("data", filters, names)

	standardplots.PlotSolutionCosts("Solution Costs", dss)
	standardplots.PlotSolutionCostsFactorOfBest("Solution Costs Factor Best", dss)


	realtimeplots.TrafficCollisionWFFiltered(dss)

	standardtables.SolutionCostSumsTable(dss)
	standardtables.SolutionCostMeansTable(dss)
}
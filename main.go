package main

import (
	"github.com/skiesel/expsys/plots/standardplots"
	"github.com/skiesel/expsys/rdb"
	//	"github.com/skiesel/expsys/plots/realtimeplots"
	//	"github.com/skiesel/expsys/tables/standardtables"
	"fmt"
	"math"
	"strconv"
)

func main() {
	filters := []map[string]string{
		map[string]string{
			"alg": "rbfs",
			"cost": "unit",
		},
		map[string]string{
			"alg": "rbfs_cr",
			"cost": "unit",
		},
		map[string]string{
			"alg": "rbfs_cr_opt",
			"cost": "unit",
		},
	}

	names := []string{"RBFS", "RBFS-CR", "RBFS-CR opt"}

	dss := rdb.GetDatasetsWithPathKeys("/Users/skiesel/Desktop/matt-java/edu.unh.ai.search.ml/bin/data", filters, names)

	solvedTest := func(solutionCost string) bool {
		val, err := strconv.ParseFloat(solutionCost, 0)
		if err != nil {
			return false
		}
		return val >= 0
	}

  dss = rdb.FilterOutUnsolved(dss, "num", solvedTest, "solution cost")



	log10 := func(v string) string {
		val, err := strconv.ParseFloat(v, 0)
		if err != nil {
			fmt.Printf("Could not convered (%s) to float64\n", v)
			panic(err)
		}
		return strconv.FormatFloat(math.Log10(val), 'f', 15, 64)
	}

	for _, ds := range dss {
		ds.AddTransformedKey("total cpu time", log10, "log10 total cpu time")
	}

	standardplots.PlotXvsY("CPU Time By Instance", dss, "log10 total cpu time", "CPU time (log10)", "num", "Instance")
	standardplots.PlotXvsFactorBestY("CPU Time By Instance (factor best)", dss, "total cpu time", "CPU time (factor best)", "num", "Instance")
	standardplots.PlotXvsY("Expansions By Instance", dss, "total nodes expanded", "expanded nodes", "num", "Instance")
	standardplots.PlotXvsFactorBestY("Expansions By Instance (factor best)", dss, "total nodes expanded", "expanded nodes", "num", "Instance")

	// realtimeplots.TrafficCollisionWFFiltered(dss)

	// standardtables.SolutionCostSumsTable(dss)
	// standardtables.SolutionCostMeansTable(dss)
}

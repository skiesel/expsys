package main

import (
	"github.com/skiesel/expsys/rdb"
	"github.com/skiesel/expsys/plots/standardplots"
//	"github.com/skiesel/expsys/plots/realtimeplots"
//	"github.com/skiesel/expsys/tables/standardtables"
	"strconv"
	"fmt"
	"math"
)

func main() {
	filters := []map[string]string{
		map[string]string{
			"alg": "rbfs",
		},
		map[string]string{
			"alg": "rbfs_cr",
		},
		map[string]string{
			"alg": "rbfs_cr_optimized",
		},
	}

	names := []string{"RBFS", "RBFS-CR", "RBFS-CR opt",}

	dss := rdb.GetDatasetsWithPathKeys("/Users/skiesel/Desktop/matt-java/edu.unh.ai.search.ml/bin/data", filters, names)

	log10 := func(v string)(string) {
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


	// realtimeplots.TrafficCollisionWFFiltered(dss)

	// standardtables.SolutionCostSumsTable(dss)
	// standardtables.SolutionCostMeansTable(dss)
}
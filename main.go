package main

import (
	"fmt"
	"github.com/skiesel/expsys/plots/standardplots"
	"github.com/skiesel/expsys/rdb"
	"strconv"
)

func main() {

	//astar  idastar  idastar_cr  KEY=alg  rbfs  rbfs_cr  rbfs_cr_global  rbfs_epsilon  rbfs_greedy  rbfs_sqrt  speedy  tmp  wastar  wrbfs_cr_global
	filters := []map[string]string{
		map[string]string{"alg": "rbfs"},
		map[string]string{"alg": "rbfs_cr"},
		map[string]string{"alg": "rbfs_cr_global"},
		/*		map[string]string{ "alg": "rbfs_epsilon", },
				map[string]string{ "alg": "rbfs_sqrt", },
				map[string]string{ "alg": "rbfs_greedy", },
				map[string]string{ "alg": "idastar", "weight": "1.0" },
				map[string]string{ "alg": "idastar_cr", }, */
	}

	for _, m := range filters {
		m["domain"] = "tiles"
		m["cost"] = "unit"
		m["rows"] = "4"
		m["cols"] = "4"
	}

	names := []string{
		"RBFS",
		"RBFS-cr",
		"RBFS-cr (global)",
		/*			"RBFS-eps",
					"RBFS-sqrt",
					"RBFS-greedy",
					"IDA*",
					"IDA*-cr", */
	}

	dss := rdb.GetDatasetsWithPathKeys("/Users/skiesel/Desktop/data", filters, names)

	solvedTest := func(solutionCost string) bool {
		val, err := strconv.ParseFloat(solutionCost, 0)
		if err != nil {
			return false
		}
		return val >= 0
	}

	for i, ds := range dss {
		fmt.Printf("%d) %s : %d\n", i, names[i], ds.GetSize())
	}

	dss = rdb.FilterOutUnsolved(dss, "num", solvedTest, "solution cost")

	for i, ds := range dss {
		fmt.Printf("%d) %s : %d\n", i, names[i], ds.GetSize())
	}

	standardplots.PlotXvsY("CPU Time By Instance", dss, "total cpu time", "CPU time", "num", "Instance")
	standardplots.PlotXvsFactorBestY("CPU Time By Instance (factor best)", dss, "total cpu time", "CPU time (factor best)", "num", "Instance")
	standardplots.PlotXvsY("Expansions By Instance", dss, "total nodes expanded", "expanded nodes", "num", "Instance")
	standardplots.PlotXvsFactorBestY("Expansions By Instance (factor best)", dss, "total nodes expanded", "expanded nodes", "num", "Instance")

}

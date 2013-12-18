package standardtables

import (
	"github.com/skiesel/expsys/rdb"
	"github.com/skiesel/expsys/tables"
	"fmt"
)

func SumsTable(dss []*rdb.Dataset, valuesKey string) {
	fmt.Printf("Name\tSum\n")
	for _, ds := range dss {
		sum := tables.Sum(ds.GetDatasetFloatValues(valuesKey))
		fmt.Printf("%s\t%f\n", ds.GetName(), sum)
	}
}

func SolutionCostSumsTable(dss []*rdb.Dataset) {
	SumsTable(dss, "final sol cost")
}

func MeansTable(dss []*rdb.Dataset, valuesKey string) {
	fmt.Printf("Name\tMean\tStdDev\n")
	for _, ds := range dss {		
		mean, stddev, _ := tables.MeanStdDevVariance(ds.GetDatasetFloatValues(valuesKey))
		fmt.Printf("%s\t%f\t%f\n", ds.GetName(), mean, stddev)
	}
}

func SolutionCostMeansTable(dss []*rdb.Dataset) {
	MeansTable(dss, "final sol cost")
}
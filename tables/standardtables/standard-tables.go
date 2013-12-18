package standardtables

import (
	"github.com/skiesel/expsys/rdb"
	"github.com/skiesel/expsys/tables"
	"fmt"
)

func SumsTable(dss []*rdb.Dataset, valuesKey string) {
	fmt.Printf("Name\tSum\n")
	for _, ds := range dss {
		sum := ds.GetDatasetSum(valuesKey)
		fmt.Printf("%s\t%f\n", ds.GetName(), sum)
	}
}

func SolutionCostSumsTable(dss []*rdb.Dataset) {
	SumsTable(dss, "final sol cost")
}

func AveragesTable(dss []*rdb.Dataset, valuesKey string) {
	fmt.Printf("Name\tAverage\tStdDev\n")
	for _, ds := range dss {		
		average := ds.GetDatasetAverage(valuesKey)
		stddev, _ := tables.StdDevAndVariance(ds.GetDatasetFloatValues(valuesKey))
		fmt.Printf("%s\t%f\t%f\n", ds.GetName(), average, stddev)
	}
}

func SolutionCostAveragesTable(dss []*rdb.Dataset) {
	AveragesTable(dss, "final sol cost")
}
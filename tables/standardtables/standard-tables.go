package standardtables

import (
	"github.com/skiesel/expsys/rdb"
	"github.com/skiesel/expsys/tables"
	"fmt"
)

func SumsTable(dss []*rdb.Dataset, valuesKey string) {
	fmt.Printf("Name\tSum\n")
	for i := range dss {
		sum := dss[i].GetDatasetSum(valuesKey)
		fmt.Printf("%s\t%f\n", dss[i].GetName(), sum)
	}
}

func SolutionCostSumsTable(dss []*rdb.Dataset) {
	SumsTable(dss, "final sol cost")
}

func AveragesTable(dss []*rdb.Dataset, valuesKey string) {
	fmt.Printf("Name\tAverage\tStdDev\n")
	for i := range dss {		
		average := dss[i].GetDatasetAverage(valuesKey)
		stddev, _ := tables.StdDevAndVariance(dss[i].GetDatasetFloatValues(valuesKey))
		fmt.Printf("%s\t%f\t%f\n", dss[i].GetName(), average, stddev)
	}
}

func SolutionCostAveragesTable(dss []*rdb.Dataset) {
	AveragesTable(dss, "final sol cost")
}
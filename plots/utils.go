package datautils

import (
	"github.com/skiesel/expsys/rdb"
)

func MatchKeys(targetIds []string, ids []string, values []float64) ([]string, []float64) {
	newids := make([]string, len(ids))
	copy(newids, ids)
	newvalues := make([]float64, len(values))
	copy(newvalues, values)

	for i := range targetIds {
		for j := i; j < len(newids); j++ {
			if targetIds[i] == newids[j] {
				if i != j {
					temp := newids[i]
					newids[i] = newids[j]
					newids[j] = temp
					temp2 := newvalues[i]
					newvalues[i] = newvalues[j]
					newvalues[j] = temp2
				}
				break
			}
		}
	}
	return newids, newvalues
}

func SortBothArrays(ids []string, values []float64) ([]string, []float64) {
	newids := make([]string, len(ids))
	copy(newids, ids)
	newvalues := make([]float64, len(values))
	copy(newvalues, values)

	for i := range newvalues {
		min := newvalues[i]
		minIndex := i
		for j := i; j < len(ids); j++ {
			if newvalues[j] < min {
				min = newvalues[j]
				minIndex = j
			}
		}
		temp := newids[i]
		newids[i] = newids[minIndex]
		newids[minIndex] = temp
		temp2 := newvalues[i]
		newvalues[i] = newvalues[minIndex]
		newvalues[minIndex] = temp2
	}

	return newids, newvalues
}

func Group(ds *rdb.Dataset, key string) map[string]*rdb.Dataset {
	values := ds.GetDatasetStringValues(key)
	valueSet := map[string]string{}
	for _, value := range values {
		_, exists := valueSet[value]
		if !exists {
			valueSet[value] = value
		}
	}

	grouped := map[string]*rdb.Dataset{}

	for value, _ := range valueSet {
		filter := func(str string) bool {
			return str == value
		}
		grouped[value] = ds.FilterDataset(filter, key).RenameDataset(ds.GetName() + " " + value)
	}
	return grouped
}

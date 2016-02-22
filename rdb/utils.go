package rdb

import (
	"fmt"
	"math"
	"strconv"
)

func FilterOutUnsolved(dss []*Dataset, identifier string,
	testSolved func(string) bool, solvedKey string) []*Dataset {
	solved := map[string]bool{}

	for i, ds := range dss {
		if i == 0 {
			for _, df := range ds.datafiles {
				if df.hasKey(solvedKey) && testSolved(df.getStringValue(solvedKey)) {
					solved[df.getStringValue(identifier)] = true
				}
			}
		} else {
			for key, _ := range solved {
				isSolved := false
				for _, df := range ds.datafiles {
					if df.hasKey(identifier) && df.getStringValue(identifier) == key {
						if df.hasKey(solvedKey) && testSolved(df.getStringValue(solvedKey)) {
							isSolved = true
						}
						break
					}
				}
				if !isSolved {
					delete(solved, key)
				}
			}
		}
	}

	filtered := make([]*Dataset, len(dss))

	includeThese := func(id string) bool {
		_, found := solved[id]
		return found
	}

	for i, ds := range dss {
		filtered[i] = ds.FilterDataset(includeThese, identifier)
	}
	return filtered
}

func AddFactorBest(dss []*Dataset, identifier string, key string, newKey string) []*Dataset {

	bests := map[string]float64{}

	for _, ds := range dss {
		values, ids := ds.GetDatasetFloatValuesPair(key, identifier)

		for j := range values {
			best, bound := bests[ids[j]]
			if !bound || values[j] < best {
				bests[ids[j]] = values[j]
			}
		}
	}

	clampValue := 1. / 1000000000.

	for key, val := range bests {
		if val == 0 {
			bests[key] = clampValue
		}
	}

	newDss := make([]*Dataset, len(dss))

	for i, ds := range dss {
		newDss[i] = ds.copyDataset()
	}

	for _, ds := range newDss {
		for _, df := range ds.datafiles {
			current := df.getFloatValue(key)

			if current == 0 {
				fmt.Printf("WARN -- CLAMPING:\n")
				fmt.Printf("%s : %s : %g -> %g\n", df.getStringValue("alg"), df.getStringValue(identifier), current, clampValue)
				current = clampValue
			}

			best := bests[df.getStringValue(identifier)]
			df.addKey(newKey, strconv.FormatFloat(current/best, 'f', 15, 64))
		}
	}

	return newDss
}

func AddLog10(dss []*Dataset, key string) []*Dataset {

	log10 := func(v string) string {
		val, err := strconv.ParseFloat(v, 0)
		if err != nil {
			fmt.Printf("Could not convert (%s) to float64\n", v)
			panic(err)
		}
		return strconv.FormatFloat(math.Log10(val), 'f', 15, 64)
	}

	new_key := key + " log10"
	for _, ds := range dss {
		ds.AddTransformedKey(key, log10, new_key)
	}

	return dss
}

func GroupByKey(dss []*Dataset, key string) map[string][]*Dataset {
	groupKeys := map[string]bool{}
	grouped := map[string][]*Dataset{}

	for _, ds := range dss {
		for _, df := range ds.datafiles {
			groupKeys[df.getStringValue(key)] = true
		}
	}

	groupKeysArray := make([]string, len(groupKeys))
	i := 0
	for key, _ := range groupKeys {
		groupKeysArray[i] = key
		i++
	}

	for _, groupKey := range groupKeysArray {
		filter := func(val string) bool { return val == groupKey }
		grouped[groupKey] = make([]*Dataset, len(dss))
		for i, ds := range dss {
			grouped[groupKey][i] = ds.FilterDataset(filter, key)
		}
	}

	return grouped
}

package rdb

import (
	"strconv"
)

func FilterOutUnsolved(dss []*Dataset, identifier string,
															testSolved func(string)bool, solvedKey string) ([]*Dataset) {
	solved := make(map[string]bool)

	for i, ds := range dss {
		
		if i == 0 {
			for _, df := range ds.datafiles {
				solved[df.getStringValue(identifier)] = true
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

	includeThese := func(id string)bool {
		_, found := solved[id]
		return found
	}

	for i, ds := range dss {
		filtered[i] = ds.FilterDataset(includeThese, identifier)
	}
	return filtered
}


func AddFactorBest(dss []*Dataset, identifier string, key string, newKey string) ([]*Dataset) {

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
	
	newDss := make([]*Dataset, len(dss))

	for i, ds := range dss {
		newDss[i] = ds.copyDataset()
	}

	for _, ds := range newDss {
		for _, df := range ds.datafiles {
			current := df.getFloatValue(key)
			best := bests[df.getStringValue(identifier)]
			df.addKey(newKey, strconv.FormatFloat(current / best, 'f', 15, 64))
		}
	}

	return newDss
}
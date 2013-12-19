package rdb

func FilterIntersectionSolved(dss []*Dataset, identifier string,
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
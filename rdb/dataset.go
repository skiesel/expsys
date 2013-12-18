package rdb

type Dataset struct {
	name string
	datafiles []*Datafile
}

func newDataset(name string, files []string) *Dataset {
	ds := new(Dataset)
	ds.name = name
	ds.datafiles = make([]*Datafile, len(files))

	for i := range files {
		ds.datafiles[i] = newDatafile(files[i])
	}

	return ds
}

func (ds Dataset) FilterDataset(filter func(string)bool, key string) (filtered *Dataset) {
	filtered.name = ds.name
	for _, df := range ds.datafiles {
		if filter(df.getStringValue(key)) {
			filtered.datafiles = append(filtered.datafiles, df)
		}
	}
	return
}

func (ds Dataset) TestDataset(test func(string)bool, key string) bool {
	for _, df := range ds.datafiles {
		if !test(df.getStringValue(key)) {
			return false
		}
	}
	return true
}

func (ds Dataset) addPathKeys(baseDirectory string) {
	for _, df := range ds.datafiles {
		df.addPathKeys(baseDirectory)
	}
}

func (ds Dataset) GetDatasetFloatValues(key string) (values []float64) {
	for _, df := range ds.datafiles {
		values = append(values, df.getFloatValue(key))
	}
	return
}

func (ds Dataset) GetDatasetStringValues(key string) (values []string) {
	for _, df := range ds.datafiles {
		values = append(values, df.getStringValue(key))
	}
	return
}

func (ds Dataset) GetDatasetIntegerValues(key string) (values []int64) {
	for _, df := range ds.datafiles {
		values = append(values, df.getIntegerValue(key))
	}
	return
}

func (ds Dataset) GetDatasetFloatValuesPair(key string, id string) (values []float64, ids []string) {
	for _, df := range ds.datafiles {
		ids = append(ids, df.getStringValue(id))
		values = append(values, df.getFloatValue(key))
	}
	return
}

func (ds Dataset) GetDatasetSum(key string) (values float64) {
	for _, df := range ds.datafiles {
		values += df.getFloatValue(key)
	}
	return
}

func (ds Dataset) GetDatasetAverage(key string) (values float64) {
	for _, df := range ds.datafiles {
		values += df.getFloatValue(key)
	}
	values /= float64(len(ds.datafiles))
	return
}

func (ds Dataset) GetName() string {
	return ds.name
}

func (ds Dataset) GetSize() int {
	return len(ds.datafiles)
}

func (ds Dataset) HasKey(key string) bool {
	hasKey := true
	for _, df := range ds.datafiles {
		if !df.hasKey(key) {
			hasKey = false
			break
		}
	}
	return hasKey
}
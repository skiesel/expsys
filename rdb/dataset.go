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

func (ds Dataset) addPathKeys(baseDirectory string) {
	for i := range ds.datafiles {
		ds.datafiles[i].addPathKeys(baseDirectory)
	}
}

func (ds Dataset) GetDatasetFloatValues(key string) (values []float64) {
	for i := range ds.datafiles {
		values = append(values, ds.datafiles[i].getFloatValue(key))
	}
	return
}

func (ds Dataset) GetDatasetFloatValuesPair(key string, id string) (values []float64, ids []string) {
	for i := range ds.datafiles {
		ids = append(ids, ds.datafiles[i].getStringValue(id))
		values = append(values, ds.datafiles[i].getFloatValue(key))
	}
	return
}

func (ds Dataset) GetDatasetSum(key string) (values float64) {
	for i := range ds.datafiles {
		values += ds.datafiles[i].getFloatValue(key)
	}
	return
}

func (ds Dataset) GetDatasetAverage(key string) (values float64) {
	for i := range ds.datafiles {
		values += ds.datafiles[i].getFloatValue(key)
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
	for i := range ds.datafiles {
		if ! ds.datafiles[i].hasKey(key) {
			hasKey = false
			break
		}
	}
	return hasKey
}
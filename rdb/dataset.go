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

func (ds Dataset) GetDatasetFloatValues(key string, id string) (values []float64, ids []string) {
	for i := range ds.datafiles {
		ids = append(ids, ds.datafiles[i].getStringValue(id))
		values = append(values, ds.datafiles[i].getFloatValue(key))
	}
	return
}

func (ds Dataset) GetName() string {
	return ds.name
}

func (ds Dataset) GetSize() int {
	return len(ds.datafiles)
}
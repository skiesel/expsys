package rdb

// You guessed it, a dataset
type Dataset struct {
	name      string
	datafiles []*Datafile
}

// Create a new dataset with this name and using these file paths
func newDataset(name string, files []string) *Dataset {
	ds := new(Dataset)
	ds.name = name
	ds.datafiles = []*Datafile{}

	for i := range files {
		df := newDatafileFromRDB(files[i])
		if len(df.values) > 0 {
			ds.datafiles = append(ds.datafiles, df)
		}
	}

	return ds
}

func (ds Dataset) copyDataset() *Dataset {
	newDs := new(Dataset)
	newDs.name = ds.name
	newDs.datafiles = make([]*Datafile, len(ds.datafiles))

	for i, df := range ds.datafiles {
		newDs.datafiles[i] = df.copyDatafile()
	}

	return newDs
}

func (ds Dataset) RenameDataset(name string) *Dataset {
	newDss := ds.copyDataset()
	newDss.name = name
	return newDss
}

func (ds Dataset) AddTransformedKey(key string, transform func(string) string, newKey string) {
	for _, df := range ds.datafiles {
		df.addKey(newKey, transform(df.getStringValue(key)))
	}
}

// Filter this dataset, returning a new one, that includes datafiles whose values
// bound to "key" cause the "include" function to return true
func (ds Dataset) FilterDataset(include func(string) bool, key string) (filtered *Dataset) {
	filtered = &Dataset{}
	filtered.name = ds.name
	for _, df := range ds.datafiles {
		if include(df.getStringValue(key)) {
			filtered.datafiles = append(filtered.datafiles, df)
		}
	}
	return
}

// Do all the datafiles in this dataset with values bound to "key"
// cause "test" to return true
func (ds Dataset) TestDataset(test func(string) bool, key string) bool {
	for _, df := range ds.datafiles {
		if !test(df.getStringValue(key)) {
			return false
		}
	}
	return true
}

// Using the the retained path value in the datafile,
// add all the path keys starting from baseDirectory
// df.path == baseDirectory/keysToBeUsed
func (ds Dataset) addPathKeys(baseDirectory string) {
	for _, df := range ds.datafiles {
		df.addRDBPathKeys(baseDirectory)
	}
}

// Accumulate a slice of float values bound to "key" across all
// datafiles in this dataset
func (ds Dataset) GetDatasetFloatValues(key string) (values []float64) {
	for _, df := range ds.datafiles {
		values = append(values, df.getFloatValue(key))
	}
	return
}

// Accumulate a slice of string values bound to "key" across all
// datafiles in this dataset
func (ds Dataset) GetDatasetStringValues(key string) (values []string) {
	for _, df := range ds.datafiles {
		values = append(values, df.getStringValue(key))
	}
	return
}

// Accumulate a slice of int values bound to "key" across all
// datafiles in this dataset
func (ds Dataset) GetDatasetIntegerValues(key string) (values []int64) {
	for _, df := range ds.datafiles {
		values = append(values, df.getIntegerValue(key))
	}
	return
}

// Accumulate a slice of float values bound to "key" across all
// datafiles in this dataset, but also include the associated string values
// bound to "id" -- this is useful when trying to match up data based on an identifier like instance
func (ds Dataset) GetDatasetFloatValuesPair(key string, id string) (values []float64, ids []string) {
	for _, df := range ds.datafiles {
		ids = append(ids, df.getStringValue(id))
		values = append(values, df.getFloatValue(key))
	}
	return
}

// Return the dataset's name
func (ds Dataset) GetName() string {
	return ds.name
}

// Return the number of datafiles in the dataset
func (ds Dataset) GetSize() int {
	return len(ds.datafiles)
}

// Checks if all datafiles in the dataset have "key" bound
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

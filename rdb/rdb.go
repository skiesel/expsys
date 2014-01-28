package rdb

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Retrieve a single dataset defined by the rdb path rooted at directory
// and filtered by keys
func GetDataset(directory string, keys map[string]string, name string, addPathKeys bool) *Dataset {
	ds := newDataset(name, getMatchingFiles(directory, keys))
	if addPathKeys {
		ds.addPathKeys(directory)
	}
	return ds
}

func GetDatasetWithPathKeys(directory string, keys map[string]string, name string) *Dataset {
	return GetDataset(directory, keys, name, true)
}

// Retrieve a set of datasets defined by the rdb path rooted at directory
// and filtered by each key set in keys individually
func GetDatasets(directory string, dss map[string]map[string]string, addPathKeys bool) []*Dataset {
	
	datasets := make([]*Dataset, len(dss))
	i := 0
	for dsName, dsKeys := range dss {
		datasets[i] = GetDataset(directory, dsKeys, dsName, addPathKeys)
		i++
	}
	return datasets
}

func GetDatasetsWithPathKeys(directory string, dss map[string]map[string]string) []*Dataset {
	return GetDatasets(directory, dss, true)
}

// A recursive function that crawls the directory structure starting at
// directory and ignores filters anything that doesn't match the keys
// passed in. The exception being, no key being specified means get everything
func getMatchingFiles(directory string, filter map[string]string) []string {
	return crawlAndCollect(directory, filter)
}

func getKeyInDirectory(directory string) string {
	fInfo, error := ioutil.ReadDir(directory)

	if error != nil { // there was an error
		panic(error)
	} else {
		for i := range fInfo {
			if strings.Contains(fInfo[i].Name(), "KEY=") {
				return strings.SplitAfter(fInfo[i].Name(), "KEY=")[1]
			}
		}
	}
	fmt.Println("No key file found in ", directory)
	panic("No key file found")
}

func crawlAndCollect(directory string, filter map[string]string) []string {
	returnPaths := []string{}

	fInfo, error := ioutil.ReadDir(directory)

	if error != nil { // there was an error
		panic(error)
	} else {
		//look for the key file in this directory
		filterKey := getKeyInDirectory(directory)

		var include func(string) bool = nil
		if filter[filterKey] == "" {
			include = func(value string) bool {
				return true
			}
		} else { //only include directories or files that match the value for the key
			include = func(value string) bool {
				return filter[filterKey] == value
			}
		}

		directories := []string{}

		for _, file := range fInfo {
			filename := file.Name()
			isKeyFile := strings.Contains(filename, "KEY=")
			isDotFile := filename[0] == '.'

			if !include(filename) || isKeyFile || isDotFile {
				continue
			}

			isDir := file.IsDir()
			relativePath := strings.Join([]string{directory, filename}, "/")

			if isDir {
				directories = append(directories, relativePath)
			} else if !isKeyFile {
				returnPaths = append(returnPaths, relativePath)
			}
		}

		for _, directory := range directories {
			returnPaths = append(returnPaths, crawlAndCollect(directory, filter)...)
		}
	}
	return returnPaths
}

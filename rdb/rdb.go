package rdb

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Retrieve a single dataset defined by the rdb path rooted at directory
// and filtered by keys
func GetDataset(directory string, keys map[string]string, name string) *Dataset {
	return newDataset(name, getMatchingFiles(directory, keys))
}

// Retrieve a set of datasets defined by the rdb path rooted at directory
// and filtered by each key set in keys individually
func GetDatasets(directory string, keys []map[string]string, names []string) []*Dataset {
	if len(keys) != len(names) {
		fmt.Printf("%d names provided for %d datasets!", len(names), len(keys))
		panic("Slice length mismatch")
	}

	datasets := make([]*Dataset, len(keys))
	for i := range keys {
		datasets[i] = GetDataset(directory, keys[i], names[i])
	}
	return datasets
}

// A recursive function that crawls the directory structure starting at
// directory and ignores filters anything that doesn't match the keys
// passed in. The exception being, no key being specified means get everything
func getMatchingFiles(directory string, filter map[string]string) []string {
	return crawlAndCollect(directory, filter)
}

func crawlAndCollect(directory string, filter map[string]string) []string {
	returnPaths := make([]string, 0)

	fInfo, error := ioutil.ReadDir(directory)

	if error != nil { // there was an error
		fmt.Printf(error.Error())
	} else {
		//look for the key file in this directory
		filterKey := ""
		for i := range fInfo {
			if strings.Contains(fInfo[i].Name(), "KEY=") {
				filterKey = strings.SplitAfter(fInfo[i].Name(), "KEY=")[1]
				break
			}
		}

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

		directories := make([]string, 0)

		for i := range fInfo {
			filename := fInfo[i].Name()
			isKeyFile := strings.Contains(filename, "KEY=")
			isDotFile := filename[0] == '.'

			if !include(filename) || isKeyFile || isDotFile {
				continue
			}

			isDir := fInfo[i].IsDir()
			relativePath := strings.Join([]string{directory, filename}, "/")

			if isDir {
				directories = append(directories, relativePath)
			} else if !isKeyFile {
				returnPaths = append(returnPaths, relativePath)
			}
		}

		for i := range directories {
			returnPaths = append(returnPaths, crawlAndCollect(directories[i], filter)...)
		}
	}
	return returnPaths
}

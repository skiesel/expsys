package rdb

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Datafile -- path is the path used to construct the datafile
type Datafile struct {
	path   string
	values map[string]string
}

// Construct a new datafile from the provided path/filename
func newDatafileFromRDB(filename string) *Datafile {
	df := new(Datafile)
	df.path = filename
	df.values = map[string]string{}

	completeDF := false

	file, err := os.Open(filename)
	if err == nil {

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "#end data file format") {
				completeDF = true
			}
			df.addRDBDataRow(line)
		}
	}

	if !completeDF {
		df.values = map[string]string{}
	}

	return df
}

func (df Datafile) copyDatafile() *Datafile {
	newDf := new(Datafile)
	newDf.path = df.path
	newDf.values = df.values
	return newDf
}

// Retroactively add keys specified in the path used to construct this datafile
// but ignore any overlap between df.path and baseDirectory
func (df Datafile) addRDBPathKeys(baseDirectory string) {
	pathPiece := strings.Replace(df.path, baseDirectory, "", -1)
	pathPiece = strings.Trim(pathPiece, "/")
	keyValues := strings.Split(pathPiece, "/")

	currentPath := baseDirectory
	for _, keyValue := range keyValues {
		key := getKeyInDirectory(currentPath)
		df.checkAndAddKeyValue(key, keyValue)
		currentPath = strings.Join([]string{currentPath, keyValue}, "/")
	}
}

// Parse the key value pair from an RDB file
func parseRDBKeyValuePair(str string) (key string, value string, ok bool) {
	tokens := strings.Split(str, "\"")
	if len(tokens) >= 4 {
		key = tokens[1]
		value = tokens[3]
		ok = true
	}
	return
}

// Add this RDB style row to the datafile
func (df Datafile) addRDBDataRow(line string) {
	if !strings.HasPrefix(line, "#pair") {
		return
	}
	key, value, ok := parseRDBKeyValuePair(line)
	if !ok {
		return
	}
	df.checkAndAddKeyValue(key, value)
}

// Dump out datafile for debugging
func (df Datafile) dump() {

	fmt.Println(df.path)
	for key, value := range df.values {
		fmt.Printf("\"%s\"\t\"%s\"\n", key, value)
	}
}

// A safety check when trying to bind new values to the values map
// It's okay if the value you're adding is already bound if the new value
// matches the old value
// Otherwise there is a problem and we don't know which value to maintain
func (df Datafile) checkAndAddKeyValue(key string, value string) {
	boundValue, keyExists := df.values[key]
	if keyExists {
		if boundValue != value {
			fmt.Printf("Trying to add mismatched values for key (\"%s\") (\"%s\", \"%s\") in %s\n",
				key, boundValue, value, df.path)
			panic("Datafile: Mismatched Key Values")
		}
	} else {
		df.values[key] = value
	}
}

func (df Datafile) addKey(key string, value string) {
	df.checkAndAddKeyValue(key, value)
}

// Does this datafile have this key bound
func (df Datafile) hasKey(key string) bool {
	_, keyExists := df.values[key]
	return keyExists
}

// Return the string value bound to "key"
func (df Datafile) getStringValue(key string) string {
	val, exists := df.values[key]

	if !exists {
		fmt.Printf("Key \"%s\" is unbound in %s\n", key, df.path)
		panic("Unbound Key")
	}

	return val
}

// Return the converted int value bound to "key"
func (df Datafile) getIntegerValue(key string) int64 {
	strVal := df.getStringValue(key)
	val, err := strconv.ParseInt(strVal, 10, 0)

	if err != nil {
		fmt.Printf("Key \"%s\" with value \"%s\" could not be converted to int64 in %s\n",
			key, strVal, df.path)
		panic(err)
	}

	return val
}

// Return the converted float value bound to "key"
func (df Datafile) getFloatValue(key string) float64 {
	strVal := df.getStringValue(key)
	val, err := strconv.ParseFloat(strVal, 0)

	if err != nil {
		fmt.Printf("Key \"%s\" with value \"%s\" could not be converted to float64 in %s\n",
			key, strVal, df.path)
		panic(err)
	}

	return val
}

// Return the converted bool value bound to "key"
func (df Datafile) getBooleanValue(key string) bool {
	strVal := df.getStringValue(key)
	val, err := strconv.ParseBool(strVal)

	if err != nil {
		fmt.Printf("Key \"%s\" with value \"%s\" could not be converted to bool in %s\n",
			key, strVal, df.path)
		panic(err)
	}

	return val
}

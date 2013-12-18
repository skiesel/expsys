package rdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

type Datafile struct {
	path string
	values map[string]string
}

func newDatafile(filename string) *Datafile {
	df := new(Datafile)
	df.path = filename
	df.values = make(map[string]string, 0)

	file, err := os.Open(filename)
	if err == nil {

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			df.addDataRow(line)
		}
	}

	return df
}

func (df Datafile) addPathKeys(baseDirectory string) {
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

func (df Datafile) hasKey(key string) bool {
	_, keyExists := df.values[key]
	return keyExists
}

func parseKeyValuePair(str string) (key string, value string, ok bool) {
	tokens := strings.Split(str, "\"")
	if len(tokens) >= 4 {
		key = tokens[1]
		value = tokens[3]
		ok = true
	}
	return
}

func (df Datafile) addDataRow(line string) {
	if !strings.HasPrefix(line, "#pair") {
		return
	}

	key, value, ok := parseKeyValuePair(line)

	if !ok {
		return
	}

	df.checkAndAddKeyValue(key, value)
}

func (df Datafile) dump() {
	for key, value := range df.values {
		fmt.Printf("\"%s\"\t\"%s\"\n", key, value)
	}
}

func (df Datafile) getStringValue(key string) string {
	val, exists := df.values[key]

	if !exists {
		fmt.Printf("Key \"%s\" is unbound in %s\n", key, df.path)
		panic("Unbound Key")
	}

	return val
}

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
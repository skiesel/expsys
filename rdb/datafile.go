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

	df.values[key] = value
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
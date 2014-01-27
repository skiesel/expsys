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
	path    string
	pairs   map[string]string
	columns map[string][][]string
}

// Construct a new datafile from the provided path/filename
func newDatafileFromRDB(filename string) *Datafile {
	df := new(Datafile)
	df.path = filename
	df.pairs = map[string]string{}
	df.columns = map[string][][]string{}

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
		df.pairs = map[string]string{}
	}

	return df
}

func (df Datafile) copyDatafile() *Datafile {
	newDf := new(Datafile)
	newDf.path = df.path
	newDf.pairs = df.pairs
	newDf.columns = df.columns
	return newDf
}

// Retroactively add keys specified in the path used to construct this datafile
// but ignore any overlap between df.path and baseDirectory
func (df Datafile) addRDBPathKeys(baseDirectory string) {
	pathPiece := strings.Replace(df.path, baseDirectory, "", -1)
	pathPiece = strings.Trim(pathPiece, "/")
	keyPairs := strings.Split(pathPiece, "/")

	currentPath := baseDirectory
	for _, keyValue := range keyPairs {
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

func parseRDBColumns(str string) (columnsName string, columns []string, ok bool) {
	tokens := strings.Split(str, "\"")
	for _, token := range tokens {
		token = strings.Trim(token, " \t")
		if token != "" && token != "#altcols" && token != "#altrow" {
			if columnsName == "" {
				columnsName = token
			} else {
				if columns == nil {
					columns = []string{token}
				} else {
					columns = append(columns, token)
				}
			}
		}
	}
	ok = true
	return
}

// Add this RDB style row to the datafile
func (df Datafile) addRDBDataRow(line string) {
	switch {
	case strings.HasPrefix(line, "#pair"):
		key, value, ok := parseRDBKeyValuePair(line)
		if !ok {
			return
		}
		df.checkAndAddKeyValue(key, value)

	case strings.HasPrefix(line, "#altcols"):
		columnsName, columnHeaders, ok := parseRDBColumns(line)
		if !ok {
			return
		}
		df.checkAndAddColumnHeaders(columnsName, columnHeaders)

	case strings.HasPrefix(line, "#altrow"):
		columnsName, columnValues, ok := parseRDBColumns(line)
		if !ok {
			return
		}
		df.checkAndAddColumnValues(columnsName, columnValues)
	}
}

// Dump out datafile for debugging
func (df Datafile) dump() {

	fmt.Println(df.path)
	for key, value := range df.pairs {
		fmt.Printf("\"%s\"\t\"%s\"\n", key, value)
	}
	for key, columnsTable := range df.columns {
		for _, row := range columnsTable {
			fmt.Printf("\"%s\"\t", key)
			for _, value := range row {
				fmt.Printf("\"%s\"\t", value)
			}
			fmt.Printf("\n")
		}
	}
}

// A safety check when trying to bind new pairs to the pairs map
// It's okay if the value you're adding is already bound if the new value
// matches the old value
// Otherwise there is a problem and we don't know which value to maintain
func (df Datafile) checkAndAddKeyValue(key string, value string) {
	boundValue, keyExists := df.pairs[key]
	if keyExists {
		if boundValue != value {
			fmt.Printf("Trying to add mismatched pairs for key (\"%s\") (\"%s\", \"%s\") in %s\n",
				key, boundValue, value, df.path)
			panic("Datafile: Mismatched Key pairs")
		}
	} else {
		df.pairs[key] = value
	}
}

func (df Datafile) checkAndAddColumnHeaders(columnsName string, columnHeaders []string) {
	columnsTable, keyExists := df.columns[columnsName]
	if keyExists {
		for i, header := range columnHeaders {
			if columnsTable[0][i] != header {
				fmt.Printf("Trying to add mismatched column for columns table (\"%s\") (\"%s\", \"%s\") in %s\n",
					columnsName, columnsTable[i], header, df.path)
				panic("Datafile: Mismatched Column headers")
			}
		}
	} else {
		df.columns[columnsName] = [][]string{}
		df.columns[columnsName] = append(df.columns[columnsName], []string{})
		for _, header := range columnHeaders {
			df.columns[columnsName][0] = append(df.columns[columnsName][0], header)
		}
	}
}

func (df Datafile) checkAndAddColumnValues(columnsName string, columnValues []string) {
	columnsTable, keyExists := df.columns[columnsName]
	rows := len(columnsTable)
	if keyExists && rows >= 1 {
		df.columns[columnsName] = append(df.columns[columnsName], []string{})
		for _, value := range columnValues {
			df.columns[columnsName][rows] = append(df.columns[columnsName][rows], value)
		}
	} else {
		fmt.Printf("Trying to add values to columns table without headers (\"%s\") in %s\n",
			columnsName, df.path)
		panic("Datafile: Missing Column Headers")
	}
}

func (df Datafile) addKey(key string, value string) {
	df.checkAndAddKeyValue(key, value)
}

// Does this datafile have this key bound
func (df Datafile) hasKey(key string) bool {
	_, keyExists := df.pairs[key]
	return keyExists
}

// Return the string value bound to "key"
func (df Datafile) getStringValue(key string) string {
	val, exists := df.pairs[key]

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

func (df Datafile) getColumnValues(tableName, key string) []string {
	_, keyExists := df.columns[tableName]
	if !keyExists {
		fmt.Printf("Could not find columns table \"%s\" in (%s)\n", tableName, df.path)
		df.dump()
		panic("Columns table not found")
	}

	columnNum := -1
	for i, header := range df.columns[tableName][0] {
		if  header == key {
			columnNum = i
			break
		}
	}
	
	if columnNum < 0 {
		fmt.Printf("Could not find column (\"%s\") in table (\"%s\") in (%s)\n", key, tableName, df.path)
		panic("Column not found")
	}

	columnVals := make([]string, len(df.columns[tableName]) - 1)
	for i := range columnVals {
		columnVals[i] = df.columns[tableName][i+1][columnNum]
	}

	return columnVals
}

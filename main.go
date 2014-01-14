package main

import (
	"flag"
	"github.com/skiesel/expsys/projects"
	"strings"
)

var (
	projectNames = flag.String("p", "", "The projects you would like to run.")
)

func main() {
	flag.Parse()

	names := strings.Split(*projectNames, " ")
	filterNames := []string{}

	for _, val := range names {
		val = strings.Trim(val, " ")
		if val != "" {
			filterNames = append(filterNames, val)
		}
	}

	projects.BuildProjectPlots(filterNames)
}

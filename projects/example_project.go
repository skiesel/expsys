package projects

import (
	"fmt"
)

var (
	_ = addProject(Project{
		name:  "example_project",
		build: buildExampleProjet})
)

func buildExampleProjet() {
	//Do whatever you want in here
	fmt.Println("Hello, World")
}
package projects

import (
	"fmt"
)

func init() {
	addProject(Project{
		name:  "example_project",
		build: buildExampleProject})
}

func buildExampleProject() {
	//Do whatever you want in here
	fmt.Println("Hello, World")
}

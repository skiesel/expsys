package projects

import (
	"fmt"
)

var (
	projects = map[string]Project{}
)

type BuildFunction func()

type Project struct {
	name  string
	build BuildFunction
}

func BuildProjectPlots(projectNames []string) {
	for _, projectName := range projectNames {
		proj, bound := projects[projectName]
		if !bound {
			fmt.Printf("Could not find project: '%s'\n", projectName)
		} else {
			proj.build()
		}
	}
}

func addProject(project Project) {
	_, bound := projects[project.name]
	if !bound {
		projects[project.name] = project
	} else {
		fmt.Printf("Project already exists: '%s'\n", project.name)
	}
}

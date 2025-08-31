package typeDefines

import (
	"fmt"
)

type Project struct {
	Name     string `json:"Name"`
	Url      string `json:"Url"`
	Auth     IAuth
	Parallel bool    `json:"Parallel"`
	Suites   []Suite `json:"Suites"`
}

func (project *Project) ExecuteProject() {
	fmt.Printf("\nExecuting Project: %s\n", project.Name)

	for i := range project.Suites {
		if project.Parallel {
			go project.Suites[i].ExecuteSuite(project.Url, project.Auth)
		} else {
			project.Suites[i].ExecuteSuite(project.Url, project.Auth)
		}
	}
}

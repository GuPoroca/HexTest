package typeDefines

import (
	"fmt"
)

type Project struct {
	Name   string  `json:"Name"`
	Url    string  `json:"Url"`
	Auth   Auth    `json:"Auth"`
	Suites []Suite `json:"Suites"`
}

func (project *Project) ExecuteProject() {
	fmt.Printf("Executing Project: %s\n", project.Name)
	for i := range project.Suites {
		project.Suites[i].ExecuteSuite(project.Url, project.Auth)
	}
}

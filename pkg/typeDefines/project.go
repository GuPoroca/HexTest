package typeDefines

import (
	"fmt"
)

type Project struct {
	Name                    string            `json:"Name"`
	Url                     string            `json:"Url"`
	Parallel                bool              `json:"Parallel"`
	Project_Headers         map[string]string `json:"Project_Headers"`
	Suites                  []Suite           `json:"Suites"`
	Auth                    IAuth
	Passed_Comparissons_num int
	Total_Comparissons_num  int
}

func (project *Project) ExecuteProject() {
	fmt.Printf("\nExecuting Project: %s\n", project.Name)
	fmt.Printf("\nCreating Authentication\n")

	project.Auth = NewoAuth2("client_credentials")

	for i := range project.Suites {
		if project.Parallel {
			go project.Suites[i].ExecuteSuite(project.Url, project.Auth)
		} else {
			project.Suites[i].ExecuteSuite(project.Url, project.Auth)
		}
		project.Passed_Comparissons_num += project.Suites[i].Passed_Comparissons_num
		project.Total_Comparissons_num += project.Suites[i].Total_Comparissons_num

	}
	fmt.Print("\n---------------------------------------\n")
	fmt.Printf("Total comparissons passed in the entire Project: %d/%d\n\n", project.Passed_Comparissons_num, project.Total_Comparissons_num)
	fmt.Print("\n---------------------------------------\n")

}

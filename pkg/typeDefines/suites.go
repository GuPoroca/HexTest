package typeDefines

import (
	"fmt"
)

type Suite struct {
	Name                    string `json:"Name"`
	Comment                 string `json:"Comment"`
	Tests                   []Test `json:"Tests"`
	Parallel                bool   `json:"Parallel"`
	Passed_Comparissons_num int
	Total_Comparissons_num  int
}

func (suite *Suite) ExecuteSuite(url string, auth IAuth) {
	fmt.Printf("Executing Suite: %s\n", suite.Name)

	fmt.Print("\n---------------------------------------\n")
	for i := range suite.Tests {
		if suite.Parallel {
			go suite.Tests[i].Execute(url, auth)
		} else {
			suite.Tests[i].Execute(url, auth)
		}
		suite.Passed_Comparissons_num += suite.Tests[i].Passed_Comparissons_num
		suite.Total_Comparissons_num += suite.Tests[i].Total_Comparissons_num

	}
	fmt.Print("\n---------------------------------------\n")
	fmt.Printf("Total comparissons passed in this Suite: %d/%d\n\n", suite.Passed_Comparissons_num, suite.Total_Comparissons_num)
	fmt.Print("\n---------------------------------------\n")
}

package typeDefines

import (
	"fmt"
)

type Suite struct {
	Name     string `json:"Name"`
	Comment  string `json:"Comment"`
	Tests    []Test `json:"Tests"`
	Parallel bool   `json:"Parallel"`
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
	}
}

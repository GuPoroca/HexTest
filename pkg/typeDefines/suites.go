package typeDefines

import (
	"fmt"
)

type Suite struct {
	Name  string `json:"Name"`
	Tests []Test `json:"Tests"`
}

func (suite *Suite) ExecuteSuite(url string, auth Auth) {
	fmt.Printf("Executing Suite: %s\n", suite.Name)
	for i := range suite.Tests {
		suite.Tests[i].Execute(url, auth)
	}
}

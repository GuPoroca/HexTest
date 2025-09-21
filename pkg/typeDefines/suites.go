package typeDefines

import ()

type Suite struct {
	Name     string `json:"Name"`
	Comment  string `json:"Comment"`
	Tests    []Test `json:"Tests"`
	Parallel bool   `json:"Parallel"`
}

func (suite *Suite) ExecuteSuite(url string) {
	for i := range suite.Tests {
		if suite.Parallel {
			go suite.Tests[i].Execute(url)
		} else {
			suite.Tests[i].Execute(url)
		}
	}
}

package typeDefines

import (
	"sync"
)

type Suite struct {
	Name     string `json:"Name"`
	Comment  string `json:"Comment"`
	Tests    []Test `json:"Tests"`
	Parallel bool   `json:"Parallel"`
}

func (suite *Suite) ExecuteSuite(url string, headermap map[string]string) {

	var wg sync.WaitGroup

	for i := range suite.Tests {
		if suite.Parallel == true {
			wg.Add(1)
			go func(idx int) {

				defer wg.Done()
				suite.Tests[idx].Execute(url, headermap)
			}(i)

		} else {
			suite.Tests[i].Execute(url, headermap)
		}
	}
	if suite.Parallel {
		wg.Wait()
	}
}

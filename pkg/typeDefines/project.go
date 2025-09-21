package typeDefines

import (
	"fmt"
	"sync"
)

type Project struct {
	Name            string            `json:"Name"`
	Url             string            `json:"Url"`
	Parallel        bool              `json:"Parallel"`
	Project_Headers map[string]string `json:"Project_Headers"`
	Suites          []Suite           `json:"Suites"`
	Auth            IAuth
}
type CheckCallback func(suiteName, testName, assertField string, checkIndex int, passed int)

func doWork(id int, wg *sync.WaitGroup, project Project) {
	var cb CheckCallback
	defer wg.Done() // Decrement counter when goroutine finishes
	project.Suites[id].ExecuteSuite(project.Url, cb)
}

func (project *Project) ExecuteProject() {

	var wg sync.WaitGroup

	for i := range project.Suites {
		if project.Parallel {
			wg.Add(1)
			go doWork(i, &wg, *project)
		} else {
			project.Suites[i].ExecuteSuite(project.Url)
		}
	}

	wg.Wait()
}

func (p Project) PrintResults() {
	checks, ch_ps, ch_fl, ch_bk := 0, 0, 0, 0
	for _, s := range p.Suites {
		for _, t := range s.Tests {
			for _, a := range t.Asserts {
				for _, c := range a.Checks {
					for i, p := range c.Passed {
						checks++
						switch p {
						case 1:
							ch_ps++
						case 0:
							ch_fl++
							fmt.Printf("Comparisson: %s\n", fmt.Sprintf("%s %s %s", StringifyMyAny(a.FieldResponseValue), c.Operand, StringifyMyAny(c.Expected[i])))
							fmt.Printf("On Assert: %s\n", a.Field)
							fmt.Printf("On Test: %s\n", t.Name)
							fmt.Printf("On Suite: %s\n", s.Name)
							fmt.Printf("FAILED\n\n")

						case -1:
							ch_bk++
							fmt.Printf("Comparisson: %s\n", fmt.Sprintf("%s %s %s", StringifyMyAny(a.FieldResponseValue), c.Operand, StringifyMyAny(c.Expected[i])))
							fmt.Printf("On Assert: %s\n", a.Field)
							fmt.Printf("On Test: %s\n", t.Name)
							fmt.Printf("On Suite: %s\n", s.Name)
							fmt.Printf("BROKEN\n\n")
						}
					}
				}
			}

		}
	}
	fmt.Printf("Number of Checks Made: %v\n", checks)
	fmt.Printf("Number of Checks Passed: %v\n", ch_ps)
	fmt.Printf("Number of Checks Failed: %v\n", ch_fl)
	fmt.Printf("Number of Checks Broken: %v\n", ch_bk)

}


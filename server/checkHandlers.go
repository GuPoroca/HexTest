package server

import (
	"github.com/GuPoroca/HexTest/front/components"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"net/http"
)

/////////////////////////////////////////////////////
// CHECK
/////////////////////////////////////////////////////

func HandleEditCheck(w http.ResponseWriter, r *http.Request) {
	operand := r.URL.Query().Get("operand")
	for _, suite := range currentProject.Suites {
		for _, test := range suite.Tests {
			for _, asrt := range test.Asserts {
				for _, check := range asrt.Checks {
					if check.Operand == operand {
						components.EditCheckForm(check).Render(r.Context(), w)
						return
					}
				}
			}
		}
	}
	http.Error(w, "Check not found", http.StatusNotFound)
}

func HandleSaveCheck(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	old := r.FormValue("oldOperand")
	operand := r.FormValue("Operand")

	for si := range currentProject.Suites {
		for ti := range currentProject.Suites[si].Tests {
			for ai := range currentProject.Suites[si].Tests[ti].Asserts {
				for ci := range currentProject.Suites[si].Tests[ti].Asserts[ai].Checks {
					c := &currentProject.Suites[si].Tests[ti].Asserts[ai].Checks[ci]
					if c.Operand == old || c.Operand == operand {
						c.Operand = operand
						components.EditCheckForm(*c).Render(r.Context(), w)
						components.ProjectSidebarOOB(r.Context(), w, currentProject)

						return
					}
				}
			}
		}
	}
	http.Error(w, "Check not found after save", http.StatusNotFound)
}

func HandleAddCheck(w http.ResponseWriter, r *http.Request) {
	assertField := r.URL.Query().Get("assert")
	newCheck := typeDefines.Check{Operand: "New Check"}

	for i := range currentProject.Suites {
		for j := range currentProject.Suites[i].Tests {
			for k := range currentProject.Suites[i].Tests[j].Asserts {
				if currentProject.Suites[i].Tests[j].Asserts[k].Field == assertField {
					currentProject.Suites[i].Tests[j].Asserts[k].Checks = append(currentProject.Suites[i].Tests[j].Asserts[k].Checks, newCheck)

					components.ProjectSidebarOOB(r.Context(), w, currentProject)
					components.EditCheckForm(newCheck).Render(r.Context(), w)
					return
				}
			}
		}
	}
	http.Error(w, "Assert not found", http.StatusNotFound)
}

func HandleDeleteCheck(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	operand := r.FormValue("operand") // in case you use hx-vals in future
	if operand == "" {
		operand = r.URL.Query().Get("operand")
	}
	for si := range currentProject.Suites {
		for ti := range currentProject.Suites[si].Tests {
			for ai := range currentProject.Suites[si].Tests[ti].Asserts {
				for ci, check := range currentProject.Suites[si].Tests[ti].Asserts[ai].Checks {
					if check.Operand == operand {
						currentProject.Suites[si].Tests[ti].Asserts[ai].Checks = append(
							currentProject.Suites[si].Tests[ti].Asserts[ai].Checks[:ci],
							currentProject.Suites[si].Tests[ti].Asserts[ai].Checks[ci+1:]...,
						)

						components.ProjectSidebarOOB(r.Context(), w, currentProject)
						w.Write([]byte(`<p class="text-gray-400">Select an item.</p></div>`))
						return
					}
				}
			}
		}
	}
	http.Error(w, "Check not found", http.StatusNotFound)
}

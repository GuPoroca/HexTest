package server

import (
	"log"
	"net/http"

	"github.com/GuPoroca/HexTest/front"
	"github.com/GuPoroca/HexTest/front/components"
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
)

var currentProject typeDefines.Project

/////////////////////////////////////////////////////
// MAIN ENTRY
/////////////////////////////////////////////////////

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	projectData := jsonOperations.ReadJSON("./one_request.json")
	currentProject = projectData

	component := front.Layout(projectData, nil)
	component.Render(r.Context(), w)
}

/////////////////////////////////////////////////////
// PROJECT
/////////////////////////////////////////////////////

func HandleEditProject(w http.ResponseWriter, r *http.Request) {
	components.EditProjectForm(currentProject).Render(r.Context(), w)
}

func HandleSaveProject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	currentProject.Name = r.FormValue("Name")
	currentProject.Url = r.FormValue("Url")

	components.EditProjectForm(currentProject).Render(r.Context(), w)
	components.ProjectSidebarOOB(r.Context(), w, currentProject)
}

func HandleAddSuite(w http.ResponseWriter, r *http.Request) {
	newSuite := typeDefines.Suite{Name: "New Suite"}
	currentProject.Suites = append(currentProject.Suites, newSuite)

	components.ProjectSidebarOOB(r.Context(), w, currentProject)
	components.EditSuiteForm(newSuite).Render(r.Context(), w)
}

/////////////////////////////////////////////////////
// SUITE
/////////////////////////////////////////////////////

func HandleEditSuite(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	for _, suite := range currentProject.Suites {
		if suite.Name == name {
			components.EditSuiteForm(suite).Render(r.Context(), w)
			return
		}
	}
	http.Error(w, "Suite not found", http.StatusNotFound)
}

func HandleSaveSuite(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	old := r.FormValue("oldName")
	new := r.FormValue("Name")

	var saved *typeDefines.Suite
	for i := range currentProject.Suites {
		if currentProject.Suites[i].Name == old || currentProject.Suites[i].Name == new {
			currentProject.Suites[i].Name = new
			saved = &currentProject.Suites[i]
			break
		}
	}
	if saved != nil {
		components.EditSuiteForm(*saved).Render(r.Context(), w)
		components.ProjectSidebarOOB(r.Context(), w, currentProject)

		return
	}
	http.Error(w, "Suite not found after save", http.StatusNotFound)
}

func HandleDeleteSuite(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	for i, suite := range currentProject.Suites {
		if suite.Name == name {
			currentProject.Suites = append(currentProject.Suites[:i], currentProject.Suites[i+1:]...)
			break
		}
	}

	components.ProjectSidebarOOB(r.Context(), w, currentProject)
	w.Write([]byte(`<p class="text-gray-400">Select an item.</p></div>`))
}

/////////////////////////////////////////////////////
// TEST
/////////////////////////////////////////////////////

func HandleEditTest(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	for _, suite := range currentProject.Suites {
		for _, test := range suite.Tests {
			if test.Name == name {
				components.EditTestForm(test).Render(r.Context(), w)
				return
			}
		}
	}
	http.Error(w, "Test not found", http.StatusNotFound)
}

func HandleSaveTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	old := r.FormValue("oldName")
	name := r.FormValue("Name")

	for si := range currentProject.Suites {
		for ti := range currentProject.Suites[si].Tests {
			t := &currentProject.Suites[si].Tests[ti]
			if t.Name == old || t.Name == name {
				t.Name = name
				t.Method = r.FormValue("Method")
				t.Request_body = r.FormValue("Request_body")
				t.Api_endpoint = r.FormValue("Api_endpoint")
				t.Comment = r.FormValue("Comment")
				components.EditTestForm(*t).Render(r.Context(), w)
				components.ProjectSidebarOOB(r.Context(), w, currentProject)
				return
			}
		}
	}

	http.Error(w, "Test not found after save", http.StatusNotFound)

}

func HandleAddTest(w http.ResponseWriter, r *http.Request) {
	suiteName := r.URL.Query().Get("suite")
	newTest := typeDefines.Test{Name: "New Test"}

	for i := range currentProject.Suites {
		if currentProject.Suites[i].Name == suiteName {
			currentProject.Suites[i].Tests = append(currentProject.Suites[i].Tests, newTest)

			components.ProjectSidebarOOB(r.Context(), w, currentProject)
			components.EditTestForm(newTest).Render(r.Context(), w)
			return
		}
	}
	http.Error(w, "Suite not found", http.StatusNotFound)
}

func HandleDeleteTest(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	for si := range currentProject.Suites {
		for ti, test := range currentProject.Suites[si].Tests {
			if test.Name == name {
				currentProject.Suites[si].Tests = append(currentProject.Suites[si].Tests[:ti], currentProject.Suites[si].Tests[ti+1:]...)

				components.ProjectSidebarOOB(r.Context(), w, currentProject)
				w.Write([]byte(`<p class="text-gray-400">Select an item.</p></div>`))
				return
			}
		}
	}
	http.Error(w, "Test not found", http.StatusNotFound)
}

/////////////////////////////////////////////////////
// ASSERT
/////////////////////////////////////////////////////

func HandleEditAssert(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	for _, suite := range currentProject.Suites {
		for _, test := range suite.Tests {
			for _, asrt := range test.Asserts {
				if asrt.Field == field {
					components.EditAssertForm(asrt).Render(r.Context(), w)
					return
				}
			}
		}
	}
	http.Error(w, "Assert not found", http.StatusNotFound)
}

func HandleSaveAssert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	old := r.FormValue("oldField")
	field := r.FormValue("Field")

	for si := range currentProject.Suites {
		for ti := range currentProject.Suites[si].Tests {
			for ai := range currentProject.Suites[si].Tests[ti].Asserts {
				a := &currentProject.Suites[si].Tests[ti].Asserts[ai]
				if a.Field == old || a.Field == field {
					a.Field = field
					components.EditAssertForm(*a).Render(r.Context(), w)
					components.ProjectSidebarOOB(r.Context(), w, currentProject)

					return
				}
			}
		}
	}
	http.Error(w, "Assert not found after save", http.StatusNotFound)
}

func HandleAddAssert(w http.ResponseWriter, r *http.Request) {
	testName := r.URL.Query().Get("test")
	newAssert := typeDefines.Assert{Field: "New Assert"}

	for i := range currentProject.Suites {
		for j := range currentProject.Suites[i].Tests {
			if currentProject.Suites[i].Tests[j].Name == testName {
				currentProject.Suites[i].Tests[j].Asserts = append(currentProject.Suites[i].Tests[j].Asserts, newAssert)

				components.ProjectSidebarOOB(r.Context(), w, currentProject)
				components.EditAssertForm(newAssert).Render(r.Context(), w)
				return
			}
		}
	}
	http.Error(w, "Test not found", http.StatusNotFound)
}

func HandleDeleteAssert(w http.ResponseWriter, r *http.Request) {
	field := r.URL.Query().Get("field")
	for si := range currentProject.Suites {
		for ti := range currentProject.Suites[si].Tests {
			for ai, asrt := range currentProject.Suites[si].Tests[ti].Asserts {
				if asrt.Field == field {
					currentProject.Suites[si].Tests[ti].Asserts = append(currentProject.Suites[si].Tests[ti].Asserts[:ai], currentProject.Suites[si].Tests[ti].Asserts[ai+1:]...)

					components.ProjectSidebarOOB(r.Context(), w, currentProject)
					w.Write([]byte(`<p class="text-gray-400">Select an item.</p></div>`))
					return
				}
			}
		}
	}
	http.Error(w, "Assert not found", http.StatusNotFound)
}

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
	log.Printf("[DEL CHECK] operand=%q", operand)
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

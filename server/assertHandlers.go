package server

import (
	"github.com/GuPoroca/HexTest/front/components"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"net/http"
)

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

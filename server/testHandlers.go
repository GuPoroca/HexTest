package server

import (
	"github.com/GuPoroca/HexTest/front/components"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"net/http"
)

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

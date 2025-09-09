package server

import (
	"github.com/GuPoroca/HexTest/front/components"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"net/http"
)

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
	comment := r.FormValue("Comment")

	var saved *typeDefines.Suite
	for i := range currentProject.Suites {
		if currentProject.Suites[i].Name == old || currentProject.Suites[i].Name == new {
			currentProject.Suites[i].Name = new
			currentProject.Suites[i].Comment = comment
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

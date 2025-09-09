package server

import (
	"github.com/GuPoroca/HexTest/front/components"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"net/http"
)

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

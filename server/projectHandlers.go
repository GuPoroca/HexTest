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

	//Headers
	keys := r.Form["HeaderKeys"]
	values := r.Form["HeaderValues"]
	headers := make(map[string]string)
	for i := range keys {
		if keys[i] != "" {
			headers[keys[i]] = values[i]
		}
	}
	currentProject.Project_Headers = headers

	components.EditProjectForm(currentProject).Render(r.Context(), w)
	components.ProjectSidebarOOB(r.Context(), w, currentProject)
}

func HandleAddSuite(w http.ResponseWriter, r *http.Request) {
	newSuite := typeDefines.Suite{Name: "New Suite"}
	currentProject.Suites = append(currentProject.Suites, newSuite)

	components.ProjectSidebarOOB(r.Context(), w, currentProject)
	components.EditSuiteForm(newSuite).Render(r.Context(), w)
}

func HandleAddProjectHeader(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
      <div class="flex space-x-2">
        <input type="text" name="HeaderKeys" placeholder="Key"
               class="flex-1 p-2 bg-gray-800 rounded" />
        <input type="text" name="HeaderValues" placeholder="Value"
               class="flex-1 p-2 bg-gray-800 rounded" />
      </div>
    `))
}

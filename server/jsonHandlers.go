package server

import (
	"encoding/json"
	"github.com/GuPoroca/HexTest/front"
	"github.com/GuPoroca/HexTest/front/components"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"io"
	"log"
	"net/http"
)

/////////////////////////////////////////////////////
// MAIN ENTRY
/////////////////////////////////////////////////////

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	currentProject = typeDefines.Project{
		Name:   "Untitled Project",
		Url:    "",
		Suites: []typeDefines.Suite{{Name: "New Suite"}},
	}
	component := front.Layout(currentProject, nil)
	component.Render(r.Context(), w)
}

/*************** MODALS ***************/
func HandleNewProjectConfirm(w http.ResponseWriter, r *http.Request) {
	components.NewConfirmModal().Render(r.Context(), w)
}

func HandleImportModal(w http.ResponseWriter, r *http.Request) {
	components.ImportModal().Render(r.Context(), w)
}

func HandleModalClose(w http.ResponseWriter, r *http.Request) {
	// Replace modal root with empty (closes modal)
	components.ModalClose().Render(r.Context(), w)
}

/*************** NEW ***************/
func HandleNewProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	currentProject = typeDefines.Project{
		Name:   "Untitled Project",
		Url:    "",
		Suites: []typeDefines.Suite{{Name: "New Suite"}},
	}

	// Refresh sidebar (OOB) and close modal
	components.ProjectSidebarOOB(r.Context(), w, currentProject)
	components.ModalClose().Render(r.Context(), w) // clears #modal-root

	// Replace main content with a simple message
	w.Write([]byte(`<div class="text-gray-400">New project created. Select an item from the sidebar to edit.</div>`))
}

/*************** IMPORT ***************/
func HandleJsonImport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		components.ImportModal().Render(r.Context(), w)
		return

	case http.MethodPost:
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
			http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file required: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "read error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var p typeDefines.Project
		if err := json.Unmarshal(data, &p); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Set in memory
		currentProject = p

		// Refresh sidebar (OOB) + close modal
		components.ProjectSidebarOOB(r.Context(), w, currentProject)
		components.ModalClose().Render(r.Context(), w)

		// Replace main content with a simple message
		w.Write([]byte(`<div class="text-gray-400">Project imported successfully. Select an item from the sidebar to edit.</div>`))
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

/*************** EXPORT ***************/
func HandleJsonExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="project.json"`)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(currentProject); err != nil {
		log.Println("export error:", err)
	}
}

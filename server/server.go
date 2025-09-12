package server

import (
	"log"
	"net/http"
)

func Run() {
	// The http.FileServer is used to serve static files like CSS and JS.
	// We will create a static directory for these files.
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// This handler will render the main page of our application.
	http.HandleFunc("/", HandleIndex)

	http.HandleFunc("/json/new/confirm", HandleNewProjectConfirm)
	http.HandleFunc("/json/new", HandleNewProject)
	http.HandleFunc("/json/import", HandleJsonImport) // GET (modal) + POST (upload)
	http.HandleFunc("/json/export", HandleJsonExport)
	http.HandleFunc("/modal/close", HandleModalClose)

	// EDIT
	http.HandleFunc("/edit/project", HandleEditProject)
	http.HandleFunc("/edit/suite", HandleEditSuite)
	http.HandleFunc("/edit/test", HandleEditTest)
	http.HandleFunc("/edit/assert", HandleEditAssert)
	http.HandleFunc("/edit/check", HandleEditCheck)

	// SAVE (in-memory updates; only write JSON if you want here)
	http.HandleFunc("/save/project", HandleSaveProject)
	http.HandleFunc("/save/suite", HandleSaveSuite)
	http.HandleFunc("/save/test", HandleSaveTest)
	http.HandleFunc("/save/assert", HandleSaveAssert)
	http.HandleFunc("/save/check", HandleSaveCheck)

	// ADD
	http.HandleFunc("/add/suite", HandleAddSuite)
	http.HandleFunc("/add/test", HandleAddTest)
	http.HandleFunc("/add/assert", HandleAddAssert)
	http.HandleFunc("/add/check", HandleAddCheck)
	http.HandleFunc("/add/check/expected", HandleAddCheckExpected)

	// DELETE
	http.HandleFunc("/delete/suite", HandleDeleteSuite)
	http.HandleFunc("/delete/test", HandleDeleteTest)
	http.HandleFunc("/delete/assert", HandleDeleteAssert)
	http.HandleFunc("/delete/check", HandleDeleteCheck)

	//Headers
	http.HandleFunc("/add/project/header", HandleAddProjectHeader)
	http.HandleFunc("/add/test/header", HandleAddTestHeader)

	log.Println("Starting frontend on :3773")
	if err := http.ListenAndServe(":3773", nil); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}

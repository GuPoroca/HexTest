package jsonOperations

import (
	"encoding/json"
	"fmt"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"io"
	"os"
)

func ReadJSON() typeDefines.Project {

	jsonFile, err := os.Open("../../project_mock.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", jsonFile.Name())

	byteValue, _ := io.ReadAll(jsonFile)
	var project typeDefines.Project
	json.Unmarshal(byteValue, &project)

	defer jsonFile.Close()
	SetDefaults(&project)
	return project
}

func PrettyPrint(project typeDefines.Project) {
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(data))
}

func SetDefaults(project *typeDefines.Project) {
	//Project part
	//Suites part

	//Tests part
	for _, suites := range project.Suites {
		for i := range suites.Tests {
			if suites.Tests[i].Content_type == "" {
				suites.Tests[i].Content_type = "application/json; charset=UTF-8"
			}
		}
	}
}

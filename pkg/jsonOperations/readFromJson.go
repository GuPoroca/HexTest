package jsonOperations

import (
	"encoding/json"
	"fmt"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"io"
	"os"
)

func ReadJSON(path string) typeDefines.Project {

	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", jsonFile.Name())

	byteValue, _ := io.ReadAll(jsonFile)
	var project typeDefines.Project

	json.Unmarshal(byteValue, &project)

	defer jsonFile.Close()
	return project
}

func PrettyPrint(any_struct any) {
	data, err := json.MarshalIndent(any_struct, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(data))
}

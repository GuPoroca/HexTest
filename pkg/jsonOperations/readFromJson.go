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
	fmt.Println("Successfully Opened users.json")

	byteValue, _ := io.ReadAll(jsonFile)
	var projeto typeDefines.Project
	json.Unmarshal(byteValue, &projeto)

	defer jsonFile.Close()

	return projeto
}

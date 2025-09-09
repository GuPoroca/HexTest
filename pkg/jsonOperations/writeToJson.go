package jsonOperations

import (
	"encoding/json"
	"fmt"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"os"
)

func WriteJSON(path string, currentProject typeDefines.Project) {
	jsonData, err := json.MarshalIndent(currentProject, "", "  ")
	if err != nil {
		fmt.Print(err)
	}
	err = os.WriteFile(currentProject.Name+".json", jsonData, 0644)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print("Saved with success!")

}

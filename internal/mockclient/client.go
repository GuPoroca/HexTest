package mockclient

import (
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
)

func MakeAllRequests() {
	project := jsonOperations.ReadJSON()

	project.ExecuteProject()
}

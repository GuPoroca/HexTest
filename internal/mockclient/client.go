package mockclient

import (
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
)

func MakeAllRequests() {
	path := "../../one_request.json"
	project := jsonOperations.ReadJSON(path)

	project.ExecuteProject()
}

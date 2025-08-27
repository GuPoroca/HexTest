package typedefines

import (
// "fmt"
// "github.com/pkg/errors"
// "io/ioutil"
// "net/http"
)

type Suite struct {
	name     string
	endpoint string
	tests    []Test
}

type Project struct {
	name   string
	url    string
	auth   Auth
	suites []Suite
}

type Auth struct {
	auth_type string
	token     string
}

package goTester

import (
	"fmt"
	"github.com/pkg/errors"
	"io/util"
	"net/http"
)

type Test struct {
	req_body       string
	exp_res_body   string
	exp_res_status int
}

type Suite struct {
	endpoint string
	tests    []Test
}

type Project struct {
	suites []Suite
	url    string
	token  string
}

func (test *Test) execute() {

}

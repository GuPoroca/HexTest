package typedefines

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Test struct {
	//Necessary {user input}
	name           string
	method         string
	req_body       string
	exp_res_body   string
	exp_res_status string
	//Necessary {defaulted values}
	content_type string
	//Response related
	last_resp_body   string
	last_resp_status string
	//Extra sugar
	comment string
}

func NewTest(name, method, req_body, exp_res_body, exp_res_status string) *Test {
	var test *Test = &Test{}
	test.name = name
	test.method = method
	test.req_body = req_body
	test.exp_res_body = exp_res_body
	test.exp_res_status = exp_res_status
	test.content_type = "application/json"
	return test
}

func (test *Test) Execute(full_url string, auth *Auth) error {
	var err error

	fmt.Printf("Executing test: %s\n", test.name)

	request, err := http.NewRequest(test.method, full_url, strings.NewReader(test.req_body))
	if err != nil {
		fmt.Printf("An error ocurred while creating the request %v\n", err)
	}

	request.Header.Add("Content-type", "application/json; charset=UTF-8")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("An error ocurred while making the request %v\n", err)
	}
	out, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("An error ocurred when reading the response %v\n", err)
	}

	test.last_resp_status = response.Status
	test.last_resp_body = string(out)

	fmt.Printf("\nStatus: %s\n", test.last_resp_status)
	fmt.Printf("Body: \n%s\n", test.last_resp_body)
	return nil
}

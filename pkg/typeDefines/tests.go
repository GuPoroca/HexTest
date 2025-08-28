package typeDefines

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Test struct {
	//Necessary {user input}
	Name                     string `json:"Name"`
	Method                   string `json:"Method"`
	Req_body                 string `json:"Req_body"`
	Expected_response_body   string `json:"Expected_response_body"`
	Expected_response_status string `json:"Expected_response_status"`
	Api_endpoint             string `json:"Api_endpoint"`
	//Necessary {defaulted values}
	Content_type string `json:"Content_type"`
	//Response related
	Last_response_body   string `json:"Last_response_body"`
	Last_response_status string `json:"Last_response_status"`
	//Extra sugar
	Comment string `json:"Comment"`
}

func NewTest(name, method, req_body, exp_res_body, exp_res_status, api_endpoint string) *Test {
	var test *Test = &Test{}
	test.Name = name
	test.Method = method
	test.Req_body = req_body
	test.Expected_response_body = exp_res_body
	test.Expected_response_status = exp_res_status
	test.Api_endpoint = api_endpoint
	test.Content_type = "application/json"
	return test
}

func (test *Test) Execute(url string, auth Auth) error {
	var err error
	full_url := url + test.Api_endpoint

	fmt.Printf("Executing test: %s\n", test.Name)

	request, err := http.NewRequest(test.Method, full_url, strings.NewReader(test.Req_body))
	if err != nil {
		fmt.Printf("An error ocurred while creating the request %v\n", err)
		return err
	}

	request.Header.Add("Content-type", "application/json; charset=UTF-8")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("An error ocurred while making the request %v\n", err)
		return err
	}
	out, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("An error ocurred when reading the response %v\n", err)
		return err
	}

	test.Last_response_status = response.Status
	test.Last_response_body = string(out)

	fmt.Printf("\nStatus: %s\n", test.Last_response_status)
	fmt.Printf("Body: \n%s\n", test.Last_response_body)
	return nil
}

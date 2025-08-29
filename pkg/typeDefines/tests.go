package typeDefines

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
	//"unsafe"
)

type Test struct {
	//Necessary {user input}
	Name                     string `json:"Name"`
	Method                   string `json:"Method"`
	Request_body             string `json:"Request_body"`
	Expected_response_body   string `json:"Expected_response_body"`
	Expected_response_status string `json:"Expected_response_status"`
	Api_endpoint             string `json:"Api_endpoint"`
	//Necessary {defaulted values}
	Content_type string `json:"Content_type"`
	//Response related
	Last_response_body   string `json:"Last_response_body"`
	Last_response_status string `json:"Last_response_status"`
	Time_to_respond      int64
	Last_response_size   int
	//Extra sugar
	Comment string `json:"Comment"`
}

func (test *Test) Execute(url string, auth Auth) error {
	var err error
	full_url := url + test.Api_endpoint

	fmt.Printf("\nExecuting Test: %s\n", test.Name)

	request, err := http.NewRequest(test.Method, full_url, strings.NewReader(test.Request_body))
	if err != nil {
		fmt.Printf("An error ocurred while creating the request %v\n", err)
		return err
	}

	test.AddAllHeaders(*request)

	start_time := time.Now()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("An error ocurred while making the request %v\n", err)
		return err
	}

	test.Time_to_respond = time.Since(start_time).Milliseconds()

	out, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("An error ocurred when reading the response %v\n", err)
		return err
	}
	response.Body.Close()
	response.Body = io.NopCloser(bytes.NewBuffer(out))

	test.checkResponseSize(*response)

	test.Last_response_status = response.Status
	test.Last_response_body = string(out)

	fmt.Printf("\nSize: %v Bytes\n", test.Last_response_size)
	fmt.Printf("Time to execute: %vms\n", test.Time_to_respond)
	fmt.Printf("Status: %s\n", test.Last_response_status)
	fmt.Printf("Body: \n%s\n", test.Last_response_body)

	fmt.Print("\n---------------------------------------\n")

	return nil
}

func (test *Test) checkResponseSize(resp http.Response) {
	dump, err := httputil.DumpResponse(&resp, true)
	if err != nil {
		fmt.Printf("Error dumping response: %v", err)
	} else {
		test.Last_response_size = len(dump)
	}
}

func (test Test) AddAllHeaders(req http.Request) {
	req.Header.Add("Content-type", test.Content_type)
}

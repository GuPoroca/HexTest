package typeDefines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type Test struct {
	//Necessary {user input}
	Name         string `json:"Name"`
	Method       string `json:"Method"`
	Request_body string `json:"Request_body"`
	Api_endpoint string `json:"Api_endpoint"`
	//Custom Request_Headers
	Request_Headers map[string]string `json:"Request_Headers"`
	//Response related
	Response_body        map[string]any
	Response_body_string string
	Response_status      string
	Time_to_respond      int64
	Response_size        int64
	//Assertion related
	Asserts []Assert `json:"Asserts"`
	Passed  []bool
	//Extra sugar
	Comment string `json:"Comment"`
}

func (test *Test) Execute(url string, auth IAuth) error {
	var err error
	full_url := url + test.Api_endpoint

	fmt.Printf("\nExecuting Test: %s\n", test.Name)

	request, err := http.NewRequest(test.Method, full_url, strings.NewReader(test.Request_body))
	if err != nil {
		log.Fatalf("An error ocurred while creating the request %v\n", err)
		return err
	}
	if auth != nil {
		token, err := auth.Authenticate()
		if err != nil {
			log.Fatalf("An error ocurred during the token request %v\n", err)
		}
		request.Header.Add("Authorization", token)
		fmt.Printf("\nAutenticado com sucesso!\n\n")
	}

	test.AddAllHeaders(*request)

	start_time := time.Now()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalf("An error ocurred while making the request %v\n", err)
		return err
	}

	test.Time_to_respond = time.Since(start_time).Milliseconds()

	//puts response.Body in a []byte
	out, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("An error ocurred when reading the response %v\n", err)
		return err
	}

	test.Response_body_string = string(out)

	if !json.Valid(out) {
		log.Fatalf("Response contains a invalid json body")
		return err
	}

	//replaces the response.Body content with a copy of the []byte
	response.Body = io.NopCloser(bytes.NewBuffer(out))

	//maps the response.body to test.Response_body (map[string]any)
	err = json.NewDecoder(response.Body).Decode(&test.Response_body)
	if err != nil {
		log.Fatalf("An error occurred when putting the response in the map %v\n", err)
		return err
	}
	//replaces the response.Body again
	response.Body = io.NopCloser(bytes.NewBuffer(out))

	response.Body.Close()
	test.checkResponseSize(*response)

	test.Response_status = response.Status

	fmt.Printf("Response Body in json: %s\n", test.Response_body_string)
	fmt.Printf("Status: %s\n", test.Response_status)
	fmt.Printf("Size: %v Bytes\n", test.Response_size)
	fmt.Printf("Time to execute: %vms\n", test.Time_to_respond)

	fmt.Print("\nRunning Assertions:\n\n")
	test.runAllAssertions()
	fmt.Print("\n---------------------------------------\n")

	return nil
}

func (test *Test) runAllAssertions() bool {
	result := 0.0
	all_passed := true
	var value any
	for i := range test.Asserts {
		switch test.Asserts[i].Field {
		case "Response Body":
			value = test.Response_body
		case "Response Status":
			value = test.Response_status
		case "Response Time":
			value = test.Time_to_respond
		case "Response Size":
			value = test.Response_size
		default:
			fmt.Printf("Assertion field \"%s\" is invalid", test.Asserts[i].Field)
			continue
		}
		fmt.Printf("\tAsserting %s\n\n", test.Asserts[i].Field)

		test.Passed = append(test.Passed, test.Asserts[i].MakeAssertions(value))
		if test.Passed[i] {
			result++
		} else {
			all_passed = false
		}
	}
	fmt.Printf("Assertions passed: %v/%v", result, len(test.Asserts))
	return all_passed
}

func (test *Test) checkResponseSize(resp http.Response) {
	dump, err := httputil.DumpResponse(&resp, true)
	if err != nil {
		fmt.Printf("Error dumping response: %v", err)
	} else {
		test.Response_size = int64(len(dump))
	}
}

func (test Test) AddAllHeaders(req http.Request) {
	for k := range test.Request_Headers {
		req.Header.Add(k, test.Request_Headers[k])
	}
}

package typeDefines

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"maps"
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
	Comment      string `json:"Comment"`
	//Custom Request_Headers
	Request_Headers map[string]string `json:"Request_Headers"`
	//Response related
	Response_body        map[string]any
	Response_body_string string
	Response_Headers     map[string][]string
	Response_status      string
	Time_to_respond      int64
	Response_size        int64
	//Assertion related
	Asserts []Assert `json:"Asserts"`
}

func (test *Test) Execute(url string, headermap map[string]string) error {
	var err error
	full_url := url + test.Api_endpoint

	request, err := http.NewRequest(test.Method, full_url, strings.NewReader(test.Request_body))
	if err != nil {
		log.Printf("An error ocurred while creating the request %v\n", err)
		return err
	}

	test.AddAllHeaders(*request, headermap)

	start_time := time.Now()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("An error ocurred while making the request %v\n", err)
		return err
	}

	test.Time_to_respond = time.Since(start_time).Milliseconds()

	//puts response.Body in a []byte
	out, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("An error ocurred when reading the response %v\n", err)
		return err
	}

	test.Response_body_string = string(out)

	if !json.Valid(out) {
		log.Printf("Response contains a invalid json body")
		return err
	}

	//replaces the response.Body content with a copy of the []byte
	response.Body = io.NopCloser(bytes.NewBuffer(out))

	//maps the response.body to test.Response_body (map[string]any)
	err = json.NewDecoder(response.Body).Decode(&test.Response_body)
	if err != nil {
		log.Printf("An error occurred when putting the response in the map %v\n", err)
		return err
	}
	//replaces the response.Body again
	response.Body = io.NopCloser(bytes.NewBuffer(out))

	response.Body.Close()
	test.checkResponseSize(*response)

	test.Response_Headers = response.Header
	test.Response_status = strings.Split(response.Status, " ")[0]

	test.runAllAssertions()

	return nil
}

func (test *Test) runAllAssertions() {
	for i := range test.Asserts {
		value := test.getValueForTest(i)
		if typeOfValue(value) == "string" {
			if value == "Invalid Assert" {
				continue
			}
		}
		test.Asserts[i].MakeAssertions(value)
	}
}

func (test *Test) getValueForTest(i int) any {
	var value any
	switch test.Asserts[i].Field {
	case "Response Body":
		value = test.Response_body
	case "Response Status":
		value = test.Response_status
	case "Response Time":
		value = test.Time_to_respond
	case "Response Size":
		value = test.Response_size
	case "Response Headers":
		value = test.Request_Headers
	case "JSON Schema Validation":
		value = test.Response_body
	default:
		//special cases
		if subFields := strings.Split(test.Asserts[i].Field, "."); subFields[0] == "Value of Body" {
			if len(subFields) == 1 {
				value = test.Response_body
			} else {
				value, _ = getSpecificVal(subFields[1:], test.Response_body)
			}
		} else if subFields := strings.Split(test.Asserts[i].Field, "."); subFields[0] == "Value of Headers" {
			if len(subFields) == 1 {
				value = test.Response_Headers
			} else {
				value = test.Response_Headers[subFields[1]][0]
			}
		} else if subFields := strings.Split(test.Asserts[i].Field, "."); subFields[0] == "Type of Body" {
			if len(subFields) == 1 {
				value = test.Response_body
				value = typeOfValue(value)
			} else {
				value, _ = getSpecificVal(subFields[1:], test.Response_body)
				value = typeOfValue(value)
			}
		} else { //end of special cases
			return "Invalid Assert"
		}
	}
	return value
}

func (test *Test) checkResponseSize(resp http.Response) {
	dump, err := httputil.DumpResponse(&resp, true)
	if err != nil {
		log.Printf("Error dumping response: %v", err)
	} else {
		test.Response_size = int64(len(dump))
	}
}

func (test Test) AddAllHeaders(req http.Request, hmps map[string]string) {
	maps.Copy(hmps, test.Request_Headers)
	for k := range hmps {
		req.Header.Add(k, hmps[k])
	}
}

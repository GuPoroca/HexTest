package typeDefines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strconv"
	"strings"
	"time"
	//"unsafe"
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
	Asserts []Assert `json:"-"`
	Result  bool
	//Extra sugar
	Comment string `json:"Comment"`
}

func (test *Test) Execute(url string, auth IAuth) error {
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

	//puts response.Body in a []byte
	out, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("An error ocurred when reading the response %v\n", err)
		return err
	}
	test.Response_body_string = string(out)
	if !json.Valid(out) {
		fmt.Printf("Response contains a invalid json body")
		return err
	}

	//replaces the response.Body content with a copy of the []byte
	response.Body = io.NopCloser(bytes.NewBuffer(out))

	//maps the response.body to test.Response_body (map[string]any)
	err = json.NewDecoder(response.Body).Decode(&test.Response_body)
	if err != nil {
		fmt.Printf("An error occurred when putting the response in the map %v\n", err)
		return err
	}
	//replaces the response.Body again
	response.Body = io.NopCloser(bytes.NewBuffer(out))

	response.Body.Close()
	test.checkResponseSize(*response)

	test.Response_status = response.Status

	fmt.Printf("Response Body in json: %s\n", test.Response_body_string)
	fmt.Printf("Size: %v Bytes\n", test.Response_size)
	fmt.Printf("Time to execute: %vms\n", test.Time_to_respond)

	fmt.Print("\n---------------------------------------\n")
	fmt.Print("Running Assertions:\n")
	test.runAllAssertions()

	return nil
}

func (test *Test) runAllAssertions() bool {
	result := true
	var value any
	for i := range test.Asserts {
		switch test.Asserts[i].FieldToCheck {
		case "Response Body":
			value = test.Response_body
		case "Response Status":
			value = test.Response_status
		case "Response Time":
			value = test.Time_to_respond
		case "Response Size":
			value = test.Response_size
		default:
			fmt.Printf("Assertion field \"%s\" is invalid", test.Asserts[i].FieldToCheck)
			continue
		}
		fmt.Printf("\tAsserting %s\n", test.Asserts[i].FieldToCheck)
		result = test.Asserts[i].MakeAssertions(value)
	}
	return result
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

// custom unmarshler
func (t *Test) UnmarshalJSON(data []byte) error {

	type TestAlias Test
	aux := &struct {
		Asserts []Assert `json:"-"`
		*TestAlias
	}{
		TestAlias: (*TestAlias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var rawAsserts struct {
		Asserts []map[string]map[string]any `json:"Asserts"`
	}
	if err := json.Unmarshal(data, &rawAsserts); err != nil {
		return err
	}

	var mappedAsserts []Assert
	for _, assertBlock := range rawAsserts.Asserts {
		for field, checksMap := range assertBlock {
			newAssert := Assert{
				FieldToCheck: field,
				Checks:       []Check{},
			}
			for operand, value := range checksMap {
				if NumberField(field) && reflect.TypeOf(value) == reflect.TypeOf("string") {
					num_value, err := strconv.ParseFloat(value.(string), 64)
					if err != nil {
						fmt.Printf("\nWARNING: There was a problem in the number conversion on:")
						fmt.Printf("\n Test: \"%s\"\n  Assertion: \"%s\"\n   Operand: \"%s\"\n    Value:\"%s\"", t.Name, newAssert.FieldToCheck, operand, value)
						fmt.Printf("\nPlease fix the JSON. The numbers should not be between quotes")
					}
					newAssert.Checks = append(newAssert.Checks, Check{
						Operand: operand,
						Value:   num_value,
					})
					continue

				}
				newAssert.Checks = append(newAssert.Checks, Check{
					Operand: operand,
					Value:   value,
				})
			}
			mappedAsserts = append(mappedAsserts, newAssert)
		}
	}
	t.Asserts = mappedAsserts

	return nil
}

func NumberField(field string) bool {
	switch field {
	case "Response Time":
		return true
	case "Response Size":
		return true
	default:
		return false
	}
}

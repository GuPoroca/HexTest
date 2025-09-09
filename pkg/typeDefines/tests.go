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
	Asserts                 []Assert `json:"Asserts"`
	Passed                  []bool
	Passed_Comparissons_num int
	Total_Comparissons_num  int
	//Extra sugar
	Comment string `json:"Comment"`
}

func (test *Test) Execute(url string, auth IAuth) error {
	var err error
	full_url := url + test.Api_endpoint

	fmt.Printf("\nExecuting Test: %s\n", test.Name)

	request, err := http.NewRequest(test.Method, full_url, strings.NewReader(test.Request_body))
	if err != nil {
		log.Printf("An error ocurred while creating the request %v\n", err)
		return err
	}
	//check how to better do authentication rn

	// if auth != nil {
	// 	token, err := auth.Authenticate()
	// 	if err != nil {
	// 		log.Fatalf("An error ocurred during the token request %v\n", err)
	// 	}
	// 	request.Header.Add("Authorization", token)
	// 	fmt.Printf("\nAutenticado com sucesso!\n\n")
	// }

	test.AddAllHeaders(*request)

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

	test.Response_status = response.Status

	fmt.Printf("Response Body in json: %s\n", test.Response_body_string)
	fmt.Printf("Status: %s\n", test.Response_status)
	fmt.Printf("Size: %v Bytes\n", test.Response_size)
	fmt.Printf("Time to execute: %vms\n", test.Time_to_respond)

	fmt.Print("\nRunning Assertions:\n\n")
	test.runAllAssertions()
	fmt.Print("\n---------------------------------------\n")
	fmt.Printf("Total comparissons passed %d/%d", test.Passed_Comparissons_num, test.Total_Comparissons_num)

	return nil
}

func (test *Test) runAllAssertions() int {
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
			//special cases
			if subFields := strings.Split(test.Asserts[i].Field, "."); subFields[0] == "Value of Body" {
				if len(subFields) == 1 {
					value = test.Response_body
				} else {
					value, _ = getSpecificVal(subFields[1:], test.Response_body)
				}
			} else {
				fmt.Printf("Assertion field \"%s\" is invalid", test.Asserts[i].Field)
				continue
			}
		}
		fmt.Printf("\tAsserting %s\n\n", test.Asserts[i].Field)

		test.Passed_Comparissons_num += test.Asserts[i].MakeAssertions(value)
		test.Total_Comparissons_num += test.Asserts[i].Total_Comparissons_num
		if (test.Asserts[i].Passed_Comparissons_num - len(test.Asserts[i].Checks)) == 0 {
			test.Passed = append(test.Passed, true)
		} else {
			test.Passed = append(test.Passed, false)
		}
	}
	return test.Passed_Comparissons_num
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

func getSpecificVal(fa []string, m any) (any, bool) {
	if len(fa) == 0 {
		return nil, false
	}
	switch val := m.(type) {
	case map[string]any:

		if len(fa) == 1 {
			return val[fa[0]], true
		} else {
			if val, ok := val[fa[0]]; ok {
				return getSpecificVal(fa[1:], val)
			} else {
				return nil, false
			}
		}
	case []any:
		for i := range val {
			if val[i].(map[string]any)[fa[0]] != nil {
				if len(fa) == 1 {
					return (val[i].(map[string]any)[fa[0]]), true
				} else {
					return getSpecificVal(fa[1:], val[i])
				}
			}
		}
	}
	return nil, false
}

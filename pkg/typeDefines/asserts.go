package typeDefines

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type Assert struct {
	FieldToCheck       string
	FieldResponseValue any
	Checks             []Check
	Result             bool
}

func (assert *Assert) MakeAssertions(fieldValue any) bool {
	assert.FieldResponseValue = fieldValue
	assert.Result = true
	str_field_value := stringifyMyAny(fieldValue)
	for i := range assert.Checks {
		assert.RunCheck(i, str_field_value)
	}
	return assert.Result
}

func (assert *Assert) RunCheck(i int, str_field_value string) {
	checkValue := reflect.ValueOf(assert.Checks[i].Value)
	if checkValue.Kind() == reflect.Slice { //multiple checks for same field
		for j := 0; j < checkValue.Len(); j++ {
			element := checkValue.Index(j).Interface()
			assert.SingleCheck(element, i, str_field_value)
		}
	} else { //single check
		assert.SingleCheck(assert.Checks[i].Value, i, str_field_value)
	}
}

func (assert *Assert) SingleCheck(element any, i int, str_field_value string) {
	assert.Checks[i].Value = matchMyType(assert.FieldToCheck, element)
	str_check_value := stringifyMyAny(assert.Checks[i].Value)
	fmt.Printf("\t%s %s %s\n", str_field_value, assert.Checks[i].Operand, str_check_value)
	if !assert.Checks[i].MakeCheck(assert.FieldResponseValue) {
		fmt.Printf("\t RESULT: FAILED\n\n")
		assert.Result = false
	} else {
		fmt.Printf("\t RESULT: PASSED\n\n")
	}
}

func stringifyMyAny(myAny any) string {
	var str string
	if _, ok := myAny.(bool); ok {
		str = fmt.Sprintf("%t", myAny)
	} else if _, ok := myAny.(string); ok {
		str = myAny.(string)
	} else if _, ok := myAny.(int64); ok {
		str = fmt.Sprintf("%d", myAny)
	} else if _, ok := myAny.(float64); ok {
		str = fmt.Sprintf("%f", myAny)
	} else if _, ok := myAny.(map[string]any); ok {
		str, err := json.Marshal(myAny) //solvethis
		if err != nil {
			log.Fatal(err)
		}
		return string(str)
	}
	return str
}

func matchMyType(field string, current any) any {
	switch field {
	case "Response Body": //string
		if _, ok := current.(string); !ok {
			fmt.Printf("WARNING: %v should be of type string on the json, Test will probably fail\n", field)
		} else {
			var responseBodyMap map[string]any
			json.Unmarshal([]byte(current.(string)), &responseBodyMap)
			if responseBodyMap != nil {
				return responseBodyMap
			}
		}
	case "Response Status": //string
		if _, ok := current.(string); !ok {
			fmt.Printf("WARNING: %v should be of type string on the json, Test will probably fail\n", field)
		}
	case "Response Time": //int64
		if _, ok := current.(int64); !ok {
			return int64(current.(float64))
		}
	case "Response Size": //int64
		if _, ok := current.(int64); !ok {
			return int64(current.(float64))
		}

	}
	return current
}

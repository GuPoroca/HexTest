package typeDefines

import (
	"encoding/json"
	"fmt"
	"log"
)

type Assert struct {
	Field              string  `json:"Field"`
	Checks             []Check `json:"Checks"`
	Results            []bool
	FieldResponseValue any
	Passed             int
}

func (assert *Assert) MakeAssertions(fieldValue any) bool {
	assert.FieldResponseValue = fieldValue
	str_field_value := stringifyMyAny(fieldValue)
	result := 0
	all_passed := true
	for i := range assert.Checks {
		assert.Results = append(assert.Results, assert.Checks[i].MakeAllChecks(str_field_value))
		if assert.Results[i] {
			result++
		} else {
			all_passed = false
		}
	}
	fmt.Printf("Checks passed: %v/%v", result, len(assert.Checks))
	assert.Passed = result
	return all_passed
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

package typeDefines

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/GuPoroca/HexTest/bus"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

type Check struct {
	Operand    string `json:"Operand"`
	Expected   []any  `json:"Expected"`
	Passed     []int
	Passed_num int
	Total_num  int
}

type MockT struct {
	errors []string
}

func (m MockT) Errorf(format string, args ...any) {
	errorMsg := fmt.Sprintf(format, args...)
	m.errors = append(m.errors, errorMsg)
}

func (check *Check) MakeAllChecks(responseVal any) int {
	check.Total_num = len(check.Expected)
	if check.Total_num == 0 {
		check.Total_num = 1
		if p, _ := check.MakeCheckWithoutExpected(responseVal); p > 0 {
			check.Passed_num++
		}
	}
	for i := range check.Expected {
		if p, _ := check.MakeCheck(responseVal, i); p > 0 {
			check.Passed_num++
		}
	}
	return check.Passed_num
}

func putValInChan() {
	bus.CheckEvents <- 1
}

func (check *Check) MakeCheck(responseVal any, i int) (int, error) {
	t := &MockT{}
	expectedVal := check.Expected[i]
	value := -2

	passed := false
	defer putValInChan()

	switch check.Operand {

	case "==", "!=", ">=", "<=", ">", "<":
		//try numbers first
		resFloat, resOk := toFloat64(responseVal)
		expFloat, expOk := toFloat64(expectedVal)
		if expFloat == math.Trunc(expFloat) && resFloat == math.Trunc(resFloat) {
		}

		if resOk && expOk {
			switch check.Operand {
			case "==":
				passed = assert.Equal(t, resFloat, expFloat)
			case "!=":
				passed = assert.NotEqual(t, resFloat, expFloat)
			case ">=":
				passed = assert.GreaterOrEqual(t, resFloat, expFloat)
			case "<=":
				passed = assert.LessOrEqual(t, resFloat, expFloat)
			case ">":
				passed = assert.Greater(t, resFloat, expFloat)
			case "<":
				passed = assert.Less(t, resFloat, expFloat)
			}
			goto checkPassed
		}
		//try dates
		resStr, resIsStr := responseVal.(string)
		expStr, expIsStr := expectedVal.(string)
		if resIsStr && expIsStr {
			resTime, resOk := tryParseTime(resStr)
			expTime, expOk := tryParseTime(expStr)

			if resOk && expOk {
				// We have two valid dates, so compare them
				switch check.Operand {
				case "==":
					passed = assert.True(t, resTime.Equal(expTime))
				case "!=":
					passed = assert.False(t, resTime.Equal(expTime))
				case ">=":
					passed = assert.True(t, resTime.Equal(expTime) || resTime.After(expTime))
				case "<=":
					passed = assert.True(t, resTime.Equal(expTime) || resTime.Before(expTime))
				case ">":
					passed = assert.True(t, resTime.After(expTime))
				case "<":
					passed = assert.True(t, resTime.Before(expTime))
				}
				goto checkPassed
			}
		}
	case "matchRegex":
		passed = assert.Regexp(t, expectedVal, responseVal)
	case "notMatchRegex":
		passed = assert.NotRegexp(t, expectedVal, responseVal)
	case "containsSubstring":
		passed = assert.Contains(t, responseVal, expectedVal)
	case "containsKey":
		passed = assert.Contains(t, responseVal, expectedVal)
	case "containsKey -R":
		_, passed = ContainsKeyRecursevely(responseVal, expectedVal.(string))
	default:
		err := OperandNotFound{check.Operand}
		return -1, err
	}

checkPassed:
	{
		if passed == true {
			value = 1
		} else {
			value = 0
		}

		check.Passed = append(check.Passed, value)
		return value, nil
	}

}

func (check *Check) MakeCheckWithoutExpected(responseVal any) (int, error) {
	defer putValInChan()
	t := &MockT{}
	passed := false

	switch check.Operand {
	case "isNull":
		passed = assert.Nil(t, responseVal)
	case "notNull":
		passed = assert.NotNil(t, responseVal)
	case "isEmpty":
		passed = assert.Empty(t, responseVal)
	case "notEmpty":
		passed = assert.NotEmpty(t, responseVal)
	default:
		err := OperandNotFound{check.Operand}
		return -1, err
	}
	value := 0
	if passed == true {
		value = 1
	} else {
		value = 0
	}
	check.Passed = append(check.Passed, value)
	return value, nil

}

func (check *Check) JsonSchema(responseVal any) int {
	defer putValInChan()
	schemaStr := check.Expected[0].(string)
	passed := false
	ok, _ := validateAgainstSchema(schemaStr, responseVal)
	if ok < 1 {
		passed = false
	} else {
		passed = true
	}
	value := 0
	if passed == true {
		value = 1
	} else {
		value = 0
	}
	return value
}

func validateAgainstSchema(schemaStr string, body any) (int, []string) {
	schemaLoader := gojsonschema.NewStringLoader(schemaStr)

	// encode body back to JSON string
	bodyBytes, _ := json.Marshal(body)
	docLoader := gojsonschema.NewBytesLoader(bodyBytes)

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return 0, []string{err.Error()}
	}

	if result.Valid() {
		return 1, nil
	}

	errs := []string{}
	for _, desc := range result.Errors() {
		errs = append(errs, desc.String())
	}
	return -1, errs
}

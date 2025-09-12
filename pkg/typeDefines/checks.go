package typeDefines

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"
)

type Check struct {
	Operand    string `json:"Operand"`
	Expected   []any  `json:"Expected"`
	Passed     []bool
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
		if check.MakeCheckWithoutExpected(responseVal) {
			check.Passed_num++
		}
	}
	for i := range check.Expected {
		if check.MakeCheck(responseVal, i) {
			check.Passed_num++
		}
	}
	return check.Passed_num
}

func (check *Check) MakeCheck(responseVal any, i int) bool {
	t := &MockT{}
	expectedVal := check.Expected[i]
	passed := false

	check.printCheckStart(responseVal, i)

	switch check.Operand {
	case "==":
		passed = assert.True(t, reflect.DeepEqual(responseVal, expectedVal))
	case "!=":
		passed = assert.False(t, reflect.DeepEqual(responseVal, expectedVal))
	case ">=", "<=", ">", "<":
		//Comparisons of ><
		//try numbers first

		resFloat, resOk := toFloat64(responseVal)
		expFloat, expOk := toFloat64(expectedVal)
		if expFloat == math.Trunc(expFloat) && resFloat == math.Trunc(resFloat) {
		}

		if resOk && expOk {
			switch check.Operand {
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
		fmt.Printf("Operand \"%s\" is not recognized", check.Operand)
	}

checkPassed:
	{
		check.Passed = append(check.Passed, passed)
		printCheckEnd(passed)
		return passed
	}

}

func (check *Check) MakeCheckWithoutExpected(responseVal any) bool {
	t := &MockT{}
	passed := false

	check.printCheckStart(responseVal, -1)

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
		fmt.Printf("Operand \"%s\" is not recognized", check.Operand)
	}
	check.Passed = append(check.Passed, passed)
	printCheckEnd(passed)
	return passed

}

func (check Check) printCheckStart(responseVal any, i int) {
	var expectedVal any
	if i != -1 {
		expectedVal = check.Expected[i]
	} else {
		expectedVal = ""
	}
	//checks if they can be floats, if so checks if they can be ints
	resFloat, resOk := toFloat64(responseVal)
	expFloat, expOk := toFloat64(expectedVal)

	if resOk && expOk {
		if expFloat == math.Trunc(expFloat) && resFloat == math.Trunc(resFloat) {
			expectedVal = int(expFloat)
			responseVal = int(resFloat)
		}
	}

	fmt.Printf("\t%s %s %s\n", StringifyMyAny(responseVal), check.Operand, StringifyMyAny(expectedVal))
}

func printCheckEnd(passed bool) {
	if passed {
		fmt.Printf("PASSED\n")
	} else {
		fmt.Printf("FAILED\n")
	}
}

func (check *Check) JsonSchema(responseVal any) bool {
	check.printCheckStart(responseVal, 0)
	schemaStr := check.Expected[0].(string)
	passed := false
	ok, _ := validateAgainstSchema(schemaStr, responseVal)
	if !ok {
		passed = false
	} else {
		passed = true
	}
	printCheckEnd(passed)
	return passed
}

func validateAgainstSchema(schemaStr string, body any) (bool, []string) {
	schemaLoader := gojsonschema.NewStringLoader(schemaStr)

	// encode body back to JSON string
	bodyBytes, _ := json.Marshal(body)
	docLoader := gojsonschema.NewBytesLoader(bodyBytes)

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return false, []string{err.Error()}
	}

	if result.Valid() {
		return true, nil
	}

	errs := []string{}
	for _, desc := range result.Errors() {
		errs = append(errs, desc.String())
	}
	return false, errs
}

func tryParseTime(s string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339, // "2006-01-02T15:04:05Z07:00"
		"2006-01-02", // "YYYY-MM-DD"
		"02/01/2006", // "DD/MM/YYYY"
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func toFloat64(v any) (float64, bool) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), true
	case reflect.Float32, reflect.Float64:
		return val.Float(), true
	case reflect.String:
		if ret, err := strconv.ParseFloat(val.String(), 64); err == nil {
			return ret, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func ContainsKeyRecursevely(responseVal any, targetVal string) (any, bool) {
	switch v := responseVal.(type) {
	case map[string]any:
		for key, val := range v {
			if key == targetVal {
				return v, true
			}
			if _, ok := ContainsKeyRecursevely(val, targetVal); ok {
				return v, true
			}
		}
	case []any:
		for _, item := range v {
			if _, ok := ContainsKeyRecursevely(item, targetVal); ok {
				return v, true
			}
		}
	}
	return nil, false
}

func StringifyMyAny(myAny any) string {
	switch v := myAny.(type) {
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case map[string]any:
		b, err := json.Marshal(v)
		if err != nil {
			log.Printf("error marshalling map: %v", err)
			return ""
		}
		return string(b)
	default:
		return fmt.Sprintf("%v", myAny)
	}
}

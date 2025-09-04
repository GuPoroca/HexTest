package typeDefines

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/stretchr/testify/assert"
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
	fmt.Printf("\t%s %s %s\n", stringifyMyAny(responseVal), check.Operand, stringifyMyAny(expectedVal))
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
	case "isNull":
		passed = assert.Empty(t, responseVal)
	case "notNull":
		passed = assert.NotEmpty(t, responseVal)
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
		if passed {
			fmt.Printf("PASSED\n")
		} else {
			fmt.Printf("FAILED\n")
		}
		return passed
	}

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

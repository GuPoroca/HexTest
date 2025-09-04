package typeDefines

import (
	"fmt"
	"github.com/stretchr/testify/assert"
)

type Check struct {
	Operand  string   `json:"Operand"`
	Expected []string `json:"Expected"`
	Passed   []bool
}

type MockT struct {
	errors []string
}

func (m MockT) Errorf(format string, args ...any) {
	errorMsg := fmt.Sprintf(format, args...)
	m.errors = append(m.errors, errorMsg)
}

func (check *Check) MakeAllChecks(responseVal any) bool {
	result := 0
	all_passed := true
	for i := range check.Expected {
		if check.MakeCheck(responseVal, i) {
			result++
		} else {
			all_passed = false
		}
	}
	fmt.Printf("Comparisons passed: %v/%v\n", result, len(check.Passed))
	return all_passed
}

func (check *Check) MakeCheck(responseVal any, i int) bool {
	t := &MockT{}

	switch check.Operand {
	case "==":
		check.Passed = append(check.Passed, assert.Equal(t, responseVal, check.Expected[i]))
	case "!=":
		check.Passed = append(check.Passed, assert.NotEqual(t, responseVal, check.Expected[i]))
	case ">=":
		check.Passed = append(check.Passed, assert.GreaterOrEqual(t, responseVal, check.Expected[i]))
	case "<=":
		check.Passed = append(check.Passed, assert.LessOrEqual(t, responseVal, check.Expected[i]))
	case "isNull":
		check.Passed = append(check.Passed, assert.Empty(t, responseVal))
	case "notNull":
		check.Passed = append(check.Passed, assert.NotEmpty(t, responseVal))
	case "containsKey":
		check.Passed = append(check.Passed, assert.Contains(t, responseVal, check.Expected[i]))
	case "containsKey -R":
		_, result := containsKeyRecursevely(responseVal, check.Expected[i])
		check.Passed = append(check.Passed, result)
	}
	return check.Passed[i]
}

func containsKeyRecursevely(responseVal any, targetVal string) (any, bool) {
	switch v := responseVal.(type) {
	case map[string]any:
		for key, val := range v {
			if key == targetVal {
				return v, true
			}
			if _, ok := containsKeyRecursevely(val, targetVal); ok {
				return v, true
			}
		}
	case []any:
		for _, item := range v {
			if _, ok := containsKeyRecursevely(item, targetVal); ok {
				return v, true
			}
		}
	}
	return nil, false
}

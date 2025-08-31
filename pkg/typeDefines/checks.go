package typeDefines

import (
	"fmt"
	"github.com/stretchr/testify/assert"
)

type Check struct {
	Operand string
	Value   any
	Result  bool
}

type MockT struct {
	errors []string
}

func (m MockT) Errorf(format string, args ...any) {
	errorMsg := fmt.Sprintf(format, args...)
	m.errors = append(m.errors, errorMsg)
}

func (check *Check) MakeCheck(responseVal any) bool {
	t := &MockT{}

	switch check.Operand {
	case "==":
		check.Result = assert.Equal(t, responseVal, check.Value)
	case "!=":
		check.Result = assert.NotEqual(t, responseVal, check.Value)
	case ">=":
		check.Result = assert.GreaterOrEqual(t, responseVal, check.Value)
	case "<=":
		check.Result = assert.LessOrEqual(t, responseVal, check.Value)
	case "isNull":
		check.Result = assert.Empty(t, responseVal)
	case "notNull":
		check.Result = assert.NotEmpty(t, responseVal)
	case "containsKey":
		check.Result = assert.Contains(t, responseVal, check.Value)
	case "containsKey -R":
		_, check.Result = containsKeyRecursevely(responseVal, check.Value.(string))

	}

	return check.Result
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

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
	case "notNull":
		check.Result = assert.NotEmpty(t, responseVal)
	case "contains":
		check.Result = assert.Contains(t, responseVal, check.Value)
	}
	return check.Result
}

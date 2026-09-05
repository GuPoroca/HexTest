package typeDefines

import "testing"

// TestMakeCheckWithExpected covers the operators that compare a response value
// against one or more expected values.
func TestMakeCheckWithExpected(t *testing.T) {
	cases := []struct {
		name       string
		operand    string
		response   any
		expected   any
		wantPassed int // 1 = passed, 0 = failed
	}{
		{"equal numbers", "==", float64(200), float64(200), 1},
		{"equal numbers from string response", "==", "200", float64(200), 1},
		{"not equal numbers", "==", float64(200), float64(404), 0},
		{"greater or equal, equal", ">=", float64(900), float64(900), 1},
		{"greater or equal, below", ">=", float64(100), float64(900), 0},
		{"less or equal", "<=", float64(120), float64(1000), 1},
		{"strictly greater", ">", float64(5), float64(1), 1},
		{"strictly less fails when equal", "<", float64(5), float64(5), 0},
		{"contains key present", "containsKey", map[string]any{"message": "hi"}, "message", 1},
		{"contains key absent", "containsKey", map[string]any{"message": "hi"}, "missing", 0},
		{"contains substring", "containsSubstring", "GET called", "called", 1},
		{"match regex", "matchRegex", "user123@example.com", `^[a-z0-9]+@[a-z]+\.[a-z]+$`, 1},
		{"not match regex", "notMatchRegex", "plain text", `^\d+$`, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Check{Operand: tc.operand, Expected: []any{tc.expected}}
			got, err := c.MakeCheck(tc.response, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantPassed {
				t.Fatalf("operand %q: got %d, want %d", tc.operand, got, tc.wantPassed)
			}
		})
	}
}

// TestMakeCheckWithoutExpected covers the presence/emptiness operators that take
// no expected value.
func TestMakeCheckWithoutExpected(t *testing.T) {
	cases := []struct {
		name       string
		operand    string
		response   any
		wantPassed int
	}{
		{"notNull on value", "notNull", "something", 1},
		{"notNull on nil", "notNull", nil, 0},
		{"isNull on nil", "isNull", nil, 1},
		{"isEmpty on empty string", "isEmpty", "", 1},
		{"notEmpty on populated map", "notEmpty", map[string]any{"a": 1}, 1},
		{"notEmpty on empty slice", "notEmpty", []any{}, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Check{Operand: tc.operand}
			got, err := c.MakeCheckWithoutExpected(tc.response)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantPassed {
				t.Fatalf("operand %q: got %d, want %d", tc.operand, got, tc.wantPassed)
			}
		})
	}
}

func TestMakeCheckUnknownOperand(t *testing.T) {
	c := &Check{Operand: "definitelyNotAnOperand", Expected: []any{"x"}}
	got, err := c.MakeCheck("x", 0)
	if got != -1 {
		t.Fatalf("value: got %d, want -1", got)
	}
	if _, ok := err.(OperandNotFound); !ok {
		t.Fatalf("error: got %T (%v), want OperandNotFound", err, err)
	}
}

func TestMakeAllChecksCountsPasses(t *testing.T) {
	body := map[string]any{"message": "GET called", "test": true}
	c := &Check{Operand: "containsKey", Expected: []any{"message", "test", "absent"}}

	passed := c.MakeAllChecks(body)
	if passed != 2 {
		t.Fatalf("passed count: got %d, want 2", passed)
	}
	if c.Total_num != 3 {
		t.Fatalf("total count: got %d, want 3", c.Total_num)
	}
}

func TestContainsKeyRecursively(t *testing.T) {
	nested := map[string]any{
		"outer": map[string]any{
			"inner": map[string]any{"needle": 1},
		},
	}
	if _, ok := ContainsKeyRecursevely(nested, "needle"); !ok {
		t.Fatal("expected to find deeply nested key")
	}
	if _, ok := ContainsKeyRecursevely(nested, "haystack"); ok {
		t.Fatal("did not expect to find a missing key")
	}
}

package typeDefines

import ()

type Assert struct {
	Field              string  `json:"Field"`
	Checks             []Check `json:"Checks"`
	FieldResponseValue any
}

func (assert *Assert) MakeAssertions(fieldValue any) {
	assert.FieldResponseValue = fieldValue
	for i := range assert.Checks {
		if assert.Field == "JSON Schema Validation" {
			assert.Checks[i].JsonSchema(assert.FieldResponseValue)
		} else {
			assert.Checks[i].MakeAllChecks(fieldValue)
		}
	}
}

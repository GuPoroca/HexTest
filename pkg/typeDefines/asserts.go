package typeDefines

import ()

type Assert struct {
	Field                   string  `json:"Field"`
	Checks                  []Check `json:"Checks"`
	Results                 []bool
	FieldResponseValue      any
	Passed_Comparissons_num int
	Total_Comparissons_num  int
}

func (assert *Assert) MakeAssertions(fieldValue any) int {
	assert.FieldResponseValue = fieldValue
	for i := range assert.Checks {
		assert.Passed_Comparissons_num += assert.Checks[i].MakeAllChecks(fieldValue)
		assert.Total_Comparissons_num += assert.Checks[i].Total_num
		if (assert.Checks[i].Passed_num - len(assert.Checks[i].Expected)) == 0 {
			assert.Results = append(assert.Results, true)
		} else {
			assert.Results = append(assert.Results, false)
		}
	}
	return assert.Passed_Comparissons_num
}

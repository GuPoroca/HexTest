package typeDefines

import "fmt"

type OperandNotFound struct {
	Operand string
}

func (e OperandNotFound) Error() string {
	return fmt.Sprintf("Operand %s not found", e.Operand)
}

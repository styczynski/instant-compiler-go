package ir

import "fmt"

type IRType string

const (
	IR_INT32   IRType = "Int32"
	IR_INT16   IRType = "Int16"
	IR_INT8    IRType = "Int8"
	IR_BIT     IRType = "Bit"
	IR_FN      IRType = "FunctionPtr"
	IR_UNKNOWN IRType = "Unknown"
)

func GetIRTypeSize(varType IRType) int {
	if varType == IR_INT32 {
		return 32
	} else if varType == IR_INT16 {
		return 16
	} else if varType == IR_INT8 {
		return 8
	} else if varType == IR_BIT {
		return 8
	} else if varType == IR_FN {
		return 32
	}
	panic(fmt.Sprintf("Invalid type was given (cannot calculate size): %s", varType))
}

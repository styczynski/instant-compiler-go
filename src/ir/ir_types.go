package ir

type IRType string

const (
	IR_INT32   IRType = "Int32"
	IR_INT16   IRType = "Int16"
	IR_INT8    IRType = "Int8"
	IR_BIT     IRType = "Bit"
	IR_FN      IRType = "FunctionPtr"
	IR_UNKNOWN IRType = "Unknown"
)

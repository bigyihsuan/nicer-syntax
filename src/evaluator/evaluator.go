package evaluator

type Evaluable interface {
	Evaluate() *NicerValue
}

type NicerValue struct {
	Type  NicerType
	Value interface{}
}

type NicerType string

// built-in types
const (
	NT_number  NicerType = "number"
	NT_boolean NicerType = "boolean"
	NT_string  NicerType = "string"
)

// maps a typename to whether it's defined
// lists and structs will need to be run-time defined
var NicerTypeList = map[NicerType]bool{
	NT_number:  true,
	NT_boolean: true,
	NT_string:  true,
}

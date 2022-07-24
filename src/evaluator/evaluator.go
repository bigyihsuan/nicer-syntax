package evaluator

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
)

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

type RuntimeError struct {
	Reason       string
	VariableName string
	Node         interface{}
}

var COLOR_ERROR = color.New(color.FgHiRed).Add(color.Underline).Add(color.Bold).Sprintf
var COLOR_KEYWORD = color.New(color.FgYellow).Sprintf
var COLOR_TOKEN = color.New(color.FgCyan).Sprintf

// for interface error.Error()
func (re *RuntimeError) Error() string {
	return fmt.Sprintf("%v %v (%v) at `%v`",
		COLOR_ERROR("RUNTIME ERROR:"),
		COLOR_KEYWORD(re.Reason),
		COLOR_TOKEN(re.VariableName),
		COLOR_TOKEN("%s", spew.Sprintf("%#v", re.Node)))
}

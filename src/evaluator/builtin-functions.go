package evaluator

import "fmt"

type NicerFunction func(parameters []NicerValue) *NicerValue

var BuiltInFunctions = map[string]NicerFunction{
	"Print":     Print,
	"PrintLine": PrintLine,
}

// TODO: Proper Expr
func Print(parameters []NicerValue) *NicerValue {
	for _, param := range parameters {
		v := param.Value
		fmt.Print(v)
	}
	return nil
}

// TODO: Proper Expr
func PrintLine(parameters []NicerValue) *NicerValue {
	for _, param := range parameters {
		v := param.Value
		fmt.Println(v)
	}
	return nil
}

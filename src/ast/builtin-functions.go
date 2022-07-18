package ast

import (
	"fmt"
	"nicer-syntax/evaluator"
)

type NicerFunction func(parameters []evaluator.Evaluable) *evaluator.NicerValue

var BuiltInFunctions = map[string]NicerFunction{
	"Print":     Print,
	"PrintLine": PrintLine,
}

// TODO: Proper Expr
func Print(parameters []evaluator.Evaluable) *evaluator.NicerValue {
	for _, param := range parameters {
		v := param.Evaluate().Value
		fmt.Print(v)
	}
	return nil
}

// TODO: Proper Expr
func PrintLine(parameters []evaluator.Evaluable) *evaluator.NicerValue {
	for _, param := range parameters {
		v := param.Evaluate().Value
		fmt.Println(v)
	}
	return nil
}

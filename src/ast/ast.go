package ast

import (
	"math/big"
	"nicer-syntax/evaluator"
	"nicer-syntax/lexer"
	"strconv"
)

type Primitive interface {
}

type NumberLiteral struct {
	*lexer.TokItem
}

func (lit NumberLiteral) Evaluate() *evaluator.NicerValue {
	var f float64
	switch l := lit.TokValue; l.(type) {
	case *big.Int:
		f, _ = new(big.Float).SetInt(l.(*big.Int)).Float64()
	case *big.Float:
		f, _ = l.(*big.Float).Float64()
	}
	return &evaluator.NicerValue{Value: f}
}

type BooleanLiteral struct {
	*lexer.TokItem
}

func (lit BooleanLiteral) Evaluate() *evaluator.NicerValue {
	b, err := strconv.ParseBool(lit.TokValue.(string))
	if err != nil {
		panic(err)
	} else {
		return &evaluator.NicerValue{Value: b}
	}
}

type StringLiteral struct {
	*lexer.TokItem
}

func (lit StringLiteral) Evaluate() *evaluator.NicerValue {
	return &evaluator.NicerValue{Value: lit.TokValue.(string)}
}

type FunctionCall struct {
	FunctionName lexer.TokItem
	// TODO: Proper Expr
	Parameters []evaluator.Evaluable
}

func (fc FunctionCall) Evaluate() *evaluator.NicerValue {
	funcname := fc.FunctionName.TokValue.(string)
	if fun, ok := BuiltInFunctions[funcname]; ok {
		return fun(fc.Parameters)
	} else {
		// TODO: User-defined functions
		return nil
	}
}

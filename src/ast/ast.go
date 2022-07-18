package ast

import (
	"fmt"
	"math/big"
	"nicer-syntax/evaluator"
	"nicer-syntax/lexer"
	"strconv"
)

type Primitive interface {
	evaluator.Evaluable
}

type NumberLiteral struct {
	*lexer.TokItem
}

func (lit NumberLiteral) Evaluate() float64 {
	var f float64
	switch l := lit.TokValue; l.(type) {
	case *big.Int:
		f, _ = new(big.Float).SetInt(l.(*big.Int)).Float64()
	case *big.Float:
		f, _ = l.(*big.Float).Float64()
	}
	return f
}

type BooleanLiteral struct {
	*lexer.TokItem
}

func (lit BooleanLiteral) Evaluate() bool {
	b, err := strconv.ParseBool(lit.TokValue.(string))
	if err != nil {
		panic(err)
	} else {
		return b
	}
}

type StringLiteral struct {
	*lexer.TokItem
}

func (lit StringLiteral) Evaluate() string {
	return lit.TokValue.(string)
}

type FunctionCall struct {
	FunctionName lexer.TokItem
	Parameters   []*lexer.TokItem
}

func (fc FunctionCall) Evaluate() evaluator.NicerValue {
	fmt.Println(fc)
	return evaluator.NicerValue{Value: nil}
}

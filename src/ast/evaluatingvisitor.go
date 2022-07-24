package ast

import (
	"nicer-syntax/evaluator"
)

type ValueStack []*evaluator.NicerValue

func (vs *ValueStack) Push(val *evaluator.NicerValue) {
	*vs = append(*vs, val)
}

func (vs *ValueStack) Pop() *evaluator.NicerValue {
	val := (*vs)[len(*vs)-1]
	*vs = (*vs)[:len(*vs)-1]
	return val
}

type EvaluatingVisitor struct {
	DefaultVisitor
	*ValueStack // stack of values
	IdentValue  map[string]*evaluator.NicerValue
}

func NewEvaluatingVisitor() *EvaluatingVisitor {
	ev := new(EvaluatingVisitor)
	ev.ValueStack = new(ValueStack)
	ev.IdentValue = make(map[string]*evaluator.NicerValue)
	return ev
}

func (v *EvaluatingVisitor) Visit(vis Visitable) {
	switch vis := vis.(type) {
	case *NumberLiteral:
		v.VisitNumberLiteral(v, vis)
	case *BooleanLiteral:
		v.VisitBooleanLiteral(v, vis)
	case *StringLiteral:
		v.VisitStringLiteral(v, vis)
	case *Identifier:
		v.VisitIdentifier(v, vis)
	}
}
func (v *EvaluatingVisitor) VisitNumberLiteral(_ Visitor, nl *NumberLiteral) {
	v.ValueStack.Push(nl.Evaluate())
}
func (v *EvaluatingVisitor) VisitBooleanLiteral(_ Visitor, bl *BooleanLiteral) {
	v.ValueStack.Push(bl.Evaluate())
}
func (v *EvaluatingVisitor) VisitStringLiteral(_ Visitor, sl *StringLiteral) {
	v.ValueStack.Push(sl.Evaluate())
}
func (v *EvaluatingVisitor) VisitIdentifier(_ Visitor, id *Identifier) {
	val, ok := v.IdentValue[id.Name]
	if ok {
		v.ValueStack.Push(val)
	} else {
		// TODO: `nothing` values
		v.ValueStack.Push(&evaluator.NicerValue{
			Type:  evaluator.NT_number,
			Value: 0.0,
		})
	}
}
func (v *EvaluatingVisitor) VisitFunctionCall(_ Visitor, fc *FunctionCall) {
	v.Visit(fc.FuncParams) // evaluate arguments first
	function, ok := evaluator.BuiltInFunctions[fc.FuncName.Name]
	if ok {
		val := v.ValueStack.Pop()
		function([]evaluator.NicerValue{*val})
	}
}
func (v *EvaluatingVisitor) VisitConstDecl(_ Visitor, cd *ConstDecl) {
	// assign to the variable map the name and value
	v.Visit(cd.Value)
	val := v.ValueStack.Pop()
	v.IdentValue[cd.ConstName.Name] = val
}

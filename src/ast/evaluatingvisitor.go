package ast

import (
	"fmt"
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
	default:
		v.ValueStack.Push(nil)
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
		// v.ValueStack.Push(&evaluator.NicerValue{
		// 	Type:  evaluator.NT_number,
		// 	Value: 0.0,
		// })
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

func (v *EvaluatingVisitor) VisitVarDecl(_ Visitor, cd *VarDecl) {
	// assign to the variable map the name and value
	v.Visit(cd.Value)
	val := v.ValueStack.Pop()
	v.IdentValue[cd.VarName.Name] = val
}

func (v *EvaluatingVisitor) VisitProgram(_ Visitor, p *Program) {
	for _, stmt := range p.Statements {
		v.VisitStatement(v, stmt)
	}
}

func (v *EvaluatingVisitor) VisitStatement(_ Visitor, s Statement) {
	switch s := s.(type) {
	case *VarAssignment:
		v.VisitVarAssignment(v, s)
	case Declaration:
		v.VisitDeclaration(v, s)
	}
}

func (v *EvaluatingVisitor) VisitDeclaration(_ Visitor, d Declaration) {
	switch d := d.(type) {
	case *VarDecl:
		v.VisitVarDecl(v, d)
	case *ConstDecl:
		v.VisitConstDecl(v, d)
	}
}

func (v *EvaluatingVisitor) VisitVarAssignment(_ Visitor, va *VarAssignment) {
	// assign the name to the new value
	v.Visit(va.Value)
	// check for the ident to exist; if not, exit
	val := v.ValueStack.Pop()
	if _, ok := v.IdentValue[va.Name.Name]; ok && val != nil {
		v.IdentValue[va.Name.Name] = val
	} else {
		err := evaluator.RuntimeError{
			Reason:       "Trying to assign to variable that does not exist",
			VariableName: va.Name.Name,
			Node:         va,
		}
		fmt.Println(err.Error())
		return
	}
}

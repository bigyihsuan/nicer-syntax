package ast

import (
	"fmt"
	"strings"
)

type stringStack []string

func (vs *stringStack) Push(val string) {
	*vs = append(*vs, val)
}

func (vs *stringStack) Pop() string {
	val := (*vs)[len(*vs)-1]
	*vs = (*vs)[:len(*vs)-1]
	return val
}

type StringVisitor struct {
	DefaultVisitor
	builder strings.Builder
	strings stringStack
}

func (v StringVisitor) String() string {
	v.builder.Reset()
	for _, s := range v.strings {
		v.builder.WriteString(s)
	}
	return v.builder.String()
}

func (v *StringVisitor) Visit(vis Visitable) {
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
		v.strings.Push("nothing")
	}
}
func (v *StringVisitor) VisitNumberLiteral(_ Visitor, nl *NumberLiteral) {
	v.strings.Push(fmt.Sprintf("%v", nl.Value))
}
func (v *StringVisitor) VisitBooleanLiteral(_ Visitor, bl *BooleanLiteral) {
	v.strings.Push(fmt.Sprintf("%v", bl.Value))
}
func (v *StringVisitor) VisitStringLiteral(_ Visitor, sl *StringLiteral) {
	v.strings.Push(fmt.Sprintf("\"%v\"", sl.Value))
}
func (v *StringVisitor) VisitIdentifier(_ Visitor, id *Identifier) {
	v.strings.Push(fmt.Sprintf("%v", id.Name))
}
func (v *StringVisitor) VisitFunctionCall(_ Visitor, fc *FunctionCall) {
	v.builder.Reset()
	v.VisitIdentifier(v, fc.FuncName)
	ident := v.strings.Pop()
	v.Visit(fc.FuncParams)
	params := v.strings.Pop()
	v.builder.WriteString(fmt.Sprintf("FunctionCall(%s %s)", ident, params))
	v.strings.Push(v.builder.String())
}

func (v *StringVisitor) VisitVarDecl(_ Visitor, vd *VarDecl) {
	v.builder.Reset()
	v.VisitIdentifier(v, vd.VarName)
	ident := v.strings.Pop()
	v.Visit(vd.TypeName)
	typeName := v.strings.Pop()
	v.Visit(vd.Value)
	value := v.strings.Pop()
	v.builder.WriteString(fmt.Sprintf("VarDecl(%s %s %s)", ident, typeName, value))
	v.strings = append(v.strings, v.builder.String())
}
func (v *StringVisitor) VisitConstDecl(_ Visitor, cd *ConstDecl) {
	v.builder.Reset()
	v.VisitIdentifier(v, cd.ConstName)
	ident := v.strings.Pop()
	v.Visit(cd.TypeName)
	typeName := v.strings.Pop()
	v.Visit(cd.Value)
	value := v.strings.Pop()
	v.builder.WriteString(fmt.Sprintf("ConstDecl(%s %s %s)", ident, typeName, value))
	v.strings = append(v.strings, v.builder.String())
}

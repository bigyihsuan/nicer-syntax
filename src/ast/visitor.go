package ast

import (
	"fmt"
	"strings"
)

type Visitable interface {
	Accept(Visitor)
}

type Visitor interface {
	Visit(vis Visitable)
	VisitNumberLiteral(v Visitor, nl *NumberLiteral)
	VisitBooleanLiteral(v Visitor, bl *BooleanLiteral)
	VisitStringLiteral(v Visitor, sl *StringLiteral)
	VisitIdentifier(v Visitor, id *Identifier)
	VisitFunctionCall(v Visitor, fc *FunctionCall)
}

type DefaultVisitor struct{}

func (*DefaultVisitor) Visit(vis Visitable)                               {}
func (*DefaultVisitor) VisitNumberLiteral(v Visitor, nl *NumberLiteral)   {}
func (*DefaultVisitor) VisitBooleanLiteral(v Visitor, bl *BooleanLiteral) {}
func (*DefaultVisitor) VisitStringLiteral(v Visitor, l *StringLiteral)    {}
func (*DefaultVisitor) VisitIdentifier(v Visitor, id *Identifier)         {}
func (*DefaultVisitor) VisitFunctionCall(v Visitor, fc *FunctionCall)     {}

type StringVisitor struct {
	DefaultVisitor
	builder strings.Builder
	strings []string
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
	}
}

func (v *StringVisitor) VisitNumberLiteral(_ Visitor, nl *NumberLiteral) {
	v.strings = append(v.strings, fmt.Sprintf("%v", nl.Value))
	// fmt.Println("VisitNumberLiteral", v.strings)
}
func (v *StringVisitor) VisitBooleanLiteral(_ Visitor, bl *BooleanLiteral) {
	v.strings = append(v.strings, fmt.Sprintf("%v", bl.Value))
	// fmt.Println("VisitBooleanLiteral", v.strings)
}
func (v *StringVisitor) VisitStringLiteral(_ Visitor, sl *StringLiteral) {
	v.strings = append(v.strings, fmt.Sprintf("\"%v\"", sl.Value))
	// fmt.Println("VisitStringLiteral", v.strings)
}
func (v *StringVisitor) VisitIdentifier(_ Visitor, id *Identifier) {
	v.strings = append(v.strings, fmt.Sprintf("%v", id.Name))
	// fmt.Println("VisitIdentifier", v.strings)
}
func (v *StringVisitor) VisitFunctionCall(_ Visitor, fc *FunctionCall) {
	v.builder.Reset()
	v.VisitIdentifier(v, fc.FuncName)
	// fmt.Println("VisitFunctionCall", v.strings)
	ident := v.strings[len(v.strings)-1]
	v.strings = append([]string{}, v.strings[:len(v.strings)-1]...)
	v.Visit(fc.FuncParams)
	// fmt.Println("VisitFunctionCall", v.strings)
	params := v.strings[len(v.strings)-1]
	v.strings = append([]string{}, v.strings[:len(v.strings)-1]...)
	v.builder.WriteString("FunctionCall(")
	v.builder.WriteString(ident + " ")
	v.builder.WriteString(params + ")")
	v.strings = append(v.strings, v.builder.String())
}

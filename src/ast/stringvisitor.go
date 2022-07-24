package ast

import (
	"fmt"
	"strings"
)

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
}
func (v *StringVisitor) VisitBooleanLiteral(_ Visitor, bl *BooleanLiteral) {
	v.strings = append(v.strings, fmt.Sprintf("%v", bl.Value))
}
func (v *StringVisitor) VisitStringLiteral(_ Visitor, sl *StringLiteral) {
	v.strings = append(v.strings, fmt.Sprintf("\"%v\"", sl.Value))
}
func (v *StringVisitor) VisitIdentifier(_ Visitor, id *Identifier) {
	v.strings = append(v.strings, fmt.Sprintf("%v", id.Name))
}
func (v *StringVisitor) VisitFunctionCall(_ Visitor, fc *FunctionCall) {
	v.builder.Reset()
	v.VisitIdentifier(v, fc.FuncName)
	ident := v.strings[len(v.strings)-1]
	v.strings = v.strings[:len(v.strings)-1]
	v.Visit(fc.FuncParams)
	params := v.strings[len(v.strings)-1]
	v.strings = v.strings[:len(v.strings)-1]
	v.builder.WriteString(fmt.Sprintf("FunctionCall(%s %s)", ident, params))
	v.strings = append(v.strings, v.builder.String())
}

func (v *StringVisitor) VisitVarDecl(_ Visitor, vd *VarDecl) {
	v.builder.Reset()
	v.VisitIdentifier(v, vd.VarName)
	ident := v.strings[len(v.strings)-1]
	v.strings = v.strings[:len(v.strings)-1]
	v.Visit(vd.Value)
	params := v.strings[len(v.strings)-1]
	v.strings = v.strings[:len(v.strings)-1]
	v.builder.WriteString(fmt.Sprintf("ConstDecl(%s %s)", ident, params))
	v.strings = append(v.strings, v.builder.String())
}
func (v *StringVisitor) VisitConstDecl(_ Visitor, cd *ConstDecl) {
	v.builder.Reset()
	v.VisitIdentifier(v, cd.ConstName)
	ident := v.strings[len(v.strings)-1]
	v.strings = v.strings[:len(v.strings)-1]
	v.Visit(cd.Value)
	params := v.strings[len(v.strings)-1]
	v.strings = v.strings[:len(v.strings)-1]
	v.builder.WriteString(fmt.Sprintf("ConstDecl(%s %s)", ident, params))
	v.strings = append(v.strings, v.builder.String())
}

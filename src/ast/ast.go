package ast

import (
	"math/big"
	"nicer-syntax/evaluator"
	"nicer-syntax/lexer"
	"strconv"
)

/* THINGS THAT RETURN VALUES
- Literals (number, boolean, string, list, map, struct)
- Ranges (List of Numbers)
- Expressions
- Functions
*/

type HasValue interface{}

type NumberLiteral struct {
	HasValue
	Value float64
}

func NewNumberLiteral(tok *lexer.TokItem) *NumberLiteral {
	nl := NumberLiteral{}
	var f float64
	switch l := tok.TokValue; l.(type) {
	case *big.Int:
		f, _ = new(big.Float).SetInt(l.(*big.Int)).Float64()
	case *big.Float:
		f, _ = l.(*big.Float).Float64()
	}
	nl.Value = f
	return &nl
}

// ast.Visitable
func (nl NumberLiteral) Accept(v Visitor) {
	v.VisitNumberLiteral(v, &nl)
}

// evaluator.Evaluable
func (nl NumberLiteral) Evaluate() *evaluator.NicerValue {
	return &evaluator.NicerValue{
		Type:  evaluator.NT_number,
		Value: nl.Value,
	}
}

type BooleanLiteral struct {
	HasValue
	Value bool
}

func NewBooleanLiteral(tok *lexer.TokItem) *BooleanLiteral {
	b, err := strconv.ParseBool(tok.TokValue.(string))
	if err != nil {
		b = false
	}
	return &BooleanLiteral{Value: b}
}

// ast.Visitable
func (bl BooleanLiteral) Accept(v Visitor) {
	v.VisitBooleanLiteral(v, &bl)
}

// evaluator.Evaluable
func (bl BooleanLiteral) Evaluate() *evaluator.NicerValue {
	return &evaluator.NicerValue{
		Type:  evaluator.NT_boolean,
		Value: bl.Value,
	}
}

type StringLiteral struct {
	HasValue
	Value string
}

func NewStringLiteral(tok *lexer.TokItem) *StringLiteral {
	return &StringLiteral{Value: tok.TokValue.(string)}
}

// ast.Visitable
func (sl StringLiteral) Accept(v Visitor) {
	v.VisitStringLiteral(v, &sl)
}

// evaluator.Evaluable
func (sl StringLiteral) Evaluate() *evaluator.NicerValue {
	return &evaluator.NicerValue{
		Type:  evaluator.NT_string,
		Value: sl.Value,
	}
}

type Identifier struct {
	HasValue
	Name string
}

func NewIdentifier(tok *lexer.TokItem) *Identifier {
	ident := &Identifier{}
	ident.Name = tok.TokValue.(string)
	return ident
}

// ast.Visitable
func (id Identifier) Accept(v Visitor) {
	v.VisitIdentifier(v, &id)
}

type FunctionCall struct {
	FuncName   *Identifier
	FuncParams Visitable
}

// TODO: Actual comma-separated list of expr
func NewFunctionCall(name *lexer.TokItem, parameters Visitable) *FunctionCall {
	return &FunctionCall{
		FuncName:   NewIdentifier(name),
		FuncParams: parameters,
	}
}

// ast.Visitable
func (fc FunctionCall) Accept(v Visitor) {
	v.VisitFunctionCall(v, &fc)
}

// type CommaSeparatedList struct {
// }

type Statement interface{}

type Declaration interface {
	Statement
}

type VarAssignment struct {
	Statement
	Name  *Identifier
	Value Visitable
}

func NewVarAssignment(varName *Identifier, val Visitable) *VarAssignment {
	varass := new(VarAssignment)
	varass.Name = varName
	varass.Value = val
	return varass
}

// ast.Visitable
func (va VarAssignment) Accept(v Visitor) {
	v.VisitVarAssignment(v, &va)
}

type VarDecl struct {
	Declaration
	VarName  *Identifier
	TypeName *Identifier
	Value    Visitable // TODO: Expr
}

// ast.Visitable
func (vd VarDecl) Accept(v Visitor) {
	v.VisitVarDecl(v, &vd)
}
func NewVarDecl(name, typeName *Identifier, value Visitable) *VarDecl {
	vardecl := new(VarDecl)
	vardecl.VarName = name
	vardecl.TypeName = typeName
	vardecl.Value = value
	return vardecl
}

type ConstDecl struct {
	Declaration
	ConstName *Identifier
	TypeName  *Identifier
	Value     Visitable // TODO: Expr
}

func NewConstDecl(name, typeName *Identifier, value Visitable) *ConstDecl {
	return &ConstDecl{
		ConstName: name,
		TypeName:  typeName,
		Value:     value,
	}
}

// ast.Visitable
func (cd ConstDecl) Accept(v Visitor) {
	v.VisitConstDecl(v, &cd)
}

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	program := new(Program)
	return program
}

// ast.Visitable
func (p Program) Accept(v Visitor) {
	v.VisitProgram(v, &p)
}

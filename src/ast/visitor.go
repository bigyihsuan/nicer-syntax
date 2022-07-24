package ast

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
	VisitDeclaration(v Visitor, d Declaration)
	VisitVarDecl(v Visitor, vd *VarDecl)
	VisitConstDecl(v Visitor, cd *ConstDecl)
	VisitProgram(v Visitor, p *Program)
	VisitStatement(v Visitor, s Statement)
	VisitVarAssignment(v Visitor, va *VarAssignment)
}

type DefaultVisitor struct{}

func (v *DefaultVisitor) Visit(vis Visitable)                             {}
func (*DefaultVisitor) VisitNumberLiteral(v Visitor, nl *NumberLiteral)   {}
func (*DefaultVisitor) VisitBooleanLiteral(v Visitor, bl *BooleanLiteral) {}
func (*DefaultVisitor) VisitStringLiteral(v Visitor, sl *StringLiteral)   {}
func (*DefaultVisitor) VisitIdentifier(v Visitor, id *Identifier)         {}
func (*DefaultVisitor) VisitFunctionCall(v Visitor, fc *FunctionCall)     {}
func (*DefaultVisitor) VisitDeclaration(v Visitor, d Declaration)         {}
func (*DefaultVisitor) VisitVarDecl(v Visitor, vd *VarDecl)               {}
func (*DefaultVisitor) VisitConstDecl(v Visitor, cd *ConstDecl)           {}
func (*DefaultVisitor) VisitProgram(v Visitor, p *Program)                {}
func (*DefaultVisitor) VisitStatement(v Visitor, s Statement)             {}
func (*DefaultVisitor) VisitVarAssignment(v Visitor, va *VarAssignment)   {}

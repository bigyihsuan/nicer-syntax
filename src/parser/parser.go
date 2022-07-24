package parser

import (
	"fmt"
	"nicer-syntax/ast"
	"nicer-syntax/lexer"

	"github.com/db47h/lex"

	"github.com/fatih/color"
)

var COLOR_ERROR = color.New(color.FgHiRed).Add(color.Underline).Add(color.Bold).Sprintf
var COLOR_KEYWORD = color.New(color.FgYellow).Sprintf
var COLOR_TOKEN = color.New(color.FgCyan).Sprintf
var COLOR_RULE = color.New(color.FgMagenta).Sprintf

type ParseError struct {
	Reason   string
	Token    lexer.TokItem
	LastRule string
}

func NewParseError(reason string, token lexer.TokItem, lastRule string) *ParseError {
	err := new(ParseError)
	err.Reason = reason
	err.Token = token
	err.LastRule = lastRule
	return err
}

func (pe *ParseError) addRule(rule string) *ParseError {
	pe.LastRule += fmt.Sprintf(" %s", rule)
	return pe
}

// for interface error.Error()
func (pe *ParseError) Error() string {
	return fmt.Sprintf("%v %v because of token `%v` within rule trace `%v`", COLOR_ERROR("PARSE ERROR:"), COLOR_KEYWORD(pe.Reason), COLOR_TOKEN("%v", pe.Token), COLOR_RULE(pe.LastRule))
}

type Parser struct {
	Tokens    []lexer.TokItem
	lastToken *lexer.TokItem
}

func NewParser(tokens []lexer.TokItem) Parser {
	return Parser{tokens, &lexer.TokItem{TokType: lexer.ItemEOF, TokName: "nothing", TokPosition: -1, TokValue: ""}}
}

// consume and return the next token in the token queue.
func (p *Parser) getNextToken() lexer.TokItem {
	tok := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	p.lastToken = &tok
	return tok
}

// peek at the front of the token queue.
func (p *Parser) peekToken() *lexer.TokItem {
	return &(p.Tokens[0])
}

// put the last-consumed token back onto the front of the token queue.
func (p *Parser) putBackToken() {
	p.Tokens = append([]lexer.TokItem{*p.lastToken}, p.Tokens...)
}

// consume a token and determine if it is of a desired token type.
func (p *Parser) expectToken(tokType lex.Token, lastRule string) (bool, *ParseError, *lexer.TokItem) {
	token := p.getNextToken()
	if token.TokType != tokType {
		return false, NewParseError(fmt.Sprintf("Expected token `%v`", lexer.TokenString[tokType]), token, lastRule), nil
	}
	return true, nil, &token
}

// peek at a token and determine if it is of a desired token type.
func (p *Parser) maybeToken(tokType lex.Token, lastRule string) (bool, *ParseError) {
	token := p.peekToken()
	if token.TokType != tokType {
		return false, NewParseError(fmt.Sprintf("Expected token `%v`", lexer.TokenString[tokType]), *token, lastRule)
	}
	return true, nil
}

//! Calling Conventions
//! - Consume tokens when expected
//! - Put back when hitting the end of a rule that consumes multiple
//! - When going to a nested rule, peek to check for the starting token

func (p *Parser) Parse() (bool, *ParseError, *ast.Program) {
	ok, err, prog := p.Program()
	if !ok {
		return false, err, nil
	}
	return true, nil, prog
}

func (p *Parser) Program() (bool, *ParseError, *ast.Program) {
	program := ast.NewProgram()
	for len(p.Tokens) > 0 {
		ok, err, stmt := p.Stmt()
		if !ok {
			return false, err.addRule("Program-Stmt"), nil
		}
		if ok, err, _ := p.expectToken(lexer.ItemSemicolon, "Program-Semicolon"); !ok {
			return false, err, nil
		}
		program.Statements = append(program.Statements, stmt)
	}
	return true, nil, program
}

func (p *Parser) Stmt() (bool, *ParseError, ast.Statement) {
	switch p.peekToken().TokType {
	case lexer.ItemIdent:
		return p.IdentAssignment()
	default:
		return p.IdentDeclaration()
	}
}

func (p *Parser) IdentAssignment() (bool, *ParseError, *ast.VarAssignment) {
	ok, err, name := p.Ident()
	if !ok {
		return false, err.addRule("IdentAssignment-Ident"), nil
	}
	if ok, err, _ := p.expectToken(lexer.KW_Is, "IdentAssignment-Is"); !ok {
		return false, err, nil
	}
	ok, err, val := p.Value()
	return ok, err, ast.NewVarAssignment(name, val)
}

func (p *Parser) IdentDeclaration() (bool, *ParseError, ast.Declaration) {
	switch p.peekToken().TokType {
	case lexer.KW_Constant:
		return p.ConstDecl()
	case lexer.KW_Variable:
		return p.VarDecl()
	default:
		return false, NewParseError("Invalid token", *p.peekToken(), "IdentDeclaration"), nil
	}
}

func (p *Parser) VarDecl() (bool, *ParseError, *ast.VarDecl) {
	ok, err, _ := p.expectToken(lexer.KW_Variable, "VarDecl-Variable")
	if !ok {
		return false, err, nil
	}
	ok, err, name, typeName := p.IdentType()
	if !ok {
		return false, err.addRule("VarDecl-IdentType"), nil
	}
	// optional value, ended with semicolon
	if p.peekToken().TokType == lexer.ItemSemicolon {
		return true, nil, ast.NewVarDecl(name, typeName, nil)
	}
	ok, err, val := p.Value()
	return ok, err, ast.NewVarDecl(name, typeName, val)
}

func (p *Parser) ConstDecl() (bool, *ParseError, *ast.ConstDecl) {
	if ok, err, _ := p.expectToken(lexer.KW_Constant, "ConstDecl-Constant"); !ok {
		return false, err, nil
	}
	ok, err, name, typeName := p.IdentType()
	if !ok {
		return false, err.addRule("ConstDecl-IdentType"), nil
	}
	ok, err, val := p.Value()
	if !ok {
		return false, err.addRule("ConstDecl-Value"), nil
	}
	return true, nil, ast.NewConstDecl(name, typeName, val)
}

func (p *Parser) IdentType() (bool, *ParseError, *ast.Identifier, *ast.Identifier) {
	ok, err, name := p.expectToken(lexer.ItemIdent, "IdentType-Ident")
	if !ok {
		return false, err, nil, nil
	}
	if ok, err, _ := p.expectToken(lexer.KW_Is, "IdentType-Is"); !ok {
		return false, err, nil, nil
	}
	ok, err, typeName := p.TypeName()
	if !ok {
		return false, err.addRule("IdentType-TypeName"), nil, nil
	}
	return true, nil, ast.NewIdentifier(name), typeName
}

func (p *Parser) TypeName() (bool, *ParseError, *ast.Identifier) {
	typeName := p.getNextToken()
	switch typeName.TokType {
	case lexer.TN_Number, lexer.TN_String, lexer.TN_Boolean:
		return true, nil, ast.NewIdentifier(&typeName)
	case lexer.TN_List:
		if ok, err, _ := p.expectToken(lexer.KW_Of, "TypeName-ListOf"); !ok {
			return false, err, nil
		}
		return p.TypeName()
	case lexer.TN_Map:
		if ok, err, _ := p.expectToken(lexer.KW_Of, "TypeName-MapOfKey"); !ok {
			return false, err, nil
		}
		if ok, err, _ := p.TypeName(); !ok {
			return false, err.addRule("TypeName-MapKey"), nil
		}
		if ok, err, _ := p.expectToken(lexer.KW_To, "TypeName-MapToValue"); !ok {
			return false, err, nil
		}
		return p.TypeName()
	case lexer.ItemIdent: // possibly undeclared typename
		return true, nil, ast.NewIdentifier(&typeName)
	default:
		return false, NewParseError("Expected type name", typeName, "TypeName"), nil
	}
}

func (p *Parser) Value() (bool, *ParseError, ast.Visitable) {
	switch p.peekToken().TokType {
	case lexer.LT_Number:
		ok, err, val := p.NumberLiteral()
		return ok, err, val
	case lexer.LT_Boolean:
		ok, err, val := p.BooleanLiteral()
		return ok, err, val
	case lexer.LT_String:
		ok, err, val := p.StringLiteral()
		return ok, err, val
	case lexer.KW_Containing:
		// TODO: Evaluating list literals
		// return p.ListLiteral()
		return true, nil, nil
	default:
		return false, NewParseError("Expected value", *p.peekToken(), "Value"), nil
	}
}

func (p *Parser) PrimitiveLiteral() (bool, *ParseError, ast.Visitable) {
	switch p.peekToken().TokType {
	case lexer.LT_Number:
		return p.NumberLiteral()
	case lexer.LT_Boolean:
		return p.BooleanLiteral()
	case lexer.LT_String:
		return p.StringLiteral()
	default:
		return false, NewParseError("Expected primitive literal", *p.peekToken(), "PrimitiveLiteral"), nil
	}
}

func (p *Parser) NumberLiteral() (bool, *ParseError, *ast.NumberLiteral) {
	if ok, err, token := p.expectToken(lexer.LT_Number, "NumberLiteral"); !ok {
		return false, err, nil
	} else {
		return true, nil, ast.NewNumberLiteral(token)
	}
}

func (p *Parser) BooleanLiteral() (bool, *ParseError, *ast.BooleanLiteral) {
	if ok, err, token := p.expectToken(lexer.LT_Boolean, "BooleanLiteral"); !ok {
		return false, err, nil
	} else {
		return true, nil, ast.NewBooleanLiteral(token)
	}
}

func (p *Parser) StringLiteral() (bool, *ParseError, *ast.StringLiteral) {
	if ok, err, token := p.expectToken(lexer.LT_String, "StringLiteral"); !ok {
		return false, err, nil
	} else {
		return true, nil, ast.NewStringLiteral(token)
	}
}

func (p *Parser) ListLiteral() (bool, *ParseError) {
	if ok, err, _ := p.expectToken(lexer.KW_Containing, "ListLiteral-Containing"); !ok {
		return false, err
	}
	element := p.peekToken()
	if element.TokType == lexer.LT_Nothing {
		p.getNextToken() // consume `nothing`
		if ok, err, _ := p.expectToken(lexer.KW_Done, "ListLiteral-NothingDone"); !ok {
			return false, err
		}
		return true, nil
	}
	if ok, err := p.ListElements(); !ok {
		return false, err.addRule("ListLiteral")
	}
	if ok, err, _ := p.expectToken(lexer.KW_Done, "ListLiteral-SomethingDone"); !ok {
		return false, err
	}
	return true, nil
}

func (p *Parser) ListElements() (bool, *ParseError) {
	if ok, err := p.ListValue(); !ok {
		return false, err.addRule("ListElements-One")
	}
	if ok, err, _ := p.expectToken(lexer.OP_Comma, "ListElements-OneComma"); !ok {
		return false, err
	}
	if p.peekToken().TokType == lexer.KW_Done {
		// single element, exit
		return true, nil
	}
	for {
		if p.peekToken().TokType == lexer.KW_And { // exit when see the last element
			p.getNextToken() // consume `and``
			break
		}
		if ok, err := p.ListValue(); !ok {
			return false, err.addRule("ListElements-MoreThan1")
		}
		if ok, err, _ := p.expectToken(lexer.OP_Comma, "ListElements-TwoComma"); !ok {
			return false, err
		}
	}
	if ok, err := p.ListValue(); !ok {
		// last element
		return false, err.addRule("ListElements-LastElement")
	}
	if ok, err, _ := p.expectToken(lexer.OP_Comma, "ListElements-LastComma"); !ok {
		return false, err
	}
	return true, nil
}

func (p *Parser) ListValue() (bool, *ParseError) {
	switch p.peekToken().TokType {
	case lexer.LT_Number:
		ok, err, _ := p.NumberLiteral()
		return ok, err
	case lexer.LT_Boolean:
		ok, err, _ := p.BooleanLiteral()
		return ok, err
	case lexer.LT_String:
		ok, err, _ := p.StringLiteral()
		return ok, err
	case lexer.ItemIdent:
		return true, nil
	case lexer.KW_From, lexer.KW_Every:
		return p.RangeLiteral()
	default:
		return false, NewParseError("TODO", *p.peekToken(), "ListValue")
	}
}

func (p *Parser) RangeLiteral() (bool, *ParseError) {
	if ok, _ := p.maybeToken(lexer.KW_Every, "RangeLiteral-Every"); ok {
		// consume `every`
		p.getNextToken()
		if ok, err := p.Nth(); !ok {
			return false, err.addRule("RangeLiteral-EveryNth")
		}
	}
	if ok, err, _ := p.expectToken(lexer.KW_From, "RangeLiteral-From"); !ok {
		return false, err
	}
	if ok, err := p.RangeStart(); !ok {
		return false, err.addRule("RangeLiteral-Start")
	}
	if ok, err, _ := p.expectToken(lexer.KW_To, "RangeLiteral-To"); !ok {
		return false, err
	}
	if ok, err := p.RangeEnd(); !ok {
		return false, err.addRule("RangeLiteral-End")
	}
	return true, nil
}

func (p *Parser) RangeStart() (bool, *ParseError) {
	if ok, _ := p.maybeToken(lexer.KW_Start, "RangeStart"); ok {
		p.getNextToken() // consume start
		return true, nil
	}
	if ok, err, _ := p.Number(); !ok {
		return false, err.addRule("RangeStart-StartN")
	}
	return true, nil

}

func (p *Parser) RangeEnd() (bool, *ParseError) {
	if ok, _ := p.maybeToken(lexer.KW_End, "RangeEnd"); ok {
		p.getNextToken() // consume end
		return true, nil
	}
	if ok, err, _ := p.Number(); !ok {
		return false, err.addRule("RangeEnd")
	}
	return true, nil
}

func (p *Parser) Ident() (bool, *ParseError, *ast.Identifier) {
	ident := p.getNextToken()
	if ident.TokType != lexer.ItemIdent {
		return false, NewParseError("Expected identifier", ident, "Ident"), nil
	}
	return true, nil, ast.NewIdentifier(&ident)
}

func (p *Parser) Number() (bool, *ParseError, ast.HasValue) {
	if ok, _ := p.maybeToken(lexer.ItemIdent, "Number-Ident"); ok {
		return p.Ident()
	}
	if ok, _ := p.maybeToken(lexer.LT_Number, "Number-NumberLiteral"); ok {
		ok, err, val := p.NumberLiteral()
		return ok, err, val
	}
	return false, NewParseError("Expected number", *p.lastToken, "Number"), nil
}

func (p *Parser) Nth() (bool, *ParseError) {
	if ok, err, _ := p.Number(); !ok {
		return false, err.addRule("Nth-Number")
	}
	ok, err, _ := p.expectToken(lexer.KW_Th, "Nth-Th")
	return ok, err
}

func (p *Parser) FunctionCall() (bool, *ParseError, *ast.FunctionCall) {
	call := ast.FunctionCall{}
	if ok, err, _ := p.expectToken(lexer.KW_Do, "FunctionCall-Do"); !ok {
		return false, err, nil
	}
	if ok, err, funcname := p.expectToken(lexer.ItemIdent, "FunctionCall-FuncName"); !ok {
		return false, err, nil
	} else {
		call.FuncName = ast.NewIdentifier(funcname)
	}
	if ok, err, _ := p.expectToken(lexer.KW_To, "FunctionCall-To"); !ok {
		return false, err, nil
	}
	// TODO: Change to proper expr
	if ok, err, primitive := p.PrimitiveLiteral(); !ok {
		return false, err.addRule("FunctionCall-Parameter1"), nil
	} else {
		call.FuncParams = primitive
	}

	return true, nil, &call
}

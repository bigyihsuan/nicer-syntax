package parser

import (
	"fmt"
	"nicer-syntax/src/lexer"
)

type Node struct {
	lexer.TokItem
}

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

// for interface error.Error()
func (pe *ParseError) Error() string {
	return fmt.Sprintf("PARSE ERROR: `%v` because token `%v` within rule `%v`", pe.Reason, pe.Token, pe.LastRule)
}

type Parser struct {
	Tokens    []lexer.TokItem
	lastToken *lexer.TokItem
}

func NewParser(tokens []lexer.TokItem) Parser {
	return Parser{tokens, &lexer.TokItem{TokType: lexer.ItemEOF, TokName: "nothing", Position: -1, Value: -1}}
}

func (p *Parser) nextToken() lexer.TokItem {
	tok := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	p.lastToken = &tok
	return tok
}

func (p *Parser) peekToken() *lexer.TokItem {
	return &(p.Tokens[0])
}

func (p *Parser) putBackToken() {
	p.Tokens = append([]lexer.TokItem{*p.lastToken}, p.Tokens...)
}

//! Calling Conventions
//! - Consume tokens when expected
//! - Put back when hitting the end of a rule that consumes multiple
//! - When going to a nested rule, peek to check for the starting token

func (p *Parser) Parse() (bool, *ParseError) {
	for len(p.Tokens) > 0 { // while there are still tokens
		if ok, err := p.Program(); !ok {
			return false, err
		}
	}
	return true, nil
}

func (p *Parser) Program() (bool, *ParseError) {
	ok, err := p.Stmt()
	if !ok {
		return false, err
	}
	semicolon := p.nextToken()
	if semicolon.TokType != lexer.ItemSemicolon {
		return false, NewParseError("Expected semicolon to end statement", semicolon, "Program")
	}
	return true, nil
}

func (p *Parser) Stmt() (bool, *ParseError) {
	if p.peekToken().TokType == lexer.ItemIdent {
		return p.IdentAssignment()
	} else {
		return p.IdentDeclaration()
	}
}

func (p *Parser) IdentAssignment() (bool, *ParseError) {
	ident := p.nextToken()
	if ident.TokType != lexer.ItemIdent {
		return false, NewParseError("Expected identifier", ident, "IdentAssignment")
	}
	is := p.nextToken()
	if is.TokType != lexer.KW_Is {
		return false, NewParseError("Expected `is`", is, "IdentAssignment")
	}
	return p.Value()
}

func (p *Parser) IdentDeclaration() (bool, *ParseError) {
	switch p.peekToken().TokType {
	case lexer.KW_Constant:
		return p.ConstDecl()
	case lexer.KW_Variable:
		return p.VarDecl()
	default:
		return false, NewParseError("Invalid token", *p.peekToken(), "IdentDeclaration")
	}
}

func (p *Parser) VarDecl() (bool, *ParseError) {
	variable := p.nextToken()
	if variable.TokType != lexer.KW_Variable {
		return false, NewParseError("Expected `variable`", variable, "VarDecl")
	}
	if ok, err := p.IdentType(); !ok {
		return false, err
	}
	// optional value, ended with semicolon
	if p.peekToken().TokType == lexer.ItemSemicolon {
		return true, nil
	}
	return p.Value()
}

func (p *Parser) ConstDecl() (bool, *ParseError) {
	constant := p.nextToken()
	if constant.TokType != lexer.KW_Constant {
		return false, NewParseError("Expected `constant`", constant, "ConstDecl")
	}
	if ok, err := p.IdentType(); !ok {
		return false, err
	}
	if ok, err := p.Value(); !ok {
		return false, err
	}
	return true, nil
}

func (p *Parser) IdentType() (bool, *ParseError) {
	ident := p.nextToken()
	if ident.TokType != lexer.ItemIdent {
		return false, NewParseError("Expected identifer", ident, "IdentType")
	}
	is := p.nextToken()
	if is.TokType != lexer.KW_Is {
		return false, NewParseError("Expected `is", is, "IdentType")
	}
	typeName := p.nextToken()
	if ok := lexer.IsTypeName[typeName.TokType]; typeName.TokType != lexer.ItemIdent && !ok {
		return false, NewParseError("Expected identifer/typename", typeName, "IdentType")
	}
	return true, nil
}

func (p *Parser) Value() (bool, *ParseError) {
	switch p.peekToken().TokType {
	case lexer.LT_Number:
		return p.NumberLiteral()
	case lexer.LT_Boolean:
		return p.BooleanLiteral()
	case lexer.LT_String:
		return p.StringLiteral()
	default:
		return false, NewParseError("Expected value", *p.peekToken(), "Value")
	}
}

func (p *Parser) NumberLiteral() (bool, *ParseError) {
	tok := p.nextToken()
	fmt.Println("NumberLiteral", tok.Value)
	if tok.TokType != lexer.LT_Number {
		p.putBackToken()
		return false, NewParseError("Expected literal", tok, "NumberLiteral")
	}
	return true, nil
}

func (p *Parser) BooleanLiteral() (bool, *ParseError) {
	tok := p.nextToken()
	fmt.Println("BooleanLiteral", tok.Value)
	if tok.TokType != lexer.LT_Boolean {
		p.putBackToken()
		return false, NewParseError("Expected literal", tok, "BooleanLiteral")
	}
	return true, nil
}

func (p *Parser) StringLiteral() (bool, *ParseError) {
	tok := p.nextToken()
	fmt.Println("StringLiteral", tok.Value)
	if tok.TokType != lexer.LT_String {
		p.putBackToken()
		return false, NewParseError("Expected literal", tok, "StringLiteral")
	}
	return true, nil
}

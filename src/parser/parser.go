package parser

import (
	"fmt"
	"nicer-syntax/src/lexer"

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
	return fmt.Sprintf("%v `%v` because of token `%v` within rule trace `%v`", COLOR_ERROR("PARSE ERROR:"), COLOR_KEYWORD(pe.Reason), COLOR_TOKEN("%v", pe.Token), COLOR_RULE(pe.LastRule))
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
		return false, err.addRule("Program-Stmt")
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
		return false, err.addRule("VarDecl-IdentType")
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
		return false, err.addRule("ConstDecl-IdentType")
	}
	if ok, err := p.Value(); !ok {
		return false, err.addRule("ConstDecl-Value")
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
	if ok, err := p.TypeName(); !ok {
		return false, err.addRule("IdentType-TypeName")
	}
	return true, nil
}

func (p *Parser) TypeName() (bool, *ParseError) {
	typeName := p.nextToken()
	switch typeName.TokType {
	case lexer.TN_Number, lexer.TN_String, lexer.TN_Boolean:
		return true, nil
	case lexer.TN_List:
		of := p.nextToken()
		if of.TokType != lexer.KW_Of {
			return false, NewParseError("Expected `of` for list element type", of, "TypeName")
		}
		return p.TypeName()
	case lexer.TN_Map:
		of := p.nextToken()
		if of.TokType != lexer.KW_Of {
			return false, NewParseError("Expected `of` for map key type", of, "TypeName")
		}
		if ok, err := p.TypeName(); !ok {
			return false, err.addRule("TypeName-MapKey")
		}
		to := p.nextToken()
		if to.TokType != lexer.KW_To {
			return false, NewParseError("Expected `to` for map value type", to, "TypeName")
		}
		return p.TypeName()
	case lexer.ItemIdent: // possibly undeclared typename
		return true, nil
	default:
		return false, NewParseError("Expected type name", typeName, "TypeName")
	}
}

func (p *Parser) Value() (bool, *ParseError) {
	switch p.peekToken().TokType {
	case lexer.LT_Number:
		return p.NumberLiteral()
	case lexer.LT_Boolean:
		return p.BooleanLiteral()
	case lexer.LT_String:
		return p.StringLiteral()
	case lexer.KW_Containing:
		return p.ListLiteral()
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

func (p *Parser) ListLiteral() (bool, *ParseError) {
	containing := p.nextToken()
	if containing.TokType != lexer.KW_Containing {
		return false, NewParseError("Expected `containing` for list literal", containing, "ListLiteral")
	}
	element := p.peekToken()
	if element.TokType == lexer.LT_Nothing {
		p.nextToken() // consume `nothing`
		done := p.nextToken()
		if done.TokType != lexer.KW_Done {
			return false, NewParseError("Expected `done` for empty list literal", done, "ListLiteral")
		}
		return true, nil
	}
	if ok, err := p.ListElements(); !ok {
		return false, err.addRule("ListLiteral")
	}
	done := p.nextToken()
	if done.TokType != lexer.KW_Done {
		return false, NewParseError("Expected `done` for end of list literal", done, "ListLiteral")
	}
	return true, nil
}

func (p *Parser) ListElements() (bool, *ParseError) {
	if ok, err := p.ListValue(); !ok {
		return false, err.addRule("ListElements-One")
	}
	comma := p.nextToken()
	if comma.TokType != lexer.OP_Comma {
		return false, NewParseError("Expected comma after list element", comma, "ListElements-Comma")
	}
	if p.peekToken().TokType == lexer.KW_Done {
		// single element, exit
		return true, nil
	}
	for {
		if p.peekToken().TokType == lexer.KW_And { // exit when see the last element
			p.nextToken() // consume `and``
			break
		}
		if ok, err := p.ListValue(); !ok {
			return false, err.addRule("ListElements-MoreThan1")
		}
		comma := p.nextToken()
		if comma.TokType != lexer.OP_Comma {
			return false, NewParseError("Expected comma after list element", comma, "ListElements-Comma-MoreThan3")
		}
	}
	if ok, err := p.ListValue(); !ok {
		// last element
		return false, err.addRule("ListElements-LastElement")
	}
	comma = p.nextToken()
	if comma.TokType != lexer.OP_Comma {
		return false, NewParseError("Expected comma after last element", comma, "ListElements-LastElementComma")
	}
	return true, nil
}

func (p *Parser) ListValue() (bool, *ParseError) {
	switch p.peekToken().TokType {
	case lexer.LT_Number:
		return p.NumberLiteral()
	case lexer.LT_Boolean:
		return p.BooleanLiteral()
	case lexer.LT_String:
		return p.StringLiteral()
	case lexer.ItemIdent:
		return true, nil
	// TODO: Ranges
	default:
		return false, NewParseError("TODO", *p.peekToken(), "ListValue")
	}
}

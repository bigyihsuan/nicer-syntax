package parser

import (
	"fmt"
	"nicer-syntax/ast"
	"nicer-syntax/evaluator"
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

func (p *Parser) Parse() (bool, *ParseError) {
	for len(p.Tokens) > 0 { // while there are still tokens
		if ok, err := p.Program(); !ok {
			return false, err
		}
	}
	return true, nil
}

func (p *Parser) Program() (bool, *ParseError) {
	if ok, err := p.Stmt(); !ok {
		return false, err.addRule("Program-Stmt")
	}
	if ok, err, _ := p.expectToken(lexer.ItemSemicolon, "Program"); !ok {
		return false, err
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
	if ok, err := p.Ident(); !ok {
		return false, err.addRule("IdentAssignment-Ident")
	}
	if ok, err, _ := p.expectToken(lexer.KW_Is, "IdentAssignment-Is"); !ok {
		return false, err
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
	if ok, err, _ := p.expectToken(lexer.KW_Variable, "VarDecl-Variable"); !ok {
		return false, err
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
	if ok, err, _ := p.expectToken(lexer.KW_Constant, "ConstDecl-Constant"); !ok {
		return false, err
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
	if ok, err, _ := p.expectToken(lexer.ItemIdent, "IdentType-Ident"); !ok {
		return false, err
	}
	if ok, err, _ := p.expectToken(lexer.KW_Is, "IdentType-Is"); !ok {
		return false, err
	}
	if ok, err := p.TypeName(); !ok {
		return false, err.addRule("IdentType-TypeName")
	}
	return true, nil
}

func (p *Parser) TypeName() (bool, *ParseError) {
	typeName := p.getNextToken()
	switch typeName.TokType {
	case lexer.TN_Number, lexer.TN_String, lexer.TN_Boolean:
		return true, nil
	case lexer.TN_List:
		if ok, err, _ := p.expectToken(lexer.KW_Of, "TypeName-ListOf"); !ok {
			return false, err
		}
		return p.TypeName()
	case lexer.TN_Map:
		if ok, err, _ := p.expectToken(lexer.KW_Of, "TypeName-MapOfKey"); !ok {
			return false, err
		}
		if ok, err := p.TypeName(); !ok {
			return false, err.addRule("TypeName-MapKey")
		}
		if ok, err, _ := p.expectToken(lexer.KW_To, "TypeName-MapToValue"); !ok {
			return false, err
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
		ok, err, _ := p.NumberLiteral()
		return ok, err
	case lexer.LT_Boolean:
		ok, err, _ := p.BooleanLiteral()
		return ok, err
	case lexer.LT_String:
		ok, err, _ := p.StringLiteral()
		return ok, err
	case lexer.KW_Containing:
		return p.ListLiteral()
	default:
		return false, NewParseError("Expected value", *p.peekToken(), "Value")
	}
}

func (p *Parser) PrimitiveLiteral() (bool, *ParseError, evaluator.Evaluable) {
	switch p.peekToken().TokType {
	case lexer.LT_Number:
		ok, err, lit := p.NumberLiteral()
		return ok, err, lit
	case lexer.LT_Boolean:
		ok, err, lit := p.BooleanLiteral()
		return ok, err, lit
	case lexer.LT_String:
		ok, err, lit := p.StringLiteral()
		return ok, err, lit
	default:
		return false, NewParseError("Expected primitive literal", *p.peekToken(), "PrimitiveLiteral"), nil
	}
}

func (p *Parser) NumberLiteral() (bool, *ParseError, *ast.NumberLiteral) {
	if ok, err, token := p.expectToken(lexer.LT_Number, "NumberLiteral"); !ok {
		return false, err, nil
	} else {
		return true, nil, &ast.NumberLiteral{TokItem: token}
	}
}

func (p *Parser) BooleanLiteral() (bool, *ParseError, *ast.BooleanLiteral) {
	if ok, err, token := p.expectToken(lexer.LT_Boolean, "BooleanLiteral"); !ok {
		return false, err, nil
	} else {
		return true, nil, &ast.BooleanLiteral{TokItem: token}
	}
}

func (p *Parser) StringLiteral() (bool, *ParseError, *ast.StringLiteral) {
	if ok, err, token := p.expectToken(lexer.LT_String, "StringLiteral"); !ok {
		return false, err, nil
	} else {
		return true, nil, &ast.StringLiteral{TokItem: token}
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
	if ok, err := p.Number(); !ok {
		return false, err.addRule("RangeStart-StartN")
	}
	return true, nil

}

func (p *Parser) RangeEnd() (bool, *ParseError) {
	if ok, _ := p.maybeToken(lexer.KW_End, "RangeEnd"); ok {
		p.getNextToken() // consume end
		return true, nil
	}
	if ok, err := p.Number(); !ok {
		return false, err.addRule("RangeEnd")
	}
	return true, nil
}

func (p *Parser) Ident() (bool, *ParseError) {
	ident := p.getNextToken()
	if ident.TokType != lexer.ItemIdent {
		return false, NewParseError("Expected identifier", ident, "Ident")
	}
	return true, nil
}

func (p *Parser) Number() (bool, *ParseError) {
	if ok, _ := p.maybeToken(lexer.ItemIdent, "Number-Ident"); ok {
		return p.Ident()
	}
	if ok, _ := p.maybeToken(lexer.LT_Number, "Number-NumberLiteral"); ok {
		ok, err, _ := p.NumberLiteral()
		return ok, err
	}
	return false, NewParseError("Expected number", *p.lastToken, "Number")
}

func (p *Parser) Nth() (bool, *ParseError) {
	if ok, err := p.Number(); !ok {
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
		call.FunctionName = *funcname
	}
	if ok, err, _ := p.expectToken(lexer.KW_To, "FunctionCall-To"); !ok {
		return false, err, nil
	}
	// TODO: Change to proper expr
	if ok, err, primitive := p.PrimitiveLiteral(); !ok {
		return false, err.addRule("FunctionCall-Parameter1"), nil
	} else {
		call.Parameters = append(call.Parameters, primitive)
	}
	return true, nil, &call
}

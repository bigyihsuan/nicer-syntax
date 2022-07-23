package tests

import (
	"bytes"
	"nicer-syntax/ast"
	"nicer-syntax/lexer"
	"nicer-syntax/parser"
	"testing"

	"github.com/db47h/lex"
)

func TestParseNumbers(t *testing.T) {
	numbers := []TestCase{
		{"1", true},
		{"12", true},
		{"123", true},
		{"123.456", true},
		{"0.123456", true},
		{".asdt123456", false},
	}

	for _, number := range numbers {
		text := []byte(number.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestNumbers "+number.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], TokPosition: pos, TokValue: val}

		p := parser.NewParser([]lexer.TokItem{token})
		_, err, _ := p.NumberLiteral()
		if err != nil && number.shouldSucceed {
			t.Errorf("failed %v, got %v", number, err)
		}
	}
}

func TestParseBooleans(t *testing.T) {
	booleans := []TestCase{
		{"false", true},
		{"true", true},
	}
	for _, boolean := range booleans {
		text := []byte(boolean.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestBooleans "+boolean.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], TokPosition: pos, TokValue: val}

		p := parser.NewParser([]lexer.TokItem{token})
		_, err, _ := p.BooleanLiteral()
		if err != nil && boolean.shouldSucceed {
			t.Errorf("failed %v, got %v", boolean, err)
		}
	}
}

func TestParseStrings(t *testing.T) {
	strings := []TestCase{
		{`"hello world!"`, true},
		{`"escaped\nstring\n"`, true},
		{`"\\\\"`, true},
		{`"line1\nline2"`, true},
		{`"'single-quoted'"`, true},
		{`"double \"quote\" inside double quotes"`, true},
		{`"should\
fail"`, false},
	}
	for _, str := range strings {
		text := []byte(str.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestStrings "+str.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], TokPosition: pos, TokValue: val}

		p := parser.NewParser([]lexer.TokItem{token})
		_, err, _ := p.StringLiteral()
		if err != nil && str.shouldSucceed {
			t.Errorf("failed %v, got %v", str, err)
		}
	}
}

func TestEvalNumbers(t *testing.T) {
	numbers := []TestCase{
		{"1", true},
		{"12", true},
		{"123", true},
		{"123.456", true},
		{"0.123456", true},
		{".asdt123456", false},
	}

	for _, number := range numbers {
		text := []byte(number.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestNumbers "+number.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], TokPosition: pos, TokValue: val}

		p := parser.NewParser([]lexer.TokItem{token})
		ok, err, numlit := p.NumberLiteral()
		if err != nil && number.shouldSucceed {
			t.Errorf("failed parsing %v, got %v", number, err)
		} else if ok {
			var visitor ast.StringVisitor
			numlit.Accept(&visitor)
		}
	}
}

func TestEvalBooleans(t *testing.T) {
	booleans := []TestCase{
		{"false", true},
		{"true", true},
	}
	for _, boolean := range booleans {
		text := []byte(boolean.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestBooleans "+boolean.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], TokPosition: pos, TokValue: val}

		p := parser.NewParser([]lexer.TokItem{token})
		ok, err, boollit := p.BooleanLiteral()
		if err != nil && boolean.shouldSucceed {
			t.Errorf("failed parsing %v, got %v", boolean, err)
		} else if ok {
			var visitor ast.StringVisitor
			boollit.Accept(&visitor)
		}
	}
}

func TestEvalStrings(t *testing.T) {
	strings := []TestCase{
		{`"hello world!"`, true},
		{`"escaped\nstring\n"`, true},
		{`"\\\\"`, true},
		{`"line1\nline2"`, true},
		{`"'single-quoted'"`, true},
		{`"double \"quote\" inside double quotes"`, true},
		{`"should\
fail"`, false},
	}
	for _, str := range strings {
		text := []byte(str.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestStrings "+str.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], TokPosition: pos, TokValue: val}

		p := parser.NewParser([]lexer.TokItem{token})
		ok, err, strlit := p.StringLiteral()
		if err != nil && str.shouldSucceed {
			t.Errorf("failed parsing %v, got %v", str, err)
		} else if ok {
			var visitor ast.StringVisitor
			strlit.Accept(&visitor)
		}
	}
}

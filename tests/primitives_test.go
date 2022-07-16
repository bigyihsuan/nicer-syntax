package tests

import (
	"bytes"
	"nicer-syntax/src/lexer"
	"nicer-syntax/src/parser"
	"testing"

	"github.com/db47h/lex"
)

func TestNumbers(t *testing.T) {
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
		// fmt.Println(file.Name())
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], Position: pos, Value: val}
		// fmt.Printf("token: %v\n", token)

		p := parser.NewParser([]lexer.TokItem{token})
		_, err := p.NumberLiteral()
		if err != nil && number.shouldSucceed {
			t.Errorf("failed %v, got %v", number, err)
		}
	}
}

func TestBooleans(t *testing.T) {
	booleans := []TestCase{
		{"false", true},
		{"true", true},
	}
	for _, boolean := range booleans {
		text := []byte(boolean.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestBooleans "+boolean.input, byteReader)
		// fmt.Println(file.Name())
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], Position: pos, Value: val}
		// fmt.Printf("token: %v\n", token)

		p := parser.NewParser([]lexer.TokItem{token})
		_, err := p.BooleanLiteral()
		if err != nil && boolean.shouldSucceed {
			t.Errorf("failed %v, got %v", boolean, err)
		}
	}
}

func TestStrings(t *testing.T) {
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
		// fmt.Println(file.Name())
		nicerLexer := lexer.NewLexer(file)
		tok, pos, val := nicerLexer.Lex()
		token := lexer.TokItem{TokType: tok, TokName: lexer.TokenString[tok], Position: pos, Value: val}
		// fmt.Printf("token: %v\n", token)

		p := parser.NewParser([]lexer.TokItem{token})
		_, err := p.StringLiteral()
		if err != nil && str.shouldSucceed {
			t.Errorf("failed %v, got %v", str, err)
		}
	}
}

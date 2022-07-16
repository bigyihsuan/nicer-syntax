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
		{"1", false},
		{"12", false},
		{"123", false},
		{"123.456", false},
		{"0.123456", false},
		{".asdt123456", true},
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
		if err != nil && !number.shouldFail {
			t.Errorf("failed %v, got %v", number, err)
		}
	}
}

func TestBooleans(t *testing.T) {
	booleans := []TestCase{
		{"true", false},
		{"false", false},
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
		if err != nil && !boolean.shouldFail {
			t.Errorf("failed %v, got %v", boolean, err)
		}
	}
}

func TestStrings(t *testing.T) {
	strings := []TestCase{
		{`"hello world!"`, false},
		{`"escaped\nstring\n"`, false},
		{`"\\\\"`, false},
		{`"line1\nline2"`, false},
		{`"'single-quoted'"`, false},
		{`"double \"quote\" inside double quotes"`, false},
		{`"should\
fail"`, true},
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
		if err != nil && !str.shouldFail {
			t.Errorf("failed %v, got %v", str, err)
		}
	}
}

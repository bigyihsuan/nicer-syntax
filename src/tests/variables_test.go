package tests

import (
	"bytes"
	"fmt"
	"nicer-syntax/ast"
	"nicer-syntax/lexer"
	"nicer-syntax/parser"
	"testing"

	"github.com/db47h/lex"
)

func TestDeclarationsConstants(t *testing.T) {
	tests := []TestCase{
		{`constant ConstNumber is number 10`, true},
		{`constant Hello is string "hello world!"`, true},
		{`constant ThisIsNotTrue is boolean false`, true},
		{`constant MissingValue is number`, false},
		{`constant MissingType is 232`, false},
	}
	for _, code := range tests {
		text := []byte(code.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestDeclarationsConstants "+code.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		p := parser.NewParser(tokens)
		var err *parser.ParseError
		ok, err, constdecl := p.ConstDecl()
		if err != nil && code.shouldSucceed {
			t.Errorf("failed `%v`, got %v", code.input, err)
		} else if ok {
			var stringVisitor ast.StringVisitor
			constdecl.Accept(&stringVisitor)
			fmt.Println(stringVisitor)
			var evaluatingVisitor = ast.NewEvaluatingVisitor()
			constdecl.Accept(evaluatingVisitor)
			fmt.Println()
		}
	}
}
func TestDeclarationsVariables(t *testing.T) {
	tests := []TestCase{
		{`variable VarNumber is number 10`, true},
		{`variable Hello is string "hello world!"`, true},
		{`variable ThisIsNotTrue is boolean false`, true},
		{`variable MissingValue is number`, true},
		{`variable MissingType is 232`, false},
	}
	for _, code := range tests {
		text := []byte(code.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestDeclarationsVariables "+code.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		p := parser.NewParser(tokens)
		var err *parser.ParseError
		ok, err, vardecl := p.VarDecl()
		if err != nil && code.shouldSucceed {
			t.Errorf("failed `%v`, got %v", code.input, err)
		} else if ok {
			var stringVisitor ast.StringVisitor
			vardecl.Accept(&stringVisitor)
			fmt.Println(stringVisitor)
			var evaluatingVisitor = ast.NewEvaluatingVisitor()
			vardecl.Accept(evaluatingVisitor)
			fmt.Println()
		}
	}
}

var varAssignments = []TestCase{
	{`variable MissingValue is number
MissingValue is 10`,
		true},
	{`variable MissingBool is boolean
MissingBool is false`,
		true},
	{`variable MissingString is string
MissingString is "no longer missing"`,
		true},
	{`variable TryingToAssignToNonExistent is number
MissingValue is 10`,
		false},
	{`variable WrongType is number
WrongType is "hello"`,
		false},
}

func TestAssignmentsVariables(t *testing.T) {
	tests := varAssignments
	for _, code := range tests {
		text := []byte(code.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestAssignmentsVariables "+code.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		p := parser.NewParser(tokens)
		var err *parser.ParseError
		ok, err, program := p.Program()
		if err != nil && code.shouldSucceed {
			t.Errorf("failed `%v`, got %v", code.input, err)
		} else if ok {
			var stringVisitor ast.StringVisitor
			program.Accept(&stringVisitor)
			fmt.Println(code.input)
			fmt.Println(stringVisitor)
			var evaluatingVisitor = ast.NewEvaluatingVisitor()
			program.Accept(evaluatingVisitor)
			fmt.Println()
		}
	}
}

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

func TestConstants(t *testing.T) {
	tests := []TestCase{
		{`constant ConstNumber is number 10`, true},
		{`constant Hello is string "hello world!"`, true},
		{`constant ThisIsNotTrue is boolean false`, true},
		{`constant MissingValue is number`, false},
		{`constant MissingType is 232`, false},
	}
	for _, print := range tests {
		text := []byte(print.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestRanges "+print.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		p := parser.NewParser(tokens)
		var err *parser.ParseError
		ok, err, constdecl := p.ConstDecl()
		if err != nil && print.shouldSucceed {
			t.Errorf("failed `%v`, got %v", print.input, err)
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

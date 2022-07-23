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

func TestEvalPrintFuncs(t *testing.T) {
	tests := []TestCase{
		{`do PrintLine to "Hello, World!"`, true},
		{`do Print to "Hello, World!"`, true},
		{`do PrintLine to 123.456`, true},
		{`do Print to 123.456`, true},
		{`do PrintLine to true`, true},
		{`do Print to true`, true},
	}
	for _, print := range tests {
		text := []byte(print.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestRanges "+print.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		p := parser.NewParser(tokens)
		var err *parser.ParseError
		ok, err, funccall := p.FunctionCall()
		if err != nil && print.shouldSucceed {
			t.Errorf("failed `%v`, got %v", print.input, err)
		} else if ok {
			var visitor ast.StringVisitor
			funccall.Accept(&visitor)
			fmt.Println(visitor)
			fmt.Println()
		}
	}
}

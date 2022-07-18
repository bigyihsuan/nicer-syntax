package tests

import (
	"bytes"
	"fmt"
	"nicer-syntax/lexer"
	"nicer-syntax/parser"
	"testing"

	"github.com/db47h/lex"
)

func TestParseRanges(t *testing.T) {
	tests := []TestCase{
		{"from start to end", true},
		{"from 10 to end", true},
		{"from start to 5", true},
		{"from 4 to 7", true},
		{"from Id1 to Id2", true},
		{"every 2-th from start to end", true},
		{"every Skip-th from start to end", true},
	}
	for _, rangelit := range tests {
		text := []byte(rangelit.input)
		fmt.Println(string(text))
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestRanges "+rangelit.input, byteReader)
		fmt.Println(file.Name())
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		p := parser.NewParser(tokens)
		var err *parser.ParseError
		_, err = p.RangeLiteral()
		if err != nil && rangelit.shouldSucceed {
			t.Errorf("failed `%v`, got %v", rangelit.input, err)
		}
	}
}

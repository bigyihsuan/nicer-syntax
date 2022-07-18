package tests

import (
	"bytes"
	"fmt"
	"nicer-syntax/lexer"
	"nicer-syntax/parser"
	"testing"

	"github.com/db47h/lex"
)

func TestParseListLiteral(t *testing.T) {
	literals := []TestCase{
		{"containing nothing done", true},
		{"containing 0, done", true},
		{"containing 0, and 1, done", true},
		{"containing 0, 1, 1, 2, 3, 5, 8, and 13, done", true},
		{"containing 0 done", false},        // no comma after last element
		{"containing 0, 1, done", false},    // no `and` before last element
		{"containing 0, and 1 done", false}, // no comma after last element
		// ranges
		{"containing from 0 to 10, done", true},                               // range only
		{"containing every 5-th from 0 to 10, done", true},                    // every-range only
		{"containing from 0 to 10, and 20, done", true},                       // range as first
		{"containing 5, from 0 to 10, and 20, done", true},                    // range as middle
		{"containing 5, 20, and from 0 to 10, done", true},                    // range as last
		{"containing from 0 to 10, and every 2-th from 20 to 40, done", true}, // multiple ranges
		// malformed ranges
		{"containing 0 to 10 done", false},                  // missing `from`
		{"containing from 0 10, done", false},               // missing `to`
		{"containing from to 10, done", false},              // missing start
		{"containing from 0 to , done", false},              // missing end
		{"containing from 0 every 3-th to 10, done", false}, // every in wrong spot
	}
	for _, list := range literals {
		text := []byte(list.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestListLiteral "+list.input, byteReader)
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		p := parser.NewParser(tokens)
		var err *parser.ParseError
		output := captureOutput(func() {
			_, err = p.ListLiteral()
		})
		if err != nil && list.shouldSucceed {
			fmt.Println(output)
			t.Errorf("failed `%v`, got %v", list.input, err)
		}
	}
}

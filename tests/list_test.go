package tests

import (
	"bytes"
	"nicer-syntax/src/lexer"
	"nicer-syntax/src/parser"
	"testing"

	"github.com/db47h/lex"
)

func TestListLiteral(t *testing.T) {
	literals := []TestCase{
		{"containing nothing done", false},
		{"containing 0, done", false},
		{"containing 0, and 1, done", false},
		{"containing 0, 1, 1, 2, 3, 5, 8, and 13, done", false},
		{"containing 0 done", true},        // no comma after last element
		{"containing 0, 1, done", true},    // no `and` before last element
		{"containing 0, and 1 done", true}, // no comma after last element
	}
	for _, list := range literals {
		text := []byte(list.input)
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestNumbers "+list.input, byteReader)
		// fmt.Println(file.Name())
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()
		// fmt.Printf("token: %v\n", token)

		p := parser.NewParser(tokens)
		_, err := p.ListLiteral()
		if err != nil && !list.shouldFail {
			t.Errorf("failed %v, got %v", list.input, err)
		}
	}
}

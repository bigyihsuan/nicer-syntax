package tests

import (
	"bytes"
	"fmt"
	"nicer-syntax/lexer"
	"nicer-syntax/parser"

	"github.com/db47h/lex"
)

func TestRanges() {
	tests := []string{
		"from start to end",
		"from 10 to end",
		"from start to 5",
		"from 4 to 7",
		"from Id1 to Id2",
		"every 2-th from start to end",
		"every Skip-th from start to end",
	}
	for _, s := range tests {
		text := []byte(s)
		fmt.Println(string(text))
		byteReader := bytes.NewBuffer(text)
		file := lex.NewFile("TestRanges "+s, byteReader)
		fmt.Println(file.Name())
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()
		// for _, tokitem := range tokens {
		// 	fmt.Printf("%v\n", tokitem)
		// }

		p := parser.NewParser(tokens)
		result, err := p.RangeLiteral()
		fmt.Printf("result: %v\n", result)
		fmt.Printf("%v\n", err)
	}
}

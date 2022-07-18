package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"nicer-syntax/lexer"
	"nicer-syntax/parser"
	"os"

	"github.com/db47h/lex"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	filename := os.Args[1]
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// tests := []string{
	// 	"from start to end",
	// 	"from 10 to end",
	// 	"from start to 5",
	// 	"from 4 to 7",
	// 	"from Id1 to Id2",
	// 	"every 2-th from start to end",
	// 	"every Skip-th from start to end",
	// }
	text = []byte("do PrintLine to 10")
	fmt.Println(string(text))
	byteReader := bytes.NewBuffer(text)
	file := lex.NewFile(filename, byteReader)
	fmt.Println(file.Name())
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	// for _, tokitem := range tokens {
	// 	fmt.Printf("%v\n", tokitem)
	// }

	p := parser.NewParser(tokens)
	result, err := p.FunctionCall()
	fmt.Printf("result: %v\n", result)
	fmt.Printf("%v\n", err)
}

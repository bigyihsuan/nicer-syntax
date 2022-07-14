package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"nicer-syntax/src/lexer"
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
	byteReader := bytes.NewBuffer(text)
	file := lex.NewFile(filename, byteReader)
	fmt.Println(file.Name())
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	for _, tokitem := range tokens {
		fmt.Printf("%v\n", tokitem)
	}
}

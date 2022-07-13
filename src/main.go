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
	// read, _ := os.ReadFile(filename)
	// reader := bytes.NewReader(read)
	// for ;reader.Len() > 0; {
	// 	ch, size, err := reader.ReadRune()
	// 	fmt.Println(string(ch), size, err)
	// }
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// text = []byte("1123")
	byteReader := bytes.NewBuffer(text)
	file := lex.NewFile(filename, byteReader)
	fmt.Println(file.Name())
	nicerLexer := lexer.NewLexer(file)

	for tok, _, v := nicerLexer.Lex(); tok != lexer.ItemEOF; tok, _, v = nicerLexer.Lex() {
		fmt.Printf("%v\t\t%v\n", lexer.TokenString[tok], v)
	}
}

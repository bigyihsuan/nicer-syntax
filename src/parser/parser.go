package parser

import "nicer-syntax/src/lexer"

type AST struct {
	Value lexer.TokItem
	Left, Right *AST
}

type Parser struct {
	TokensLeft []lexer.TokItem
}
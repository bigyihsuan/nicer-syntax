package lexer

// https://pkg.go.dev/github.com/db47h/lex/state#example-package-Go

import (
	"fmt"
	"unicode"

	"github.com/db47h/lex"
	"github.com/db47h/lex/state"
)

const (
	ItemError lex.Token = iota
	ItemEOF
	ItemComment
	ItemIdent
	ItemSemicolon
	// literals
	LT_Number
	LT_String
	LT_True
	LT_False
	LT_Null
	// built-in type names
	TN_Number
	TN_String
	TN_List
	TN_Map
	TN_Struct
	// declaration keywords
	KW_Variable
	KW_Constant
	KW_Type
	KW_Function
	// random keywords
	KW_Is
	KW_Of
	KW_Can
	KW_Where
	KW_Do
	KW_Done
	KW_Does
	KW_Th
	KW_Containing
	KW_Taking
	KW_Returning
	KW_Return
	// loops
	KW_For
	KW_While
	KW_Loop
	// conditionals
	KW_If
	KW_Then
	KW_Else
	// ranges
	KW_From
	KW_To
	KW_Start
	KW_End
	KW_Every
	// boolean operators
	KW_And
	KW_Or
	KW_Not
	// symbolic operators
	OP_Eq
	OP_Neq
	OP_Gt
	OP_Lt
	OP_GtEq
	OP_LtEq
	OP_Plus
	OP_Minus
	OP_Star
	OP_Slash
	OP_Percent
	OP_Caret
	OP_Lparen
	OP_Rparen
	OP_Comma
)

var keywords = map[string]lex.Token{
	// literals
	"true": LT_True,
	"false": LT_False,
	"null": LT_Null,
	// built-in type names
	"number": TN_Number,
	"string": TN_String,
	"list": TN_List,
	"map": TN_Map,
	"struct": TN_Struct,
	// declaration keywords
	"variable": KW_Variable,
	"constant": KW_Constant,
	"type": KW_Type,
	"function": KW_Function,
	// random keywords
	"is": KW_Is,
	"of": KW_Of,
	"can": KW_Can,
	"where": KW_Where,
	"do": KW_Do,
	"done": KW_Done,
	"does": KW_Does,
	"-th": KW_Th,
	"containing": KW_Containing,
	"from": KW_From,
	"to": KW_To,
	"start": KW_Start,
	"end": KW_End,
	"every": KW_Every,
	"taking": KW_Taking,
	"returning": KW_Returning,
	"return": KW_Return,
	// loops
	"for": KW_For,
	"while": KW_While,
	"loop": KW_Loop,
	// conditionals
	"if": KW_If,
	"then": KW_Then,
	"else": KW_Else,
	// boolean operators
	"and": KW_And,
	"or": KW_Or,
	"not": KW_Not,
}

var TokenString = map[lex.Token]string{
	ItemError:     "\t\tItemError",
	ItemEOF: "ItemEOF",
	ItemComment: "ItemComment",
	ItemIdent: "ItemIdent",
	ItemSemicolon: "ItemSemicolon",
	// literals
	LT_Number: "LT_Number",
	LT_String: "LT_String",
	LT_True: "LT_True",
	LT_False: "LT_False",
	LT_Null: "LT_Null",
	// built-in type names
	TN_Number: "TN_Number",
	TN_String: "TN_String",
	TN_List: "TN_List",
	TN_Map: "TN_Map",
	TN_Struct: "TN_Struct",
	// declaration keywords
	KW_Variable: "KW_Variable",
	KW_Constant: "KW_Constant",
	KW_Type: "KW_Type",
	KW_Function: "KW_Function",
	// random keywords
	KW_Is: "KW_Is",
	KW_Of: "KW_Of",
	KW_Can: "KW_Can",
	KW_Where: "KW_Where",
	KW_Do: "KW_Do",
	KW_Done: "KW_Done",
	KW_Does: "KW_Does",
	KW_Th: "KW_Th",
	KW_Containing: "KW_Containing",
	KW_Taking: "KW_Taking",
	KW_Returning: "KW_Returning",
	KW_Return: "KW_Return",
	// loops
	KW_For: "KW_For",
	KW_While: "KW_While",
	KW_Loop: "KW_Loop",
	// conditionals
	KW_If: "KW_If",
	KW_Then: "KW_Then",
	KW_Else: "KW_Else",
	// ranges
	KW_From: "KW_From",
	KW_To: "KW_To",
	KW_Start: "KW_Start",
	KW_End: "KW_End",
	KW_Every: "KW_Every",
	// boolean operators
	KW_And: "KW_And",
	KW_Or: "KW_Or",
	KW_Not: "KW_Not",
	// symbolic operators
	OP_Eq: "OP_Eq",
	OP_Neq: "OP_Neq",
	OP_Gt: "OP_Gt",
	OP_Lt: "OP_Lt",
	OP_GtEq: "OP_GtEq",
	OP_LtEq: "OP_LtEq",
	OP_Plus: "OP_Plus",
	OP_Minus: "OP_Minus",
	OP_Star: "OP_Star",
	OP_Slash: "OP_Slash",
	OP_Percent: "OP_Percent",
	OP_Caret: "OP_Caret",
	OP_Lparen: "OP_Lparen",
	OP_Rparen: "OP_Rparen",
	OP_Comma: "OP_Comma",
}

type TokItem struct {
	TokType lex.Token
	TokName string
	Position int
	Value interface{}
}

// for String() string
func (ti TokItem) String() string {
	return fmt.Sprintf("{%s @ %v: `%s`}", ti.TokName, ti.Position, ti.Value)
}

type NicerLexer struct {
	lex.Lexer
}

func NewLexer(file *lex.File) *NicerLexer {
	l := &NicerLexer{}
	l.Lexer = *lex.NewLexer(file, l.program)
	return l
}

func (nl *NicerLexer) LexAll() []TokItem {
	var tokens []TokItem
	for tok, pos, v := nl.Lex(); tok != ItemEOF; tok, pos, v = nl.Lex() {
		tokens = append(tokens, TokItem{tok, TokenString[tok], pos, v})
	}
	return tokens
}

func (nl *NicerLexer) program(s *lex.State) lex.StateFn {
	r := s.Next()
	pos := s.Pos()
	switch r { // single-character tokens
	case lex.EOF:
		s.Emit(pos, ItemEOF, nil)
		return nil
	case '\n': // newlines separate statements
		s.Emit(pos, ItemSemicolon, ";")
		return nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return state.Number(LT_Number, LT_Number, '.')
		case '"': // strings
		return state.QuotedString(LT_String)
	case ',':
		s.Emit(pos, OP_Comma, string(r))
		return nil
	case '+':
		s.Emit(pos, OP_Plus, string(r))
		return nil
	case '-': // either binary minus, unary minus, or -th operator
		if s.Peek() == 't' {
			return nl.keyword
		}
		s.Emit(pos, OP_Minus, string(r))
		return nil
	case '*':
		s.Emit(pos, OP_Star, string(r))
		return nil
	case '/':
		s.Emit(pos, OP_Slash, string(r))
		return nil
	case '%':
		s.Emit(pos, OP_Percent, string(r))
		return nil
	case '^':
		s.Emit(pos, OP_Caret, string(r))
		return nil
	case '=':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, OP_Eq, "==")
		} else {
			s.Emit(pos, ItemError, string(r))
			s.Backup()
		}
		return nil
	case '!':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, OP_Eq, "!=")
		} else {
			s.Emit(pos, ItemError, string(r))
			s.Backup()
		}
		return nil
	case '>':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, OP_GtEq, ">=")
		} else {
			s.Emit(pos, OP_Gt, ">")
			s.Backup()
		}
		return nil
	case '<':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, OP_LtEq, "<=")
		} else {
			s.Emit(pos, OP_Lt, "<")
			s.Backup()
		}
		return nil
	case '(':
		s.Emit(pos, OP_Lparen, "(")
		return nil
	case ')':
		s.Emit(pos, OP_Rparen, ")")
		return nil
	case '#': // comments
		return nl.comment
	}
	switch {
	case unicode.IsSpace(r):
		// consume spaces
		for r = s.Next(); unicode.IsSpace(r); r = s.Next() {
			// nop
		}
		s.Backup()
		return nil
	case unicode.IsUpper(r): // identifier
		return nl.ident
	case unicode.IsLower(r): // keyword
		return nl.keyword
	}
	return nil
}

func (nl *NicerLexer) ident(s *lex.State) lex.StateFn {
	identName := make([]rune, 0, 64)
	return func(l *lex.State) lex.StateFn {
		pos := l.Pos()
		identName = append(identName[:0], l.Current())
		for r := l.Next(); unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'; r = l.Next() {
			identName = append(identName, r)
		}
		l.Backup()
		l.Emit(pos, ItemIdent, string(identName))
		return nil
	}
}

func (nl *NicerLexer) comment(s *lex.State) lex.StateFn {
	comment := make([]rune, 0, 64)
	return func(l *lex.State) lex.StateFn {
		comment = append(comment[:0], l.Current())
		for r := l.Next(); r != '\n'; r = l.Next() {
			comment = append(comment, r)
		}
		l.Backup()
		// completely ignore comments
		return nil
	}
}

func (nl *NicerLexer) keyword(s *lex.State) lex.StateFn {
	kw := make([]rune, 0, 64)
	return func(l *lex.State) lex.StateFn {
		pos := l.Pos()
		kw = append(kw[:0], l.Current())
		for r := l.Next(); unicode.IsLetter(r); r = l.Next() {
			kw = append(kw, r)
		}
		l.Backup()
		l.Emit(pos, keywords[string(kw)], string(kw))
		return nil
	}
}

// func (nl *NicerLexer) VarDecl(s *lex.State) lex.StateFn {
// 	r := s.Next()
// 	if r == lex.EOF {
// 		return nl.eof
// 	}
// 	// eat chars until hit space
// 	for string(nl.buffer) != "variable" && !unicode.IsSpace(r) {
// 		nl.buffer = append(nl.buffer, r)
// 		r = s.Next()
// 	}
// 	s.Emit(s.TokenPos(), keywordVariable, string(nl.buffer))
// 	// get identifier

// 	return nil
// }

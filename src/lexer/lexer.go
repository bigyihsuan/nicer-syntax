package lexer

// https://pkg.go.dev/github.com/db47h/lex/state#example-package-Go

import (
	"unicode"

	"github.com/db47h/lex"
	"github.com/db47h/lex/state"
)

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
		s.Emit(pos, ItemSemicolon, nil)
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

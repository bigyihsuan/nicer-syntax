package lexer

import (
	"fmt"

	"github.com/db47h/lex"
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
	LT_Boolean
	LT_Nothing
	// built-in type names
	TN_Number
	TN_Boolean
	TN_String
	TN_List
	TN_Map
	TN_Struct
	// struct-specific keywords
	KW_Containing
	KW_Can
	// declaration keywords
	KW_Variable
	KW_Constant
	KW_Type
	KW_Function
	KW_Is
	// random keywords
	KW_Of
	KW_Where
	KW_Do
	KW_Done
	KW_Does
	KW_Th
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
	"true":    LT_Boolean,
	"false":   LT_Boolean,
	"nothing": LT_Nothing,
	// built-in type names
	"number":  TN_Number,
	"boolean": TN_Boolean,
	"string":  TN_String,
	"list":    TN_List,
	"map":     TN_Map,
	"struct":  TN_Struct,
	// declaration keywords
	"variable": KW_Variable,
	"constant": KW_Constant,
	"type":     KW_Type,
	"function": KW_Function,
	// random keywords
	"is":         KW_Is,
	"of":         KW_Of,
	"can":        KW_Can,
	"where":      KW_Where,
	"do":         KW_Do,
	"done":       KW_Done,
	"does":       KW_Does,
	"-th":        KW_Th,
	"containing": KW_Containing,
	"from":       KW_From,
	"to":         KW_To,
	"start":      KW_Start,
	"end":        KW_End,
	"every":      KW_Every,
	"taking":     KW_Taking,
	"returning":  KW_Returning,
	"return":     KW_Return,
	// loops
	"for":   KW_For,
	"while": KW_While,
	"loop":  KW_Loop,
	// conditionals
	"if":   KW_If,
	"then": KW_Then,
	"else": KW_Else,
	// boolean operators
	"and": KW_And,
	"or":  KW_Or,
	"not": KW_Not,
}

var TokenString = map[lex.Token]string{
	ItemError:     "\t\tItemError",
	ItemEOF:       "ItemEOF",
	ItemComment:   "ItemComment",
	ItemIdent:     "ItemIdent",
	ItemSemicolon: "ItemSemicolon",
	// literals
	LT_Number:  "LT_Number",
	LT_String:  "LT_String",
	LT_Boolean: "LT_Boolean",
	LT_Nothing: "LT_Nothing",
	// built-in type names
	TN_Number:  "TN_Number",
	TN_Boolean: "TN_Boolean",
	TN_String:  "TN_String",
	TN_List:    "TN_List",
	TN_Map:     "TN_Map",
	TN_Struct:  "TN_Struct",
	// declaration keywords
	KW_Variable: "KW_Variable",
	KW_Constant: "KW_Constant",
	KW_Type:     "KW_Type",
	KW_Function: "KW_Function",
	// random keywords
	KW_Is:         "KW_Is",
	KW_Of:         "KW_Of",
	KW_Can:        "KW_Can",
	KW_Where:      "KW_Where",
	KW_Do:         "KW_Do",
	KW_Done:       "KW_Done",
	KW_Does:       "KW_Does",
	KW_Th:         "KW_Th",
	KW_Containing: "KW_Containing",
	KW_Taking:     "KW_Taking",
	KW_Returning:  "KW_Returning",
	KW_Return:     "KW_Return",
	// loops
	KW_For:   "KW_For",
	KW_While: "KW_While",
	KW_Loop:  "KW_Loop",
	// conditionals
	KW_If:   "KW_If",
	KW_Then: "KW_Then",
	KW_Else: "KW_Else",
	// ranges
	KW_From:  "KW_From",
	KW_To:    "KW_To",
	KW_Start: "KW_Start",
	KW_End:   "KW_End",
	KW_Every: "KW_Every",
	// boolean operators
	KW_And: "KW_And",
	KW_Or:  "KW_Or",
	KW_Not: "KW_Not",
	// symbolic operators
	OP_Eq:      "OP_Eq",
	OP_Neq:     "OP_Neq",
	OP_Gt:      "OP_Gt",
	OP_Lt:      "OP_Lt",
	OP_GtEq:    "OP_GtEq",
	OP_LtEq:    "OP_LtEq",
	OP_Plus:    "OP_Plus",
	OP_Minus:   "OP_Minus",
	OP_Star:    "OP_Star",
	OP_Slash:   "OP_Slash",
	OP_Percent: "OP_Percent",
	OP_Caret:   "OP_Caret",
	OP_Lparen:  "OP_Lparen",
	OP_Rparen:  "OP_Rparen",
	OP_Comma:   "OP_Comma",
}

type TokItem struct {
	TokType  lex.Token
	TokName  string
	Position int
	Value    interface{}
}

// for String() string
func (ti TokItem) String() string {
	return fmt.Sprintf("{%s @ %v: %#v}", ti.TokName, ti.Position, ti.Value)
}

package token

import (
	"fmt"
	"strings"
)

// Type represents a type of token, such as an integer or a boolean.
type Type string

// Token represents a small, easily categorizable chunk of text within the source code.
// For example, the number "5" is an integer, so we need to know that it's type is INTEGER and its value is "5".
// It also contains information about the location of the token within the source as to give more meaningul error messages.
type Token struct {
	Type    Type   // The type of token.
	Literal string // The literal value of the token, such as "5" or "true".

	Line     int // The line the token is located at.
	StartCol int // The location of the start of the token.
	EndCol   int // The location of the end of the token.
}

// String returns a string representation of the token.
// The format is LITERAL<TYPE>(line=LINE_NUM,col=COL_NUM)
func (t Token) String() string {
	if t.Literal == "\n" {
		return fmt.Sprintf("(Lit: '%s', Type: '%s', line=%d, startcol=%d, endcol=%d)", "\\n", t.Type, t.Line, t.StartCol, t.EndCol)
	}

	return fmt.Sprintf("(Lit: '%s', Type: '%s', line=%d, startcol=%d, endcol=%d)", t.Literal, t.Type, t.Line, t.StartCol, t.EndCol)
}

// Definitions of token types.
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers/Literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	FLOAT  = "FLOAT"

	// Operators
	ASSIGN   = "<-"
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "MOD"
	DIV      = "DIV"
	LSHIFT   = "<<"
	RSHIFT   = ">>"

	LT     = "<"
	GT     = ">"
	LT_EQ  = "<="
	GT_EQ  = ">="
	EQ     = "=="
	NOT_EQ = "!="

	DOT = "."

	NOT = "NOT"
	AND = "AND"
	OR  = "OR"
	XOR = "XOR"

	// Delimeters
	COMMA   = ","
	NEWLINE = "\\n"
	COLON   = ":"

	// Brackets/Braces/Parenthesis
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	SUBROUTINE = "SUBROUTINE"
	CONSTANT   = "CONSTANT"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	IF         = "IF"
	ELSE       = "ELSE"
	RETURN     = "RETURN"

	WHILE  = "WHILE"
	REPEAT = "REPEAT"
	UNTIL  = "UNTIL"
	FOR    = "FOR"
	TO     = "TO"

	ENDWHILE      = "ENDWHILE"
	ENDFOR        = "ENDFOR"
	ENDIF         = "ENDIF"
	ENDSUBROUTINE = "ENDSUBROUTINE"

	THEN = "THEN"
	MAP  = "MAP"

	// Importing
	// IMPORT "file" (directly accessible)
	// IMPORT "folder" (accessible using `folder.functionName`
	//
	// The file needs to be run
	IMPORT = "IMPORT"
	AS     = "AS"
	FROM   = "FROM"

	// Special 'functions'
	OUTPUT    = "OUTPUT"
	USERINPUT = "USERINPUT"
)

// NewToken returns a new token from a given Type, Literal and position in the source.
func NewToken(tokenType Type, lit string, line, startCol, endCol int) Token {
	return Token{
		Type:     tokenType,
		Literal:  lit,
		Line:     line,
		StartCol: startCol,
		EndCol:   endCol,
	}
}

// Keywords maps the lowercase name of a keyword to the associated token.Type.
var Keywords = map[string]Type{
	"constant": CONSTANT,
	"return":   RETURN,

	"true":  TRUE,
	"false": FALSE,

	"if":         IF,
	"else":       ELSE,
	"subroutine": SUBROUTINE,

	"while":    WHILE,
	"endwhile": ENDWHILE,

	"repeat": REPEAT,
	"until":  UNTIL,

	"for":    FOR,
	"endfor": ENDFOR,
	"to":     TO,

	"then":          THEN,
	"endif":         ENDIF,
	"endsubroutine": ENDSUBROUTINE,

	"output":    OUTPUT,
	"userinput": USERINPUT,

	"div": DIV,
	"mod": MOD,

	"not": NOT,
	"and": AND,
	"or":  OR,
	"xor": XOR,

	"import": IMPORT,
	"as":     AS,
	"from":   FROM,

	"map": MAP,
}

// LookupKeyword converts a keyword name into a keyword.
// When checking if a given ident is a keyword, we only want to accept keywords that are either all UPPERCASE or all lowercase,
// not mIxEdCaSe.
func LookupKeyword(ident string) Type {
	if keyword, ok := Keywords[strings.ToLower(ident)]; ok {
		return keyword
	}

	return ILLEGAL
}

// LookupIdent converts an ident name into a token.Type. If it cannot find a keyword or associated token type, it will return
// token.IDENT.
func LookupIdent(ident string) Type {
	upper := strings.ToUpper(ident)
	lower := strings.ToLower(ident)

	// Mixed case is always an ident.
	if lower != ident && upper != ident {
		return IDENT
	}

	if val := LookupKeyword(ident); val != ILLEGAL {
		return val
	}

	return IDENT
}

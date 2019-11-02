package lexer

import (
	"testing"

	"github.com/ollybritton/aqa/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `five <- 5
ten <- 10
!+-*/5
!= == = @
5 < 10 > 5
({[1, potato, 3]})
# I shouldn't affect anything
"foobar"
"foo bar"
'foobar'
'foo bar'
'foo\'s bar'
"foo\"s bar"
WHILE true
	a
ENDWHILE
REPEAT
	a
UNTIL false
FOR i <- 10 TO 20
	a
ENDFOR
# Twice
# Twice
OUTPUT USERINPUT
1.234
0xBADA55
0b10
[1, 2]
1<=2
1>=2
1>>2
2<<1
1 MOD 2
2 DIV 3
NOT true OR false AND false XOR true
a123
MAP {
    "a": 10,
}`

	tests := []token.Token{
		{Type: token.IDENT, Literal: "five", Line: 0, Column: 0},
		{Type: token.ASSIGN, Literal: "<-", Line: 0, Column: 5},
		{Type: token.INT, Literal: "5", Line: 0, Column: 8},
		{Type: token.NEWLINE, Literal: "\n", Line: 0, Column: 9},
		{Type: token.IDENT, Literal: "ten", Line: 1, Column: 0},
		{Type: token.ASSIGN, Literal: "<-", Line: 1, Column: 4},
		{Type: token.INT, Literal: "10", Line: 1, Column: 7},
		{Type: token.NEWLINE, Literal: "\n", Line: 1, Column: 9},
		{Type: token.BANG, Literal: "!", Line: 2, Column: 0},
		{Type: token.PLUS, Literal: "+", Line: 2, Column: 1},
		{Type: token.MINUS, Literal: "-", Line: 2, Column: 2},
		{Type: token.ASTERISK, Literal: "*", Line: 2, Column: 3},
		{Type: token.SLASH, Literal: "/", Line: 2, Column: 4},
		{Type: token.INT, Literal: "5", Line: 2, Column: 5},
		{Type: token.NEWLINE, Literal: "\n", Line: 2, Column: 6},
		{Type: token.NOT_EQ, Literal: "!=", Line: 3, Column: 0},
		{Type: token.EQ, Literal: "==", Line: 3, Column: 3},
		{Type: token.EQ, Literal: "=", Line: 3, Column: 6},
		{Type: token.ILLEGAL, Literal: "@", Line: 3, Column: 8},
		{Type: token.NEWLINE, Literal: "\n", Line: 3, Column: 9},
		{Type: token.INT, Literal: "5", Line: 4, Column: 0},
		{Type: token.LT, Literal: "<", Line: 4, Column: 2},
		{Type: token.INT, Literal: "10", Line: 4, Column: 4},
		{Type: token.GT, Literal: ">", Line: 4, Column: 7},
		{Type: token.INT, Literal: "5", Line: 4, Column: 9},
		{Type: token.NEWLINE, Literal: "\n", Line: 4, Column: 10},
		{Type: token.LPAREN, Literal: "(", Line: 5, Column: 0},
		{Type: token.LBRACE, Literal: "{", Line: 5, Column: 1},
		{Type: token.LBRACKET, Literal: "[", Line: 5, Column: 2},
		{Type: token.INT, Literal: "1", Line: 5, Column: 3},
		{Type: token.COMMA, Literal: ",", Line: 5, Column: 4},
		{Type: token.IDENT, Literal: "potato", Line: 5, Column: 6},
		{Type: token.COMMA, Literal: ",", Line: 5, Column: 12},
		{Type: token.INT, Literal: "3", Line: 5, Column: 14},
		{Type: token.RBRACKET, Literal: "]", Line: 5, Column: 15},
		{Type: token.RBRACE, Literal: "}", Line: 5, Column: 16},
		{Type: token.RPAREN, Literal: ")", Line: 5, Column: 17},
		{Type: token.NEWLINE, Literal: "\n", Line: 5, Column: 18},

		{Type: token.STRING, Literal: "foobar", Line: 7, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 7, Column: 8},
		{Type: token.STRING, Literal: "foo bar", Line: 8, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 8, Column: 9},
		{Type: token.STRING, Literal: "foobar", Line: 9, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 9, Column: 8},
		{Type: token.STRING, Literal: "foo bar", Line: 10, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 10, Column: 9},
		{Type: token.STRING, Literal: "foo's bar", Line: 11, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 11, Column: 12},
		{Type: token.STRING, Literal: "foo\"s bar", Line: 12, Column: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 12, Column: 12},
		{Type: token.WHILE, Literal: "WHILE", Line: 13, Column: 0},
		{Type: token.TRUE, Literal: "true", Line: 13, Column: 6},
		{Type: token.NEWLINE, Literal: "\n", Line: 13, Column: 10},
		{Type: token.IDENT, Literal: "a", Line: 14, Column: 1},
		{Type: token.NEWLINE, Literal: "\n", Line: 14, Column: 2},
		{Type: token.ENDWHILE, Literal: "ENDWHILE", Line: 15, Column: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 15, Column: 8},
		{Type: token.REPEAT, Literal: "REPEAT", Line: 16, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 16, Column: 6},
		{Type: token.IDENT, Literal: "a", Line: 17, Column: 1},
		{Type: token.NEWLINE, Literal: "\n", Line: 17, Column: 2},
		{Type: token.UNTIL, Literal: "UNTIL", Line: 18, Column: 0},
		{Type: token.FALSE, Literal: "false", Line: 18, Column: 6},

		{Type: token.NEWLINE, Literal: "\n", Line: 18, Column: 11},
		{Type: token.FOR, Literal: "FOR", Line: 19, Column: 0},
		{Type: token.IDENT, Literal: "i", Line: 19, Column: 4},
		{Type: token.ASSIGN, Literal: "<-", Line: 19, Column: 6},
		{Type: token.INT, Literal: "10", Line: 19, Column: 9},
		{Type: token.TO, Literal: "TO", Line: 19, Column: 12},
		{Type: token.INT, Literal: "20", Line: 19, Column: 15},
		{Type: token.NEWLINE, Literal: "\n", Line: 19, Column: 17},
		{Type: token.IDENT, Literal: "a", Line: 20, Column: 1},
		{Type: token.NEWLINE, Literal: "\n", Line: 20, Column: 2},
		{Type: token.ENDFOR, Literal: "ENDFOR", Line: 21, Column: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 21, Column: 6},
		// ...skip comments...
		{Type: token.OUTPUT, Literal: "OUTPUT", Line: 24, Column: 0},
		{Type: token.USERINPUT, Literal: "USERINPUT", Line: 24, Column: 7},

		{Type: token.NEWLINE, Literal: "\n", Line: 24, Column: 16},
		{Type: token.FLOAT, Literal: "1.234", Line: 25, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 25, Column: 5},
		{Type: token.INT, Literal: "0xBADA55", Line: 26, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 26, Column: 8},
		{Type: token.INT, Literal: "0b10", Line: 27, Column: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 27, Column: 4},
		{Type: token.LBRACKET, Literal: "[", Line: 28, Column: 0},
		{Type: token.INT, Literal: "1", Line: 28, Column: 1},
		{Type: token.COMMA, Literal: ",", Line: 28, Column: 2},
		{Type: token.INT, Literal: "2", Line: 28, Column: 4},
		{Type: token.RBRACKET, Literal: "]", Line: 28, Column: 5},

		{Type: token.NEWLINE, Literal: "\n", Line: 28, Column: 6},
		{Type: token.INT, Literal: "1", Line: 29, Column: 0},
		{Type: token.LT_EQ, Literal: "<=", Line: 29, Column: 1},
		{Type: token.INT, Literal: "2", Line: 29, Column: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 29, Column: 4},
		{Type: token.INT, Literal: "1", Line: 30, Column: 0},
		{Type: token.GT_EQ, Literal: ">=", Line: 30, Column: 1},
		{Type: token.INT, Literal: "2", Line: 30, Column: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 30, Column: 4},
		{Type: token.INT, Literal: "1", Line: 31, Column: 0},
		{Type: token.RSHIFT, Literal: ">>", Line: 31, Column: 1},
		{Type: token.INT, Literal: "2", Line: 31, Column: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 31, Column: 4},
		{Type: token.INT, Literal: "2", Line: 32, Column: 0},
		{Type: token.LSHIFT, Literal: "<<", Line: 32, Column: 1},
		{Type: token.INT, Literal: "1", Line: 32, Column: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 32, Column: 4},
		{Type: token.INT, Literal: "1", Line: 33, Column: 0},
		{Type: token.MOD, Literal: "MOD", Line: 33, Column: 2},
		{Type: token.INT, Literal: "2", Line: 33, Column: 6},

		{Type: token.NEWLINE, Literal: "\n", Line: 33, Column: 7},
		{Type: token.INT, Literal: "2", Line: 34, Column: 0},
		{Type: token.DIV, Literal: "DIV", Line: 34, Column: 2},
		{Type: token.INT, Literal: "3", Line: 34, Column: 6},

		{Type: token.NEWLINE, Literal: "\n", Line: 34, Column: 7},
		{Type: token.NOT, Literal: "NOT", Line: 35, Column: 0},
		{Type: token.TRUE, Literal: "true", Line: 35, Column: 4},
		{Type: token.OR, Literal: "OR", Line: 35, Column: 9},
		{Type: token.FALSE, Literal: "false", Line: 35, Column: 12},
		{Type: token.AND, Literal: "AND", Line: 35, Column: 18},
		{Type: token.FALSE, Literal: "false", Line: 35, Column: 22},
		{Type: token.XOR, Literal: "XOR", Line: 35, Column: 28},
		{Type: token.TRUE, Literal: "true", Line: 35, Column: 32},

		{Type: token.NEWLINE, Literal: "\n", Line: 35, Column: 36},
		{Type: token.IDENT, Literal: "a123", Line: 36, Column: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 36, Column: 4},
		{Type: token.MAP, Literal: "MAP", Line: 37, Column: 0},
		{Type: token.LBRACE, Literal: "{", Line: 37, Column: 4},

		{Type: token.NEWLINE, Literal: "\n", Line: 37, Column: 5},
		{Type: token.STRING, Literal: "a", Line: 38, Column: 4},
		{Type: token.COLON, Literal: ":", Line: 38, Column: 7},
		{Type: token.INT, Literal: "10", Line: 38, Column: 9},
		{Type: token.COMMA, Literal: ",", Line: 38, Column: 11},

		{Type: token.NEWLINE, Literal: "\n", Line: 38, Column: 12},
		{Type: token.RBRACE, Literal: "}", Line: 39, Column: 0},

		{Type: token.EOF, Literal: "", Line: 39, Column: 0},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.Type, tok.Type, "token type wrong for token %s, expecting %s", tok, tt.String())
		assert.Equal(t, tt.Literal, tok.Literal, "token literal wrong for token %s, expecting %s", tok, tt.String())
		assert.Equal(t, tt.Line, tok.Line, "token line number wrong for token %s, expecting %s", tok, tt.String())
		assert.Equal(t, tt.Column, tok.Column, "token column number wrong for token %s, expecting %s", tok, tt.String())
	}

	assert.Equal(t, byte(0), l.peekChar(), "lexer should have read all input before tests finish, not enough test cases")
}

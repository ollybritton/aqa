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
		{Type: token.IDENT, Literal: "five", Line: 0, StartCol: 0, EndCol: 3},
		{Type: token.ASSIGN, Literal: "<-", Line: 0, StartCol: 5, EndCol: 6},
		{Type: token.INT, Literal: "5", Line: 0, StartCol: 8, EndCol: 8},
		{Type: token.NEWLINE, Literal: "\n", Line: 0, StartCol: 9, EndCol: 9},
		{Type: token.IDENT, Literal: "ten", Line: 1, StartCol: 0, EndCol: 2},
		{Type: token.ASSIGN, Literal: "<-", Line: 1, StartCol: 4, EndCol: 5},
		{Type: token.INT, Literal: "10", Line: 1, StartCol: 7, EndCol: 8},
		{Type: token.NEWLINE, Literal: "\n", Line: 1, StartCol: 9, EndCol: 9},
		{Type: token.BANG, Literal: "!", Line: 2, StartCol: 0, EndCol: 0},
		{Type: token.PLUS, Literal: "+", Line: 2, StartCol: 1, EndCol: 1},
		{Type: token.MINUS, Literal: "-", Line: 2, StartCol: 2, EndCol: 2},
		{Type: token.ASTERISK, Literal: "*", Line: 2, StartCol: 3, EndCol: 3},
		{Type: token.SLASH, Literal: "/", Line: 2, StartCol: 4, EndCol: 4},
		{Type: token.INT, Literal: "5", Line: 2, StartCol: 5, EndCol: 5},
		{Type: token.NEWLINE, Literal: "\n", Line: 2, StartCol: 6, EndCol: 6},
		{Type: token.NOT_EQ, Literal: "!=", Line: 3, StartCol: 0, EndCol: 1},
		{Type: token.EQ, Literal: "==", Line: 3, StartCol: 3, EndCol: 4},
		{Type: token.EQ, Literal: "=", Line: 3, StartCol: 6, EndCol: 6},
		{Type: token.ILLEGAL, Literal: "@", Line: 3, StartCol: 8, EndCol: 8},
		{Type: token.NEWLINE, Literal: "\n", Line: 3, StartCol: 9, EndCol: 9},
		{Type: token.INT, Literal: "5", Line: 4, StartCol: 0, EndCol: 0},
		{Type: token.LT, Literal: "<", Line: 4, StartCol: 2, EndCol: 2},
		{Type: token.INT, Literal: "10", Line: 4, StartCol: 4, EndCol: 5},
		{Type: token.GT, Literal: ">", Line: 4, StartCol: 7, EndCol: 7},
		{Type: token.INT, Literal: "5", Line: 4, StartCol: 9, EndCol: 9},
		{Type: token.NEWLINE, Literal: "\n", Line: 4, StartCol: 10, EndCol: 10},
		{Type: token.LPAREN, Literal: "(", Line: 5, StartCol: 0, EndCol: 0},
		{Type: token.LBRACE, Literal: "{", Line: 5, StartCol: 1, EndCol: 1},
		{Type: token.LBRACKET, Literal: "[", Line: 5, StartCol: 2, EndCol: 2},
		{Type: token.INT, Literal: "1", Line: 5, StartCol: 3, EndCol: 3},
		{Type: token.COMMA, Literal: ",", Line: 5, StartCol: 4, EndCol: 4},
		{Type: token.IDENT, Literal: "potato", Line: 5, StartCol: 6, EndCol: 11},
		{Type: token.COMMA, Literal: ",", Line: 5, StartCol: 12, EndCol: 12},
		{Type: token.INT, Literal: "3", Line: 5, StartCol: 14, EndCol: 14},
		{Type: token.RBRACKET, Literal: "]", Line: 5, StartCol: 15, EndCol: 15},
		{Type: token.RBRACE, Literal: "}", Line: 5, StartCol: 16, EndCol: 16},
		{Type: token.RPAREN, Literal: ")", Line: 5, StartCol: 17, EndCol: 17},
		{Type: token.NEWLINE, Literal: "\n", Line: 5, StartCol: 18, EndCol: 18},

		{Type: token.STRING, Literal: "foobar", Line: 7, StartCol: 0, EndCol: 7},
		{Type: token.NEWLINE, Literal: "\n", Line: 7, StartCol: 8},
		{Type: token.STRING, Literal: "foo bar", Line: 8, StartCol: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 8, StartCol: 9},
		{Type: token.STRING, Literal: "foobar", Line: 9, StartCol: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 9, StartCol: 8},
		{Type: token.STRING, Literal: "foo bar", Line: 10, StartCol: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 10, StartCol: 9},
		{Type: token.STRING, Literal: "foo's bar", Line: 11, StartCol: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 11, StartCol: 12},
		{Type: token.STRING, Literal: "foo\"s bar", Line: 12, StartCol: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 12, StartCol: 12},
		{Type: token.WHILE, Literal: "WHILE", Line: 13, StartCol: 0, EndCol: 4},
		{Type: token.TRUE, Literal: "true", Line: 13, StartCol: 6},
		{Type: token.NEWLINE, Literal: "\n", Line: 13, StartCol: 10},
		{Type: token.IDENT, Literal: "a", Line: 14, StartCol: 1},
		{Type: token.NEWLINE, Literal: "\n", Line: 14, StartCol: 2},
		{Type: token.ENDWHILE, Literal: "ENDWHILE", Line: 15, StartCol: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 15, StartCol: 8},
		{Type: token.REPEAT, Literal: "REPEAT", Line: 16, StartCol: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 16, StartCol: 6},
		{Type: token.IDENT, Literal: "a", Line: 17, StartCol: 1},
		{Type: token.NEWLINE, Literal: "\n", Line: 17, StartCol: 2},
		{Type: token.UNTIL, Literal: "UNTIL", Line: 18, StartCol: 0},
		{Type: token.FALSE, Literal: "false", Line: 18, StartCol: 6},

		{Type: token.NEWLINE, Literal: "\n", Line: 18, StartCol: 11},
		{Type: token.FOR, Literal: "FOR", Line: 19, StartCol: 0},
		{Type: token.IDENT, Literal: "i", Line: 19, StartCol: 4},
		{Type: token.ASSIGN, Literal: "<-", Line: 19, StartCol: 6, EndCol: 7},
		{Type: token.INT, Literal: "10", Line: 19, StartCol: 9},
		{Type: token.TO, Literal: "TO", Line: 19, StartCol: 12},
		{Type: token.INT, Literal: "20", Line: 19, StartCol: 15},
		{Type: token.NEWLINE, Literal: "\n", Line: 19, StartCol: 17},
		{Type: token.IDENT, Literal: "a", Line: 20, StartCol: 1},
		{Type: token.NEWLINE, Literal: "\n", Line: 20, StartCol: 2},
		{Type: token.ENDFOR, Literal: "ENDFOR", Line: 21, StartCol: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 21, StartCol: 6},
		// ...skip comments...
		{Type: token.OUTPUT, Literal: "OUTPUT", Line: 24, StartCol: 0},
		{Type: token.USERINPUT, Literal: "USERINPUT", Line: 24, StartCol: 7},

		{Type: token.NEWLINE, Literal: "\n", Line: 24, StartCol: 16},
		{Type: token.FLOAT, Literal: "1.234", Line: 25, StartCol: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 25, StartCol: 5},
		{Type: token.INT, Literal: "0xBADA55", Line: 26, StartCol: 0, EndCol: 7},
		{Type: token.NEWLINE, Literal: "\n", Line: 26, StartCol: 8},
		{Type: token.INT, Literal: "0b10", Line: 27, StartCol: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 27, StartCol: 4},
		{Type: token.LBRACKET, Literal: "[", Line: 28, StartCol: 0},
		{Type: token.INT, Literal: "1", Line: 28, StartCol: 1},
		{Type: token.COMMA, Literal: ",", Line: 28, StartCol: 2},
		{Type: token.INT, Literal: "2", Line: 28, StartCol: 4},
		{Type: token.RBRACKET, Literal: "]", Line: 28, StartCol: 5},

		{Type: token.NEWLINE, Literal: "\n", Line: 28, StartCol: 6},
		{Type: token.INT, Literal: "1", Line: 29, StartCol: 0},
		{Type: token.LT_EQ, Literal: "<=", Line: 29, StartCol: 1},
		{Type: token.INT, Literal: "2", Line: 29, StartCol: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 29, StartCol: 4},
		{Type: token.INT, Literal: "1", Line: 30, StartCol: 0},
		{Type: token.GT_EQ, Literal: ">=", Line: 30, StartCol: 1},
		{Type: token.INT, Literal: "2", Line: 30, StartCol: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 30, StartCol: 4},
		{Type: token.INT, Literal: "1", Line: 31, StartCol: 0},
		{Type: token.RSHIFT, Literal: ">>", Line: 31, StartCol: 1},
		{Type: token.INT, Literal: "2", Line: 31, StartCol: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 31, StartCol: 4},
		{Type: token.INT, Literal: "2", Line: 32, StartCol: 0},
		{Type: token.LSHIFT, Literal: "<<", Line: 32, StartCol: 1},
		{Type: token.INT, Literal: "1", Line: 32, StartCol: 3},

		{Type: token.NEWLINE, Literal: "\n", Line: 32, StartCol: 4},
		{Type: token.INT, Literal: "1", Line: 33, StartCol: 0},
		{Type: token.MOD, Literal: "MOD", Line: 33, StartCol: 2},
		{Type: token.INT, Literal: "2", Line: 33, StartCol: 6},

		{Type: token.NEWLINE, Literal: "\n", Line: 33, StartCol: 7},
		{Type: token.INT, Literal: "2", Line: 34, StartCol: 0},
		{Type: token.DIV, Literal: "DIV", Line: 34, StartCol: 2},
		{Type: token.INT, Literal: "3", Line: 34, StartCol: 6},

		{Type: token.NEWLINE, Literal: "\n", Line: 34, StartCol: 7},
		{Type: token.NOT, Literal: "NOT", Line: 35, StartCol: 0},
		{Type: token.TRUE, Literal: "true", Line: 35, StartCol: 4},
		{Type: token.OR, Literal: "OR", Line: 35, StartCol: 9},
		{Type: token.FALSE, Literal: "false", Line: 35, StartCol: 12},
		{Type: token.AND, Literal: "AND", Line: 35, StartCol: 18},
		{Type: token.FALSE, Literal: "false", Line: 35, StartCol: 22},
		{Type: token.XOR, Literal: "XOR", Line: 35, StartCol: 28},
		{Type: token.TRUE, Literal: "true", Line: 35, StartCol: 32},

		{Type: token.NEWLINE, Literal: "\n", Line: 35, StartCol: 36},
		{Type: token.IDENT, Literal: "a123", Line: 36, StartCol: 0},

		{Type: token.NEWLINE, Literal: "\n", Line: 36, StartCol: 4},
		{Type: token.MAP, Literal: "MAP", Line: 37, StartCol: 0},
		{Type: token.LBRACE, Literal: "{", Line: 37, StartCol: 4},

		{Type: token.NEWLINE, Literal: "\n", Line: 37, StartCol: 5},
		{Type: token.STRING, Literal: "a", Line: 38, StartCol: 4},
		{Type: token.COLON, Literal: ":", Line: 38, StartCol: 7},
		{Type: token.INT, Literal: "10", Line: 38, StartCol: 9},
		{Type: token.COMMA, Literal: ",", Line: 38, StartCol: 11},

		{Type: token.NEWLINE, Literal: "\n", Line: 38, StartCol: 12},
		{Type: token.RBRACE, Literal: "}", Line: 39, StartCol: 0},

		{Type: token.EOF, Literal: "", Line: 39, StartCol: 0},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.Type, tok.Type, "token type wrong for token %s, expecting %s", tok, tt.String())
		assert.Equal(t, tt.Literal, tok.Literal, "token literal wrong for token %s, expecting %s", tok, tt.String())
		assert.Equal(t, tt.Line, tok.Line, "token line number wrong for token %s, expecting %s", tok, tt.String())
		assert.Equal(t, tt.StartCol, tok.StartCol, "token StartCol number wrong for token %s, expecting %s", tok, tt.String())

		// TODO: Finish adding test cases for endline. Remove this zero check to see if the struct has been
		// given a EndCol value.
		if tt.EndCol != 0 {
			assert.Equal(t, tt.EndCol, tok.EndCol, "token EndCol number wrong for token %s, expecting %s", tok, tt.String())
		}
	}

	assert.Equal(t, byte(0), l.peekChar(), "lexer should have read all input before tests finish, not enough test cases")
}

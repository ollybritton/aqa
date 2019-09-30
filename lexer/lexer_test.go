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
"foobar"
"foo bar"
'foobar'
'foo bar'
'foo\'s bar'
"foo\"s bar"`

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
		{Type: token.ILLEGAL, Literal: "=", Line: 3, Column: 6},
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
		{Type: token.STRING, Literal: "foobar", Line: 6, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 6, Column: 8},
		{Type: token.STRING, Literal: "foo bar", Line: 7, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 7, Column: 9},
		{Type: token.STRING, Literal: "foobar", Line: 8, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 8, Column: 8},
		{Type: token.STRING, Literal: "foo bar", Line: 9, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 9, Column: 9},
		{Type: token.STRING, Literal: "foo's bar", Line: 10, Column: 0},
		{Type: token.NEWLINE, Literal: "\n", Line: 10, Column: 12},
		{Type: token.STRING, Literal: "foo\"s bar", Line: 11, Column: 0},
		{Type: token.EOF, Literal: "", Line: 11, Column: 11},
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

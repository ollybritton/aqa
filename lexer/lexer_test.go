package lexer

import (
	"testing"

	"github.com/ollybritton/aqa++/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `five <- 5
ten <- 10
!+-*/5
!= ==
5 < 10 > 5`

	tests := []token.Token{
		{Type: token.IDENT, Literal: "five", Line: 0, Column: 0},
		{Type: token.ASSIGN, Literal: "<-", Line: 0, Column: 5},
		{Type: token.INT, Literal: "5", Line: 0, Column: 8},
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

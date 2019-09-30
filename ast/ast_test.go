package ast

import (
	"testing"

	"github.com/ollybritton/aqa/token"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VariableAssignment{
				Tok: token.NewToken(token.IDENT, "a", 0, 0),
				Name: &Identifier{
					Tok:   token.NewToken(token.IDENT, "a", 0, 0),
					Value: "a",
				},
				Value: &Identifier{
					Tok:   token.NewToken(token.IDENT, "b", 0, 0),
					Value: "b",
				},
			},
		},
	}

	expected := `a <- b`

	assert.Equal(t, expected, program.String(), "expected program string and actual program string differ")
}

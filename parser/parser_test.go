package parser

import (
	"testing"

	"github.com/ollybritton/aqa++/ast"
	"github.com/ollybritton/aqa++/lexer"
)

// private test methods/utilities.
func checkParserErrors(t *testing.T, p *Parser) bool {
	errors := p.Errors()

	if len(errors) == 0 {
		return true
	}

	t.Errorf("parser has %d errors:", len(errors))
	for i, err := range errors {
		t.Errorf("parser error <%d>: %v", i+1, err)
	}

	t.FailNow()

	return false
}

func parseProgram(t *testing.T, input string) (*Parser, *ast.Program) {
	l := lexer.New(input)
	p := New(l)

	program := p.Parse()
	if program == nil {
		t.Fatalf(".Parse() returned nil")
	}

	checkParserErrors(t, p)

	return p, program
}

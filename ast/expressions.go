package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ollybritton/aqa++/token"
)

// Identifier represents an identifier in the AST. Idents are expressions because they produce values (the value they represent)
// Example: `a`
// General: `{ident}`
type Identifier struct {
	Tok   token.Token // the token.IDENT token.
	Value string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) Token() token.Token { return i.Tok }
func (i *Identifier) String() string {
	return i.Value
}

// IntegerLiteral represents an integer value in the AST.
// Example: `5`
// General: `{token.INT}`
type IntegerLiteral struct {
	Tok   token.Token // the token.INT token.
	Value int64
}

func (il *IntegerLiteral) expressionNode()    {}
func (il *IntegerLiteral) Token() token.Token { return il.Tok }
func (il *IntegerLiteral) String() string {
	return fmt.Sprint(il.Value)
}

// BooleanLiteral represents a boolean in the AST.
// Example: `true`
// General: `{token.TRUE or token.FALSE}`
type BooleanLiteral struct {
	Tok   token.Token // the boolean token (token.TRUE or token.FALSE)
	Value bool
}

func (bl *BooleanLiteral) expressionNode()    {}
func (bl *BooleanLiteral) Token() token.Token { return bl.Tok }
func (bl *BooleanLiteral) String() string {
	return fmt.Sprint(bl.Value)
}

// PrefixExpression represents an expression involving a prefix operator.
// Example: `-10`
// General: `{- or !}{expression}`
type PrefixExpression struct {
	Tok      token.Token // the token of the prefix operator.
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()    {}
func (pe *PrefixExpression) Token() token.Token { return pe.Tok }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an expression involving an operator sandwiched between two other expressions.
// Example: `10-5`
// General: `{expression}{opeator}{expression}`
type InfixExpression struct {
	Tok token.Token

	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()    {}
func (ie *InfixExpression) Token() token.Token { return ie.Tok }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// SubroutineCall represents a call to a subroutine within the AST.
// Example: `add(1,2)`
// General: `{IDENT}({expression}, {expression}...)`
type SubroutineCall struct {
	Tok        token.Token // The '(' token
	Subroutine *Identifier
	Arguments  []Expression
}

func (sc *SubroutineCall) expressionNode()    {}
func (sc *SubroutineCall) Token() token.Token { return sc.Tok }
func (sc *SubroutineCall) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range sc.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(sc.Subroutine.Value)
	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}

package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ollybritton/aqa/token"
)

// Identifier represents an identifier in the AST. Idents are expressions because they produce values (the value they represent)
// Example: `a`
// General: `{ident}`
type Identifier struct {
	Tok      token.Token // the token.IDENT token.
	Constant bool
	Value    string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) Token() token.Token { return i.Tok }
func (i *Identifier) String() string {
	var out bytes.Buffer

	if i.Constant {
		out.WriteString("constant ")
	}

	out.WriteString(i.Value)

	return out.String()
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

// FloatLiteral represents an float value in the AST.
// Example: `5.5`
// General: `{token.INT}`
type FloatLiteral struct {
	Tok   token.Token // the token.FLOAT token.
	Value float64
}

func (fl *FloatLiteral) expressionNode()    {}
func (fl *FloatLiteral) Token() token.Token { return fl.Tok }
func (fl *FloatLiteral) String() string {
	return fmt.Sprint(fl.Value)
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
	Subroutine Expression
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

	out.WriteString(sc.Subroutine.String())
	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}

// StringLiteral represents a string inside the program.
// Example: `"hello"`, `'Dave\'s mom was sad'`
// General: `{'|"}{characters}{'|"}`
type StringLiteral struct {
	Tok   token.Token // The token.STRING token.
	Value string
}

func (sl *StringLiteral) expressionNode()    {}
func (sl *StringLiteral) Token() token.Token { return sl.Tok }
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("\"")
	out.WriteString(sl.Value)
	out.WriteString("\"")

	return out.String()
}

// ArrayLiteral represents an array inside the AST.
type ArrayLiteral struct {
	Tok      token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()    {}
func (al *ArrayLiteral) Token() token.Token { return al.Tok }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// IndexExpression represents an access to an array or map within the AST.
type IndexExpression struct {
	Tok   token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()    {}
func (ie *IndexExpression) Token() token.Token { return ie.Tok }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")

	return out.String()
}

// HashLiteral represents a hashmap inside the AST.
type HashLiteral struct {
	Tok   token.Token // The token.MAP token.
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()    {}
func (hl *HashLiteral) Token() token.Token { return hl.Tok }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// DotExpression is the use of the dot (.) operator on two operator. It gets the member associated object
// from a module.
type DotExpression struct {
	Tok    token.Token // the token.DOT token.
	Parent Identifier
	Child  Identifier
}

func (de *DotExpression) expressionNode()    {}
func (de *DotExpression) Token() token.Token { return de.Tok }
func (de *DotExpression) String() string {
	var out bytes.Buffer

	out.WriteString(de.Parent.Value)
	out.WriteString(".")
	out.WriteString(de.Child.Value)

	return out.String()
}

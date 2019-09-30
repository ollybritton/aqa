package ast

import (
	"bytes"
	"strings"

	"github.com/ollybritton/aqa/token"
)

// VariableAssignment represents the process of assignment to a variable in the AST.
// Example: `a <- 10`.
// General: `{ident} <- {expression}`
type VariableAssignment struct {
	Tok   token.Token // the token.ASSIGN token.
	Name  *Identifier
	Value Expression
}

func (va *VariableAssignment) statementNode()     {}
func (va *VariableAssignment) Token() token.Token { return va.Tok }
func (va *VariableAssignment) String() string {
	var out bytes.Buffer

	out.WriteString(va.Name.String())
	out.WriteString(" <- ")
	out.WriteString(va.Value.String())

	return out.String()
}

// ReturnStatement represents a return statement from a function or subroutine within a program.
// Example: `return a`
// General: `return {expression}`
type ReturnStatement struct {
	Tok         token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()     {}
func (rs *ReturnStatement) Token() token.Token { return rs.Tok }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString("return ")
	out.WriteString(rs.ReturnValue.String())

	return out.String()
}

// ExpressionStatement is a single expression by itself on one line.
// Example: `{start} a+10 {end}` (where start & end are the start and end of the line)
// General: `{start}{expression}{end}`
type ExpressionStatement struct {
	Tok        token.Token // The first token of the expression.
	Expression Expression
}

func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) Token() token.Token { return es.Tok }
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

// BlockStatement is a series of statements, wrapped in a block.
// Example: `IF 1 == 1 THEN {I'm a block statement} ENDIF`
// General: `{START}{list of statements}{END}`
type BlockStatement struct {
	Tok        token.Token // the start token, such as token.THEN
	Statements []Statement
}

func (bs *BlockStatement) statementNode()     {}
func (bs *BlockStatement) Token() token.Token { return bs.Tok }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// IfStatement represents an if-elseif-else statement inside the AST.
// Example:
//   IF a == 0 THEN
//     RETURN a
//   ELSE IF a < 0 THEN
//     a <- a + 1
//   ELSE
//     a <- a - 1
//   ENDIF
// General: IF {Expression} THEN {Statements} ELSE IF {Expression} THEN {STATEMENTS} ELSE {STATEMENTS} ENDIF
// The ELSE & ELSE IF are optional.
type IfStatement struct {
	Tok       token.Token // the token.IF token.
	Condition Expression

	Consequence *BlockStatement
	Else        *BlockStatement
	ElseIf      *IfStatement // Many else-ifs are represented as the .ElseIf of this .ElseIf
}

func (is *IfStatement) statementNode()     {}
func (is *IfStatement) Token() token.Token { return is.Tok }
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("IF")
	out.WriteString(is.Condition.String())
	out.WriteString(" THEN\n")

	for _, s := range is.Consequence.Statements {
		out.WriteString("  " + s.String())
	}

	current := is.ElseIf
	for current != nil {
		out.WriteString("\nELSE IF ")
		out.WriteString(current.Condition.String())
		out.WriteString(" THEN\n")

		for _, s := range current.Consequence.Statements {
			out.WriteString("  " + s.String())
		}

		current = current.ElseIf
	}

	if is.Else != nil {
		out.WriteString("\nELSE\n")

		for _, s := range is.Else.Statements {
			out.WriteString("  " + s.String())
		}
	}

	out.WriteString("\nENDIF")

	return out.String()
}

// Subroutine represents a subroutine inside the program. Subroutines are here for compliance with the spec, and are not
// expressions like FUNC will be.
// Example:
//   SUBROUTINE
//   show_add(a, b)
//     result <- a + b
//     OUTPUT result
//   ENDSUBROUTINE
// General:
//   SUBROUTINE
//   {ident}({ident}, {ident}...)
//     {statements}
//   ENDSUBROUTINE
type Subroutine struct {
	Tok        token.Token // the token.SUBROUTINE token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (s *Subroutine) statementNode()     {}
func (s *Subroutine) Token() token.Token { return s.Tok }
func (s *Subroutine) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range s.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(s.Tok.Literal + "\n")
	out.WriteString(s.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString("\n")

	for _, s := range s.Body.Statements {
		out.WriteString("  " + s.String() + "\n")
	}

	out.WriteString("ENDSUBROUTINE")

	return out.String()
}

package object

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/ollybritton/aqa/ast"
)

// Type represents a type of object, such as an integer or a subroutine.
type Type string

// Definition of object types.
const (
	INTEGER_OBJ  = "INTEGER"
	FLOAT_OBJ    = "FLOAT"
	BOOLEAN_OBJ  = "BOOLEAN"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ   = "STRING"
	ARRAY_OBJ    = "ARRAY"
	HASH_OBJ     = "HASH"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
	NULL_OBJ         = "NULL"
)

// Object is an interface which allows different objects to be represented.
type Object interface {
	Type() Type      // Type reveals an object's type
	Inspect() string // Inspect gets the value of the object as a string.
}

// BuiltinFunction represents an external function that is avaliable inside an AQA++ program.
type BuiltinFunction func(args ...Object) Object

// Builtin represents a builtin inside the program.
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type      { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string { return "<builtin>" }

// Integer represents an integer within the program.
type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Float represents an Float within the program.
type Float struct {
	Value float64
}

func (f *Float) Type() Type      { return FLOAT_OBJ }
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

// Boolean represents a boolean value, such as true or false, within the program.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null represents the lack/absence of a value. It is like nil.
type Null struct{}

func (n *Null) Type() Type      { return NULL_OBJ }
func (n *Null) Inspect() string { return "null" }

// ReturnValue represents a value that is being returned from a subroutine or from a program as a whole.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type      { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error represents an error that occurs during the evalutation of the programming language.
type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Subroutine represents a subroutine within the evaluator.
type Subroutine struct {
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (s *Subroutine) Type() Type { return FUNCTION_OBJ }
func (s *Subroutine) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range s.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("SUBROUTINE\n")
	out.WriteString(s.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")\n")

	for _, stmt := range s.Body.Statements {
		out.WriteString("  " + stmt.String())
	}
	out.WriteString("\n")

	out.WriteString("ENDSUBROUTINE")

	return out.String()
}

// String represents a string within the evaluator.
type String struct {
	Value string
}

func (s *String) Type() Type      { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }

// Array represents an arrau within the evaluator.
type Array struct {
	Elements []Object
}

func (a *Array) Type() Type { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

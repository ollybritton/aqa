package object

import "fmt"

// Type represents a type of object, such as an integer or a subroutine.
type Type string

// Definition of object types.
const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"

	RETURN_OBJ = "RETURN_VALUE"
	NULL_OBJ   = "NULL"
)

// Object is an interface which allows different objects to be represented.
type Object interface {
	Type() Type      // Type reveals an object's type
	Inspect() string // Inspect gets the value of the object as a string.
}

// Integer represents an integer within the program.
type Integer struct {
	Value int64
}

func (i *Integer) Object() Type    { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean represents a boolean value, such as true or false, within the program.
type Boolean struct {
	Value bool
}

func (b *Boolean) Object() Type    { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null represents the lack/absence of a value. It is like nil.
type Null struct{}

func (n *Null) Object() Type    { return NULL_OBJ }
func (n *Null) Inspect() string { return "null" }

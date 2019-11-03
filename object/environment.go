package object

import (
	"fmt"
)

// Environment represents the variables and identifiers inside their program, mapped to their actual Object values.
type Environment struct {
	store     map[string]Object
	constants map[string]Object
	outer     *Environment
	imported  map[string]*Environment
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]Object)
	i := make(map[string]*Environment)
	return &Environment{store: s, constants: c, imported: i, outer: nil}
}

// NewEnclosedEnvironment creates a new enclosed environment, extending from a previous.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

// Get gets an object by name.
func (e *Environment) Get(name string) (Object, bool) {

	obj, ok := e.store[name]
	if !ok {
		obj, ok = e.constants[name]
		if !ok && e.outer != nil {
			obj, ok = e.outer.Get(name)
		}
	}

	return obj, ok
}

// Set sets an object by name.
func (e *Environment) Set(name string, value Object) Object {

	if _, ok := e.constants[name]; ok {
		return &Error{Message: fmt.Sprintf("cannot assign to constant %s", name)}
	}

	e.store[name] = value
	return value
}

// Keys gets the list of all symbols.
func (e *Environment) Keys() map[string]bool {
	symbols := make(map[string]bool)

	for k := range e.store {
		symbols[k] = true
	}

	for k := range e.constants {
		symbols[k] = true
	}

	return symbols
}

// SetConstant sets a constant.
func (e *Environment) SetConstant(name string, value Object) Object {
	if _, ok := e.constants[name]; ok {
		return &Error{Message: fmt.Sprintf("cannot assign to constant %s", name)}
	}

	e.constants[name] = value
	return value
}

// AddEnvironment adds the objects from one environment into another.
func (e *Environment) AddEnvironment(env *Environment) {
	for ident, value := range env.store {
		e.store[ident] = value
	}

	for ident, value := range env.constants {
		e.constants[ident] = value
	}
}

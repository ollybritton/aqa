package object

import (
	"fmt"
)

// Environment represents the variables and identifiers inside their program, mapped to their actual Object values.
type Environment struct {
	store     map[string]Object
	constants map[string]Object
	outer     *Environment
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]Object)
	return &Environment{store: s, constants: c, outer: nil}
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

// SetConstant sets a constant.
func (e *Environment) SetConstant(name string, value Object) Object {
	if _, ok := e.constants[name]; ok {
		return &Error{Message: fmt.Sprintf("cannot assign to constant %s", name)}
	}

	e.constants[name] = value
	return value
}

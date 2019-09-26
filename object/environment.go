package object

// Environment represents the variables and identifiers inside their program, mapped to their actual Object values.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
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
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set sets an object by name.
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

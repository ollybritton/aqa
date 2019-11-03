package evaluator

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/object"
	"github.com/ollybritton/aqa/parser"
)

// EvalString will execute a string of aqa++ code.
func EvalString(str string, env *object.Environment) (object.Object, []error) {
	l := lexer.New(str)
	p := parser.New(l)

	program := p.Parse()
	if len(p.Errors()) != 0 {
		return &object.Null{}, p.Errors()
	}

	eval := Eval(program, env)
	if eval == nil {
		return &object.Null{}, []error{}
	}

	if eval.Type() == object.ERROR_OBJ {
		return &object.Null{}, []error{errors.New(eval.Inspect())}
	}

	return eval, []error{}
}

// EvalFile will execute a file containing aqa++ code.
func EvalFile(f *os.File, env *object.Environment) (object.Object, []error) {
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return &object.Null{}, []error{err}
	}

	return EvalString(string(bytes), env)
}

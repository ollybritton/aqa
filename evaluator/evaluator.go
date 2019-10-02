package evaluator

import (
	"bufio"
	"os"

	"github.com/ollybritton/aqa/ast"
	"github.com/ollybritton/aqa/object"
)

// Eval evaluates a node, and returns its representation as an object.Object.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfStatement:
		return evalIfStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	case *ast.RepeatStatement:
		return evalRepeatStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}

	case *ast.VariableAssignment:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

	case *ast.Subroutine:
		params := node.Parameters
		body := node.Body
		name := node.Name

		env.Set(name.Value, &object.Subroutine{Parameters: params, Env: env, Body: body, Name: name})

	case *ast.SubroutineCall:
		ident := evalIdentifier(node.Subroutine, env)
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applySubroutine(ident, args)

	// Literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	// Expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)

		if isError(left) {
			return left
		} else if isError(right) {
			return right
		}

		return evalInfixExpression(left, node.Operator, right)

	case *ast.Identifier:
		return evalIdentifier(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()

			if rt == object.ERROR_OBJ || rt == object.RETURN_VALUE_OBJ {
				return result
			}
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		// TODO: I want to avoid null, should error here
		return TRUE
	default:
		// TODO: I don't want the bang operator to work on things that aren't booleans, should error
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(left, operator, right)

	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(left, operator, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(left, operator, right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfStatement(node *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	}

	if node.ElseIf != nil {
		elseIf := evalIfStatement(node.ElseIf, env)
		if isError(elseIf) {
			return elseIf
		}

		if elseIf != NULL {
			return elseIf
		}
	}

	if node.Else != nil {
		return Eval(node.Else, env)
	}

	return NULL
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	if node.Value == "USERINPUT" || node.Value == "userinput" {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		return &object.String{Value: text}
	}

	return newError("identifier not found: " + node.Value)
}

func evalWhileStatement(node *ast.WhileStatement, env *object.Environment) object.Object {
	val := Eval(node.Condition, env)
	if isError(val) {
		return val
	}

	cond, ok := val.(*object.Boolean)
	if !ok {
		return newError("need a boolean for while loop, not %T", val)
	}

	var result object.Object

	for cond.Value {
		result = Eval(node.Body, env)
		if isError(result) {
			return result
		}

		cond, ok = Eval(node.Condition, env).(*object.Boolean)
		if !ok {
			return newError("need a boolean for while loop, not %T", cond)
		}
	}

	return result
}

func evalForStatement(node *ast.ForStatement, env *object.Environment) object.Object {
	extended := object.NewEnclosedEnvironment(env)

	lower, ok := Eval(node.Lower, env).(*object.Integer)
	if !ok {
		return newError("expected integer expression for `for` loop lower bound, got=%T", node.Lower)
	}

	upper, ok := Eval(node.Upper, env).(*object.Integer)
	if !ok {
		return newError("expected integer expression for `for` loop upper bounds, got=%T", node.Upper)
	}

	var val object.Object

	for i := lower.Value; i <= upper.Value; i++ {
		extended.Set(node.Ident.Value, &object.Integer{Value: i})
		val = Eval(node.Body, extended)

		if isError(val) {
			return val
		}
	}

	return val
}

func evalRepeatStatement(node *ast.RepeatStatement, env *object.Environment) object.Object {
	val := Eval(node.Condition, env)
	if isError(val) {
		return val
	}

	cond, ok := val.(*object.Boolean)
	if !ok {
		return newError("need a boolean for while loop, not %T", val)
	}

	var result object.Object

	for !cond.Value {
		result = Eval(node.Body, env)
		if isError(result) {
			return result
		}

		cond, ok = Eval(node.Condition, env).(*object.Boolean)
		if !ok {
			return newError("need a boolean for while loop, not %T", cond)
		}
	}

	return result
}

func applySubroutine(sub object.Object, args []object.Object) object.Object {
	switch sub := sub.(type) {
	case *object.Subroutine:
		extended := extendSubroutineEnv(sub, args)
		evaluated := Eval(sub.Body, extended)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return sub.Fn(args...)

	default:
		return newError("not a subroutine, function or builtin: %s", sub.Type())
	}
}

func extendSubroutineEnv(sub *object.Subroutine, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(sub.Env)

	for paramIDx, param := range sub.Parameters {
		env.Set(param.Value, args[paramIDx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

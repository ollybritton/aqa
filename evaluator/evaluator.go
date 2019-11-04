package evaluator

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ollybritton/aqa/ast"
	"github.com/ollybritton/aqa/builtins"
	"github.com/ollybritton/aqa/object"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

		if isBuiltin(node.Name.Value) {
			return newError("cannot assign to builtin: %s", node.Name.Value)
		}

		if node.Name.Constant {
			err := env.SetConstant(node.Name.Value, val)
			if isError(err) {
				return err
			}
		} else {
			err := env.Set(node.Name.Value, val)
			if isError(err) {
				return err
			}
		}

	case *ast.Subroutine:
		params := node.Parameters
		body := node.Body
		name := node.Name

		if isBuiltin(name.Value) {
			return newError("cannot assign to builtin: %s", node.Name.Value)
		}

		err := env.Set(name.Value, &object.Subroutine{Parameters: params, Env: env, Body: body, Name: name})
		if isError(err) {
			return err
		}

	case *ast.SubroutineCall:
		expression := Eval(node.Subroutine, env)
		if isError(expression) {
			return expression
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applySubroutine(expression, args)

	case *ast.ImportStatement:
		err := evalImport(node, env)

		if isError(err) {
			return err
		}

	// Literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	// Expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		if node.Operator == "." {
			ident, ok := node.Right.(*ast.Identifier)
			if !ok {
				return newError("right-hand side of dot expression is not an identifier, got %T.", node.Right)
			}

			return evalDotExpression(left, ident.Value)
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(left, node.Operator, right)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
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
	case "!", "NOT":
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
		return newError("unknown operator: !null")
	default:
		return newError("unknown operator: !%s", right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch val := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -val.Value}
	case *object.Float:
		return &object.Float{Value: -val.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	left, right = coerceInfix(left, operator, right)

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(left, operator, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(left, operator, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(left, operator, right)

	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(left, operator, right)

	case operator == "=" || operator == "==":
		if l, ok := left.(object.Hashable); ok {
			if r, ok := right.(object.Hashable); ok {
				return nativeBoolToBooleanObject(l.HashKey() == r.HashKey())
			}
		}

		fallthrough

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftInt := left.(*object.Integer)
	rightInt := right.(*object.Integer)

	switch operator {
	case "+":
		return &object.Integer{Value: leftInt.Value + rightInt.Value}
	case "-":
		return &object.Integer{Value: leftInt.Value - rightInt.Value}
	case "*":
		return &object.Integer{Value: leftInt.Value * rightInt.Value}
	case "/":
		if rightInt.Value == 0 {
			return newError("division error: division by zero")
		}

		if leftInt.Value%rightInt.Value == 0 {
			return &object.Integer{Value: leftInt.Value / rightInt.Value}
		}

		lf := object.IntegerToFloat(leftInt)
		rf := object.IntegerToFloat(rightInt)

		return &object.Float{Value: lf.Value / rf.Value}
	case ">>":
		if rightInt.Value < 0 {
			return newError("cannot perform bit shift using negative number: %d >> %d", leftInt.Value, rightInt.Value)
		}

		return &object.Integer{Value: leftInt.Value >> uint64(rightInt.Value)}
	case "<<":
		if rightInt.Value < 0 {
			return newError("cannot perform bit shift using negative number: %d << %d", leftInt.Value, rightInt.Value)
		}

		return &object.Integer{Value: leftInt.Value << uint64(rightInt.Value)}

	case "DIV", "div":
		return &object.Integer{Value: int64(
			math.Floor(float64(leftInt.Value / rightInt.Value)),
		)}

	case "MOD", "mod":
		return &object.Integer{Value: int64(
			leftInt.Value % rightInt.Value,
		)}

	case "==", "=":
		return nativeBoolToBooleanObject(leftInt.Value == rightInt.Value)
	case "!=":
		return nativeBoolToBooleanObject(leftInt.Value != rightInt.Value)
	case ">":
		return nativeBoolToBooleanObject(leftInt.Value > rightInt.Value)
	case "<":
		return nativeBoolToBooleanObject(leftInt.Value < rightInt.Value)
	case ">=":
		return nativeBoolToBooleanObject(leftInt.Value >= rightInt.Value)
	case "<=":
		return nativeBoolToBooleanObject(leftInt.Value <= rightInt.Value)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	lf := left.(*object.Float)
	rf := right.(*object.Float)

	switch operator {
	case "+":
		return &object.Float{Value: lf.Value + rf.Value}
	case "-":
		return &object.Float{Value: lf.Value - rf.Value}
	case "*":
		return &object.Float{Value: lf.Value * rf.Value}
	case "/":
		if rf.Value == 0 {
			return newError("division error: division by zero")
		}

		return &object.Float{Value: lf.Value / rf.Value}

	case "DIV":
		return &object.Integer{Value: int64(
			math.Floor(lf.Value / rf.Value),
		)}

	case "==", "=":
		return nativeBoolToBooleanObject(lf.Value == rf.Value)
	case "!=":
		return nativeBoolToBooleanObject(lf.Value != rf.Value)
	case ">":
		return nativeBoolToBooleanObject(lf.Value > rf.Value)
	case "<":
		return nativeBoolToBooleanObject(lf.Value < rf.Value)
	case ">=":
		return nativeBoolToBooleanObject(lf.Value >= rf.Value)
	case "<=":
		return nativeBoolToBooleanObject(lf.Value <= rf.Value)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==", "=":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	case "OR":
		return nativeBoolToBooleanObject(leftVal || rightVal)
	case "XOR":
		return nativeBoolToBooleanObject((leftVal || rightVal) && !(leftVal && rightVal))
	case "AND":
		return nativeBoolToBooleanObject(leftVal && rightVal)

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

	case "==", "=":
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		// Translate the strings into integers
		l, err := strconv.ParseInt(leftVal, 0, 64)
		if err != nil {
			break
		}

		r, err := strconv.ParseInt(rightVal, 0, 64)
		if err != nil {
			break
		}

		lobj := &object.Integer{Value: l}
		robj := &object.Integer{Value: r}

		return evalIntegerInfixExpression(lobj, operator, robj)
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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

	if builtin, ok := builtins.Builtins[strings.ToUpper(node.Value)]; ok {
		return builtin
	}

	if node.Value == "USERINPUT" || node.Value == "userinput" {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
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
		if isBuiltin(node.Ident.Value) {
			return newError("cannot assign to builtin: %s", node.Ident.Value)
		}

		err := env.Set(node.Ident.Value, &object.Integer{Value: i})
		if isError(err) {
			return err
		}

		val = Eval(node.Body, env)

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

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(left, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(array.Elements) - 1)

	if idx < 0 || idx > max {
		return newError("index out of bounds: %d", idx)
	}

	return array.Elements[idx]
}

func evalStringIndexExpression(left, index object.Object) object.Object {
	str := left.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return newError("index out of bounds: %d", idx)
	}

	return &object.String{Value: string(str.Value[idx])}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hask key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalDotExpression(parent object.Object, child string) object.Object {
	module, ok := parent.(*object.Module)
	if !ok {
		return newError("cannot use dot operator on %T object", parent)
	}

	val, exists := module.Env.Get(child)
	if !exists {
		return newError("unknown child %q in %s", child, module.Inspect())
	}

	if !module.Exposed[child] {
		return newError("unexposed child %q in %s", child, module.Inspect())
	}

	return val
}

func applySubroutine(sub object.Object, args []object.Object) object.Object {
	switch sub := sub.(type) {
	case *object.Subroutine:
		extended, err := extendSubroutineEnv(sub, args)
		if err != nil {
			return err
		}

		evaluated := Eval(sub.Body, extended)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return sub.Fn(args...)

	default:
		return newError("not a subroutine, function or builtin: %s", sub.Type())
	}
}

func extendSubroutineEnv(sub *object.Subroutine, args []object.Object) (*object.Environment, *object.Error) {
	env := object.NewEnclosedEnvironment(sub.Env)

	for paramIDx, param := range sub.Parameters {
		if isBuiltin(param.Value) {
			return &object.Environment{}, newError("cannot assign to builtin: %s", param.Value)
		}

		env.Set(param.Value, args[paramIDx])
	}

	return env, nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

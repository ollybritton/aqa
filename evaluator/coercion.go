package evaluator

import "github.com/ollybritton/aqa/object"

// coerceInfix will convert the types in an infix expression automatically so that objects can be used with one another
// without having to deal with converting types yourself.
//
// Rules:
// string + int/float => string + string
// int/float + string => string + string
// int, float => float & float
// float, int => float & float
func coerceInfix(left object.Object, operator string, right object.Object) (object.Object, object.Object) {
	switch {
	case left.Type() == object.STRING_OBJ && operator == "+" && right.Type() == object.FLOAT_OBJ:
		x := left.(*object.String)
		y := right.(*object.Float)

		return x, object.FloatToString(y)

	case left.Type() == object.STRING_OBJ && operator == "+" && right.Type() == object.INTEGER_OBJ:
		x := left.(*object.String)
		y := right.(*object.Integer)

		return x, object.IntegerToString(y)

	case left.Type() == object.FLOAT_OBJ && operator == "+" && right.Type() == object.STRING_OBJ:
		x := left.(*object.Float)
		y := right.(*object.String)

		return object.FloatToString(x), y

	case left.Type() == object.INTEGER_OBJ && operator == "+" && right.Type() == object.STRING_OBJ:
		x := left.(*object.Integer)
		y := right.(*object.String)

		return object.IntegerToString(x), y

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		leftStr := left.(*object.String)
		rightStr := right.(*object.String)

		var x object.Object
		var y object.Object
		var err error

		switch operator {
		case "+", "-", "*", "/", "MOD", "DIV", "<<", ">>":
			x, err = object.StringToInteger(leftStr)
			if err != nil {
				x, err = object.StringToFloat(leftStr)
				if err != nil {
					break
				}
			}

			y, err = object.StringToInteger(rightStr)
			if err != nil {
				y, err = object.StringToFloat(rightStr)
				if err != nil {
					break
				}
			}

			return coerceInfix(x, operator, y)
		}

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		x := left.(*object.Float)
		y := right.(*object.Integer)

		return x, object.IntegerToFloat(y)

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		x := left.(*object.Integer)
		y := right.(*object.Float)

		return object.IntegerToFloat(x), y

	default:
		return left, right
	}

	return left, right
}

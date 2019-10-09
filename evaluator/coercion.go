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
	if left.Type() == right.Type() {
		return left, right
	}

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
}

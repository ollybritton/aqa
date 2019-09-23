package evaluator

import (
	"github.com/ollybritton/monkey/object"
)

// Reused objects as to prevent allocations
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}

	NULL = &object.Null{}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	// TODO: get rid of this idead
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

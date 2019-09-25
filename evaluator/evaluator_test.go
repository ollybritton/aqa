package evaluator

import (
	"testing"

	"github.com/ollybritton/aqa++/lexer"
	"github.com/ollybritton/aqa++/parser"
	"github.com/ollybritton/monkey/object"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"false", false},
		{"FALSE", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"IF true THEN 10 ENDIF", 10},
		{"IF false THEN 10 ENDIF", nil},
		{"IF 1 THEN 10 ENDIF", 10},
		{"IF 1 < 2 THEN 10 ENDIF", 10},
		{"IF 1 > 2 THEN 10 ENDIF", nil},
		{"IF 1 > 2 THEN 10 ELSE 20 ENDIF", 20},
		{"IF 1 < 2 THEN 10 ELSE 20 ENDIF", 10},
		{"IF 1 == 0 THEN 1 ELSE IF 1 == 2 THEN 2 ELSE 3 ENDIF", 3},
		{"IF 1 == 0 THEN 1 ELSE IF 1 == 2 THEN 2 ELSE IF 1 == 1 THEN 4 ELSE 3 ENDIF", 4},
		{`IF false THEN 10 ELSE
		20
		ENDIF`, 20},
		{`IF false THEN 10 ELSE
			IF true THEN
				20
			ENDIF
		ENDIF`, 20},
		{`IF false THEN 10 ELSE
			IF false THEN
				20
			ELSE
				30
			ENDIF
		ENDIF`, 30},
		{`IF false THEN 10 ELSE
			IF true THEN
				20
			ELSE
				30
			ENDIF
		ENDIF`, 20},
		{`IF false THEN 10 ELSE IF 1 == 0 THEN IF 1 == 1 THEN 12 ENDIF ELSE 100 ENDIF`, 100},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10", 10},
		{`return 10
		9`, 10},
		{`return 2*5
		 9`, 10},
		{`9
		return 10
		9`, 10},
		{
			`IF 10 > 1 THEN
				IF 10 > 1 THEN
					return 10
				ENDIF

				return 1
			ENDIF`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if returnVal, ok := evaluated.(*object.ReturnValue); ok {
			testIntegerObject(t, returnVal.Value, tt.expected)
		}
	}

}

// private testing methods/functions
func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()

	if len(p.Errors()) != 0 {
		t.Errorf("parser errors:")
		t.Errorf(input)

		for _, e := range p.Errors() {
			t.Errorf("%T occured while parsing: %v", e, e)
		}

		t.FailNow()
	}

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("object is not Null. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

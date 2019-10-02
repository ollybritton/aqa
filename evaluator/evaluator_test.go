package evaluator

import (
	"testing"

	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/object"
	"github.com/ollybritton/aqa/parser"
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

func TestWhileStatement(t *testing.T) {
	input := `
a <- 0
WHILE a < 10
	a <- a + 1
ENDWHILE

return a`

	evaluated := testEval(t, input)
	if returnVal, ok := evaluated.(*object.ReturnValue); ok {
		testIntegerObject(t, returnVal.Value, 10)
	}
}

func TestRepeatStatement(t *testing.T) {
	input := `
a <- 0

REPEAT
	a <- a + 1
UNTIL a > 10

return a`

	evaluated := testEval(t, input)
	if returnVal, ok := evaluated.(*object.ReturnValue); ok {
		testIntegerObject(t, returnVal.Value, 10)
	}
}

func TestForStatement(t *testing.T) {
	input := `
a <- 0
FOR i <- 1 TO 10
	a <- a + i
ENDFOR

return a`

	evaluated := testEval(t, input)
	if returnVal, ok := evaluated.(*object.ReturnValue); ok {
		if val, ok := returnVal.Value.(*object.Integer); ok {
			testIntegerObject(t, val, 55)
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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`5 + true`,
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			`5 + true
			5`,
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`5
			true + false
			5`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if 10 > 1 then true + false endif",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
IF 10 > 1 THEN
	IF 10 > 1 THEN
		RETURN true + false
	ENDIF
	RETURN 1
ENDIF`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`foobar`,
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestVariableAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a <- 5\na", 5},
		{"a <- 5*5\na", 25},
		{"a <- 5\nb <- a\nb", 5},
		{"a <- 5\nb <- a\nc <- a+b+5\nc", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestSubroutines(t *testing.T) {
	add := `SUBROUTINE add(x,y)
return x+y
ENDSUBROUTINE
`

	identity := `SUBROUTINE identity(x)
return x
ENDSUBROUTINE
`

	double := `SUBROUTINE double(x)
return x*2
ENDSUBROUTINE
`

	tests := []struct {
		input    string
		expected int64
	}{
		{identity + "identity(5)", 5},
		{double + "double(5)", 10},
		{add + "add(1,2)", 3},
		{add + "add(-1, 5*4)", 19},
		{add + "add(add(1,2),add(5,3))", 11},
		{add + "add(5+5, add(5,5))", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not a String. got=%T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not a String. got=%T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`LEN("")`, 0},
		{`LEN("four")`, 4},
		{`LEN("hello world!")`, 12},
		{`LEN(1)`, "argument to `LEN` not supported, got=INTEGER"},
		{`LEN("one", "two")`, "wrong number of arguments. got=2, want=1"},

		{`POSITION("computer science", "m")`, 2},
		{`POSITION("comp")`, "wrong number of arguments. got=1, want=2"},
		{`POSITION(5, "ten")`, "argument to `POSITION` not supported, got=INTEGER"},

		{`SUBSTRING(2, 9, "computer science")`, "mputer s"},
		{`SUBSTRING("oops", "oops", 5)`, "argument 1 to `SUBSTRING` not supported, got=STRING"},
		{`SUBSTRING(2, "oops", 5)`, "argument 2 to `SUBSTRING` not supported, got=STRING"},
		{`SUBSTRING(2, 3, 5)`, "argument 3 to `SUBSTRING` not supported, got=INTEGER"},
		{`SUBSTRING(2, -1, "hello")`, "invalid bounds [2:-1] in call to SUBSTRING"},

		{`STRING_TO_INT('16')`, 16},
		{`STRING_TO_INT('16', '16')`, "wrong number of arguments. got=2, want=1"},
		{`STRING_TO_INT(10)`, "argument to `STRING_TO_INT` not supported, got=INTEGER"},
		{`STRING_TO_INT('abc')`, "failed to convert \"abc\" to integer in call to `STRING_TO_INT`"},

		{`INT_TO_STRING(16)`, "16"},
		{`INT_TO_STRING(16, 16)`, "wrong number of arguments. got=2, want=1"},
		{`INT_TO_STRING("10")`, "argument to `INT_TO_STRING` not supported, got=STRING"},

		{`CHAR_TO_CODE('a')`, 97},
		{`CHAR_TO_CODE('a', 'a')`, "wrong number of arguments. got=2, want=1"},
		{`CHAR_TO_CODE(97)`, "argument to `CHAR_TO_CODE` not supported, got=INTEGER"},
		{`CHAR_TO_CODE('abc')`, "argument to `CHAR_TO_CODE` not supported, cannot convert multiple characters, got=abc"},

		{`CODE_TO_CHAR(97)`, "a"},
		{`CODE_TO_CHAR(97, 97)`, "wrong number of arguments. got=2, want=1"},
		{`CODE_TO_CHAR("a")`, "argument to `CODE_TO_CHAR` not supported, got=STRING"},

		{`RANDOM_INT(1)`, "wrong number of arguments. got=1, want=2"},
		{`RANDOM_INT('a', 2)`, "argument 1 to `RANDOM_INT` not supported, got=STRING"},
		{`RANDOM_INT(2, 'a')`, "argument 2 to `RANDOM_INT` not supported, got=STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			switch evaluated := evaluated.(type) {
			case *object.Error:
				if evaluated.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, evaluated.Message)
				}
			case *object.String:
				if evaluated.Value != expected {
					t.Errorf("wrong string value. expected=%q, got=%q", expected, evaluated.Value)
				}
			}
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

	return Eval(program, object.NewEnvironment())
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

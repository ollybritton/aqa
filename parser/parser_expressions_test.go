package parser

import (
	"fmt"
	"testing"

	"github.com/ollybritton/aqa/ast"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"
	_, program := parseProgram(t, input)

	ok := assert.Equal(t, 1, len(program.Statements), "program should contain exactly 1 statement. got=%d", len(program.Statements))
	if !ok {
		t.FailNow()
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt)

	ident, ok := stmt.Expression.(*ast.Identifier)
	assert.True(t, ok, "expression not *ast.Identifier. got=%T", ident)

	assert.Equal(t, "foobar", ident.Value, "ident.Value does not equal 'foobar'")
	assert.Equal(t, "foobar", ident.Token().Literal, "ident.Tok.Literal does not equal 'foobar'")

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	_, program := parseProgram(t, input)

	ok := assert.Equal(t, 1, len(program.Statements), "program should contain exactly 1 statement. got=%d", len(program.Statements))
	if !ok {
		t.FailNow()
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	assert.Equal(t, int64(5), literal.Value, "literal.Value should equal 5")
	assert.Equal(t, "5", literal.Tok.Literal, "literal.Tok.Literal should equal '5'")

}

func TestFloatLiteralExpression(t *testing.T) {
	input := "5.5"

	_, program := parseProgram(t, input)

	ok := assert.Equal(t, 1, len(program.Statements), "program should contain exactly 1 statement. got=%d", len(program.Statements))
	if !ok {
		t.FailNow()
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FloatLiteral. got=%T", stmt.Expression)
	}

	assert.Equal(t, 5.5, literal.Value, "literal.Value should equal 5.5")
	assert.Equal(t, "5.5", literal.Tok.Literal, "literal.Tok.Literal should equal '5.5'")

}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true"

	_, program := parseProgram(t, input)

	ok := assert.Equal(t, 1, len(program.Statements), "program should contain exactly 1 statement. got=%d", len(program.Statements))
	if !ok {
		t.FailNow()
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Expression)
	}

	assert.Equal(t, true, literal.Value, "literal.Value should equal 5")
	assert.Equal(t, "true", literal.Tok.Literal, "literal.Tok.Literal should equal '5'")

}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		_, program := parseProgram(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 >= 5", 5, ">=", 5},
		{"5 <= 5", 5, "<=", 5},
		{"5 >> 5", 5, ">>", 5},
		{"5 << 5", 5, "<<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"5 DIV 5", 5, "DIV", 5},
		{"5 MOD 5", 5, "MOD", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		_, program := parseProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.InfixExpression. got=%T", exp)
		}

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
		{
			"NOT true OR false AND false XOR true",
			"(NOT(((true OR false) AND false) XOR true))",
		},
	}

	for i, tt := range tests {
		_, program := parseProgram(t, tt.input)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("<%d> input=%q :: expected=%q, got=%q", i, tt.input, tt.expected, actual)
		}
	}
}

func TestSubroutineCallParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	_, program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.SubroutineCall)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.SubroutineCall. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Subroutine, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong number of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `'hello\'s world!'`

	_, program := parseProgram(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", literal)
	}

	if literal.Value != "hello's world!" {
		t.Errorf("literal.Value not %q. got=%q", "hello's world!", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := `[1, 2*2, 3+3]`

	_, program := parseProgram(t, input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := `myArray[1+1]`

	_, program := parseProgram(t, input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralWithStringKeys(t *testing.T) {
	input := `MAP { 'one': 1, 'two': 2, 'three': 3 }`

	_, program := parseProgram(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	mapStmt, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral, got=%T", stmt.Expression)
	}

	if len(mapStmt.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d, want=3", len(mapStmt.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range mapStmt.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not *ast.StringLiteral. got=%T", literal)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, v, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := `MAP {}`

	_, program := parseProgram(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	mapStmt, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral, got=%T", stmt.Expression)
	}

	if len(mapStmt.Pairs) != 0 {
		t.Fatalf("hash.Pairs has wrong length. got=%d, want=0", len(mapStmt.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `MAP { 'one': 0+1, 'two': 8/4, 'three': 27/9 }`

	_, program := parseProgram(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	mapStmt, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not *ast.HashLiteral, got=%T", stmt.Expression)
	}

	if len(mapStmt.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d, want=3", len(mapStmt.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 8, "/", 4)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 27, "/", 9)
		},
	}

	for k, v := range mapStmt.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", k)
			continue
		}

		testFunc, ok := tests[literal.Value]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(v)
	}
}

// private functions for testing
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.Tok.Literal != value {
		t.Errorf("ident.Tok.Literal not %s. got=%s", value, ident.Tok.Literal)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.Token().Literal != fmt.Sprintf("%d", value) {
		t.Errorf("integ.Tok.Literal not %d. got=%s", value, integ.Token().Literal)
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	float, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
		return false
	}

	if float.Value != value {
		t.Errorf("float.Value not %f. got=%f", value, float.Value)
		return false
	}

	if float.Token().Literal != fmt.Sprintf("%f", value) {
		t.Errorf("float.Tok.Literal not %f. got=%s", value, float.Token().Literal)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
	b, ok := bl.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("bl not *ast.BooleanLitertal. got=%T", bl)
		return false
	}

	if b.Value != value {
		t.Errorf("b.Value not %t. got=%t", value, b.Value)
		return false
	}

	if b.Token().Literal != fmt.Sprintf("%t", value) {
		t.Errorf("b.Tok.Literal not %t. got=%s", value, b.Token().Literal)
		return false
	}

	return true
}

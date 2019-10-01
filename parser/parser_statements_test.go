package parser

import (
	"testing"

	"github.com/ollybritton/aqa/ast"
	"github.com/stretchr/testify/assert"
)

func TestVariableAssignments(t *testing.T) {
	input := `x <- 5
y <- 10
foobar <- 838383`

	_, program := parseProgram(t, input)
	assert.Equal(t, 3, len(program.Statements), "program should contain exactly 3 statements. got=%d", len(program.Statements))

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		if !testVariableAssignment(t, stmt, test.expectedIdent) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5
return 10
return 993322`

	_, program := parseProgram(t, input)
	assert.Equal(t, 3, len(program.Statements), "program should contain exactly 3 statements. got=%d", len(program.Statements))

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.Token().Literal != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got=%q", returnStmt.Token().Literal)
		}
	}
}

func TestIfStatement(t *testing.T) {
	input := `IF 1 == 1 THEN a ENDIF`
	_, program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, 1, "==", 1) {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", stmt.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "a") {
		return
	}

	if stmt.Else != nil {
		t.Errorf("stmt.Else wasn't nil. got=%+v", stmt.Else)
	}

	if stmt.ElseIf != nil {
		t.Errorf("unexpected else ifs in if stmt. got=%+v", stmt.ElseIf)
	}
}

func TestSubroutineStatement(t *testing.T) {
	input := `SUBROUTINE add(a, b) 1+1 ENDSUBROUTINE`

	_, program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.Subroutine)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Subroutine. got=%T", program.Statements[0])
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("subroutine body does not contain %d. got=%d", 1, len(stmt.Body.Statements))
	}

	expr, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("subroutine body is not ast.ExpressionStatement. got=%T", stmt.Body.Statements[0])
	}

	if !testInfixExpression(t, expr.Expression, 1, "+", 1) {
		return
	}

}

func TestSubroutineParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"SUBROUTINE add() ENDSUBROUTINE", []string{}},
		{"SUBROUTINE add(a) ENDSUBROUTINE", []string{"a"}},
		{"SUBROUTINE add(a, b) ENDSUBROUTINE", []string{"a", "b"}},
	}

	for _, tt := range tests {
		_, program := parseProgram(t, tt.input)

		sub := program.Statements[0].(*ast.Subroutine)

		if len(sub.Parameters) != len(tt.expectedParams) {
			t.Errorf("length of parameters wrong. want=%d, got=%d", len(tt.expectedParams), len(sub.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, sub.Parameters[i], ident)
		}
	}
}

// Oh god. I suggest you minimise this one
func TestIfElseElseIfStatement(t *testing.T) {
	input := `IF 1 == 1 THEN
	a 
ELSE IF 2 == 2 THEN
	b
ELSE IF 3 == 3 THEN
	c
ELSE
	d
ENDIF`
	_, program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, 1, "==", 1) {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", stmt.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "a") {
		return
	}

	if len(stmt.Else.Statements) != 1 {
		t.Errorf("else consequence is not 1 statement. got=%d", len(stmt.Else.Statements))
	}

	elseConsequence, ok := stmt.Else.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("else.Statements[0] not *ast.ExpressionStatement. got=%T", stmt.Else.Statements[0])
	}

	if !testIdentifier(t, elseConsequence.Expression, "d") {
		return
	}

	if stmt.ElseIf == nil {
		t.Errorf("expected else if in stmt. got=<nil>")
	}

	if !testInfixExpression(t, stmt.ElseIf.Condition, 2, "==", 2) {
		return
	}

	if len(stmt.ElseIf.Consequence.Statements) != 1 {
		t.Errorf("stmt.ElseIf.Statements is not 1 statement. got=%d", len(stmt.ElseIf.Consequence.Statements))
	}

	c, ok := stmt.ElseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.ElseIf.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", stmt.ElseIf.Consequence.Statements[0])
	}

	if !testIdentifier(t, c.Expression, "b") {
		return
	}

	if stmt.ElseIf.ElseIf == nil {
		t.Errorf("expected nested else-if in stmt. got=<nil>")
	}

	if !testInfixExpression(t, stmt.ElseIf.ElseIf.Condition, 3, "==", 3) {
		return
	}

	if len(stmt.ElseIf.ElseIf.Consequence.Statements) != 1 {
		t.Errorf("stmt.ElseIf.ElseIf.Statements is not 1 statement. got=%d", len(stmt.ElseIf.ElseIf.Consequence.Statements))
	}

	c, ok = stmt.ElseIf.ElseIf.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.ElseIf.ElseIf.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", stmt.ElseIf.ElseIf.Consequence.Statements[0])
	}

	if !testIdentifier(t, c.Expression, "c") {
		return
	}
}

func TestWhileStatement(t *testing.T) {
	input := `WHILE a != 10
	a <- a + 1
ENDWHILE`

	_, program := parseProgram(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	while, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T", program.Statements[0])
	}

	if while.Condition == nil {
		t.Fatalf("while.Condition is nil.")
	}

	if while.Condition.String() != "(a != 10)" {
		t.Fatalf("while.Condition is not equal to `(a != 10)`. got=%s", while.Condition.String())
	}

	if len(while.Body.Statements) != 1 {
		t.Fatalf("while.Body.Statements does not contain %d statements. got=%d", 1, len(while.Body.Statements))
	}

	_, ok = while.Body.Statements[0].(*ast.VariableAssignment)
	if !ok {
		t.Fatalf("while.Body.Statements[0] is not *ast.VariableAssignment. got=%T", while.Body.Statements[0])
	}

}

func TestRepeatStatement(t *testing.T) {
	input := `REPEAT
	a <- a + 1
UNTIL a == 10`

	_, program := parseProgram(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	repeat, ok := program.Statements[0].(*ast.RepeatStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T", program.Statements[0])
	}

	if repeat.Condition == nil {
		t.Fatalf("repeat.Condition is nil.")
	}

	if repeat.Condition.String() != "(a == 10)" {
		t.Fatalf("repeat.Condition is not equal to `(a == 10)`. got=%s", repeat.Condition.String())
	}

	if len(repeat.Body.Statements) != 1 {
		t.Fatalf("repeat.Body.Statements does not contain %d statements. got=%d", 1, len(repeat.Body.Statements))
	}

	_, ok = repeat.Body.Statements[0].(*ast.VariableAssignment)
	if !ok {
		t.Fatalf("repeat.Body.Statements[0] is not *ast.VariableAssignment. got=%T", repeat.Body.Statements[0])
	}

}
func TestForStatement(t *testing.T) {
	input := `FOR i <- 10 TO 20
  a <- a + 1
ENDFOR`

	_, program := parseProgram(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T", program.Statements[0])
	}

	if stmt.Ident.Value != "i" {
		t.Fatalf("incorrect for loop variable. got=%s, want=%s", stmt.Ident.Value, "i")
	}

	if stmt.Lower.Value != 10 {
		t.Fatalf("incorrect for loop lower bound. got=%d, want=%d", stmt.Lower.Value, 10)
	}

	if stmt.Upper.Value != 20 {
		t.Fatalf("incorrect for loop upper bound. got=%d, want=%d", stmt.Upper.Value, 10)
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("stmt.Body.Statements does not contain %d statements. got=%d", 1, len(stmt.Body.Statements))
	}

	_, ok = stmt.Body.Statements[0].(*ast.VariableAssignment)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not *ast.VariableAssignment. got=%T", stmt.Body.Statements[0])
	}

}

// private methods to help with statement tests
func testVariableAssignment(t *testing.T, s ast.Statement, expectedName string) bool {
	if s.Token().Literal != expectedName {
		t.Errorf("s.Token.Literal not %q. got=%q", expectedName, s.Token().Literal)
		return false
	}

	varStmt, ok := s.(*ast.VariableAssignment)
	if !ok {
		t.Errorf("s not *ast.VariableAssignment. got=%T", s)
		return false
	}

	if varStmt.Name.Value != expectedName {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", expectedName, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.Token().Literal != expectedName {
		t.Errorf("varStmt.Name not '%s'. got=%s", expectedName, varStmt.Name.Token().Literal)
		return false
	}

	return true

}

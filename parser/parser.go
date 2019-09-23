package parser

import (
	"fmt"
	"strconv"

	"github.com/ollybritton/aqa++/ast"
	"github.com/ollybritton/aqa++/lexer"
	"github.com/ollybritton/aqa++/token"
)

// Definition of parsing functions.
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser represents a parser for an aqa++ program.
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []error

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New returns a new parser from a given lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = map[token.Type]prefixParseFn{
		token.IDENT: p.parseIdentifier,
		token.INT:   p.parseIntegerLiteral,
		token.TRUE:  p.parseBooleanLiteral,
		token.FALSE: p.parseBooleanLiteral,

		token.BANG:  p.parsePrefixExpression,
		token.MINUS: p.parsePrefixExpression,

		token.LPAREN: p.parseGroupedExpression,
	}

	p.infixParseFns = map[token.Type]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NOT_EQ:   p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,

		token.LPAREN: p.parseCallExpression,
	}

	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns the errors that occured during parsing.
func (p *Parser) Errors() []error {
	return p.errors
}

// addError adds an error to the parser's internal error list.
func (p *Parser) addError(err error) {
	p.errors = append(p.errors, err)
}

// Parse parses the input program into a ast.Program.
func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	for p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}

	switch {
	case p.curToken.Type == token.IDENT && p.peekTokenIs(token.ASSIGN):
		return p.parseVariableAssignment()
	case p.curToken.Type == token.RETURN:
		return p.parseReturnStatement()
	case p.curToken.Type == token.IF:
		return p.parseIfStatement()
	case p.curToken.Type == token.SUBROUTINE:
		return p.parseSubroutineDefinition()
	default:
		return p.parseExpressionStatement()
	}
}

// Individual Statement Parsing
func (p *Parser) parseVariableAssignment() *ast.VariableAssignment {
	stmt := &ast.VariableAssignment{Tok: p.curToken}
	stmt.Name = &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		p.addError(
			NewUnexpectedTokenError(p.curToken, p.peekToken, token.ASSIGN),
		)

		return nil
	}

	p.nextToken()

	// TODO: Skipping expressions until we can parse them
	// for !p.curTokenIs(token.NEWLINE) && !p.curTokenIs(token.EOF) {
	// 	p.nextToken()
	// }
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Tok: p.curToken}

	p.nextToken()

	// TODO: Skipping expressions until we can parse them
	// for !p.curTokenIs(token.NEWLINE) && !p.curTokenIs(token.EOF) {
	// 	p.nextToken()
	// }
	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Tok: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.BLOCK_START) {
		p.addError(
			NewUnexpectedTokenError(p.curToken, p.peekToken, token.BLOCK_START),
		)

		return nil
	}

	stmt.Consequence = p.parseBlockStatement([]token.Type{token.BLOCK_END, token.ELSE})

	switch {
	case p.curTokenIs(token.BLOCK_END):
		return stmt

	case p.curTokenIs(token.ELSE) && p.peekTokenIs(token.IF):

		stmt.ElseIf = p.parseElseIfStatement()
		fallthrough

	case p.curTokenIs(token.ELSE):

		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}

		stmt.Else = p.parseBlockStatement([]token.Type{token.BLOCK_END})

		return stmt

	default:
		p.addError(
			NewInvalidTokenError(p.curToken, p.peekToken, p.curToken),
		)

		return nil
	}
}

func (p *Parser) parseElseIfStatement() *ast.IfStatement {

	stmt := &ast.IfStatement{Tok: p.curToken}

	// Skip token.ELSE and token.IF
	p.nextToken()
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.BLOCK_START) {

		return nil
	}

	stmt.Consequence = p.parseBlockStatement([]token.Type{token.BLOCK_END, token.ELSE})

	switch {
	case p.curTokenIs(token.ELSE) && p.peekTokenIs(token.IF):

		stmt.ElseIf = p.parseElseIfStatement()

	case p.curTokenIs(token.BLOCK_END) || p.curTokenIs(token.ELSE):
		return stmt

	default:

		return nil
	}

	return stmt

}

func (p *Parser) parseSubroutineDefinition() *ast.Subroutine {
	sub := &ast.Subroutine{Tok: p.curToken}
	p.nextToken()

	for p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}

	sub.Name = p.parseIdentifier().(*ast.Identifier)
	if !p.expectPeek(token.LPAREN) {
		p.addError(
			NewUnexpectedTokenError(p.curToken, p.peekToken, token.LPAREN),
		)

		return nil
	}

	sub.Parameters = p.parseParameters()
	sub.Body = p.parseBlockStatement([]token.Type{token.BLOCK_END})

	p.nextToken()

	return sub
}

func (p *Parser) parseParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatement(until []token.Type) *ast.BlockStatement {
	block := &ast.BlockStatement{Tok: p.curToken}

	p.nextToken()
	for !p.curTokenIs(token.EOF) {
		for _, stopToken := range until {
			if p.curTokenIs(stopToken) {
				return block
			}
		}

		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Tok: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}

	return stmt
}

// Expression Parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addError(
			NewNoPrefixParseFnError(p.curToken, p.peekToken, p.curToken.Type),
		)

		return nil
	}
	leftExp := prefix()

	for !(p.peekTokenIs(token.NEWLINE) || p.peekTokenIs(token.EOF)) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			fmt.Println("sup beach")
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Tok: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(
			NewIntegerParseError(p.curToken, p.peekToken, p.curToken.Literal),
		)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Tok: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	switch t := left.(type) {
	case *ast.Identifier:
		return p.parseSubroutineCall(t)
	default:
		// TODO: add support for function expressions
		p.addError(
			NewInvalidTokenError(left.Token(), p.curToken, left.Token()),
		)
		return nil
	}
}

func (p *Parser) parseSubroutineCall(ident *ast.Identifier) ast.Expression {
	exp := &ast.SubroutineCall{Tok: p.curToken, Subroutine: ident}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Tok:      p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	var expression = &ast.InfixExpression{
		Tok:      p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
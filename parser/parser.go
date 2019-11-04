package parser

import (
	"errors"
	"strconv"

	"github.com/ollybritton/aqa/ast"
	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/token"
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
		token.FLOAT: p.parseFloatLiteral,
		token.TRUE:  p.parseBooleanLiteral,
		token.FALSE: p.parseBooleanLiteral,

		token.BANG:  p.parsePrefixExpression,
		token.MINUS: p.parsePrefixExpression,

		token.NOT: p.parsePrefixExpression,

		token.LPAREN:   p.parseGroupedExpression,
		token.LBRACKET: p.parseArrayLiteral,
		token.STRING:   p.parseStringLiteral,

		token.OUTPUT:    p.parseOutput,
		token.USERINPUT: p.parseUserinput,

		token.MAP:    p.parseHashLiteral,
		token.LBRACE: p.parseHashLiteral,
	}

	p.infixParseFns = map[token.Type]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NOT_EQ:   p.parseInfixExpression,

		token.LT:     p.parseInfixExpression,
		token.GT:     p.parseInfixExpression,
		token.LT_EQ:  p.parseInfixExpression,
		token.GT_EQ:  p.parseInfixExpression,
		token.LSHIFT: p.parseInfixExpression,
		token.RSHIFT: p.parseInfixExpression,
		token.DIV:    p.parseInfixExpression,
		token.MOD:    p.parseInfixExpression,
		token.DOT:    p.parseInfixExpression,

		token.AND: p.parseInfixExpression,
		token.OR:  p.parseInfixExpression,
		token.XOR: p.parseInfixExpression,

		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
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

	for {
		if p.curTokenIs(token.EOF) {
			break
		}

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

	if p.curTokenIs(token.EOF) {
		return nil
	}

	switch {
	case p.curToken.Type == token.IDENT && p.peekTokenIs(token.ASSIGN):
		return p.parseVariableAssignment()
	case p.curToken.Type == token.CONSTANT && p.peekTokenIs(token.IDENT):
		return p.parseConstantAssignment()
	case p.curToken.Type == token.RETURN:
		return p.parseReturnStatement()
	case p.curToken.Type == token.IF:
		return p.parseIfStatement()
	case p.curToken.Type == token.SUBROUTINE:
		return p.parseSubroutineDefinition()
	case p.curToken.Type == token.WHILE:
		return p.parseWhileStatement()
	case p.curToken.Type == token.FOR:
		return p.parseForStatement()
	case p.curToken.Type == token.REPEAT:
		return p.parseRepeatStatement()
	case p.curToken.Type == token.IMPORT:
		return p.parseImportStatement()
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

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Tok: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Tok: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		p.addError(
			NewUnexpectedTokenError(p.curToken, p.peekToken, token.THEN),
		)

		return nil
	}

	stmt.Consequence = p.parseBlockStatement([]token.Type{token.ENDIF, token.ELSE})

	switch {
	case p.curTokenIs(token.ENDIF):
		return stmt

	case p.curTokenIs(token.ELSE) && p.peekTokenIs(token.IF):

		stmt.ElseIf = p.parseElseIfStatement()
		fallthrough

	case p.curTokenIs(token.ELSE):

		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}

		stmt.Else = p.parseBlockStatement([]token.Type{token.ENDIF})

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

	if !p.expectPeek(token.THEN) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement([]token.Type{token.ENDIF, token.ELSE})

	switch {
	case p.curTokenIs(token.ELSE) && p.peekTokenIs(token.IF):
		stmt.ElseIf = p.parseElseIfStatement()

	case p.curTokenIs(token.ENDIF) || p.curTokenIs(token.ELSE):
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
	sub.Body = p.parseBlockStatement([]token.Type{token.ENDSUBROUTINE})

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
		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}

		for _, stopToken := range until {
			if p.curTokenIs(stopToken) {
				return block
			}
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
			NewFloatParseError(p.curToken, p.peekToken, p.curToken.Literal),
		)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Tok: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(
			NewFloatParseError(p.curToken, p.peekToken, p.curToken.Literal),
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
	// some expression ( 1, 2 )
	//                 ^

	return p.parseSubroutineCall(left)
	// switch t := left.(type) {
	// // case *a÷st.Identifier:
	// default:
	// TODO: add support for function expressions
	// log.Printf("Im here: %T", left)
	// p.addError(
	// 	NewInvalidTokenError(left.Token(), p.curToken, left.Token()),
	// )
	// 	return nil
	// }
}

func (p *Parser) parseSubroutineCall(expression ast.Expression) ast.Expression {
	exp := &ast.SubroutineCall{Tok: p.curToken, Subroutine: expression}
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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Tok: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	while := &ast.WhileStatement{Tok: p.curToken}
	p.nextToken()

	while.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.NEWLINE) {
		p.addError(NewUnexpectedTokenError(p.curToken, p.peekToken, token.NEWLINE))
	}
	while.Body = p.parseBlockStatement([]token.Type{token.ENDWHILE})

	return while
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Tok: p.curToken}
	p.nextToken()

	stmt.Ident = p.parseIdentifier().(*ast.Identifier)

	if !p.expectPeek(token.ASSIGN) {
		p.addError(NewInvalidTokenError(p.curToken, p.peekToken, p.curToken))
		return &ast.ForStatement{}
	}
	p.nextToken()

	lower := p.parseExpression(LOWEST)

	if !p.expectPeek(token.TO) {
		p.addError(NewInvalidTokenError(p.curToken, p.peekToken, p.curToken))
		return &ast.ForStatement{}
	}
	p.nextToken()

	upper := p.parseExpression(LOWEST)

	stmt.Lower = lower
	stmt.Upper = upper

	if !p.expectPeek(token.NEWLINE) {
		p.addError(NewUnexpectedTokenError(p.curToken, p.peekToken, token.NEWLINE))
	}

	stmt.Body = p.parseBlockStatement([]token.Type{token.ENDFOR})

	return stmt
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	repeat := &ast.RepeatStatement{Tok: p.curToken}

	repeat.Body = p.parseBlockStatement([]token.Type{token.UNTIL})
	p.nextToken()

	repeat.Condition = p.parseExpression(LOWEST)

	return repeat
}

func (p *Parser) parseOutput() ast.Expression {

	call := &ast.SubroutineCall{Tok: p.curToken}

	call.Subroutine = &ast.Identifier{Tok: p.curToken, Value: "OUTPUT"}

	p.nextToken()

	call.Arguments = []ast.Expression{p.parseExpression(LOWEST)}

	return call
}

func (p *Parser) parseUserinput() ast.Expression {
	return &ast.Identifier{Tok: p.curToken, Value: "USERINPUT"}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{Tok: p.curToken, Elements: p.parseExpressionList(token.RBRACKET)}
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		p.addError(NewUnexpectedTokenError(p.curToken, p.peekToken, end))
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Tok: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		p.addError(NewUnexpectedTokenError(p.curToken, p.peekToken, token.RBRACKET))
		return nil
	}

	return exp
}

func (p *Parser) parseConstantAssignment() *ast.VariableAssignment {
	p.nextToken()

	stmt := p.parseVariableAssignment()
	stmt.Name.Constant = true

	return stmt
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Tok: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	p.skipNewlines()

	if p.curTokenIs(token.MAP) {
		p.nextToken()
	}

	p.nextToken()

	for !p.peekTokenIs(token.RBRACE) {
		p.skipNewlines()

		if p.curTokenIs(token.RBRACE) {
			return hash
		}

		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		p.nextToken()

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
			p.skipNewlines()
		} else {
			break
		}
	}

	p.skipNewlines()

	return hash
}

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	// import "file.aqa"
	// import "file.aqa" as otherName
	// import abc, def from "file.aqa"
	// import * from "file.aqa"

	// import "dir"
	// import "dir" as otherName

	stmt := &ast.ImportStatement{Tok: p.curToken}
	p.nextToken()

	switch p.curToken.Type {
	case token.STRING:
		stmt.Path = p.curToken.Literal

		if p.peekTokenIs(token.AS) {
			p.nextToken()
		} else {
			break
		}

		if !p.expectPeek(token.IDENT) {
			p.addError(errors.New("missing name after 'as' in import statement"))
		}

		stmt.As = p.curToken.Literal

	case token.IDENT, token.ASTERISK:
		from := []string{p.curToken.Literal}
		p.nextToken()

		for p.curTokenIs(token.COMMA) {
			p.nextToken()

			if p.curTokenIs(token.IDENT) || p.curTokenIs(token.ASTERISK) {
				from = append(from, p.curToken.Literal)
				p.nextToken()
			} else {
				p.addError(errors.New("unknown import syntax"))
				return nil
			}
		}

		if !p.expectPeek(token.STRING) {
			p.addError(errors.New("unknown import syntax"))
			return nil
		}

		stmt.Path = p.curToken.Literal
		stmt.From = from

	default:
		p.addError(errors.New("unknown import syntax"))
		return nil
	}

	return stmt
}

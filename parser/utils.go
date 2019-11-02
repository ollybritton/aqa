package parser

import "github.com/ollybritton/aqa/token"

func (p *Parser) curTokenIs(tt token.Type) bool {
	return p.curToken.Type == tt
}

func (p *Parser) peekTokenIs(tt token.Type) bool {
	return p.peekToken.Type == tt
}

// expectPeek only advances the parser if the next token is correct.
func (p *Parser) expectPeek(tt token.Type) bool {
	if p.peekTokenIs(tt) {
		p.nextToken()
		return true
	}

	return false
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) skipNewlines() {
	for p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}
}

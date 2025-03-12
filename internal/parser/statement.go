package parser

import (
	"github.com/timur-makarov/monkey-interpreter/internal/ast"
	"github.com/timur-makarov/monkey-interpreter/internal/token"
)

func (p *Parser) parseLetStatement() ast.Statement {
	statement := ast.LetStatement{Token: p.token}

	if !p.expectRead(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.token, Value: p.token.Literal}

	if !p.expectRead(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() ast.Statement {
	statement := ast.ReturnStatement{Token: p.token}

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := ast.ExpressionStatement{Token: p.token}

	statement.Expression = p.parseExpression(LOWEST)

	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseBlockStatement() ast.BlockStatement {
	statement := ast.BlockStatement{Token: p.token}

	p.nextToken()

	for p.token.Type != token.RBRACE && p.token.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			statement.Statements = append(statement.Statements, stmt)
		}
		p.nextToken()
	}

	return statement
}

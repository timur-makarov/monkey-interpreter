package parser

import (
	"strconv"

	"github.com/timur-makarov/monkey-interpreter/internal/ast"
	"github.com/timur-makarov/monkey-interpreter/internal/token"
)

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.Identifier{Token: p.token, Value: p.token.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	value, err := strconv.ParseBool(p.token.Literal)
	if err != nil {
		p.pushError(invalidValue("boolean", err))
	}

	return ast.Boolean{Token: p.token, Value: value}
}

func (p *Parser) parseInteger() ast.Expression {
	value, err := strconv.Atoi(p.token.Literal)
	if err != nil {
		p.pushError(invalidValue("integer", err))
	}

	return ast.Integer{Token: p.token, Value: value}
}

func (p *Parser) parseString() ast.Expression {
	return ast.String{Token: p.token, Value: p.token.Literal}
}

func (p *Parser) parseGroup() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectRead(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIf() ast.Expression {
	expression := ast.If{Token: p.token}

	if !p.expectRead(token.LPAREN) {
		return nil
	}

	condition := p.parseGroup()
	if condition == nil {
		return nil
	}

	expression.Conditions = append(expression.Conditions, condition)

	if !p.expectRead(token.LBRACE) {
		return nil
	}

	expression.Consequences = append(expression.Consequences, p.parseBlockStatement())

	if p.readToken.Type == token.ELSE {
		p.nextToken()

		if p.readToken.Type == token.IF {
			p.nextToken()
			exp := p.parseIf().(ast.If)
			expression.Conditions = append(expression.Conditions, exp.Conditions...)
			expression.Consequences = append(expression.Consequences, exp.Consequences...)
			expression.Alternative = exp.Alternative
			return expression
		}

		if !p.expectRead(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseWhile() ast.Expression {
	expression := ast.While{Token: p.token}

	if !p.expectRead(token.LPAREN) {
		return nil
	}

	expression.Condition = p.parseGroup()
	if expression.Condition == nil {
		return nil
	}

	if !p.expectRead(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseFunction() ast.Expression {
	expression := ast.Function{Token: p.token}

	if !p.expectRead(token.LPAREN) {
		return nil
	}

	p.nextToken()

	if p.token.Type != token.RPAREN {
		expression.Parameters = p.parseFunctionParameters()
		p.nextToken()
	}

	if !p.expectRead(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseFunctionParameters() []ast.Identifier {
	var parameters []ast.Identifier

	for {
		ident := p.parseIdentifier()
		parameters = append(parameters, ident.(ast.Identifier))

		if p.readToken.Type == token.RPAREN {
			break
		}

		if !p.expectRead(token.COMMA) {
			break
		}

		p.nextToken()
	}

	return parameters
}

func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	expression := ast.Call{Token: p.token, Function: function}
	p.nextToken()
	if p.token.Type != token.RPAREN {
		expression.Arguments = p.parseCallArguments()
	}
	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var arguments []ast.Expression

	for {
		argument := p.parseExpression(LOWEST)
		if argument == nil {
			break
		}

		arguments = append(arguments, argument)

		if p.readToken.Type == token.RPAREN {
			p.nextToken()
			break
		}

		if !p.expectRead(token.COMMA) {
			break
		}

		p.nextToken()
	}

	return arguments
}

func (p *Parser) parseArray() ast.Expression {
	expression := ast.Array{Token: p.token}

	p.nextToken()

	for p.token.Type != token.RBRACKET && p.token.Type != token.EOF {
		exp := p.parseExpression(LOWEST)
		expression.Items = append(expression.Items, exp)

		if p.readToken.Type == token.COMMA {
			p.nextToken()
		} else if p.readToken.Type != token.RBRACKET {
			p.pushError(unexpectedTypeError(token.RBRACKET, p.readToken.Type))
		}
		p.nextToken()
	}

	return expression
}

func (p *Parser) parseHashTable() ast.Expression {
	expression := ast.HashTable{Token: p.token, Items: make(map[ast.Expression]ast.Expression)}

	for p.readToken.Type != token.RBRACE && p.readToken.Type != token.EOF {
		if p.readToken.Type != token.IDENT && p.readToken.Type != token.STRING {
			return nil
		}
		p.nextToken()

		keyExp := p.parseExpression(LOWEST)

		if !p.expectRead(token.COLON) {
			return nil
		}

		p.nextToken()

		valExp := p.parseExpression(LOWEST)

		expression.Items[keyExp] = valExp

		if p.readToken.Type == token.COMMA {
			p.nextToken()
		} else if p.readToken.Type != token.RBRACE {
			p.pushError(unexpectedTypeError(token.RBRACE, p.readToken.Type))
		}
	}

	p.nextToken()

	return expression
}

func (p *Parser) parseAccessByIndexOrKey(leftExp ast.Expression) ast.Expression {
	expression := ast.AccessByExpression{Token: p.token, Left: leftExp}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectRead(token.RBRACKET) {
		return nil
	}

	return expression
}

func (p *Parser) parsePrefix() ast.Expression {
	expression := ast.Prefix{Token: p.token, Operator: p.token.Literal}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfix(leftExp ast.Expression) ast.Expression {
	expression := ast.Infix{Token: p.token, Operator: p.token.Literal, Left: leftExp}

	precedence := precedences[p.token.Type]
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

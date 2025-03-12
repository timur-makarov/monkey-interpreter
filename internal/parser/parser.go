package parser

import (
	"github.com/timur-makarov/monkey-interpreter/internal/ast"
	"github.com/timur-makarov/monkey-interpreter/internal/lexer"
	"github.com/timur-makarov/monkey-interpreter/internal/token"
)

type Parser struct {
	l *lexer.Lexer

	token     token.Token
	readToken token.Token

	errors []Error

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

type prefixParseFn = func() ast.Expression
type infixParseFn = func(ast.Expression) ast.Expression

type Precedence = int

const (
	LOWEST       Precedence = 1
	ASSIGN       Precedence = 2
	EQUALS       Precedence = 3
	COMPARISON   Precedence = 4
	SUM          Precedence = 5
	PRODUCT      Precedence = 6
	PREFIX       Precedence = 7
	CALL         Precedence = 8
	INDEX_OR_KEY Precedence = 8
)

var precedences = map[token.Type]Precedence{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.GT:       COMPARISON,
	token.LT:       COMPARISON,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.MULTIPLY: PRODUCT,
	token.DIVIDE:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX_OR_KEY,
	token.ASSIGN:   ASSIGN,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)

	p.registerPrefixFn(token.INT, p.parseInteger)
	p.registerPrefixFn(token.STRING, p.parseString)
	p.registerPrefixFn(token.TRUE, p.parseBoolean)
	p.registerPrefixFn(token.FALSE, p.parseBoolean)

	p.registerPrefixFn(token.BANG, p.parsePrefix)
	p.registerPrefixFn(token.MINUS, p.parsePrefix)

	p.registerPrefixFn(token.LPAREN, p.parseGroup)
	p.registerPrefixFn(token.LBRACKET, p.parseArray)
	p.registerPrefixFn(token.LBRACE, p.parseHashTable)

	p.registerPrefixFn(token.IF, p.parseIf)
	p.registerPrefixFn(token.WHILE, p.parseWhile)
	p.registerPrefixFn(token.FUNCTION, p.parseFunction)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfixFn(token.PLUS, p.parseInfix)
	p.registerInfixFn(token.MINUS, p.parseInfix)
	p.registerInfixFn(token.DIVIDE, p.parseInfix)
	p.registerInfixFn(token.MULTIPLY, p.parseInfix)
	p.registerInfixFn(token.GT, p.parseInfix)
	p.registerInfixFn(token.LT, p.parseInfix)
	p.registerInfixFn(token.EQ, p.parseInfix)
	p.registerInfixFn(token.NEQ, p.parseInfix)
	p.registerInfixFn(token.ASSIGN, p.parseInfix)
	p.registerInfixFn(token.LPAREN, p.parseCall)
	p.registerInfixFn(token.LBRACKET, p.parseAccessByIndexOrKey)

	return p
}

func (p *Parser) Errors() []Error {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	var program ast.Program

	for p.token.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return &program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.token.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix, ok := p.prefixParseFns[p.token.Type]
	if !ok {
		p.pushError(parseFnNotImplemented(p.token.Type))
		return nil
	}

	leftExp := prefix()

	for p.readToken.Type != token.SEMICOLON && precedence < precedences[p.readToken.Type] {
		infix, ok := p.infixParseFns[p.readToken.Type]
		if !ok {
			p.pushError(parseFnNotImplemented(p.token.Type))
			return nil
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) nextToken() {
	p.token = p.readToken
	p.readToken = p.l.NextToken()
}

func (p *Parser) pushError(error Error) {
	p.errors = append(p.errors, error)
}

func (p *Parser) expectRead(expectedType token.Type) bool {
	if p.readToken.Type == expectedType {
		p.nextToken()
		return true
	} else {
		p.pushError(unexpectedTypeError(expectedType, p.readToken.Type))
		return false
	}
}

func (p *Parser) registerPrefixFn(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

package lexer

import (
	"iter"
	"unicode"
	"unicode/utf8"

	"github.com/timur-makarov/monkey-interpreter/internal/token"
)

type inputIterator struct {
	next func() (rune, bool)
}

type Lexer struct {
	input        string
	iterator     inputIterator
	position     int
	readPosition int
	character    rune
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	ii := createInputIterator(input)
	l.iterator.next, _ = iter.Pull(ii)
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	l.character, _ = l.iterator.next()
	l.position = l.readPosition
	l.readPosition += utf8.RuneLen(l.character)
}

func (l *Lexer) peekChar() rune {
	if l.readPosition > len(l.input) {
		return 0
	} else {
		return []rune(l.input[l.readPosition : l.readPosition+1])[0]
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.position

	for unicode.IsLetter(l.character) || l.character == '_' {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readInteger() string {
	pos := l.position

	for unicode.IsDigit(l.character) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readString() string {
	pos := l.position + 1

	for {
		l.readChar()

		if l.character == '"' || l.character == 0 {
			break
		}
	}

	return l.input[pos:l.position]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.character) {
		l.readChar()
	}
}

func (l *Lexer) determineTokenType(fType, sType token.Type, lookFor rune) token.Token {
	if l.peekChar() == lookFor {
		ch := l.character
		l.readChar()
		return token.Token{Type: sType, Literal: string(ch) + string(l.character)}
	} else {
		return newToken(fType, l.character)
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.character {
	case '=':
		tok = l.determineTokenType(token.ASSIGN, token.EQ, '=')
	case '+':
		tok = newToken(token.PLUS, l.character)
	case '-':
		tok = newToken(token.MINUS, l.character)
	case '/':
		tok = newToken(token.DIVIDE, l.character)
	case '*':
		tok = newToken(token.MULTIPLY, l.character)
	case '(':
		tok = newToken(token.LPAREN, l.character)
	case ')':
		tok = newToken(token.RPAREN, l.character)
	case '{':
		tok = newToken(token.LBRACE, l.character)
	case '}':
		tok = newToken(token.RBRACE, l.character)
	case '[':
		tok = newToken(token.LBRACKET, l.character)
	case ']':
		tok = newToken(token.RBRACKET, l.character)
	case ',':
		tok = newToken(token.COMMA, l.character)
	case ':':
		tok = newToken(token.COLON, l.character)
	case ';':
		tok = newToken(token.SEMICOLON, l.character)
	case '!':
		tok = l.determineTokenType(token.BANG, token.NEQ, '=')
	case '>':
		tok = newToken(token.GT, l.character)
	case '<':
		tok = newToken(token.LT, l.character)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if unicode.IsLetter(l.character) || l.character == '_' {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIndent(tok.Literal)
			return tok
		}

		if unicode.IsDigit(l.character) {
			tok.Type = token.INT
			tok.Literal = l.readInteger()
			return tok
		}

		tok = newToken(token.ILLEGAL, l.character)
	}

	l.readChar()

	return tok
}

func newToken(tokenType token.Type, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func createInputIterator(input string) iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for _, ch := range input {
			if !yield(ch) {
				return
			}
		}
	}
}

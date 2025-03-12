package parser

import (
	"fmt"

	"github.com/timur-makarov/monkey-interpreter/internal/token"
)

type Error struct {
	Message string
}

func unexpectedTypeError(expected token.Type, actual token.Type) Error {
	message := fmt.Sprintf("expected next token to be '%s', got %s instead", expected, actual)
	return Error{Message: message}
}

func parseFnNotImplemented(expected token.Type) Error {
	message := fmt.Sprintf("parse function for token type '%s' is not implemented", expected)
	return Error{Message: message}
}

func invalidValue(expected string, err error) Error {
	message := fmt.Sprintf("error parsing %s value: %v", expected, err)
	return Error{Message: message}
}

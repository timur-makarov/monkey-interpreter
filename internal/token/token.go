package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	DIVIDE   = "/"
	MULTIPLY = "*"

	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"
	BANG      = "!"
	GT        = ">"
	LT        = "<"
	EQ        = "=="
	NEQ       = "!="

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	WHILE    = "WHILE"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"while":  WHILE,
}

func LookupIndent(ident string) Type {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}

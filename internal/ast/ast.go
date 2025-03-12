package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/timur-makarov/monkey-interpreter/internal/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement = Node
type Expression = Node

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i Identifier) String() string {
	return i.Value
}

type Integer struct {
	Token token.Token
	Value int
}

func (i Integer) TokenLiteral() string {
	return i.Token.Literal
}

func (i Integer) String() string {
	return strconv.Itoa(i.Value)
}

type String struct {
	Token token.Token
	Value string
}

func (s String) TokenLiteral() string {
	return s.Token.Literal
}

func (s String) String() string {
	return "\"" + s.Value + "\""
}

type Array struct {
	Token token.Token
	Items []Expression
}

func (a Array) TokenLiteral() string {
	return a.Token.Literal
}

func (a Array) String() string {
	return fmt.Sprintf("%+v", a.Items)
}

type HashTable struct {
	Token token.Token
	Items map[Expression]Expression
}

func (ht HashTable) TokenLiteral() string {
	return ht.Token.Literal
}

func (ht HashTable) String() string {
	return fmt.Sprintf("%+v", ht.Items)
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b Boolean) String() string {
	return b.TokenLiteral()
}

type If struct {
	Token        token.Token
	Conditions   []Expression
	Consequences []BlockStatement
	Alternative  BlockStatement
}

func (i If) TokenLiteral() string {
	return i.Token.Literal
}

func (i If) String() string {
	return fmt.Sprintf(
		"if %+v {%+v} else {%+v}", i.Conditions, i.Consequences, i.Alternative,
	)
}

type While struct {
	Token     token.Token
	Condition Expression
	Body      BlockStatement
}

func (w While) TokenLiteral() string {
	return w.Token.Literal
}

func (w While) String() string {
	return fmt.Sprintf(
		"for (%s) {%s}", w.Condition, w.Body,
	)
}

type Function struct {
	Token      token.Token
	Parameters []Identifier
	Body       BlockStatement
}

func (f Function) TokenLiteral() string {
	return f.Token.Literal
}

func (f Function) String() string {
	return fmt.Sprintf("fn (%+v) {%s}", f.Parameters, f.Body)
}

type Call struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (f Call) TokenLiteral() string {
	return f.Token.Literal
}

func (f Call) String() string {
	return fmt.Sprintf("call fn %s with args (%+v)", f.Function, f.Arguments)
}

type AccessByExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (f AccessByExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f AccessByExpression) String() string {
	return fmt.Sprintf("(%s)[%s]", f.Left, f.Index)
}

type Prefix struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p Prefix) TokenLiteral() string {
	return p.Token.Literal
}

func (p Prefix) String() string {
	return fmt.Sprintf("%s %v", p.Operator, p.Right)
}

type Infix struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i Infix) TokenLiteral() string {
	return i.Token.Literal
}

func (i Infix) String() string {
	return fmt.Sprintf("%s %s %s", i.Left, i.Operator, i.Right)
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls LetStatement) String() string {
	return fmt.Sprintf("%s %s = %+v", ls.Token.Literal, ls.Name, ls.Value)
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs ReturnStatement) String() string {
	return fmt.Sprintf("%s %v", rs.Token.Literal, rs.Value)
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es ExpressionStatement) String() string {
	return es.Expression.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

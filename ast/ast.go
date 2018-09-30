package ast

import "github.com/nrtkbb/go-MEL/token"

// Node is top of AST interface
type Node interface {
	TokenLiteral() string
}

// Statement have some expression
type Statement interface {
	Node
	statementNode()
}

// Expression ...
type Expression interface {
	Node
	expressionNode()
}

// Program is represent the entire program
type Program struct {
	Statements []Statement
}

// TokenLiteral ...
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// StringStatement ...
type StringStatement struct {
	Token token.Token // token.String
	Name  *Identifier
	Value Expression
}

func (ls *StringStatement) statementNode() {}

// TokenLiteral ...
func (ls *StringStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Identifier is token.Ident
type Identifier struct {
	Token token.Token // token.Ident
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral ...
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

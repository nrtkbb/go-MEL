package parser

import (
	"fmt"
	"strconv"

	"github.com/nrtkbb/go-MEL/ast"
	"github.com/nrtkbb/go-MEL/lexer"
	"github.com/nrtkbb/go-MEL/token"
)

// precedence const
const (
	LOWEST      int = 1 + iota
	TERNARY         // ? and :
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	CREMENT         // X++ or X--
	PREFIX          // -X or !X
	CALL            // myFunction(X
)

var precedences = map[token.Type]int{
	token.Eq:        EQUALS,
	token.NotEq:     EQUALS,
	token.Lt:        LESSGREATER,
	token.Gt:        LESSGREATER,
	token.Plus:      SUM,
	token.Minus:     SUM,
	token.Slash:     PRODUCT,
	token.Asterisk:  PRODUCT,
	token.Increment: CREMENT,
	token.Decrement: CREMENT,
	token.Question:  TERNARY,
	token.Coron:     TERNARY,
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
	ternaryParseFn func(ast.Expression) ast.Expression
)

// Parser use Lexer and Token
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.Type]prefixParseFn
	infixParseFns   map[token.Type]infixParseFn
	postfixParseFns map[token.Type]postfixParseFn
	ternaryParseFns map[token.Type]ternaryParseFn
}

// Errors return parsing error strings..
func (p *Parser) Errors() []string {
	return p.errors
}

// New make Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// set prefix parse func.
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.Decrement, p.parsePrefixExpression)
	p.registerPrefix(token.Increment, p.parsePrefixExpression)

	// set infix parse func.
	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Eq, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.Lt, p.parseInfixExpression)
	p.registerInfix(token.Gt, p.parseInfixExpression)

	// set postfix parse func.
	p.postfixParseFns = make(map[token.Type]postfixParseFn)
	p.registerPostfix(token.Increment, p.parsePostfixExpression)
	p.registerPostfix(token.Decrement, p.parsePostfixExpression)

	// set ternary parse func.
	p.ternaryParseFns = make(map[token.Type]ternaryParseFn)
	p.registerTernary(token.Question, p.parseTernaryExpression)

	// Read two token. Set both curToken and peekToken.
	p.nextToken()
	p.nextToken()

	return p
}

// ParseProgram is to read token and build the Program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.StringDec:
		return p.parseStringStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		ternary := p.ternaryParseFns[p.peekToken.Type]
		if ternary != nil {
			leftExp = ternary(leftExp)
		}

		postfix := p.postfixParseFns[p.peekToken.Type]
		if postfix != nil {
			leftExp = postfix(leftExp)
		}
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedences := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedences)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	p.nextToken()

	expression := &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	return expression
}

func (p *Parser) parseTernaryExpression(conditional ast.Expression) ast.Expression {
	p.nextToken()

	expression := &ast.TernaryExpression{
		Conditional: conditional,
		Token1:      p.curToken,
		Operator1:   p.curToken.Literal,
	}
	precedences := p.curPrecedence()

	p.nextToken()
	expression.TrueExp = p.parseExpression(precedences)
	p.nextToken()

	expression.Token2 = p.curToken
	expression.Operator2 = p.curToken.Literal
	precedences = p.curPrecedence()

	p.nextToken()
	expression.FalseExp = p.parseExpression(precedences)

	return expression
}

func (p *Parser) parseStringStatement() ast.Statement {
	stmt := &ast.StringStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	// TODO: skip to semicolon.
	for !p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: skip to semicolon.
	for !p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("line:%d.%d expected next token to be %s, got %s instead",
		p.peekToken.Row, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.Type, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *Parser) registerTernary(tokenType token.Type, fn ternaryParseFn) {
	p.ternaryParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if pre, ok := precedences[p.peekToken.Type]; ok {
		return pre
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if pre, ok := precedences[p.curToken.Type]; ok {
		return pre
	}

	return LOWEST
}

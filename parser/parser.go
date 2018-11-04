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
	OR              // ||
	AND             // &&
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	CREMENT         // X++ or X--
	PREFIX          // -X or !X
	CALL            // myFunction(X
	INDEX           // array[index]
)

var precedences = map[token.Type]int{
	token.Or:         OR,
	token.And:        AND,
	token.Eq:         EQUALS,
	token.NotEq:      EQUALS,
	token.Lt:         LESSGREATER,
	token.Gt:         LESSGREATER,
	token.Plus:       SUM,
	token.Minus:      SUM,
	token.Slash:      PRODUCT,
	token.Asterisk:   PRODUCT,
	token.Increment:  CREMENT,
	token.Decrement:  CREMENT,
	token.Question:   TERNARY,
	token.Coron:      TERNARY,
	token.Lparen:     CALL,
	token.BackQuotes: CALL,
	token.Lbracket:   INDEX,
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

	// Function Style
	prefixParseFns  map[token.Type]prefixParseFn
	infixParseFns   map[token.Type]infixParseFn
	postfixParseFns map[token.Type]postfixParseFn
	ternaryParseFns map[token.Type]ternaryParseFn

	commandStyleMode bool
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
	p.registerPrefix(token.ProcIdent, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)
	p.registerPrefix(token.Float, p.parseFloatLiteral)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.Decrement, p.parsePrefixExpression)
	p.registerPrefix(token.Increment, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.On, p.parseBoolean)
	p.registerPrefix(token.Off, p.parseBoolean)
	p.registerPrefix(token.Ltensor, p.parseTensorLiteral)
	p.registerPrefix(token.Lbrace, p.parseArrayLiteral)
	p.registerPrefix(token.Lparen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Proc, p.parseFunctionLiteral)
	p.registerPrefix(token.BackQuotes, p.parseBackQuotesCallExpression)

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
	p.registerInfix(token.And, p.parseInfixExpression)
	p.registerInfix(token.Or, p.parseInfixExpression)
	p.registerInfix(token.Lparen, p.parseCallExpression)
	p.registerInfix(token.Lbracket, p.parseIndexExpression)

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
	case token.IntDec:
		return p.parseIntegerStatement()
	case token.FloatDec:
		return p.parseFloatStatement()
	case token.VectorDec:
		return p.parseVectorStatement()
	case token.MatrixDec:
		return p.parseMatrixStatement()
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

	if p.curTokenIs(token.ProcIdent) &&
		!p.peekTokenIs(token.Lparen) &&
		p.commandStyleMode == false {
		leftExp = p.parseCommandCallExpression(leftExp)
	}

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		ternary := p.ternaryParseFns[p.peekToken.Type]
		if ternary != nil {
			leftExp = ternary(leftExp)
		}

		postfix := p.postfixParseFns[p.peekToken.Type]
		if postfix != nil {
			leftExp = postfix(leftExp)
		}
		if p.commandStyleMode && p.peekTokenIs(token.Lparen) {
			return leftExp
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

// like this:
//   add 1 (2 + 3) a "b";
func (p *Parser) parseCommandCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	preCommandMode := p.commandStyleMode
	p.commandStyleMode = true
	defer func() { p.commandStyleMode = preCommandMode }()

	exp.Arguments = p.parseCommandCallArguments()

	return exp
}

func (p *Parser) parseCommandCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for !p.peekTokenIs(token.Semicolon) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.peekTokenIs(token.Semicolon) {
		return nil
	}
	p.nextToken()

	return args
}

// like this:
//   `add 1 (2 + 3) a "b"`;
func (p *Parser) parseBackQuotesCallExpression() ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken}
	if p.peekTokenIs(token.ProcIdent) {
		p.nextToken()
		exp.Function = p.parseIdentifier()
	} else {
		return nil
	}

	preCommandMode := p.commandStyleMode
	p.commandStyleMode = true
	defer func() { p.commandStyleMode = preCommandMode }()

	exp.Arguments = p.parseBackQuotesCallArguments()

	return exp
}

func (p *Parser) parseBackQuotesCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.BackQuotes) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for !p.peekTokenIs(token.BackQuotes) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.peekTokenIs(token.BackQuotes) {
		return nil
	}
	p.nextToken()

	return args
}

// like this:
//   add(1, (2 + 3), "a", "b");
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.Rparen)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	if p.curTokenIs(token.Rbracket) {
		exp.Index = nil
		return exp
	}

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.Rbracket) {
		return nil
	}

	return exp
}

func (p *Parser) parseMatrixStatement() ast.Statement {
	stmt := &ast.MatrixStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Names, stmt.Values = p.parseBulkDefinition()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVectorStatement() ast.Statement {
	stmt := &ast.VectorStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Names, stmt.Values = p.parseBulkDefinition()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIntegerStatement() ast.Statement {
	stmt := &ast.IntegerStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Names, stmt.Values = p.parseBulkDefinition()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseFloatStatement() ast.Statement {
	stmt := &ast.FloatStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Names, stmt.Values = p.parseBulkDefinition()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseStringStatement() ast.Statement {
	stmt := &ast.StringStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Names, stmt.Values = p.parseBulkDefinition()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBulkDefinition() ([]ast.Expression, []ast.Expression) {
	var names []ast.Expression
	var values []ast.Expression

	name := p.parseExpression(LOWEST)
	names = append(names, name)
	if !p.peekTokenIs(token.Assign) {
		values = append(values, nil)
	} else {
		p.nextToken()
		p.nextToken()
		value := p.parseExpression(LOWEST)
		values = append(values, value)
	}

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		name := p.parseExpression(LOWEST)
		names = append(names, name)

		if !p.peekTokenIs(token.Assign) {
			values = append(values, nil)
			continue
		}
		p.nextToken()
		p.nextToken()

		value := p.parseExpression(LOWEST)
		values = append(values, value)
	}

	return names, values
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.Rparen) {
		return nil
	}

	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	lit.ReturnType = p.parseTypeDeclaration()

	if p.peekTokenIs(token.ProcIdent) {
		p.nextToken()
		lit.Name = p.curToken
	} else {
		return nil
	}

	if !p.expectPeek(token.Lparen) {
		return nil
	}

	lit.ParamTypes, lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.Lbrace) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() ([]*ast.TypeDeclaration, []*ast.Identifier) {
	var typeDeclarations []*ast.TypeDeclaration
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.Rparen) {
		p.nextToken()
		return typeDeclarations, identifiers
	}

	typeDeclaration := p.parseTypeDeclaration()
	typeDeclarations = append(typeDeclarations, typeDeclaration)
	p.nextToken()
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		typeDeclaration := p.parseTypeDeclaration()
		typeDeclarations = append(typeDeclarations, typeDeclaration)
		p.nextToken()
		identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.Rparen) {
		return nil, nil
	}

	return typeDeclarations, identifiers
}

func (p *Parser) parseTypeDeclaration() *ast.TypeDeclaration {
	switch p.peekToken.Type {
	case token.StringDec:
		fallthrough
	case token.FloatDec:
		fallthrough
	case token.VectorDec:
		fallthrough
	case token.IntDec:
		p.nextToken()
		td := &ast.TypeDeclaration{Token: p.curToken}
		if p.peekTokenIs(token.Lbracket) {
			p.nextToken()
			if p.peekTokenIs(token.Rbracket) {
				p.nextToken()
				td.IsArray = true
			} else {
				return nil
			}
		} else {
			return nil
		}
		return td
	case token.MatrixDec:
		p.nextToken()
		td := &ast.TypeDeclaration{Token: p.curToken}
		return td
	}
	return nil
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.Lparen) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.Rparen) {
		return nil
	}

	if p.peekTokenIs(token.Lbrace) {
		p.nextToken()
		expression.Consequence = p.parseBlockStatement()
	} else {
		p.nextToken()
		expression.Consequence = p.parseSingleBlockStatement()
	}

	if p.peekTokenIs(token.Else) {
		p.nextToken()

		if p.peekTokenIs(token.Lbrace) {
			p.nextToken()
			expression.Alternative = p.parseBlockStatement()
		} else {
			p.nextToken()
			expression.Alternative = p.parseSingleBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.Rbrace) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseSingleBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: token.Token{
		Type:    token.Lbrace,
		Literal: "{",
		Row:     p.curToken.Row,
		Column:  p.curToken.Column,
	}}
	block.Statements = []ast.Statement{}

	stmt := p.parseStatement()
	if stmt != nil {
		block.Statements = append(block.Statements, stmt)
	}

	return block
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

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseTensorLiteral() ast.Expression {
	tl := &ast.TensorLiteral{Token: p.curToken}
	p.nextToken()

	var value []ast.Expression
	value = append(value, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		value = append(value, p.parseExpression(LOWEST))
	}
	tl.Values = append(tl.Values, value)
	value = nil

	for p.peekTokenIs(token.Semicolon) {
		p.nextToken()
		value = append(value, p.parseExpression(LOWEST))
		for p.peekTokenIs(token.Comma) {
			p.nextToken()
			p.nextToken()
			value = append(value, p.parseExpression(LOWEST))
		}
		tl.Values = append(tl.Values, value)
		value = nil
	}

	if !p.expectPeek(token.Rtensor) {
		return nil
	}

	return tl
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.Rbrace)

	return array
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.True)}
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

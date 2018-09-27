package lexer

import (
	"github.com/nrtkbb/go-MEL/token"
)

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	rune         rune
}

func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readRune()
	return l
}

func (l *Lexer) readRune() {
	if l.readPosition >= len(l.input) {
		l.rune = 0
	} else {
		l.rune = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, r rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(r)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.rune {
	case '=':
		tok = newToken(token.ASSIGN, l.rune)
	case ';':
		tok = newToken(token.SEMICOLON, l.rune)
	case '(':
		tok = newToken(token.LPAREN, l.rune)
	case ')':
		tok = newToken(token.RPAREN, l.rune)
	case '{':
		tok = newToken(token.LBRACE, l.rune)
	case '}':
		tok = newToken(token.RBRACE, l.rune)
	case ',':
		tok = newToken(token.COMMA, l.rune)
	case '+':
		tok = newToken(token.PLUS, l.rune)
	case '$':
		tok.Literal = l.readIdentifier()
		tok.Type = token.IDENT
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.rune) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.rune) {
			tok.Type = token.INT_DATA
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.rune)
		}
	}

	l.readRune()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	l.readRune() // はじめの $ をスキップする
	for isIdentifier(l.rune) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '_' == r
}

func isIdentifier(r rune) bool {
	return isLetter(r) || '0' <= r && r <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.rune) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func (l *Lexer) skipWhitespace() {
	for ' ' == l.rune || '\t' == l.rune || '\n' == l.rune || '\r' == l.rune {
		l.readRune()
	}
}

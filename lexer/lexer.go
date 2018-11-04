package lexer

import (
	"github.com/nrtkbb/go-MEL/token"
)

// Lexer は字句解析を行うための構造体
type Lexer struct {
	input        []rune // 字句解析対象のすべてのrune配列
	position     int    // 字句解析中のinputのインデックス
	readPosition int    // 字句解析中の一つ先のinputのインデックス
	rune         rune   // positionの位置にあるrune
	row          int    // 行数 1行はじまり
	column       int    // 列数 1列はじまり
}

// New はMELの文字列を受け取りLexerを生成して返す
func New(input string) *Lexer {
	l := &Lexer{
		input: []rune(input),
		row:   1, // 1行はじまり
	}
	l.readRune()
	return l
}

func (l *Lexer) readRune() {
	if l.readPosition >= len(l.input) {
		l.rune = 0
	} else {
		if '\n' == l.rune {
			l.row++
			l.column = 0
		}
		if '\r' == l.rune && '\n' != l.input[l.readPosition] {
			// '\r\n' の文章の '\r' の時はまだ改行しない
			l.row++
			l.column = 0
		}
		l.rune = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekRune() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func newToken(tokenType token.Type, r rune, row, column int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(r),
		Row:     row,
		Column:  column,
	}
}

// NextToken は実行される度に一つずつTokenを生成して返す
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.rune {
	case '=':
		tok = l.peekRuneCheck('=', token.Eq, token.Assign)
	case '!':
		tok = l.peekRuneCheck('=', token.NotEq, token.Bang)
	case '<':
		tok = l.peekRuneCheck('<', token.Ltensor, token.Lt)
	case '>':
		tok = l.peekRuneCheck('>', token.Rtensor, token.Gt)
	case '+':
		tok = l.peekRuneCheck('+', token.Increment, token.Plus)
	case '?':
		tok = newToken(token.Question, l.rune, l.row, l.column)
	case ':':
		tok = newToken(token.Coron, l.rune, l.row, l.column)
	case ';':
		tok = newToken(token.Semicolon, l.rune, l.row, l.column)
	case '`':
		tok = newToken(token.BackQuotes, l.rune, l.row, l.column)
	case '(':
		tok = newToken(token.Lparen, l.rune, l.row, l.column)
	case ')':
		tok = newToken(token.Rparen, l.rune, l.row, l.column)
	case '{':
		tok = newToken(token.Lbrace, l.rune, l.row, l.column)
	case '}':
		tok = newToken(token.Rbrace, l.rune, l.row, l.column)
	case '[':
		tok = newToken(token.Lbracket, l.rune, l.row, l.column)
	case ']':
		tok = newToken(token.Rbracket, l.rune, l.row, l.column)
	case ',':
		tok = newToken(token.Comma, l.rune, l.row, l.column)
	case '-':
		if 'a' <= l.peekRune() && l.peekRune() <= 'z' {
			tok.Type = token.Flag
			tok.Row = l.row
			tok.Column = l.column
			tok.Literal = l.readFlag()
			return tok
		}
		tok = l.peekRuneCheck('-', token.Decrement, token.Minus)
	case '/':
		if '/' == l.peekRune() {
			tok.Row = l.row
			tok.Column = l.column
			tok.Type = token.Comment
			tok.Literal = l.readLineComment()
			return tok
		} else if '*' == l.peekRune() {
			tok.Row = l.row
			tok.Column = l.column
			tok.Type = token.Comment
			tok.Literal = l.readComment()
			return tok
		}
		tok = newToken(token.Slash, l.rune, l.row, l.column)
	case '*':
		tok = newToken(token.Asterisk, l.rune, l.row, l.column)
	case '$':
		tok.Type = token.Ident
		tok.Row = l.row
		tok.Column = l.column
		tok.Literal = l.readIdentifier()
		return tok
	case '"':
		tok.Type = token.String
		tok.Row = l.row
		tok.Column = l.column
		tok.Literal = l.readString()
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Row = l.row
		tok.Column = l.column
		tok.Literal = ""
	default:
		if isLetter(l.rune) {
			tok.Row = l.row
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.rune) || '.' == l.rune && isDigit(l.peekRune()) {
			if '0' == l.rune && 'x' == l.peekRune() {
				tok.Type = token.Int16
				tok.Row = l.row
				tok.Column = l.column
				tok.Literal = l.readHexadecimalNumber()
			} else {
				tok.Row = l.row
				tok.Column = l.column
				tok.Type, tok.Literal = l.readNumber()
			}
			return tok
		} else {
			tok = newToken(token.Illegal, l.rune, l.row, l.column)
		}
	}

	l.readRune()
	return tok
}

func (l *Lexer) readLineComment() string {
	l.readRune() // '/'
	l.readRune() // ?
	position := l.position
	for !isNewLine(l.rune) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readComment() string {
	l.readRune() // '*'
	l.readRune() // ?
	position := l.position
	for !('*' == l.rune && '/' == l.peekRune()) {
		l.readRune()
		if l.position == len(l.input) {
			break
		}
	}
	comment := string(l.input[position:l.position])
	l.readRune() // '/'
	return comment
}

func isNewLine(r rune) bool {
	return '\n' == r || '\r' == r
}

func (l *Lexer) readFlag() string {
	position := l.position
	l.readRune() // '-'
	for isFlag(l.rune) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) peekRuneCheck(peek rune, trueType, falseType token.Type) token.Token {
	row := l.row
	column := l.column
	if peek == l.peekRune() {
		ch := l.rune
		l.readRune()
		literal := string(ch) + string(l.rune)
		return token.Token{
			Type:    trueType,
			Literal: literal,
			Row:     row,
			Column:  column,
		}
	}
	return newToken(falseType, l.rune, row, column)
}

func (l *Lexer) readString() string {
	position := l.position
	l.readRune() // '"'
	for '"' != l.rune && 0 != l.rune {
		if '\\' == l.rune {
			l.readRune() // '\\'
		}
		l.readRune()
	}
	l.readRune() // '"'
	return string(l.input[position:l.position])
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	l.readRune() // '$'
	for isIdentifier(l.rune) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func isFlag(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9'
}

func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '_' == r
}

func isIdentifier(r rune) bool {
	return isLetter(r) || '0' <= r && r <= '9'
}

func (l *Lexer) readHexadecimalNumber() string {
	position := l.position
	l.readRune() // '0'
	l.readRune() // 'x'
	for isHexadecimalDigit(l.rune) {
		l.readRune()
	}
	return string(l.input[position:l.position])
}

func isHexadecimalDigit(r rune) bool {
	return '0' <= r && r <= '9' || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F'
}

func (l *Lexer) readNumber() (token.Type, string) {
	var typ token.Type
	typ = token.Int
	position := l.position
	for isDigit(l.rune) {
		l.readRune()
	}
	if '.' == l.rune {
		typ = token.Float
		l.readRune()
		for isDigit(l.rune) {
			l.readRune()
		}
	}
	if 'e' == l.rune || 'E' == l.rune {
		if '-' == l.peekRune() || '+' == l.peekRune() {
			l.readRune() // 'e' or 'E'
			l.readRune() // '-' or '+'
			for isDigit(l.rune) {
				l.readRune()
			}
		}
	}
	return typ, string(l.input[position:l.position])
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func (l *Lexer) skipWhitespace() {
	for ' ' == l.rune || '\t' == l.rune || '\n' == l.rune || '\r' == l.rune {
		l.readRune()
	}
}

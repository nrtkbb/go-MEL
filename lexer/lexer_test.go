package lexer

import (
	"github.com/nrtkbb/go-MEL/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `int $five = 5;
int $ten = 10;

global proc add ( int $x, int $y ) {
	return $x + $y;
}
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "int"},
		{token.IDENT, "$five"},
		{token.ASSIGN, "="},
		{token.INT_DATA, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "int"},
		{token.IDENT, "$ten"},
		{token.ASSIGN, "="},
		{token.INT_DATA, "10"},
		{token.SEMICOLON, ";"},
		{token.GLOBAL, "global"},
		{token.PROC, "proc"},
		{token.PROC_IDENT, "add"},
		{token.LPAREN, "("},
		{token.INT, "int"},
		{token.IDENT, "$x"},
		{token.COMMA, ","},
		{token.INT, "int"},
		{token.IDENT, "$y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "$x"},
		{token.PLUS, "+"},
		{token.IDENT, "$y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT_DATA, "5"},
		{token.SEMICOLON, ";"},
		{token.INT_DATA, "5"},
		{token.LT, "<"},
		{token.INT_DATA, "10"},
		{token.GT, ">"},
		{token.INT_DATA, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT_DATA, "5"},
		{token.LT, "<"},
		{token.INT_DATA, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q %q",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

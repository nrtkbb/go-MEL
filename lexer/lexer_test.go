package lexer

import (
	"testing"

	"github.com/nrtkbb/go-MEL/token"
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

10 == 10;
10 != 9;

int $int16 = 0x123;
$int16 = 0xff;
$int16 = 0xFF;

float $float = 1.0;
$float = 1.0e-3;
$float = 1.0e+3;
$float = .01;

vector $vec = <<1, 2, 3.0>>;
matrix $mat[1][2] = <<1, 2; 3.0, 4.0>>;
string $str = "test \"test\"";
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
		{token.INT_DATA, "10"},
		{token.EQ, "=="},
		{token.INT_DATA, "10"},
		{token.SEMICOLON, ";"},
		{token.INT_DATA, "10"},
		{token.NOT_EQ, "!="},
		{token.INT_DATA, "9"},
		{token.SEMICOLON, ";"},
		{token.INT, "int"},
		{token.IDENT, "$int16"},
		{token.ASSIGN, "="},
		{token.INT_16DATA, "0x123"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "$int16"},
		{token.ASSIGN, "="},
		{token.INT_16DATA, "0xff"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "$int16"},
		{token.ASSIGN, "="},
		{token.INT_16DATA, "0xFF"},
		{token.SEMICOLON, ";"},
		{token.FLOAT, "float"},
		{token.IDENT, "$float"},
		{token.ASSIGN, "="},
		{token.FLOAT_DATA, "1.0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "$float"},
		{token.ASSIGN, "="},
		{token.FLOAT_DATA, "1.0e-3"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "$float"},
		{token.ASSIGN, "="},
		{token.FLOAT_DATA, "1.0e+3"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "$float"},
		{token.ASSIGN, "="},
		{token.FLOAT_DATA, ".01"},
		{token.SEMICOLON, ";"},
		{token.VECTOR, "vector"},
		{token.IDENT, "$vec"},
		{token.ASSIGN, "="},
		{token.LTENSOR, "<<"},
		{token.INT_DATA, "1"},
		{token.COMMA, ","},
		{token.INT_DATA, "2"},
		{token.COMMA, ","},
		{token.FLOAT_DATA, "3.0"},
		{token.RTENSOR, ">>"},
		{token.SEMICOLON, ";"},
		{token.MATRIX, "matrix"},
		{token.IDENT, "$mat"},
		{token.LBRACKET, "["},
		{token.INT_DATA, "1"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.INT_DATA, "2"},
		{token.RBRACKET, "]"},
		{token.ASSIGN, "="},
		{token.LTENSOR, "<<"},
		{token.INT_DATA, "1"},
		{token.COMMA, ","},
		{token.INT_DATA, "2"},
		{token.SEMICOLON, ";"},
		{token.FLOAT_DATA, "3.0"},
		{token.COMMA, ","},
		{token.FLOAT_DATA, "4.0"},
		{token.RTENSOR, ">>"},
		{token.SEMICOLON, ";"},
		{token.STRING, "string"},
		{token.IDENT, "$str"},
		{token.ASSIGN, "="},
		{token.STRING_DATA, "\"test \\\"test\\\"\""},
		{token.SEMICOLON, ";"},
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

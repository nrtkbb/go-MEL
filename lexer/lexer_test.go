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
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Int, "int"},
		{token.Ident, "$five"},
		{token.Assign, "="},
		{token.IntData, "5"},
		{token.Semicolon, ";"},
		{token.Int, "int"},
		{token.Ident, "$ten"},
		{token.Assign, "="},
		{token.IntData, "10"},
		{token.Semicolon, ";"},
		{token.Global, "global"},
		{token.Proc, "proc"},
		{token.ProcIdent, "add"},
		{token.Lparen, "("},
		{token.Int, "int"},
		{token.Ident, "$x"},
		{token.Comma, ","},
		{token.Int, "int"},
		{token.Ident, "$y"},
		{token.Rparen, ")"},
		{token.Lbrace, "{"},
		{token.Return, "return"},
		{token.Ident, "$x"},
		{token.Plus, "+"},
		{token.Ident, "$y"},
		{token.Semicolon, ";"},
		{token.Rbrace, "}"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.IntData, "5"},
		{token.Semicolon, ";"},
		{token.IntData, "5"},
		{token.Lt, "<"},
		{token.IntData, "10"},
		{token.Gt, ">"},
		{token.IntData, "5"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.Lparen, "("},
		{token.IntData, "5"},
		{token.Lt, "<"},
		{token.IntData, "10"},
		{token.Rparen, ")"},
		{token.Lbrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.Rbrace, "}"},
		{token.Else, "else"},
		{token.Lbrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.Rbrace, "}"},
		{token.IntData, "10"},
		{token.Eq, "=="},
		{token.IntData, "10"},
		{token.Semicolon, ";"},
		{token.IntData, "10"},
		{token.NotEq, "!="},
		{token.IntData, "9"},
		{token.Semicolon, ";"},
		{token.Int, "int"},
		{token.Ident, "$int16"},
		{token.Assign, "="},
		{token.Int16Data, "0x123"},
		{token.Semicolon, ";"},
		{token.Ident, "$int16"},
		{token.Assign, "="},
		{token.Int16Data, "0xff"},
		{token.Semicolon, ";"},
		{token.Ident, "$int16"},
		{token.Assign, "="},
		{token.Int16Data, "0xFF"},
		{token.Semicolon, ";"},
		{token.Float, "float"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.FloatData, "1.0"},
		{token.Semicolon, ";"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.FloatData, "1.0e-3"},
		{token.Semicolon, ";"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.FloatData, "1.0e+3"},
		{token.Semicolon, ";"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.FloatData, ".01"},
		{token.Semicolon, ";"},
		{token.Vector, "vector"},
		{token.Ident, "$vec"},
		{token.Assign, "="},
		{token.Ltensor, "<<"},
		{token.IntData, "1"},
		{token.Comma, ","},
		{token.IntData, "2"},
		{token.Comma, ","},
		{token.FloatData, "3.0"},
		{token.Rtensor, ">>"},
		{token.Semicolon, ";"},
		{token.Matrix, "matrix"},
		{token.Ident, "$mat"},
		{token.Lbracket, "["},
		{token.IntData, "1"},
		{token.Rbracket, "]"},
		{token.Lbracket, "["},
		{token.IntData, "2"},
		{token.Rbracket, "]"},
		{token.Assign, "="},
		{token.Ltensor, "<<"},
		{token.IntData, "1"},
		{token.Comma, ","},
		{token.IntData, "2"},
		{token.Semicolon, ";"},
		{token.FloatData, "3.0"},
		{token.Comma, ","},
		{token.FloatData, "4.0"},
		{token.Rtensor, ">>"},
		{token.Semicolon, ";"},
		{token.String, "string"},
		{token.Ident, "$str"},
		{token.Assign, "="},
		{token.StringData, "\"test \\\"test\\\"\""},
		{token.Semicolon, ";"},
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

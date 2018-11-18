package lexer

import (
	"testing"

	"github.com/nrtkbb/go-MEL/token"
)

func TestLexerExample1(t *testing.T) {
	input := `setAttr ($tforms[0] + ".translateZ") (-1.0 * $val);`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ProcIdent, "setAttr"},
		{token.Lparen, "("},
		{token.Ident, "$tforms"},
		{token.Lbracket, "["},
		{token.Int, "0"},
		{token.Rbracket, "]"},
		{token.Plus, "+"},
		{token.String, `".translateZ"`},
		{token.Rparen, ")"},
		{token.Lparen, "("},
		{token.Minus, "-"},
		{token.Float, "1.0"},
		{token.Asterisk, "*"},
		{token.Ident, "$val"},
		{token.Rparen, ")"},
		{token.Semicolon, ";"},
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

func TestIsLetterToken(t *testing.T) {
	input := `
|all|body|;
($ident1||$ident2);
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ProcIdent, "|all|body|"},
		{token.Semicolon, ";"},
		{token.Lparen, "("},
		{token.Ident, "$ident1"},
		{token.Or, "||"},
		{token.Ident, "$ident2"},
		{token.Rparen, ")"},
		{token.Semicolon, ";"},
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

func TestBooleanToken(t *testing.T) {
	input := `
true;
false;
on;
off;
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.On, "on"},
		{token.Semicolon, ";"},
		{token.Off, "off"},
		{token.Semicolon, ";"},
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

func TestTernaryOperator(t *testing.T) {
	input := `
int $a = 1;
int $b = 2;
int $c = $a < $b ? 10 : 20;
`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IntDec, "int"},
		{token.Ident, "$a"},
		{token.Assign, "="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.IntDec, "int"},
		{token.Ident, "$b"},
		{token.Assign, "="},
		{token.Int, "2"},
		{token.Semicolon, ";"},
		{token.IntDec, "int"},
		{token.Ident, "$c"},
		{token.Assign, "="},
		{token.Ident, "$a"},
		{token.Lt, "<"},
		{token.Ident, "$b"},
		{token.Question, "?"},
		{token.Int, "10"},
		{token.Coron, ":"},
		{token.Int, "20"},
		{token.Semicolon, ";"},
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

func TestIncrementDecrement(t *testing.T) {
	input := `
int $i = 0;
--$i;
$i--;
++$i;
$i++;
`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IntDec, "int"},
		{token.Ident, "$i"},
		{token.Assign, "="},
		{token.Int, "0"},
		{token.Semicolon, ";"},
		{token.Decrement, "--"},
		{token.Ident, "$i"},
		{token.Semicolon, ";"},
		{token.Ident, "$i"},
		{token.Decrement, "--"},
		{token.Semicolon, ";"},
		{token.Increment, "++"},
		{token.Ident, "$i"},
		{token.Semicolon, ";"},
		{token.Ident, "$i"},
		{token.Increment, "++"},
		{token.Semicolon, ";"},
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

func TestStringLiteralToken(t *testing.T) {
	input := `
"s" + "s";
`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.String, "\"s\""},
		{token.Plus, "+"},
		{token.String, "\"s\""},
		{token.Semicolon, ";"},
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

func TestCommentToken(t *testing.T) {
	input := `
// test
/* test */
`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.EOF {
		t.Fatalf("input is not EOF. got=%T", tok.Type)
	}
}

func TestNextToken(t *testing.T) {
	input := `int $five = 5;
int $ten = 10;

global proc add ( int $x, int $y ) {
	return $x + $y;
}
!-/ *5;
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

getAttr -s "pCubeShape1.fc";
string $ls[] = ` + "`ls -sl`" + `;

{1, 2};

int $t = true && false || true;
$t += 1;
$t -= 1;
$t *= 2;
$t /= 2;
$t <= 1;
$t >= 1;
1e-6;
$t.x;
setParent ..;
$t ^ $t;
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IntDec, "int"},
		{token.Ident, "$five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.IntDec, "int"},
		{token.Ident, "$ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Global, "global"},
		{token.Proc, "proc"},
		{token.ProcIdent, "add"},
		{token.Lparen, "("},
		{token.IntDec, "int"},
		{token.Ident, "$x"},
		{token.Comma, ","},
		{token.IntDec, "int"},
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
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Int, "5"},
		{token.Lt, "<"},
		{token.Int, "10"},
		{token.Gt, ">"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.Lparen, "("},
		{token.Int, "5"},
		{token.Lt, "<"},
		{token.Int, "10"},
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
		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.NotEq, "!="},
		{token.Int, "9"},
		{token.Semicolon, ";"},
		{token.IntDec, "int"},
		{token.Ident, "$int16"},
		{token.Assign, "="},
		{token.Int16, "0x123"},
		{token.Semicolon, ";"},
		{token.Ident, "$int16"},
		{token.Assign, "="},
		{token.Int16, "0xff"},
		{token.Semicolon, ";"},
		{token.Ident, "$int16"},
		{token.Assign, "="},
		{token.Int16, "0xFF"},
		{token.Semicolon, ";"},
		{token.FloatDec, "float"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.Float, "1.0"},
		{token.Semicolon, ";"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.Float, "1.0e-3"},
		{token.Semicolon, ";"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.Float, "1.0e+3"},
		{token.Semicolon, ";"},
		{token.Ident, "$float"},
		{token.Assign, "="},
		{token.Float, ".01"},
		{token.Semicolon, ";"},
		{token.VectorDec, "vector"},
		{token.Ident, "$vec"},
		{token.Assign, "="},
		{token.Ltensor, "<<"},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.Comma, ","},
		{token.Float, "3.0"},
		{token.Rtensor, ">>"},
		{token.Semicolon, ";"},
		{token.MatrixDec, "matrix"},
		{token.Ident, "$mat"},
		{token.Lbracket, "["},
		{token.Int, "1"},
		{token.Rbracket, "]"},
		{token.Lbracket, "["},
		{token.Int, "2"},
		{token.Rbracket, "]"},
		{token.Assign, "="},
		{token.Ltensor, "<<"},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.Semicolon, ";"},
		{token.Float, "3.0"},
		{token.Comma, ","},
		{token.Float, "4.0"},
		{token.Rtensor, ">>"},
		{token.Semicolon, ";"},
		{token.StringDec, "string"},
		{token.Ident, "$str"},
		{token.Assign, "="},
		{token.String, "\"test \\\"test\\\"\""},
		{token.Semicolon, ";"},
		{token.ProcIdent, "getAttr"},
		{token.Flag, "-s"},
		{token.String, "\"pCubeShape1.fc\""},
		{token.Semicolon, ";"},
		{token.StringDec, "string"},
		{token.Ident, "$ls"},
		{token.Lbracket, "["},
		{token.Rbracket, "]"},
		{token.Assign, "="},
		{token.BackQuotes, "`"},
		{token.ProcIdent, "ls"},
		{token.Flag, "-sl"},
		{token.BackQuotes, "`"},
		{token.Semicolon, ";"},
		{token.Lbrace, "{"},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.Rbrace, "}"},
		{token.Semicolon, ";"},
		{token.IntDec, "int"},
		{token.Ident, "$t"},
		{token.Assign, "="},
		{token.True, "true"},
		{token.And, "&&"},
		{token.False, "false"},
		{token.Or, "||"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.PAssign, "+="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.MAssign, "-="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.AAssign, "*="},
		{token.Int, "2"},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.SAssign, "/="},
		{token.Int, "2"},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.LtEq, "<="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.GtEq, ">="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.Float, "1e-6"},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.Dot, "."},
		{token.ProcIdent, "x"},
		{token.Semicolon, ";"},
		{token.ProcIdent, "setParent"},
		{token.ProcIdent, ".."},
		{token.Semicolon, ";"},
		{token.Ident, "$t"},
		{token.Hat, "^"},
		{token.Ident, "$t"},
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

func TestLineCount(t *testing.T) {
	input := `int $five = 5;
int $ten = 10;
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
		expectedRow     int
		expectedColumn  int
	}{
		{token.IntDec, "int", 1, 1},
		{token.Ident, "$five", 1, 5},
		{token.Assign, "=", 1, 11},
		{token.Int, "5", 1, 13},
		{token.Semicolon, ";", 1, 14},
		{token.IntDec, "int", 2, 1},
		{token.Ident, "$ten", 2, 5},
		{token.Assign, "=", 2, 10},
		{token.Int, "10", 2, 12},
		{token.Semicolon, ";", 2, 14},
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

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%q, got=%q",
				i, tt.expectedColumn, tok.Column)
		}

		if tok.Row != tt.expectedRow {
			t.Fatalf("tests[%d] - row wrong. expected=%q, got=%q",
				i, tt.expectedRow, tok.Row)
		}
	}
}

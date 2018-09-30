package parser

import (
	"github.com/nrtkbb/go-MEL/ast"
	"testing"

	"github.com/nrtkbb/go-MEL/lexer"
)

func TestStringStatement(t *testing.T) {
	input := `
string $x = "x";
string $y = "y";
string $foobar = "foobar";
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"$x"},
		{"$y"},
		{"$foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testStringStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testStringStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "string" {
		t.Errorf("s.TokenLiteral not 'string'. got=%q", s.TokenLiteral())
		return false
	}

	stringStmt, ok := s.(*ast.StringStatement)
	if !ok {
		t.Errorf("s not *ast.StringStatement. got=%T", s)
		return false
	}

	if stringStmt.Name.Value != name {
		t.Errorf("stringStmt.Name.Value not '%s'. got=%s", name, stringStmt.Name.Value)
		return false
	}

	if stringStmt.Name.TokenLiteral() != name {
		t.Errorf("stringStmt.Name.Value not '%s'. got=%s", name, stringStmt.Name.Value)
		return false
	}

	return true
}

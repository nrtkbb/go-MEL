package parser

import (
	"testing"

	"github.com/nrtkbb/go-MEL/ast"
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
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
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

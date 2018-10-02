package parser

import (
	"testing"

	"github.com/nrtkbb/go-MEL/ast"
	"github.com/nrtkbb/go-MEL/lexer"
)

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s",
			"5", literal.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `$foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "$foobar" {
		t.Errorf("ident.Value not %s. got=%s",
			"$foobar", ident.Value)
	}
	if ident.TokenLiteral() != "$foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s",
			"$foobar", ident.TokenLiteral())
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
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

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("resutnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

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

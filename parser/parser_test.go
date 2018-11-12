package parser

import (
	"fmt"
	"testing"

	"github.com/nrtkbb/go-MEL/ast"
	"github.com/nrtkbb/go-MEL/lexer"
)

func TestCommentParsing(t *testing.T) {
	input := `// test
/*
	test test
*/`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 0 {
		t.Fatalf("len(program.Statements) does not 0. got=%d",
			len(program.Statements))
	}
}

func TestSwitchStatement(t *testing.T) {
	input := `
switch ($x + 1) {
	case 1:
		string $x = "x";
		break;
	default:
		string $y = "y";
		break;
}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) does not 1. got=%d",
			len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statemtns[0] does not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	se, ok := exp.Expression.(*ast.SwitchExpression)
	if !ok {
		t.Fatalf("exp.Expression does not *ast.SwitchExpression. got=%T",
			exp.Expression)
	}

	if !testInfixExpression(t, se.Condition, "$x", "+", 1) {
		return
	}

	if len(se.Cases) != 2 {
		t.Fatalf("len(se.Cases) is not 2. got=%d", len(se.Cases))
	}

	if len(se.CaseStatements) != 2 {
		t.Fatalf("len(se.CaseStatements) is not 2. got=%d", len(se.CaseStatements))
	}

	if !testLiteralExpression(t, se.Cases[0], 1) {
		return
	}

	if se.Cases[1] != nil {
		t.Fatalf("se.Cases[1] is not nil. got=%q", se.Cases[1])
	}
}

func TestParsingContinueStatement(t *testing.T) {
	input := "continue;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) does not 1. got=%d",
			len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("program.Statemtns[0] does not *ast.ContinueStatement. got=%T",
			program.Statements[0])
	}
}

func TestParsingBreakStatement(t *testing.T) {
	input := "break;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) does not 1. got=%d",
			len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("program.Statemtns[0] does not *ast.BreakStatemnt. got=%T",
			program.Statements[0])
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "$myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "$myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingArrayLiteral(t *testing.T) {
	input := "{1, 2 * 2, 3 + 3};"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3, got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestCallExpressionParsing3(t *testing.T) {
	input := "add 1 (2 + 3) `add 1 2 a \"b\"` a \"b\";"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T\n",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 5 {
		t.Fatalf("wrong length of arguments. got=%d\n", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "+", 3)
	testIdentifier(t, exp.Arguments[3], "a")
	testLiteralExpression(t, exp.Arguments[4], `"b"`)
}

func TestCallExpressionParsing2(t *testing.T) {
	input := "`add 1 (2 + 3) x $y`;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T\n",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 4 {
		t.Fatalf("wrong length of arguments. got=%d\n", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "+", 3)
	testIdentifier(t, exp.Arguments[2], "x")
	testLiteralExpression(t, exp.Arguments[3], "$y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T\n",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d\n", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestGlobalStatementParsing(t *testing.T) {
	input := `
global proc Proc(string $x, string $y) {
    $x + $y;
}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
			1, len(program.Statements))
	}

	gs, ok := program.Statements[0].(*ast.GlobalStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.GlobalStatement. got=%T\n",
			program.Statements[0])
	}

	stmt, ok := gs.Statement.(*ast.ProcStatement)
	if !ok {
		t.Fatalf("gs.Statement is not ast.ProcStatement. got=%T\n",
			program.Statements[0])
	}

	if stmt.Name.Literal != "Proc" {
		t.Fatalf("function name is not %s. got=%s\n",
			"Proc", stmt.Name.Literal)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong, want 2, got=%d\n",
			len(stmt.Parameters))
	}

	testLiteralExpression(t, stmt.Parameters[0], "$x")
	testLiteralExpression(t, stmt.Parameters[1], "$y")

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExressionStatement. got=%T\n",
			stmt.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "$x", "+", "$y")
}

func TestProcStatementParsing(t *testing.T) {
	input := `
proc Proc(string $x, string $y) {
	$x + $y;
}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ProcStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ProcStatement. got=%T\n",
			program.Statements[0])
	}

	if stmt.Name.Literal != "Proc" {
		t.Fatalf("proc name is not %s. got=%s\n",
			"Proc", stmt.Name.Literal)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("proc statement parameters wrong, want 2, got=%d\n",
			len(stmt.Parameters))
	}

	testLiteralExpression(t, stmt.Parameters[0], "$x")
	testLiteralExpression(t, stmt.Parameters[1], "$y")

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("stmt.Body.Statements has not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt body stmt is not ast.ExressionStatement. got=%T\n",
			stmt.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "$x", "+", "$y")
}

func TestForInExpressionSingleBlock(t *testing.T) {
	input := `for ($i in $array) string $x = "x";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.ForInExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T\n",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Element, "$i") {
		return
	}

	if !testIdentifier(t, exp.ArrayElement, "$array") {
		return
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.StringStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.StringStatement, got=%T\n",
			exp.Consequence.Statements[0])
	}

	if !testStringStatement(t, consequence, "$x") {
		return
	}

	val := consequence.Values[0]
	if !testStringLiteral(t, val, `"x"`) {
		return
	}

	return
}

func TestForExpressionSingleBlock(t *testing.T) {
	input := `for ($i = 0; $i < 10; $i++) string $x = "x";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T\n",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.InitNames[0], "$i") {
		return
	}

	if !testLiteralExpression(t, exp.InitValues[0], 0) {
		return
	}

	if !testInfixExpression(t, exp.Condition, "$i", "<", 10) {
		return
	}

	if exp.ChangeOf.String() != "($i++)" {
		t.Errorf("exp.ChangeOf is not '%s'. got=%s", "($i++)", exp.ChangeOf.String())
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.StringStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.StringStatement, got=%T\n",
			exp.Consequence.Statements[0])
	}

	if !testStringStatement(t, consequence, "$x") {
		return
	}

	val := consequence.Values[0]
	if !testStringLiteral(t, val, `"x"`) {
		return
	}

	return
}

func TestDoWhileExpression(t *testing.T) {
	inputs := []string{
		`
do {
    string $x = "x";
} while ($x < $y);`,
		`
do
    string $x = "x";
while ($x < $y);`,
	}

	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.DoWhileExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.DoWhileExpression. got=%T\n",
				stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "$x", "<", "$y") {
			return
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.StringStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.StringStatement, got=%T\n",
				exp.Consequence.Statements[0])
		}

		if !testStringStatement(t, consequence, "$x") {
			return
		}

		val := consequence.Values[0]
		if !testStringLiteral(t, val, `"x"`) {
			return
		}
	}
}

func TestWhileExpressionSingleBlock(t *testing.T) {
	input := `while ($x < $y) string $x = "x";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.WhileExpression. got=%T\n",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "$x", "<", "$y") {
		return
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.StringStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.StringStatement, got=%T\n",
			exp.Consequence.Statements[0])
	}

	if !testStringStatement(t, consequence, "$x") {
		return
	}

	val := consequence.Values[0]
	if !testStringLiteral(t, val, `"x"`) {
		return
	}

	return
}

func TestIfExpressionSingleBlock(t *testing.T) {
	input := `if ($x < $y) string $x = "x";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T\n",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "$x", "<", "$y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.StringStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.StringStatement, got=%T\n",
			exp.Consequence.Statements[0])
	}

	if !testStringStatement(t, consequence, "$x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v\n",
			exp.Alternative)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if ($x < $y) { string $x = "x"; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T\n",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "$x", "<", "$y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.StringStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.StringStatement, got=%T\n",
			exp.Consequence.Statements[0])
	}

	if !testStringStatement(t, consequence, "$x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v\n",
			exp.Alternative)
	}
}

func TestOperatorPrecendenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 * (2 + 3)", "(1 * (2 + 3))"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"-(1 + 2)", "(-(1 + 2))"},
		{"!(true == true)", "(!(true == true))"},
		{"!(on == on)", "(!(on == on))"},
		{"!(on != off)", "(!(on != off))"},
		{"!(on != off) && true == false", "((!(on != off)) && (true == false))"},
		{"$a + add($b * $c) + $d", "(($a + add(($b * $c))) + $d)"},
		{"add($a, $b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add($a, $b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add($a + $b + $c * $d / $f + $g)",
			"add(((($a + $b) + (($c * $d) / $f)) + $g))"},
		{"add($a * $b[2], $b[1], 2 * 1)",
			"add(($a * ($b[2])), ($b[1]), (2 * 1))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestTernaryOperatorParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 < 2 ? 1 : 2", "((1 < 2) ? 1 : 2)"},
		{"1 < 2 ? 1 + 1 : 2 * 2", "((1 < 2) ? (1 + 1) : (2 * 2))"},
		{"1 < 2 ? $i++ : --$i", "((1 < 2) ? ($i++) : (--$i))"},
		{"1 < 2 ? -1 : 1 + -2", "((1 < 2) ? (-1) : (1 + (-2)))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestPostfixParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"--$i", "(--$i)"},
		{"$i-- + $a", "(($i--) + $a)"},
		{"$i-- * $a", "(($i--) * $a)"},
		{"-$i-- * $a", "(((-$i)--) * $a)"},
		{"$i++ + $a", "(($i++) + $a)"},
		{"$i++ * $a", "(($i++) * $a)"},
		{"-$i++ * $a", "(((-$i)++) * $a)"},
		{"$i++ + $a * $b", "(($i++) + ($a * $b))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-$a * $b", "((-$a) * $b)"},
		{"!-$a", "(!(-$a))"},
		{"$a + $b + $c", "(($a + $b) + $c)"},
		{"$a + $b - $c", "(($a + $b) - $c)"},
		{"$a * $b * $c", "(($a * $b) * $c)"},
		{"$a * $b / $c", "(($a * $b) / $c)"},
		{"$a + $b * $c + $d / $e - $f", "((($a + ($b * $c)) + ($d / $e)) - $f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParseInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"\"s\" + \"a\"", "\"s\"", "+", "\"a\""},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
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

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
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

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

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
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{`return 5;`, 5},
		{`return true;`, true},
		{`return "foobar"`, `"foobar"`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]

		val := stmt.(*ast.ReturnStatement).ReturnValue
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestMatrixStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		i1                 int64
		i2                 int64
		expectedValue      [][]float64
	}{
		{`matrix $x[2][3] = <<1, 2>>;`, "$x", 3, 2, nil},
		{`matrix $y[1][1] = <<1>>;`, "$y", 1, 1, nil},
		{`matrix $foobar[1][1] = <<123123>>`, "$foobar", 1, 1, nil},
	}

	tests[0].expectedValue = append(tests[0].expectedValue, []float64{1, 2})
	tests[1].expectedValue = append(tests[1].expectedValue, []float64{1})
	tests[2].expectedValue = append(tests[2].expectedValue, []float64{123123})

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testMatrixStatement(t, stmt, tt.expectedIdentifier, tt.i1, tt.i2) {
			return
		}

		val := stmt.(*ast.MatrixStatement).Values[0]
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestMatrixStatement2(t *testing.T) {
	tests := []struct {
		input               string
		expectedIdentifiers []string
		expectedValues      []interface{}
	}{
		{`matrix $x = <<1, 2, 3>>, $y, $z = <<2, 3, 4>>;`, []string{"$x", "$y", "$z"}, nil},
		{`matrix $x, $y = <<1, 2, 3>>, $z = <<2, 3, 4>>;`, []string{"$x", "$y", "$z"}, nil},
	}

	tests[0].expectedValues = append(tests[0].expectedValues, [][]float64{{1, 2, 3}})
	tests[0].expectedValues = append(tests[0].expectedValues, nil)
	tests[0].expectedValues = append(tests[0].expectedValues, [][]float64{{2, 3, 4}})
	tests[1].expectedValues = append(tests[1].expectedValues, nil)
	tests[1].expectedValues = append(tests[1].expectedValues, [][]float64{{1, 2, 3}})
	tests[1].expectedValues = append(tests[1].expectedValues, [][]float64{{2, 3, 4}})

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testMatrixStatement2(t, stmt, tt.expectedIdentifiers) {
			return
		}

		vals := stmt.(*ast.MatrixStatement).Values
		for i, val := range vals {
			if !testLiteralExpression(t, val, tt.expectedValues[i]) {
				return
			}
		}
	}
}

func TestVectorStatement2(t *testing.T) {
	tests := []struct {
		input               string
		expectedIdentifiers []string
		expectedValues      []interface{}
	}{
		{`vector $x = <<1, 2, 3>>, $y, $z = <<2, 3, 4>>;`, []string{"$x", "$y", "$z"}, nil},
		{`vector $x, $y = <<1, 2, 3>>, $z = <<2, 3, 4>>;`, []string{"$x", "$y", "$z"}, nil},
	}

	tests[0].expectedValues = append(tests[0].expectedValues, [][]float64{{1, 2, 3}})
	tests[0].expectedValues = append(tests[0].expectedValues, nil)
	tests[0].expectedValues = append(tests[0].expectedValues, [][]float64{{2, 3, 4}})
	tests[1].expectedValues = append(tests[1].expectedValues, nil)
	tests[1].expectedValues = append(tests[1].expectedValues, [][]float64{{1, 2, 3}})
	tests[1].expectedValues = append(tests[1].expectedValues, [][]float64{{2, 3, 4}})

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVectorStatement2(t, stmt, tt.expectedIdentifiers) {
			return
		}

		vals := stmt.(*ast.VectorStatement).Values
		for i, val := range vals {
			if !testLiteralExpression(t, val, tt.expectedValues[i]) {
				return
			}
		}
	}
}

func TestFloatStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{`float $x = 123.123;`, "$x", 123.123},
		{`float $y = 0.1;`, "$y", 0.1},
		{`float $foobar = 0;`, "$foobar", 0},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testFloatStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.FloatStatement).Values[0]
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestIntegerStatement2(t *testing.T) {
	tests := []struct {
		input               string
		expectedIdentifiers []string
		expectedValues      []interface{}
	}{
		{`int $x = 5, $y, $z = 6`, []string{"$x", "$y", "$z"}, nil},
		{`int $x, $y = 5, $z = 6`, []string{"$x", "$y", "$z"}, nil},
	}

	tests[0].expectedValues = append(tests[0].expectedValues, 5)
	tests[0].expectedValues = append(tests[0].expectedValues, nil)
	tests[0].expectedValues = append(tests[0].expectedValues, 6)
	tests[1].expectedValues = append(tests[1].expectedValues, nil)
	tests[1].expectedValues = append(tests[1].expectedValues, 5)
	tests[1].expectedValues = append(tests[1].expectedValues, 6)

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testIntegerStatement2(t, stmt, tt.expectedIdentifiers) {
			return
		}

		vals := stmt.(*ast.IntegerStatement).Values
		for i, val := range vals {
			if !testLiteralExpression(t, val, tt.expectedValues[i]) {
				return
			}
		}
	}
}

func TestIntegerStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{`int $x = 5;`, "$x", 5},
		{`int $y = 0;`, "$y", 0},
		{`int $foobar = $y;`, "$foobar", "$y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testIntegerStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.IntegerStatement).Values[0]
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestStringStatement2(t *testing.T) {
	tests := []struct {
		input               string
		expectedIdentifiers []string
		expectedValues      []interface{}
	}{
		{`string $x = "5", $y, $z = "6"`, []string{"$x", "$y", "$z"}, nil},
		{`string $x, $y = "5", $z = "6"`, []string{"$x", "$y", "$z"}, nil},
	}

	tests[0].expectedValues = append(tests[0].expectedValues, `"5"`)
	tests[0].expectedValues = append(tests[0].expectedValues, nil)
	tests[0].expectedValues = append(tests[0].expectedValues, `"6"`)
	tests[1].expectedValues = append(tests[1].expectedValues, nil)
	tests[1].expectedValues = append(tests[1].expectedValues, `"5"`)
	tests[1].expectedValues = append(tests[1].expectedValues, `"6"`)

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testStringStatement2(t, stmt, tt.expectedIdentifiers) {
			return
		}

		vals := stmt.(*ast.StringStatement).Values
		for i, val := range vals {
			if !testLiteralExpression(t, val, tt.expectedValues[i]) {
				return
			}
		}
	}
}

func TestStringStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{`string $x = "5";`, "$x", `"5"`},
		{`string $y = "true";`, "$y", `"true"`},
		{`string $foobar = $y;`, "$foobar", `$y`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testStringStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.StringStatement).Values[0]
		if !testLiteralExpression(t, val, tt.expectedValue) {
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

func testMatrixStatement2(t *testing.T, s ast.Statement, names []string) bool {
	if s.TokenLiteral() != "matrix" {
		t.Errorf("s.TokenLiteral not 'matrix'. got=%q", s.TokenLiteral())
		return false
	}

	matStmt, ok := s.(*ast.MatrixStatement)
	if !ok {
		t.Errorf("s not *ast.MatrixStatement. got=%T", s)
		return false
	}

	for i, stmtName := range matStmt.Names {
		nameIdent, ok := stmtName.(*ast.Identifier)
		if !ok {
			t.Errorf("matStmt.Name does not ast.Identifier. got=%T", stmtName)
			return false
		}

		if nameIdent.Value != names[i] {
			t.Errorf("nameIdent.Value not '%s'. got=%s", names[i], nameIdent.Value)
			return false
		}

		if nameIdent.TokenLiteral() != names[i] {
			t.Errorf("nameIdent.TokenLiteral not '%s'. got=%s", names[i], nameIdent.TokenLiteral())
			return false
		}
	}

	return true
}

func testVectorStatement2(t *testing.T, s ast.Statement, names []string) bool {
	if s.TokenLiteral() != "vector" {
		t.Errorf("s.TokenLiteral not 'vector'. got=%q", s.TokenLiteral())
		return false
	}

	vecStmt, ok := s.(*ast.VectorStatement)
	if !ok {
		t.Errorf("s not *ast.VectorStatement. got=%T", s)
		return false
	}

	for i, stmtName := range vecStmt.Names {
		nameIdent, ok := stmtName.(*ast.Identifier)
		if !ok {
			t.Errorf("vecStmt.Name does not ast.Identifier. got=%T", stmtName)
			return false
		}

		if nameIdent.Value != names[i] {
			t.Errorf("nameIdent.Value not '%s'. got=%s", names[i], nameIdent.Value)
			return false
		}

		if nameIdent.TokenLiteral() != names[i] {
			t.Errorf("nameIdent.TokenLiteral not '%s'. got=%s", names[i], nameIdent.TokenLiteral())
			return false
		}
	}

	return true
}

func testIntegerStatement2(t *testing.T, s ast.Statement, names []string) bool {
	if s.TokenLiteral() != "int" {
		t.Errorf("s.TokenLiteral not 'int'. got=%q", s.TokenLiteral())
		return false
	}

	intStmt, ok := s.(*ast.IntegerStatement)
	if !ok {
		t.Errorf("s not *ast.IntegerStatement. got=%T", s)
		return false
	}

	for i, stmtName := range intStmt.Names {
		nameIdent, ok := stmtName.(*ast.Identifier)
		if !ok {
			t.Errorf("stmtName does not ast.Identifier. got=%T", stmtName)
			return false
		}

		if nameIdent.Value != names[i] {
			t.Errorf("nameIndent.Value not '%s'. got=%s", names[i], nameIdent.Value)
			return false
		}

		if nameIdent.TokenLiteral() != names[i] {
			t.Errorf("nameIndent.TokenLiteral not '%s'. got=%s", names[i], nameIdent.TokenLiteral())
			return false
		}
	}

	return true
}

func testMatrixStatement(t *testing.T, s ast.Statement, name string, i1 int64, i2 int64) bool {
	if s.TokenLiteral() != "matrix" {
		t.Errorf("s.TokenLiteral not 'matrix'. got=%q", s.TokenLiteral())
		return false
	}

	matStmt, ok := s.(*ast.MatrixStatement)
	if !ok {
		t.Errorf("s not *ast.MatrixStatement. got=%T", s)
		return false
	}

	idxExp1, ok := matStmt.Names[0].(*ast.IndexExpression)
	if !ok {
		t.Errorf("matStmt.Names[0] does not ast.IndexExpression. got=%T", matStmt.Names[0])
		return false
	}

	idx1, ok := idxExp1.Index.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("idxExp1.Index does not ast.IntegerLiteral. got=%T", idxExp1.Index)
		return false
	}

	if idx1.Value != i1 {
		t.Errorf("idx1.Value not %d. got=%d", i1, idx1.Value)
		return false
	}

	idxExp2, ok := idxExp1.Left.(*ast.IndexExpression)
	if !ok {
		t.Errorf("idxExp1.Left does not ast.IndexExpression. got=%T", idxExp1.Left)
		return false
	}

	idx2, ok := idxExp2.Index.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("idxExp2.Index does not ast.IntegerLiteral. got=%T", idxExp2.Index)
		return false
	}

	if idx2.Value != i2 {
		t.Errorf("idx2.Value not %d. got=%d", i2, idx2.Value)
		return false
	}

	nameIdent, ok := idxExp2.Left.(*ast.Identifier)
	if !ok {
		t.Errorf("idxExp2.Left does not ast.Identifier. got=%T", idxExp2.Left)
		return false
	}

	if nameIdent.Value != name {
		t.Errorf("nameIdent.Value not '%s'. got=%s", name, nameIdent.Value)
		return false
	}

	if nameIdent.TokenLiteral() != name {
		t.Errorf("nameIdent.TokenLiteral() not '%s'. got=%s", name, nameIdent.Value)
		return false
	}

	return true
}

func testIntegerStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "int" {
		t.Errorf("s.TokenLiteral not 'int'. got=%q", s.TokenLiteral())
		return false
	}

	intStmt, ok := s.(*ast.IntegerStatement)
	if !ok {
		t.Errorf("s not *ast.IntegerStatement. got=%T", s)
		return false
	}

	nameIdent, ok := intStmt.Names[0].(*ast.Identifier)
	if !ok {
		t.Errorf("intStmt.Names[0] does not ast.Identifier. got=%T", intStmt.Names[0])
		return false
	}

	if nameIdent.Value != name {
		t.Errorf("nameIdent.Value not '%s'. got=%s", name, nameIdent.Value)
		return false
	}

	if nameIdent.TokenLiteral() != name {
		t.Errorf("nameIdent.TokenLiteral not '%s'. got=%s", name, nameIdent.Value)
		return false
	}

	return true
}

func testFloatStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "float" {
		t.Errorf("s.TokenLiteral not 'float'. got=%q", s.TokenLiteral())
		return false
	}

	intStmt, ok := s.(*ast.FloatStatement)
	if !ok {
		t.Errorf("s not *ast.FloatStatement. got=%T", s)
		return false
	}

	nameIdent, ok := intStmt.Names[0].(*ast.Identifier)
	if !ok {
		t.Errorf("intStmt.Names[0] does not ast.Identifier. got=%T", intStmt.Names[0])
		return false
	}

	if nameIdent.Value != name {
		t.Errorf("nameIdent.Value not '%s'. got=%s", name, nameIdent.Value)
		return false
	}

	if nameIdent.TokenLiteral() != name {
		t.Errorf("nameIdent.TokenLiteral not '%s'. got=%s", name, nameIdent.Value)
		return false
	}

	return true
}

func testStringStatement2(t *testing.T, s ast.Statement, names []string) bool {
	if s.TokenLiteral() != "string" {
		t.Errorf("s.TokenLiteral not 'string'. got=%q", s.TokenLiteral())
		return false
	}

	stringStmt, ok := s.(*ast.StringStatement)
	if !ok {
		t.Errorf("s not *ast.StringStatement. got=%T", s)
		return false
	}

	for i, stmtName := range stringStmt.Names {
		nameIdent, ok := stmtName.(*ast.Identifier)
		if !ok {
			t.Errorf("stmtName does not ast.Identifier. got=%T", stmtName)
			return false
		}

		if nameIdent.Value != names[i] {
			t.Errorf("nameIdent.Value not '%s'. got=%s", names[i], nameIdent.Value)
			return false
		}

		if nameIdent.TokenLiteral() != names[i] {
			t.Errorf("nameIdent.TokenLiteral not '%s'. got=%s", names[i], nameIdent.TokenLiteral())
			return false
		}
	}

	return true
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

	nameIdent, ok := stringStmt.Names[0].(*ast.Identifier)
	if !ok {
		t.Errorf("stringStmt.Names[0] does not ast.Identifier. got=%T", stringStmt.Names[0])
		return false
	}
	if nameIdent.Value != name {
		t.Errorf("nameIdent.Value not '%s'. got=%s", name, nameIdent.Value)
		return false
	}

	if nameIdent.TokenLiteral() != name {
		t.Errorf("nameIdent.TokenLiteral not '%s'. got=%s", name, nameIdent.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Eperator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, v)
	case string:
		if v[0] == '"' {
			return testStringLiteral(t, exp, v)
		}
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case nil:
		return testNil(t, exp, v)
	case [][]float64:
		return testTensorLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testTensorLiteral(t *testing.T, vl ast.Expression, value [][]float64) bool {
	vector, ok := vl.(*ast.TensorLiteral)
	if !ok {
		t.Errorf("il not *ast.TensorLiteral. got=%T", vl)
		return false
	}

	for i, vv := range vector.Values {
		for ii, vvv := range vv {
			if vvv.String() != fmt.Sprint(value[i][ii]) {
				t.Errorf("vector.Values[%d][%d] not %f, got=%f", i, ii, value[i][ii], vvv)
				return false
			}
		}
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d, got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s",
			value, integ.TokenLiteral())
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, flexp ast.Expression, value float64) bool {
	flLit, ok := flexp.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("flexp not *ast.IntegerLiteral. got=%T", flexp)
		return false
	}

	if flLit.Value != value {
		t.Errorf("flLit.Value not %f, got=%f", value, flLit.Value)
		return false
	}

	if flLit.TokenLiteral() != fmt.Sprint(value) {
		t.Errorf("flLit.TokenLiteral not %f. got=%s",
			value, flLit.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, sl ast.Expression, value string) bool {
	st, ok := sl.(*ast.StringLiteral)
	if !ok {
		t.Errorf("sl not *ast.StringLiteral. got=%T", sl)
		return false
	}

	if st.Value != value {
		t.Errorf("st.Value not %s, got=%s", value, st.Value)
		return false
	}

	if st.TokenLiteral() != fmt.Sprintf("%s", value) {
		t.Errorf("st.TokenLiteral not %s. got=%s",
			value, st.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s, got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t, got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testNil(t *testing.T, n ast.Expression, _ interface{}) bool {
	if n != nil {
		t.Errorf("n not nil. got=%T", n)
		return false
	}

	return true
}

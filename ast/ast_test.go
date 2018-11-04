package ast

import (
	"testing"

	"github.com/nrtkbb/go-MEL/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&StringStatement{
				Token: token.Token{Type: token.StringDec, Literal: "string"},
				Names: []*Identifier{{
					Token: token.Token{Type: token.Ident, Literal: "$myVar"},
					Value: "$myVar",
				}},
				Values: []Expression{&Identifier{
					Token: token.Token{Type: token.Ident, Literal: "$anotherVar"},
					Value: "$anotherVar",
				}},
			},
		},
	}

	if program.String() != "string $myVar = $anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

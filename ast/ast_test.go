package ast

import (
	"github.com/nrtkbb/go-MEL/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&StringStatement{
				Token: token.Token{Type: token.String, Literal: "string"},
				Name: &Identifier{
					Token: token.Token{Type: token.Ident, Literal: "$myVar"},
					Value: "$myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.Ident, Literal: "$anotherVar"},
					Value: "$anotherVar",
				},
			},
		},
	}

	if program.String() != "string $myVar = $anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

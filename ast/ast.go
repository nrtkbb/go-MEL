package ast

import (
	"bytes"
	"strings"

	"github.com/nrtkbb/go-MEL/token"
)

// Node is top of AST interface
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement have some expression
type Statement interface {
	Node
	statementNode()
}

// Expression ...
type Expression interface {
	Node
	expressionNode()
}

// Literal ...
type Literal interface {
	Expression
	literalNode()
}

// Program is represent the entire program
type Program struct {
	Statements []Statement
}

// TokenLiteral ...
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String ...
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ExpressionStatement ...
type ExpressionStatement struct {
	Token      token.Token // first token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral ...
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String ...
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// InfixExpression ...
type InfixExpression struct {
	Token    token.Token // infix token. ex) '+' or '-' ...
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral ...
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String ...
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// PrefixExpression ...
type PrefixExpression struct {
	Token    token.Token // prefix token. ex) '-' or '!'
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral ...
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String ...
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// PostfixExpression ...
type PostfixExpression struct {
	Token    token.Token // postfix token. ex) '--' or '++'
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode() {}

// TokenLiteral ...
func (pe *PostfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String ...
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")

	return out.String()
}

// TernaryExpression ...
type TernaryExpression struct {
	Conditional Expression
	Token1      token.Token // '?'
	Operator1   string      // "?"
	TrueExp     Expression
	Token2      token.Token // ':'
	Operator2   string      // ":"
	FalseExp    Expression
}

func (te *TernaryExpression) expressionNode() {}

// TokenLiteral ...
func (te *TernaryExpression) TokenLiteral() string {
	return te.Token1.Literal
}

// String ...
func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(te.Conditional.String())
	out.WriteString(" " + te.Operator1 + " ")
	out.WriteString(te.TrueExp.String())
	out.WriteString(" " + te.Operator2 + " ")
	out.WriteString(te.FalseExp.String())
	out.WriteString(")")

	return out.String()
}

// TypeDeclaration ...
type TypeDeclaration struct {
	Token   token.Token // string or int or float or matrix or vector
	IsArray bool
}

func (td *TypeDeclaration) expressionNode() {}

// TokenLiteral ...
func (td *TypeDeclaration) TokenLiteral() string {
	if td.IsArray {
		return td.Token.Literal + "[]"
	}
	return td.Token.Literal
}

// String ...
func (td *TypeDeclaration) String() string {
	if td.IsArray {
		return td.Token.Literal + "[]"
	}
	return td.Token.Literal
}

// CallExpression ...
type CallExpression struct {
	Token     token.Token // '(' token or '`' token or function
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral ...
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String ...
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// FunctionLiteral ...
type FunctionLiteral struct {
	Token      token.Token // proc
	Name       token.Token // ProcIdent
	IsGlobal   bool
	ReturnType *TypeDeclaration
	ParamTypes []*TypeDeclaration
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral ...
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String ...
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	if fl.IsGlobal {
		out.WriteString("global ")
	}

	out.WriteString(fl.TokenLiteral() + " ")

	if fl.ReturnType != nil {
		out.WriteString(fl.ReturnType.String() + " ")
	}

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// ForExpression ...
type ForExpression struct {
	Token       token.Token // for
	InitNames   []Expression
	InitValues  []Expression
	Condition   Expression
	ChangeOf    Expression
	Consequence *BlockStatement
}

func (fe *ForExpression) expressionNode() {}

// TokenLiteral ...
func (fe *ForExpression) TokenLiteral() string {
	return fe.Token.Literal
}

// String ...
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fe.TokenLiteral())
	out.WriteString(" (")

	var outNames []string
	for i, name := range fe.InitNames {
		if fe.InitValues[i] != nil {
			outNames = append(outNames, name.String()+" = "+fe.InitValues[i].String())
		} else {
			outNames = append(outNames, name.String())
		}
	}
	out.WriteString(strings.Join(outNames, ", "))

	out.WriteString("; ")
	out.WriteString(fe.Condition.String())
	out.WriteString("; ")
	out.WriteString(fe.ChangeOf.String())
	out.WriteString(") ")
	out.WriteString(fe.Consequence.String())

	return out.String()
}

// ForInExpression ...
type ForInExpression struct {
	Token        token.Token // for
	Element      *Identifier
	ArrayElement Expression
	Consequence  *BlockStatement
}

func (fie *ForInExpression) expressionNode() {}

// TokenLiteral ...
func (fie *ForInExpression) TokenLiteral() string {
	return fie.Token.Literal
}

// String ...
func (fie *ForInExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fie.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(fie.Element.String())
	out.WriteString(" in ")
	out.WriteString(fie.ArrayElement.String())
	out.WriteString(") ")
	out.WriteString(fie.Consequence.String())

	return out.String()
}

// DoWhileExpression ...
type DoWhileExpression struct {
	Token       token.Token // do
	Condition   Expression
	Consequence *BlockStatement
}

func (dwe *DoWhileExpression) expressionNode() {}

// TokenLiteral ...
func (dwe *DoWhileExpression) TokenLiteral() string {
	return dwe.Token.Literal
}

// String ...
func (dwe *DoWhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString(dwe.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(dwe.Consequence.String())
	out.WriteString(" while ")
	out.WriteString(dwe.Condition.String())
	out.WriteString(";")

	return out.String()
}

// WhileExpression ...
type WhileExpression struct {
	Token       token.Token // while
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode() {}

// TokenLiteral ...
func (we *WhileExpression) TokenLiteral() string {
	return we.Token.Literal
}

// String ...
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString(we.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Consequence.String())

	return out.String()
}

// IfExpression ...
type IfExpression struct {
	Token       token.Token // if
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral ...
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String ...
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// BlockStatement ...
type BlockStatement struct {
	Token      token.Token // '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral ...
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String ...
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// SwitchExpression ...
type SwitchExpression struct {
	Token          token.Token // switch
	Condition      Expression
	Cases          []Literal
	CaseStatements []*CaseStatement
}

func (se *SwitchExpression) expressionNode() {}

// TokenLiteral ...
func (se *SwitchExpression) TokenLiteral() string {
	return se.Token.Literal
}

// String ...
func (se *SwitchExpression) String() string {
	var out bytes.Buffer

	out.WriteString(se.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(se.Condition.String())
	out.WriteString(" {")
	for i, cas := range se.Cases {
		out.WriteString("case ")
		out.WriteString(cas.String())
		out.WriteString(se.CaseStatements[i].String())
	}
	out.WriteString(" }")

	return out.String()
}

// CaseStatement ...
type CaseStatement struct {
	Token      token.Token // ':' token
	Statements []Statement
}

func (cs *CaseStatement) statementNode() {}

// TokenLiteral ...
func (cs *CaseStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// String ...
func (cs *CaseStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.TokenLiteral())
	for _, stmt := range cs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

// VectorStatement ...
type VectorStatement struct {
	Token  token.Token // token.VectorDec
	Names  []Expression
	Values []Expression
}

func (vs *VectorStatement) statementNode() {}

// TokenLiteral ...
func (vs *VectorStatement) TokenLiteral() string {
	return vs.Token.Literal
}

// String ...
func (vs *VectorStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral())
	out.WriteString(" ")

	var outNames []string
	for i, name := range vs.Names {
		if vs.Values[i] != nil {
			outNames = append(outNames, name.String()+" = "+vs.Values[i].String())
		} else {
			outNames = append(outNames, name.String())
		}
	}
	out.WriteString(strings.Join(outNames, ", "))

	out.WriteString(";")

	return out.String()
}

// MatrixStatement ...
type MatrixStatement struct {
	Token  token.Token // token.MatrixDec
	Names  []Expression
	Values []Expression
}

func (ms *MatrixStatement) statementNode() {}

// TokenLiteral ...
func (ms *MatrixStatement) TokenLiteral() string {
	return ms.Token.Literal
}

// String ...
func (ms *MatrixStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ms.TokenLiteral())
	out.WriteString(" ")

	var outNames []string
	for i, name := range ms.Names {
		if ms.Values[i] != nil {
			outNames = append(outNames, name.String()+" = "+ms.Values[i].String())
		} else {
			outNames = append(outNames, name.String())
		}
	}
	out.WriteString(strings.Join(outNames, ", "))

	out.WriteString(";")

	return out.String()
}

// IntegerStatement ...
type IntegerStatement struct {
	Token  token.Token // token.IntDec
	Names  []Expression
	Values []Expression
}

func (is *IntegerStatement) statementNode() {}

// TokenLiteral ...
func (is *IntegerStatement) TokenLiteral() string {
	return is.Token.Literal
}

// String ...
func (is *IntegerStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.TokenLiteral())
	out.WriteString(" ")

	var outNames []string
	for i, name := range is.Names {
		if is.Values[i] != nil {
			outNames = append(outNames, name.String()+" = "+is.Values[i].String())
		} else {
			outNames = append(outNames, name.String())
		}
	}
	out.WriteString(strings.Join(outNames, ", "))

	out.WriteString(";")

	return out.String()
}

// FloatStatement ...
type FloatStatement struct {
	Token  token.Token // token.FloatDec
	Names  []Expression
	Values []Expression
}

func (fs *FloatStatement) statementNode() {}

// TokenLiteral ...
func (fs *FloatStatement) TokenLiteral() string {
	return fs.Token.Literal
}

// String ...
func (fs *FloatStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")

	var outNames []string
	for i, name := range fs.Names {
		if fs.Values[i] != nil {
			outNames = append(outNames, name.String()+" = "+fs.Values[i].String())
		} else {
			outNames = append(outNames, name.String())
		}
	}
	out.WriteString(strings.Join(outNames, ", "))

	out.WriteString(";")

	return out.String()
}

// StringStatement ...
type StringStatement struct {
	Token  token.Token // token.StringDec
	Names  []Expression
	Values []Expression
}

func (ss *StringStatement) statementNode() {}

// TokenLiteral ...
func (ss *StringStatement) TokenLiteral() string {
	return ss.Token.Literal
}

// String ...
func (ss *StringStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ss.TokenLiteral())
	out.WriteString(" ")

	var outNames []string
	for i, name := range ss.Names {
		if ss.Values[i] != nil {
			outNames = append(outNames, name.String()+" = "+ss.Values[i].String())
		} else {
			outNames = append(outNames, name.String())
		}
	}
	out.WriteString(strings.Join(outNames, ", "))

	out.WriteString(";")

	return out.String()
}

// ArrayLiteral ...
type ArrayLiteral struct {
	Token    token.Token // '{' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) literalNode()    {}

// TokenLiteral ...
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// String ...
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

// IndexExpression ...
type IndexExpression struct {
	Token token.Token // token.Lbracket
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral ...
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String ...
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// Identifier is token.Ident
type Identifier struct {
	Token token.Token // token.Ident
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral ...
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String ...
func (i *Identifier) String() string {
	return i.Value
}

// BreakStatement ...
type BreakStatement struct {
	Token token.Token // token.Break
}

func (bs *BreakStatement) statementNode() {}

// TokenLiteral ...
func (bs *BreakStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String ...
func (bs *BreakStatement) String() string {
	return bs.Token.Literal + ";"
}

// ContinueStatement ...
type ContinueStatement struct {
	Token token.Token // token.Break
}

func (bs *ContinueStatement) statementNode() {}

// TokenLiteral ...
func (bs *ContinueStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String ...
func (bs *ContinueStatement) String() string {
	return bs.Token.Literal + ";"
}

// ReturnStatement ...
type ReturnStatement struct {
	Token       token.Token // token.Return
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral ...
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String ...
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// IntegerLiteral ...
type IntegerLiteral struct {
	Token token.Token // token.Int
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) literalNode()    {}

// TokenLiteral ...
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String ...
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// FloatLiteral ...
type FloatLiteral struct {
	Token token.Token // token.Float
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) literalNode()    {}

// TokenLiteral ...
func (fl *FloatLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String ...
func (fl *FloatLiteral) String() string {
	return fl.Token.Literal
}

// StringLiteral ...
type StringLiteral struct {
	Token token.Token // token.String
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) literalNode()    {}

// TokenLiteral ...
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String ...
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// BooleanLiteral ...
type BooleanLiteral struct {
	Token token.Token // true or false
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) literalNode()    {}

// TokenLiteral ...
func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

// String ...
func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

// TensorLiteral ...
type TensorLiteral struct {
	Token  token.Token // token.Ltensor
	Values [][]Expression
}

func (vl *TensorLiteral) expressionNode() {}
func (vl *TensorLiteral) literalNode()    {}

// TokenLiteral ...
func (vl *TensorLiteral) TokenLiteral() string {
	return vl.Token.Literal
}

// String ...
func (vl *TensorLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("<<")

	var outValues []string
	for _, v := range vl.Values {
		var inValues []string
		for _, vv := range v {
			inValues = append(inValues, vv.String())
		}
		outValues = append(outValues, strings.Join(inValues, ", "))
		inValues = nil
	}
	out.WriteString(strings.Join(outValues, ";\n  "))

	out.WriteString(">>")

	return out.String()
}

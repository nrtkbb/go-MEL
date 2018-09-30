package token

// Type ...
type Type string

// Token ...
type Token struct {
	Type    Type
	Literal string
	Row     int // 行数 1行はじまり
	Column  int // 列数 1列はじまり
}

// Type strings.
const (
	Illegal = "Illegal"
	EOF     = "EOF"

	// 識別子 + リテラル
	Ident      = "Ident"      // $add, $foobar, $x, $y, ...
	ProcIdent  = "ProcIdent"  // add, FuncName, ...
	IntData    = "IntData"    // 1343456
	Int16Data  = "Int16Data"  // 0xA0, 0xfff, ...
	FloatData  = "FloatData"  // 1.1, 1e-3, 1e+3, ...
	StringData = "StringData" // "node.attr", ...
	Flag       = "Flag"       // -size, -s, ...

	// 演算子
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Slash    = "/"
	Asterisk = "*"
	Bang     = "!"

	Lt = "<"
	Gt = ">"

	Eq    = "=="
	NotEq = "!="

	// デリミタ
	Comma      = ","
	Semicolon  = ";"
	BackQuotes = "`"

	Lparen   = "("
	Rparen   = ")"
	Lbrace   = "{"
	Rbrace   = "}"
	Lbracket = "["
	Rbracket = "]"
	Ltensor  = "<<"
	Rtensor  = ">>"

	// For LookupIdent
	Global = "Global"
	Proc   = "Proc"
	String = "String"
	Int    = "Int"
	Float  = "Float"
	Vector = "Vector"
	Matrix = "Matrix"
	True   = "True"
	False  = "False"
	If     = "If"
	Else   = "Else"
	Return = "Return"
)

var keywords = map[string]Type{
	"global": Global,
	"proc":   Proc,
	"string": String,
	"int":    Int,
	"float":  Float,
	"vector": Vector,
	"matrix": Matrix,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

// LookupIdent ...
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ProcIdent
}

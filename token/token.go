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
	Ident     = "Ident"     // $add, $foobar, $x, $y, ...
	ProcIdent = "ProcIdent" // add, FuncName, ...
	Int       = "Int"       // 1343456
	Int16     = "Int16"     // 0xA0, 0xfff, ...
	Float     = "Float"     // 1.1, 1e-3, 1e+3, ...
	String    = "String"    // "node.attr", ...
	Flag      = "Flag"      // -size, -s, ...
	True      = "True"
	On        = "On"
	False     = "False"
	Off       = "Off"

	Comment = "Comment" // // ...

	// 演算子
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Slash    = "/"
	Asterisk = "*"
	Bang     = "!"

	Lt = "<"
	Gt = ">"

	Question = "?"
	Coron    = ":"

	Eq        = "=="
	NotEq     = "!="
	And       = "&&"
	Or        = "||"
	Increment = "++"
	Decrement = "--"

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
	Global    = "Global"
	Proc      = "Proc"
	StringDec = "StringDec"
	IntDec    = "IntDec"
	FloatDec  = "FloatDec"
	VectorDec = "VectorDec"
	MatrixDec = "MatrixDec"
	BoolDec   = "BoolDec"
	If        = "If"
	While     = "While"
	Else      = "Else"
	Return    = "Return"
)

var keywords = map[string]Type{
	"global": Global,
	"proc":   Proc,
	"string": StringDec,
	"int":    IntDec,
	"float":  FloatDec,
	"vector": VectorDec,
	"matrix": MatrixDec,
	"true":   True,
	"false":  False,
	"on":     On,
	"off":    Off,
	"bool":   BoolDec,
	"if":     If,
	"while":  While,
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

package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 + リテラル
	IDENT    = "IDENT"    // $add, $foobar, $x, $y, ...
	INT_DATA = "INT_DATA" // 1343456
	PROC_IDENT = "PROC_IDENT" // add, FuncName, ...

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	BANG     = "!"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// キーワード
	GLOBAL = "GLOBAL"
	PROC   = "PROC"
	STRING = "STRING"
	INT    = "INT"
	FLOAT  = "FLOAT"
	VECTOR = "VECTOR"
	MATRIX = "MATRIX"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
)

var keywords = map[string]TokenType{
	"global": GLOBAL,
	"proc":   PROC,
	"string": STRING,
	"int":    INT,
	"float":  FLOAT,
	"vector": VECTOR,
	"matrix": MATRIX,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return PROC_IDENT
}

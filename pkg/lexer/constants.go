package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	GT     TokenType = ">"
	LT     TokenType = "<"
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	COMMA  TokenType = ","
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	QWORD  TokenType = "QWORD"
	STR_KW TokenType = "STR_KW"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	WHILE  TokenType = "WHILE"
	IMPORT TokenType = "IMPORT"
)

var keywords = map[string]TokenType{
	"qword":  QWORD,
	"string": STR_KW,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"import": IMPORT,
}

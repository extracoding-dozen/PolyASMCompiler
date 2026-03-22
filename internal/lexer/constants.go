package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Идентификаторы и литералы
	IDENT  TokenType = "IDENT"  // my_var, target_dir
	NUMBER TokenType = "NUMBER" // 15, 0x400
	STRING TokenType = "STRING" // "hello world"

	// Операторы
	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	GT     TokenType = ">"
	LT     TokenType = "<"
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	// Разделители
	COMMA  TokenType = ","
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Ключевые слова
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

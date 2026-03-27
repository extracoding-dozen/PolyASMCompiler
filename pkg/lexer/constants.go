// Package lexer реализует лексический анализатор для полиморфного компилятора.
package lexer

// TokenType представляет строковый тип для классификации различных лексем языка.
type TokenType string

// Список констант для всех поддерживаемых типов токенов.
const (
	ILLEGAL TokenType = "ILLEGAL" // Неизвестный символ/токен
	EOF     TokenType = "EOF"     // Конец файла

	// Литералы
	IDENT  TokenType = "IDENT"  // Идентификаторы (названия переменных, функций)
	NUMBER TokenType = "NUMBER" // Числовые литералы (десятичные и шестнадцатеричные)
	STRING TokenType = "STRING" // Строковые литералы

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
	QWORD  TokenType = "QWORD"  // Тип данных 64-бит
	STR_KW TokenType = "STR_KW" // Тип данных string
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	WHILE  TokenType = "WHILE"
	IMPORT TokenType = "IMPORT"
)

// keywords сопоставляет зарезервированные слова языка с их типами TokenType.
var keywords = map[string]TokenType{
	"qword":  QWORD,
	"string": STR_KW,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"import": IMPORT,
}

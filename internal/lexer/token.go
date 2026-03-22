package lexer

type Token struct {
	Type    TokenType
	Literal string // Само текстовое значение (например "15" или "my_var")
	Line    int
	Column  int
}

package lexer

type Lexer interface {
	NextToken() Token
}

package lexer

// Lexer определяет интерфейс для лексического сканера,
// который последовательно возвращает токены из исходного кода.
type Lexer interface {
	// NextToken анализирует входные данные и возвращает следующий найденный токен.
	NextToken() Token
	SetInput(input string)
}

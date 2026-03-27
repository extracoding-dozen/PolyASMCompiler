package parser

import "go.mod/pkg/ast"

// Parser определяет интерфейс для взаимодействия с синтаксическим анализатором.
type Parser interface {
	// Errors возвращает накопленный список ошибок парсинга.
	Errors() []string
	// ParseProgram запускает процесс синтаксического анализа кода.
	ParseProgram() *ast.Program
}

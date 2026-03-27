// Package parser реализует синтаксический анализатор (парсер), преобразующий поток токенов
// в абстрактное синтаксическое дерево (AST) методом рекурсивного спуска (Pratt Parser).
package parser

import (
	"go.mod/pkg/ast"
	"go.mod/pkg/lexer"
)

// Приоритеты операторов.
const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	EQUALS      // ==
	LESSGREATER // > или <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X или !X
	CALL        // myFunction(X)
)

// precedences сопоставляет типы токенов с их приоритетами.
var precedences = map[lexer.TokenType]int{
	lexer.ASSIGN: ASSIGN,
	lexer.EQ:     EQUALS,
	lexer.NOT_EQ: EQUALS,
	lexer.LT:     LESSGREATER,
	lexer.GT:     LESSGREATER,
	lexer.PLUS:   SUM,
	lexer.MINUS:  SUM,
	lexer.LPAREN: CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

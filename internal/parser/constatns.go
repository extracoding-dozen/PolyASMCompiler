package parser

import (
	"go.mod/internal/ast"
	"go.mod/internal/lexer"
)

const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	EQUALS      // == or !=
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // execute(...)
)

// Таблица соответствия токена и его приоритета
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

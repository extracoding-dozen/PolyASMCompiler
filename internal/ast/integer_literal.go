package ast

import "go.mod/internal/lexer"

type IntegerLiteral struct {
	Token lexer.Token
	Value int64 // Сохраним пока как int64, для 0x400 это будет распарсено
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

package ast

import (
	"strings"

	"go.mod/internal/lexer"
)

type CallExpression struct {
	Token     lexer.Token // Токен '('
	Function  Expression  // Identifier
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}

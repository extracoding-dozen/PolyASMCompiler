package ast

import (
	"bytes"

	"go.mod/internal/lexer"
)

type IfStatement struct {
	Token       lexer.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("if " + is.Condition.String() + " " + is.Consequence.String())
	if is.Alternative != nil {
		out.WriteString(" else " + is.Alternative.String())
	}
	return out.String() + "\n"
}

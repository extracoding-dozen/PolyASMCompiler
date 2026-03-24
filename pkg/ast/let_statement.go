package ast

import "go.mod/pkg/lexer"

type LetStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	return ls.TokenLiteral() + " " + ls.Name.String() + " = " + ls.Value.String() + "\n"
}

package ast

import "go.mod/internal/lexer"

type ImportStatement struct {
	Token lexer.Token
	Path  *StringLiteral
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	return "import " + is.Path.String() + "\n"
}

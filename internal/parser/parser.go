package parser

import "go.mod/internal/ast"

type Parser interface {
	Errors() []string
	ParseProgram() *ast.Program
}

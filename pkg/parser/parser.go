package parser

import "go.mod/pkg/ast"

type Parser interface {
	Errors() []string
	ParseProgram() *ast.Program
}

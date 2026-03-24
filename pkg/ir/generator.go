package ir

import (
	"fmt"

	"go.mod/pkg/ast"
	"go.mod/pkg/ir/ir_constants_and_types"
)

type Generator struct {
	Instructions []Instruction
	vregCounter  int
	labelCounter int
	env          map[string]VReg
	ctx          *IRContext
}

func NewGenerator(ctr *IRContext) *Generator {
	return &Generator{
		Instructions: []Instruction{},
		vregCounter:  0,
		labelCounter: 0,
		env:          make(map[string]VReg),
		ctx:          ctr,
	}
}

func (g *Generator) emit(inst Instruction) {
	g.Instructions = append(g.Instructions, inst)
}

func (g *Generator) Generate(program *ast.Program) {
	for _, stmt := range program.Statements {
		g.generateStatement(stmt)
	}
}

func (g *Generator) generateStatement(node ast.Statement) {
	switch n := node.(type) {

	case *ast.LetStatement:
		valReg := g.generateExpression(n.Value)
		varReg := g.ctx.NextVReg()
		g.env[n.Name.Value] = varReg
		g.emit(Instruction{Op: ir_constants_and_types.MOV, Dst: varReg, Src1: valReg})

	case *ast.IfStatement:
		condReg := g.generateExpression(n.Condition)

		lblFalse := g.ctx.NextLabel()
		lblEnd := g.ctx.NextLabel()

		g.emit(Instruction{Op: ir_constants_and_types.CMP, Dst: condReg, Src1: Imm{Value: 0}})
		g.emit(Instruction{Op: ir_constants_and_types.JE, Dst: lblFalse})

		for _, stmt := range n.Consequence.Statements {
			g.generateStatement(stmt)
		}
		g.emit(Instruction{Op: ir_constants_and_types.JMP, Dst: lblEnd})

		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblFalse})
		if n.Alternative != nil {
			for _, stmt := range n.Alternative.Statements {
				g.generateStatement(stmt)
			}
		}

		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblEnd})

	case *ast.WhileStatement:
		lblStart := g.ctx.NextLabel()
		lblEnd := g.ctx.NextLabel()

		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblStart})

		condReg := g.generateExpression(n.Condition)
		g.emit(Instruction{Op: ir_constants_and_types.CMP, Dst: condReg, Src1: Imm{Value: 0}})
		g.emit(Instruction{Op: ir_constants_and_types.JE, Dst: lblEnd})

		for _, stmt := range n.Body.Statements {
			g.generateStatement(stmt)
		}
		g.emit(Instruction{Op: ir_constants_and_types.JMP, Dst: lblStart})

		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblEnd})

	case *ast.ExpressionStatement:
		g.generateExpression(n.Expression)
	}
}

func (g *Generator) generateExpression(node ast.Expression) VReg {
	switch n := node.(type) {

	case *ast.IntegerLiteral:
		reg := g.ctx.NextVReg()
		g.emit(Instruction{Op: ir_constants_and_types.MOV, Dst: reg, Src1: Imm{Value: n.Value}})
		return reg

	case *ast.StringLiteral:
		reg := g.ctx.NextVReg()
		g.emit(Instruction{Op: ir_constants_and_types.LOAD_STR, Dst: reg, Src1: Str{Value: n.Value}})
		return reg

	case *ast.Identifier:
		if reg, ok := g.env[n.Value]; ok {
			return reg
		}
		panic(fmt.Sprintf("Неизвестная переменная: %s", n.Value))

	case *ast.InfixExpression:
		leftReg := g.generateExpression(n.Left)
		rightReg := g.generateExpression(n.Right)
		resultReg := g.ctx.NextVReg()

		g.emit(Instruction{Op: ir_constants_and_types.MOV, Dst: resultReg, Src1: leftReg})

		switch n.Operator {
		case "+":
			g.emit(Instruction{Op: ir_constants_and_types.ADD, Dst: resultReg, Src1: rightReg})
		case "-":
			g.emit(Instruction{Op: ir_constants_and_types.SUB, Dst: resultReg, Src1: rightReg})
		case "==":

			g.emit(Instruction{Op: ir_constants_and_types.CMP, Dst: resultReg, Src1: rightReg})

		}
		return resultReg

	case *ast.CallExpression:
		funcName := n.Function.String()
		resultReg := g.ctx.NextVReg()

		var argRegs []Value
		for _, arg := range n.Arguments {
			argReg := g.generateExpression(arg)
			argRegs = append(argRegs, argReg)
		}

		switch funcName {

		case "exit":
			g.emit(Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_EXIT}, Args: argRegs})

		case "fork":
			g.emit(Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_FORK}, Args: argRegs})

		case "execute":
			g.emit(Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_EXECVE}, Args: argRegs})

		case "chmod":
			g.emit(Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_CHMOD}, Args: argRegs})
		case "delete":
			g.emit(Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_UNLINK}, Args: argRegs})

		case "copy":

			g.emit(Instruction{Op: ir_constants_and_types.MACRO_COPY, Dst: resultReg, Args: argRegs})
		case "useradd":
			g.emit(Instruction{Op: ir_constants_and_types.MACRO_USERADD, Dst: resultReg, Args: argRegs})
		case "write":
			g.emit(Instruction{Op: ir_constants_and_types.MACRO_WRITE, Dst: resultReg, Args: argRegs})
		case "get_file_size":
			g.emit(Instruction{Op: ir_constants_and_types.MACRO_GET_FILE_SIZE, Dst: resultReg, Args: argRegs})
		case "rename":
			g.emit(Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_RENAME}, Args: argRegs})
		case "sleep":
			g.emit(Instruction{Op: ir_constants_and_types.MACRO_SLEEP, Dst: resultReg, Args: argRegs})
		default:

			panic(fmt.Sprintf("Неизвестная встроенная функция: %s", funcName))
		}

		return resultReg
	}

	return VReg{ID: 0}
}

// Print - выводит сгенерированный код на экран
func (g *Generator) Print() {
	fmt.Println("=== Сгенерированный IR код ===")
	for _, inst := range g.Instructions {
		fmt.Println(inst.String())
	}
	fmt.Println("==============================")
}

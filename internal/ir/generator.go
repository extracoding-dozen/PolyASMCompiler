package ir

import (
	"fmt"

	"go.mod/internal/ast"
	"go.mod/internal/ir/ir_constants_and_types"
)

type Generator struct {
	Instructions []Instruction
	vregCounter  int
	labelCounter int
	env          map[string]VReg // Хранит привязку переменных к регистрам: x -> v1
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

// Главный метод входа
func (g *Generator) Generate(program *ast.Program) {
	for _, stmt := range program.Statements {
		g.generateStatement(stmt)
	}
}

// Обработка Statements (Утверждений)
func (g *Generator) generateStatement(node ast.Statement) {
	switch n := node.(type) {

	// let x = 5
	case *ast.LetStatement:
		valReg := g.generateExpression(n.Value)
		varReg := g.ctx.NextVReg()
		g.env[n.Name.Value] = varReg // Запоминаем, что переменная x лежит в varReg
		g.emit(Instruction{Op: ir_constants_and_types.MOV, Dst: varReg, Src1: valReg})

	// if (x > 5) { ... } else { ... }
	case *ast.IfStatement:
		condReg := g.generateExpression(n.Condition)

		lblFalse := g.ctx.NextLabel()
		lblEnd := g.ctx.NextLabel()

		// Если условие ложно (0), прыгаем в блок False
		g.emit(Instruction{Op: ir_constants_and_types.CMP, Dst: condReg, Src1: Imm{Value: 0}})
		g.emit(Instruction{Op: ir_constants_and_types.JE, Dst: lblFalse})

		// Блок True
		for _, stmt := range n.Consequence.Statements {
			g.generateStatement(stmt)
		}
		g.emit(Instruction{Op: ir_constants_and_types.JMP, Dst: lblEnd}) // Пропускаем Else

		// Блок False (Else)
		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblFalse})
		if n.Alternative != nil {
			for _, stmt := range n.Alternative.Statements {
				g.generateStatement(stmt)
			}
		}

		// Конец
		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblEnd})

	// while (x < 10) { ... }
	case *ast.WhileStatement:
		lblStart := g.ctx.NextLabel()
		lblEnd := g.ctx.NextLabel()

		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblStart})

		condReg := g.generateExpression(n.Condition)
		g.emit(Instruction{Op: ir_constants_and_types.CMP, Dst: condReg, Src1: Imm{Value: 0}})
		g.emit(Instruction{Op: ir_constants_and_types.JE, Dst: lblEnd}) // Если ложь - выходим из цикла

		for _, stmt := range n.Body.Statements {
			g.generateStatement(stmt)
		}
		g.emit(Instruction{Op: ir_constants_and_types.JMP, Dst: lblStart}) // Прыжок в начало цикла

		g.emit(Instruction{Op: ir_constants_and_types.LABEL, Dst: lblEnd})

	case *ast.ExpressionStatement:
		g.generateExpression(n.Expression)
	}
}

// Обработка Expressions (Выражений) - Возвращает регистр, в котором лежит результат
func (g *Generator) generateExpression(node ast.Expression) VReg {
	switch n := node.(type) {

	// Число: 5
	case *ast.IntegerLiteral:
		reg := g.ctx.NextVReg()
		g.emit(Instruction{Op: ir_constants_and_types.MOV, Dst: reg, Src1: Imm{Value: n.Value}})
		return reg

	// Строка: "/bin/sh"
	case *ast.StringLiteral:
		reg := g.ctx.NextVReg()
		g.emit(Instruction{Op: ir_constants_and_types.LOAD_STR, Dst: reg, Src1: Str{Value: n.Value}})
		return reg

	// Использование переменной: x
	case *ast.Identifier:
		if reg, ok := g.env[n.Value]; ok {
			return reg
		}
		panic(fmt.Sprintf("Неизвестная переменная: %s", n.Value)) // В реальном коде лучше возвращать ошибку

	// Математика и логика: x + y, x == y
	case *ast.InfixExpression:
		leftReg := g.generateExpression(n.Left)
		rightReg := g.generateExpression(n.Right)
		resultReg := g.ctx.NextVReg()

		// В x86 обычно результат пишется в левый операнд, поэтому сначала копируем:
		g.emit(Instruction{Op: ir_constants_and_types.MOV, Dst: resultReg, Src1: leftReg})

		switch n.Operator {
		case "+":
			g.emit(Instruction{Op: ir_constants_and_types.ADD, Dst: resultReg, Src1: rightReg})
		case "-":
			g.emit(Instruction{Op: ir_constants_and_types.SUB, Dst: resultReg, Src1: rightReg})
		case "==":
			// Для сравнения вернем 1 если равно, и 0 если нет. (Слишком сложная логика для IR, пока делаем хак - пишем результат CMP)
			g.emit(Instruction{Op: ir_constants_and_types.CMP, Dst: resultReg, Src1: rightReg})
			// В будущем здесь будет логика SETE / SETNE
		}
		return resultReg

	// Вызов функции: execute(path, flags)
	case *ast.CallExpression:
		funcName := n.Function.String()
		resultReg := g.ctx.NextVReg()

		// Генерируем значения для всех переданных аргументов
		var argRegs []Value
		for _, arg := range n.Arguments {
			argReg := g.generateExpression(arg)
			argRegs = append(argRegs, argReg)
		}

		// --- СТАНДАРТНАЯ БИБЛИОТЕКА POLYASM ---
		switch funcName {

		// 1. Простые системные вызовы (сразу подставляем номера Syscall для x86_64)
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
		// 2. Сложные макросы (Отправляем в Распутыватель)
		case "copy":
			// Обрати внимание: Op теперь MACRO_COPY, а не SYSCALL!
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
			// Если мы допишем импорт стороннего ASM (например my_func()),
			// то здесь будет вызов CALL
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

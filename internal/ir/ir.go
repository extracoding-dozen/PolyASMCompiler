package ir

import (
	"fmt"
	"strings"

	"go.mod/internal/ir/ir_constants_and_types"
)

// Opcode - код операции

// Value - интерфейс для операндов (регистр, число, строка или метка)
type Value interface {
	String() string
}

// VReg - Виртуальный регистр (v0, v1, v2)
type VReg struct {
	ID int
}

func (v VReg) String() string {
	return fmt.Sprintf("v%d", v.ID)
}

// Imm - Числовой литерал (10, 0x400)
type Imm struct {
	Value int64
}

func (i Imm) String() string {
	return fmt.Sprintf("%d", i.Value)
}

// Str - Строковый литерал
type Str struct {
	Value string
}

func (s Str) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

// Lbl - Метка для переходов (.L1, .L2)
type Lbl struct {
	Name string
}

func (l Lbl) String() string {
	return l.Name
}

// Instruction - одна команда в нашем плоском IR-коде
type Instruction struct {
	Op   ir_constants_and_types.Opcode
	Dst  Value   // Куда записываем результат (обычно VReg)
	Src1 Value   // Откуда берем данные (VReg или Imm)
	Src2 Value   // Второй аргумент (для ADD, SUB, CMP)
	Args []Value // Дополнительные аргументы (нужно для SYSCALL)
}

func (inst Instruction) String() string {
	switch inst.Op {
	case ir_constants_and_types.LABEL:
		return fmt.Sprintf("%s:", inst.Dst)
	case ir_constants_and_types.SYSCALL:
		args := []string{}
		for _, a := range inst.Args {
			args = append(args, a.String())
		}
		return fmt.Sprintf("  %s %s, [%s] -> %s", inst.Op, inst.Src1, strings.Join(args, ", "), inst.Dst)
	case ir_constants_and_types.JMP, ir_constants_and_types.JE, ir_constants_and_types.JNE:
		return fmt.Sprintf("  %s %s", inst.Op, inst.Dst)
	case ir_constants_and_types.LOAD_STR:
		return fmt.Sprintf("  %s %s, %s", inst.Op, inst.Dst, inst.Src1)
	default:
		if inst.Src2 != nil {
			return fmt.Sprintf("  %s %s, %s, %s", inst.Op, inst.Dst, inst.Src1, inst.Src2)
		}
		if inst.Src1 != nil {
			return fmt.Sprintf("  %s %s, %s", inst.Op, inst.Dst, inst.Src1)
		}
		return fmt.Sprintf("  %s %s", inst.Op, inst.Dst)
	}
}

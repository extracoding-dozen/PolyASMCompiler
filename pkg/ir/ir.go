package ir

import (
	"fmt"
	"strings"

	"go.mod/pkg/ir/ir_constants_and_types"
)

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
	Bytes []byte
}

func (s Str) String() string {
	if len(s.Bytes) > 0 {
		return fmt.Sprintf("<enc_bytes_len:%d>", len(s.Bytes))
	}
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
	Dst  Value
	Src1 Value
	Src2 Value
	Args []Value
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

		syscallName := inst.Src1.String()
		if immVal, ok := inst.Src1.(Imm); ok {

			syscallName = ir_constants_and_types.SyscallNames[immVal.Value]
		}

		return fmt.Sprintf("  SYSCALL %s, [%s] -> %s", syscallName, strings.Join(args, ", "), inst.Dst)
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

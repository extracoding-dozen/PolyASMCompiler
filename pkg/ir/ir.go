// Package ir определяет базовые структуры и типы данных для работы с промежуточным представлением кода.
package ir

import (
	"fmt"
	"strings"

	"go.mod/pkg/ir/ir_constants_and_types"
)

// Value представляет собой общий интерфейс для всех типов операндов в IR-коде.
type Value interface {
	String() string
}

// VReg представляет виртуальный регистр (например, v0, v1).
type VReg struct {
	ID int
}

// String возвращает строковое представление виртуального регистра.
func (v VReg) String() string {
	return fmt.Sprintf("v%d", v.ID)
}

// Imm представляет непосредственное числовое значение (целочисленный литерал).
type Imm struct {
	Value int64
}

// String возвращает строковое представление числового литерала.
func (i Imm) String() string {
	return fmt.Sprintf("%d", i.Value)
}

// Str представляет строковый литерал, который может содержать как исходный текст, так и зашифрованные байты.
type Str struct {
	Value string
	Bytes []byte
}

// String возвращает описание строки или информацию о длине зашифрованных байт.
func (s Str) String() string {
	if len(s.Bytes) > 0 {
		return fmt.Sprintf("<enc_bytes_len:%d>", len(s.Bytes))
	}
	return fmt.Sprintf(`"%s"`, s.Value)
}

// Lbl представляет метку в коде для выполнения переходов.
type Lbl struct {
	Name string
}

// String возвращает имя метки.
func (l Lbl) String() string {
	return l.Name
}

// Instruction описывает отдельную команду промежуточного представления.
type Instruction struct {
	Op   ir_constants_and_types.Opcode
	Dst  Value
	Src1 Value
	Src2 Value
	Args []Value
}

// String генерирует человекочитаемое текстовое представление IR-инструкции.
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

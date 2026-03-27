// Package ir предоставляет инструменты для работы с промежуточным представлением программы и генерации машинного кода.
package ir

import (
	"fmt"
	"math/rand"
	"strings"

	"go.mod/pkg/ir/ir_constants_and_types"
)

var syscallArgsMap = []string{ir_constants_and_types.RDI, ir_constants_and_types.RSI, ir_constants_and_types.RDX, ir_constants_and_types.R10, ir_constants_and_types.R8, ir_constants_and_types.R9}

// RegisterAllocator преобразует список промежуточных инструкций (IR) в ассемблерный код x86_64.
// Использует стековую модель выделения памяти под виртуальные регистры.
type RegisterAllocator struct {
	asmCode strings.Builder
}

// NewRegisterAllocator создает новый экземпляр аллокатора регистров.
func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{}
}

func (a *RegisterAllocator) emitAsm(format string, args ...interface{}) {
	a.asmCode.WriteString(fmt.Sprintf(format+"\n", args...))
}

func (a *RegisterAllocator) vRegToStack(v VReg) string {
	offset := v.ID * 8
	return fmt.Sprintf("qword [rbp - %d]", offset)
}

func getRandomJunk() string {
	junks := []string{
		"nop",
		"xchg rax, rax",
		"add rsp, 0",
		"sub rsp, 0",
	}
	return junks[rand.Intn(len(junks))]
}

func (a *RegisterAllocator) loadValue(val Value, physReg string) {
	switch v := val.(type) {
	case Imm:
		a.emitAsm(" mov %s, %d", physReg, v.Value)
	case VReg:
		a.emitAsm(" mov %s, %s", physReg, a.vRegToStack(v))
	case Str:
		a.emitAsm(" lea %s, [rel %s]", physReg, v.Value)
	}
}

// Allocate транслирует массив IR-инструкций в итоговую строку ассемблерного кода x86_64.
// Формирует пролог и эпилог функции start, вычисляет размер стека и обрабатывает макросы.
func (a *RegisterAllocator) Allocate(instructions []Instruction, totalVRegs int) string {
	a.asmCode.Reset()
	stackSize := (totalVRegs + 1) * 8
	if stackSize%16 != 0 {
		stackSize += 8
	}
	a.emitAsm(".intel_syntax noprefix")
	a.emitAsm(".global _start")
	a.emitAsm("_start:")
	a.emitAsm(" push rbp")
	a.emitAsm(" mov rbp, rsp")
	a.emitAsm(" sub rsp, %d", stackSize)
	for _, inst := range instructions {
		switch inst.Op {
		case ir_constants_and_types.LABEL:
			a.emitAsm("%s:", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.MOV:
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.loadValue(inst.Src1, ir_constants_and_types.R12)
			a.emitAsm(" mov %s, %s", dstStack, ir_constants_and_types.R12)
		case ir_constants_and_types.ADD, ir_constants_and_types.SUB:
			a.loadValue(inst.Src1, ir_constants_and_types.R12)
			a.loadValue(inst.Src2, ir_constants_and_types.R13)
			if inst.Op == ir_constants_and_types.ADD {
				a.emitAsm(" add %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)
			} else {
				a.emitAsm(" sub %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)
			}
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.emitAsm(" mov %s, %s", dstStack, ir_constants_and_types.R12)
		case ir_constants_and_types.CMP:
			a.loadValue(inst.Dst, ir_constants_and_types.R12)
			a.loadValue(inst.Src1, ir_constants_and_types.R13)
			a.emitAsm(" cmp %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)
		case ir_constants_and_types.JMP:
			a.emitAsm(" jmp %s", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.JE:
			a.emitAsm(" je %s", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.JNE:
			a.emitAsm(" jne %s", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.JL:
			a.emitAsm(" jl %s", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.LEA:
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			srcStack := a.vRegToStack(inst.Src1.(VReg))
			addrStr := strings.TrimPrefix(srcStack, "qword ")
			a.emitAsm(" lea r12, %s", addrStr)
			a.emitAsm(" mov %s, r12", dstStack)
		case ir_constants_and_types.SYSCALL:
			a.loadValue(inst.Src1, ir_constants_and_types.RAX)
			for i, arg := range inst.Args {
				if i >= len(syscallArgsMap) {
					panic("Слишком много аргументов для syscall (>6)")
				}
				targetReg := syscallArgsMap[i]
				a.loadValue(arg, targetReg)
			}
			a.emitAsm(" syscall")
			if inst.Dst != nil {
				dstStack := a.vRegToStack(inst.Dst.(VReg))
				a.emitAsm(" mov %s, %s", dstStack, ir_constants_and_types.RAX)
			}
		case ir_constants_and_types.LOAD_STR:
			strObj := inst.Src1.(Str)
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			vID := inst.Dst.(VReg).ID
			var rawBytes []byte
			if len(strObj.Bytes) > 0 {
				rawBytes = strObj.Bytes
			} else {
				rawBytes = []byte(strObj.Value)
			}
			a.emitAsm(" jmp .L_str%d_skip", vID)
			a.emitAsm(".L_str%d_data:", vID)
			hexStrs := []string{}
			for _, b := range rawBytes {
				hexStrs = append(hexStrs, fmt.Sprintf("0x%02x", b))
			}
			hexStrs = append(hexStrs, "0x00")
			a.emitAsm(" .byte %s", strings.Join(hexStrs, ", "))
			a.emitAsm(".L_str%d_skip:", vID)
			a.emitAsm(" lea r12, [rip + .L_str%d_data]", vID)
			a.emitAsm(" mov %s, r12", dstStack)
		case ir_constants_and_types.MACRO_STRLEN:
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.loadValue(inst.Src1, ir_constants_and_types.RDI)
			a.emitAsm(" sub rcx, rcx")
			a.emitAsm(" not rcx")
			a.emitAsm(" sub rax, rax")
			a.emitAsm(" cld")
			a.emitAsm(" repne scasb")
			a.emitAsm(" not rcx")
			a.emitAsm(" dec rcx")
			a.emitAsm(" mov %s, rcx", dstStack)
		case ir_constants_and_types.MACRO_DECRYPT_STR:
			ptrStack := a.vRegToStack(inst.Dst.(VReg))
			xorKey := inst.Src1.(Imm).Value
			strLen := inst.Src2.(Imm).Value
			allocLen := (strLen + 15) &^ 15
			tmpSrc := "r14"
			a.emitAsm(" mov %s, %s", tmpSrc, ptrStack)
			a.emitAsm(" sub rsp, %d", allocLen)
			a.emitAsm(" mov %s, rsp", ptrStack)
			srcRegs := []string{"rsi", "r8", "r10"}
			dstRegs := []string{"rdi", "r9", "r11"}
			cntRegs := []string{"rcx", "rbx", "rdx"}
			dataRegs := []string{"rax", "r12", "r13", "r15"}
			srcReg := srcRegs[rand.Intn(len(srcRegs))]
			dstReg := dstRegs[rand.Intn(len(dstRegs))]
			cntReg := cntRegs[rand.Intn(len(cntRegs))]
			dataReg64 := dataRegs[rand.Intn(len(dataRegs))]
			dataReg8 := func(reg string) string {
				res, exists := ir_constants_and_types.Map8BitRegs[reg]
				if !exists {
					return "al"
				}
				return res
			}(dataReg64)
			a.emitAsm(" mov %s, %s", srcReg, tmpSrc)
			a.emitAsm(" mov %s, rsp", dstReg)
			a.emitAsm(" mov %s, %d", cntReg, strLen)
			loopLbl := fmt.Sprintf(".L_decrypt%d", inst.Dst.(VReg).ID)
			a.emitAsm("%s:", loopLbl)
			if rand.Intn(100) < 30 {
				a.emitAsm(" %s", getRandomJunk())
			}
			a.emitAsm(" mov %s, byte ptr [%s]", dataReg8, srcReg)
			if rand.Intn(2) == 0 {
				a.emitAsm(" xor %s, %d", dataReg8, xorKey)
			} else {
				a.emitAsm(" mov al, %d", xorKey)
				a.emitAsm(" xor %s, al", dataReg8)
			}
			a.emitAsm(" mov byte ptr [%s], %s", dstReg, dataReg8)
			if rand.Intn(2) == 0 {
				a.emitAsm(" inc %s", srcReg)
				a.emitAsm(" inc %s", dstReg)
			} else {
				a.emitAsm(" add %s, 1", srcReg)
				a.emitAsm(" add %s, 1", dstReg)
			}
			if rand.Intn(100) < 30 {
				a.emitAsm(" %s", getRandomJunk())
			}
			if rand.Intn(2) == 0 {
				a.emitAsm(" dec %s", cntReg)
				a.emitAsm(" test %s, %s", cntReg, cntReg)
				a.emitAsm(" jnz %s", loopLbl)
			} else {
				a.emitAsm(" sub %s, 1", cntReg)
				a.emitAsm(" cmp %s, 0", cntReg)
				a.emitAsm(" jg %s", loopLbl)
			}
			a.emitAsm(" mov byte ptr [%s], 0", dstReg)
		}
	}
	a.emitAsm(" mov rsp, rbp")
	a.emitAsm(" pop rbp")
	a.emitAsm(" ret")
	return a.asmCode.String()
}

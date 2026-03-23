package ir

import (
	"fmt"
	"strings"

	"go.mod/internal/ir/ir_constants_and_types"
)

// Конвенция вызовов для Syscall (Linux x86_64)
var syscallArgsMap = []string{ir_constants_and_types.RDI, ir_constants_and_types.RSI, ir_constants_and_types.RDX, ir_constants_and_types.R10, ir_constants_and_types.R8, ir_constants_and_types.R9}

type RegisterAllocator struct {
	asmCode strings.Builder
}

func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{}
}

// Вспомогательная функция для добавления строки asm
func (a *RegisterAllocator) emitAsm(format string, args ...interface{}) {
	a.asmCode.WriteString(fmt.Sprintf(format+"\n", args...))
}

// vRegToStack конвертирует виртуальный регистр (v5) в адрес на стеке ([rbp - 40])
func (a *RegisterAllocator) vRegToStack(v VReg) string {
	// Каждый qword занимает 8 байт
	offset := v.ID * 8
	return fmt.Sprintf("qword [rbp - %d]", offset)
}

// loadValue загружает Value (Число или vReg) в физический регистр (например R12)
func (a *RegisterAllocator) loadValue(val Value, physReg string) {
	switch v := val.(type) {
	case Imm:
		a.emitAsm("  mov %s, %d", physReg, v.Value)
	case VReg:
		a.emitAsm("  mov %s, %s", physReg, a.vRegToStack(v))
	case Str:
		// Для строк генерируем хитрый трюк (RIP-relative адресация или push на стек)
		// Пока для упрощения предполагаем, что строка уже лежит где-то, или используем трюк JMP-CALL-POP
		// Это будет доработано в модуле генерации строк. Оставим заглушку:
		a.emitAsm("  lea %s, [rel %s]", physReg, v.Value)
	}
}

// Главный метод: трансляция IR -> ASM
func (a *RegisterAllocator) Allocate(instructions []Instruction, totalVRegs int) string {
	a.asmCode.Reset()

	// 1. Пролог функции (Установка фрейма стека)
	stackSize := (totalVRegs + 1) * 8
	// Выравниваем по 16 байт (требование ABI x86_64)
	if stackSize%16 != 0 {
		stackSize += 8
	}
	a.emitAsm(".intel_syntax noprefix")
	a.emitAsm(".global _start")
	//a.emitAsm("section .text")
	a.emitAsm("_start:")
	a.emitAsm("  push rbp")
	a.emitAsm("  mov rbp, rsp")
	a.emitAsm("  sub rsp, %d", stackSize)

	// 2. Трансляция инструкций
	for _, inst := range instructions {
		switch inst.Op {

		case ir_constants_and_types.LABEL:
			a.emitAsm("%s:", inst.Dst.(Lbl).Name)

		case ir_constants_and_types.MOV:
			// mov dst, src
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.loadValue(inst.Src1, ir_constants_and_types.R12)
			a.emitAsm("  mov %s, %s", dstStack, ir_constants_and_types.R12)

		case ir_constants_and_types.ADD, ir_constants_and_types.SUB:
			// add dst, src1, src2 -> dst = src1 + src2
			a.loadValue(inst.Src1, ir_constants_and_types.R12) // Левый операнд
			a.loadValue(inst.Src2, ir_constants_and_types.R13) // Правый операнд
			if inst.Op == ir_constants_and_types.ADD {
				a.emitAsm("  add %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)
			} else {
				a.emitAsm("  sub %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)
			}
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.emitAsm("  mov %s, %s", dstStack, ir_constants_and_types.R12) // Сохраняем результат

		case ir_constants_and_types.CMP:
			// cmp src1, src2
			a.loadValue(inst.Dst, ir_constants_and_types.R12) // Dst здесь выступает как левый операнд сравнения (хак IR)
			a.loadValue(inst.Src1, ir_constants_and_types.R13)
			a.emitAsm("  cmp %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)

		case ir_constants_and_types.JMP:
			a.emitAsm("  jmp %s", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.JE:
			a.emitAsm("  je %s", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.JNE:
			a.emitAsm("  jne %s", inst.Dst.(Lbl).Name)
		case ir_constants_and_types.JL:
			a.emitAsm("  jl %s", inst.Dst.(Lbl).Name)

		case ir_constants_and_types.SYSCALL:
			// Самая сложная часть: подготовка регистров для ядра Linux
			// 1. Номер сисколла идет в RAX
			a.loadValue(inst.Src1, ir_constants_and_types.RAX)

			// 2. Аргументы идут в RDI, RSI, RDX, R10, R8, R9
			for i, arg := range inst.Args {
				if i >= len(syscallArgsMap) {
					panic("Слишком много аргументов для syscall (>6)")
				}
				targetReg := syscallArgsMap[i]
				a.loadValue(arg, targetReg)
			}

			// 3. Вызываем ядро
			a.emitAsm("  syscall")

			// 4. Результат возвращается в RAX. Сохраняем его в виртуальный регистр (если он указан)
			if inst.Dst != nil {
				dstStack := a.vRegToStack(inst.Dst.(VReg))
				a.emitAsm("  mov %s, %s", dstStack, ir_constants_and_types.RAX)
			}
			// ОБРАБОТКА СТРОК (Shellcode Inline Data Trick)
		case ir_constants_and_types.LOAD_STR:
			strVal := inst.Src1.(Str).Value
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			vID := inst.Dst.(VReg).ID

			// 1. Прыгаем через байты строки, чтобы процессор их не выполнил
			a.emitAsm("  jmp .L_str_%d_skip", vID)

			// 2. Кладем саму строку (с нуль-терминатором)
			a.emitAsm(".L_str_%d_data:", vID)
			a.emitAsm("  .asciz \"%s\"", strVal)

			// 3. Метка после строки
			a.emitAsm(".L_str_%d_skip:", vID)

			// 4. Загружаем адрес этой строки в R12 через RIP-relative (PIE)
			a.emitAsm("  lea r12, [rip + .L_str_%d_data]", vID)

			// 5. Сохраняем указатель на строку в наш виртуальный регистр на стеке
			a.emitAsm("  mov %s, r12", dstStack)
		case ir_constants_and_types.MACRO_STRLEN:
			// Алгоритм поиска длины строки в x86_64:
			// rdi = указатель на строку
			// al = 0 (ищем нуль-терминатор)
			// rcx = -1 (максимальный счетчик)
			// repne scasb (сканируем строку до нуля)
			// результат: длина = (-1) - rcx - 1

			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.loadValue(inst.Src1, ir_constants_and_types.RDI) // Грузим указатель на строку в RDI

			a.emitAsm("  sub rcx, rcx") // rcx = 0
			a.emitAsm("  not rcx")      // rcx = -1 (0xFFFFFFFFFFFFFFFF)
			a.emitAsm("  sub rax, rax") // rax = 0 (ищем \x00)
			a.emitAsm("  cld")          // Направление поиска - вперед
			a.emitAsm("  repne scasb")  // Искать AL в [RDI], уменьшая RCX

			a.emitAsm("  not rcx") // Инвертируем обратно
			a.emitAsm("  dec rcx") // Вычитаем 1 (сам нулевой байт)

			// Сохраняем длину в виртуальный регистр
			a.emitAsm("  mov %s, rcx", dstStack)
		}
	}
	// 3. Эпилог (обычно до него не доходит из-за sys_exit, но для чистоты)
	a.emitAsm("  mov rsp, rbp")
	a.emitAsm("  pop rbp")
	a.emitAsm("  ret")

	return a.asmCode.String()
}

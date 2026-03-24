package ir

import (
	"fmt"
	"math/rand"
	"strings"

	"go.mod/pkg/ir/ir_constants_and_types"
)

var syscallArgsMap = []string{ir_constants_and_types.RDI, ir_constants_and_types.RSI, ir_constants_and_types.RDX, ir_constants_and_types.R10, ir_constants_and_types.R8, ir_constants_and_types.R9}

type RegisterAllocator struct {
	asmCode strings.Builder
}

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
		a.emitAsm("  mov %s, %d", physReg, v.Value)
	case VReg:
		a.emitAsm("  mov %s, %s", physReg, a.vRegToStack(v))
	case Str:
		a.emitAsm("  lea %s, [rel %s]", physReg, v.Value)
	}
}

func (a *RegisterAllocator) Allocate(instructions []Instruction, totalVRegs int) string {
	a.asmCode.Reset()

	stackSize := (totalVRegs + 1) * 8

	if stackSize%16 != 0 {
		stackSize += 8
	}
	a.emitAsm(".intel_syntax noprefix")
	a.emitAsm(".global _start")

	a.emitAsm("_start:")
	a.emitAsm("  push rbp")
	a.emitAsm("  mov rbp, rsp")
	a.emitAsm("  sub rsp, %d", stackSize)

	for _, inst := range instructions {
		switch inst.Op {

		case ir_constants_and_types.LABEL:
			a.emitAsm("%s:", inst.Dst.(Lbl).Name)

		case ir_constants_and_types.MOV:

			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.loadValue(inst.Src1, ir_constants_and_types.R12)
			a.emitAsm("  mov %s, %s", dstStack, ir_constants_and_types.R12)

		case ir_constants_and_types.ADD, ir_constants_and_types.SUB:

			a.loadValue(inst.Src1, ir_constants_and_types.R12)
			a.loadValue(inst.Src2, ir_constants_and_types.R13)
			if inst.Op == ir_constants_and_types.ADD {
				a.emitAsm("  add %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)
			} else {
				a.emitAsm("  sub %s, %s", ir_constants_and_types.R12, ir_constants_and_types.R13)
			}
			dstStack := a.vRegToStack(inst.Dst.(VReg))
			a.emitAsm("  mov %s, %s", dstStack, ir_constants_and_types.R12)

		case ir_constants_and_types.CMP:

			a.loadValue(inst.Dst, ir_constants_and_types.R12)
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
		case ir_constants_and_types.LEA:

			dstStack := a.vRegToStack(inst.Dst.(VReg))
			srcStack := a.vRegToStack(inst.Src1.(VReg))

			addrStr := strings.TrimPrefix(srcStack, "qword ")

			a.emitAsm("  lea r12, %s", addrStr)

			a.emitAsm("  mov %s, r12", dstStack)

		case ir_constants_and_types.SYSCALL:

			a.loadValue(inst.Src1, ir_constants_and_types.RAX)

			for i, arg := range inst.Args {
				if i >= len(syscallArgsMap) {
					panic("Слишком много аргументов для syscall (>6)")
				}
				targetReg := syscallArgsMap[i]
				a.loadValue(arg, targetReg)
			}

			a.emitAsm("  syscall")

			if inst.Dst != nil {
				dstStack := a.vRegToStack(inst.Dst.(VReg))
				a.emitAsm("  mov %s, %s", dstStack, ir_constants_and_types.RAX)
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

			a.emitAsm("  jmp .L_str_%d_skip", vID)

			a.emitAsm(".L_str_%d_data:", vID)

			hexStrs := []string{}
			for _, b := range rawBytes {
				hexStrs = append(hexStrs, fmt.Sprintf("0x%02x", b))
			}

			// Проще добавить его всегда, он не помешает).
			hexStrs = append(hexStrs, "0x00")

			a.emitAsm("  .byte %s", strings.Join(hexStrs, ", "))

			// 4. Метка конца
			a.emitAsm(".L_str_%d_skip:", vID)

			// 5. Загружаем адрес данных (RIP-relative)
			a.emitAsm("  lea r12, [rip + .L_str_%d_data]", vID)

			// 6. Сохраняем указатель
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
		case ir_constants_and_types.MACRO_DECRYPT_STR:
			ptrStack := a.vRegToStack(inst.Dst.(VReg)) // qword [rbp - X]
			xorKey := inst.Src1.(Imm).Value
			strLen := inst.Src2.(Imm).Value

			a.emitAsm("  // --- POLYMORPHIC RUNTIME DECRYPT (STACK COPY) ---")

			// 1. ВЫДЕЛЕНИЕ ПАМЯТИ НА АППАРАТНОМ СТЕКЕ
			// Выравниваем длину по 16 байт, чтобы не сломать ABI (иначе syscall/SSE упадут)
			allocLen := (strLen + 15) &^ 15

			// Для копирования нам нужен временный регистр (чтобы достать исходный указатель на .text)
			// Берем R14, так как он не используется в syscall-конвенциях
			tmpSrc := "r14"
			a.emitAsm("  mov %s, %s", tmpSrc, ptrStack)

			// Выделяем память на стеке
			a.emitAsm("  sub rsp, %d", allocLen)

			// ОБНОВЛЯЕМ виртуальный регистр! Теперь он указывает на стек (RW-память), а не на .text (R-O)
			a.emitAsm("  mov %s, rsp", ptrStack)

			// --- 2. ПОЛИМОРФИЗМ: Выбираем случайные регистры ---
			// Важно: мы не берем регистры, которые могут быть заняты (RSP, RBP, R14)
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
			}(dataReg64) // Получаем младший байт (например, r12 -> r12b)

			// Инициализируем выбранные регистры
			a.emitAsm("  mov %s, %s", srcReg, tmpSrc) // Указатель на зашифрованные данные (.text)
			a.emitAsm("  mov %s, rsp", dstReg)        // Указатель на выделенную память (Стек)
			a.emitAsm("  mov %s, %d", cntReg, strLen) // Счетчик длины строки

			// --- 3. ПОЛИМОРФНЫЙ ЦИКЛ КОПИРОВАНИЯ И РАСШИФРОВКИ ---
			loopLbl := fmt.Sprintf(".L_decrypt_%d", inst.Dst.(VReg).ID)
			a.emitAsm("%s:", loopLbl)

			// Шанс 30% вставить безопасный мусор в начало итерации
			if rand.Intn(100) < 30 {
				a.emitAsm("  %s", getRandomJunk())
			}

			// ЧТЕНИЕ БАЙТА ИЗ .TEXT
			a.emitAsm("  mov %s, byte ptr [%s]", dataReg8, srcReg)

			// РАСШИФРОВКА (XOR)
			// Вариант 1: Прямой XOR (xor r12b, 0x42)
			// Вариант 2: XOR через временный регистр AL
			if rand.Intn(2) == 0 {
				a.emitAsm("  xor %s, %d", dataReg8, xorKey)
			} else {
				a.emitAsm("  mov al, %d", xorKey)
				a.emitAsm("  xor %s, al", dataReg8)
			}

			// ЗАПИСЬ БАЙТА НА СТЕК (Тут больше нет SIGSEGV!)
			a.emitAsm("  mov byte ptr [%s], %s", dstReg, dataReg8)

			// СДВИГ УКАЗАТЕЛЕЙ (inc reg ИЛИ add reg, 1)
			if rand.Intn(2) == 0 {
				a.emitAsm("  inc %s", srcReg)
				a.emitAsm("  inc %s", dstReg)
			} else {
				a.emitAsm("  add %s, 1", srcReg)
				a.emitAsm("  add %s, 1", dstReg)
			}

			// Шанс 30% вставить мусор перед концом итерации
			if rand.Intn(100) < 30 {
				a.emitAsm("  %s", getRandomJunk())
			}

			// ПРОВЕРКА УСЛОВИЯ (Конец цикла)
			// Вариант 1: dec rcx; jnz. Вариант 2: sub rcx, 1; cmp rcx, 0; jg
			if rand.Intn(2) == 0 {
				a.emitAsm("  dec %s", cntReg)
				a.emitAsm("  test %s, %s", cntReg, cntReg)
				a.emitAsm("  jnz %s", loopLbl)
			} else {
				a.emitAsm("  sub %s, 1", cntReg)
				a.emitAsm("  cmp %s, 0", cntReg)
				a.emitAsm("  jg %s", loopLbl)
			}

			// 4. ГАРАНТИРУЕМ НУЛЬ-ТЕРМИНАТОР (\x00)
			// Так как dstReg после цикла указывает на байт СРАЗУ ПОСЛЕ строки,
			// мы просто пишем туда ноль. Это делает строку валидной для C-функций ядра (open, execve).
			a.emitAsm("  mov byte ptr [%s], 0", dstReg)
			a.emitAsm("  // -------------------------------------")
		}

	}
	// 3. Эпилог (обычно до него не доходит из-за sys_exit, но для чистоты)
	a.emitAsm("  mov rsp, rbp")
	a.emitAsm("  pop rbp")
	a.emitAsm("  ret")

	return a.asmCode.String()
}

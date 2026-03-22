package ir

// Новые опкоды для макросов (High-Level IR)
const (
	MACRO_COPY    Opcode = "MACRO_COPY"
	MACRO_USERADD Opcode = "MACRO_USERADD"
)

type MacroExpander struct {
	ctx *IRContext
}

func NewMacroExpander(ctx *IRContext) *MacroExpander {
	return &MacroExpander{ctx: ctx}
}

// Главная функция. Принимает High-Level IR, возвращает Low-Level IR.
func (e *MacroExpander) Expand(input []Instruction) []Instruction {
	var output []Instruction

	for _, inst := range input {
		switch inst.Op {

		case MACRO_COPY:
			// Распутываем copy(src, dst)
			expandedBlock := e.expandCopy(inst.Args[0], inst.Args[1])
			// Вставляем результат распутывания вместо оригинальной инструкции
			output = append(output, expandedBlock...)

		//case MACRO_USERADD:
		//	expandedBlock := e.expandUserAdd(inst.Args[0], inst.Args[1])
		//	output = append(output, expandedBlock...)

		default:
			// Обычные инструкции (MOV, ADD, SYSCALL 59) пробрасываем без изменений
			output = append(output, inst)
		}
	}

	return output
}

func (e *MacroExpander) expandCopy(srcReg Value, dstReg Value) []Instruction {
	var block []Instruction

	// Запрашиваем новые виртуальные регистры для файловых дескрипторов
	fdIn := e.ctx.NextVReg()
	fdOut := e.ctx.NextVReg()

	// 1. fd_in = sys_open(src, 0 /* O_RDONLY */)
	block = append(block, Instruction{
		Op:   SYSCALL,
		Dst:  fdIn,
		Src1: Imm{Value: 2}, // номер sys_open
		Args: []Value{srcReg, Imm{Value: 0}, Imm{Value: 0}},
	})

	// 2. fd_out = sys_open(dst, 65 /* O_WRONLY|O_CREAT */, 0777)
	block = append(block, Instruction{
		Op:   SYSCALL,
		Dst:  fdOut,
		Src1: Imm{Value: 2},
		Args: []Value{dstReg, Imm{Value: 65}, Imm{Value: 511}}, // 511 = 0777 в восьмеричной
	})

	// 3. Вызываем sys_sendfile (номер 40 в x86_64).
	// Он копирует данные напрямую в ядре, не гоняя их в UserSpace!
	// sys_sendfile(out_fd, in_fd, offset, count)
	// count сделаем огромным числом, например 0x7fffffff, чтобы скопировать за 1 раз.
	resultReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   SYSCALL,
		Dst:  resultReg,
		Src1: Imm{Value: 40}, // номер sys_sendfile
		Args: []Value{fdOut, fdIn, Imm{Value: 0}, Imm{Value: 2147483647}},
	})

	// 4. Закрываем дескрипторы
	block = append(block, Instruction{Op: SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: 3}, Args: []Value{fdIn}})
	block = append(block, Instruction{Op: SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: 3}, Args: []Value{fdOut}})

	return block
}

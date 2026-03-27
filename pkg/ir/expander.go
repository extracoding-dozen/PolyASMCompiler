// Package ir предоставляет инструменты для трансформации и расширения промежуточного представления.
package ir

import "go.mod/pkg/ir/ir_constants_and_types"

// MacroExpander выполняет раскрытие высокоуровневых макрокоманд в последовательности базовых IR-инструкций.
type MacroExpander struct {
	ctx *IRContext
}

// NewMacroExpander создает новый экземпляр расширителя макросов с заданным контекстом.
func NewMacroExpander(ctx *IRContext) *MacroExpander {
	return &MacroExpander{ctx: ctx}
}

// GetCtx возвращает текущий контекст IR.
func (e *MacroExpander) GetCtx() *IRContext {
	return e.ctx
}

// Expand обходит список инструкций и заменяет макросы на соответствующие блоки низкоуровневых команд.
func (e *MacroExpander) Expand(input []Instruction) []Instruction {
	var output []Instruction
	for _, inst := range input {
		switch inst.Op {
		case ir_constants_and_types.MACRO_COPY:
			output = append(output, e.expandCopy(inst.Dst, inst.Args[0], inst.Args[1])...)
		case ir_constants_and_types.MACRO_USERADD:
			output = append(output, e.expandUserAdd(inst.Dst, inst.Args[0], inst.Args[1])...)
		case ir_constants_and_types.MACRO_WRITE:
			output = append(output, e.expandWrite(inst.Dst, inst.Args[0], inst.Args[1], inst.Args[2])...)
		case ir_constants_and_types.MACRO_GET_FILE_SIZE:
			output = append(output, e.expandGetFileSize(inst.Dst, inst.Args[0])...)
		case ir_constants_and_types.MACRO_SLEEP:
			output = append(output, e.expandSleep(inst.Dst, inst.Args[0])...)
		default:
			output = append(output, inst)
		}
	}
	return output
}

// expandCopy реализует копирование файла через системный вызов sendfile.
func (e *MacroExpander) expandCopy(resReg, srcReg Value, dstReg Value) []Instruction {
	var block []Instruction
	fdIn := e.ctx.NextVReg()
	fdOut := e.ctx.NextVReg()
	errLabel := e.ctx.NextLabel()
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: resReg, Src1: Imm{Value: ir_constants_and_types.RETURN_ERROR}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: fdIn, Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, Args: []Value{srcReg, Imm{Value: 0}, Imm{Value: 0}}})
	block = append(block, Instruction{Op: ir_constants_and_types.CMP, Dst: fdIn, Src1: Imm{Value: 0}})
	block = append(block, Instruction{Op: ir_constants_and_types.JL, Dst: errLabel})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: fdOut, Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, Args: []Value{dstReg, Imm{Value: 65}, Imm{Value: 511}}})
	block = append(block, Instruction{Op: ir_constants_and_types.CMP, Dst: fdOut, Src1: Imm{Value: 0}})
	block = append(block, Instruction{Op: ir_constants_and_types.JL, Dst: errLabel})
	resultReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_SEND_FILE}, Args: []Value{fdOut, fdIn, Imm{Value: 0}, Imm{Value: 2147483647}}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE}, Args: []Value{fdIn}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE}, Args: []Value{fdOut}})
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.RETURN_SUCCESS}})
	block = append(block, Instruction{Op: ir_constants_and_types.LABEL, Dst: errLabel})
	return block
}

// expandWrite реализует запись данных в файл по заданному смещению.
func (e *MacroExpander) expandWrite(resReg, pathReg, offsetReg, dataReg Value) []Instruction {
	var block []Instruction
	fdOut := e.ctx.NextVReg()
	errLabel := e.ctx.NextLabel()
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: resReg, Src1: Imm{Value: ir_constants_and_types.RETURN_ERROR}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: fdOut, Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, Args: []Value{pathReg, Imm{Value: 65}, Imm{Value: 420}}})
	block = append(block, Instruction{Op: ir_constants_and_types.CMP, Dst: fdOut, Src1: Imm{Value: 0}})
	block = append(block, Instruction{Op: ir_constants_and_types.JL, Dst: errLabel})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_LSEEK}, Args: []Value{fdOut, offsetReg, Imm{Value: 0}}})
	dataSizeReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.MACRO_STRLEN, Dst: dataSizeReg, Src1: dataReg})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_WRITE}, Args: []Value{fdOut, dataReg, dataSizeReg}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE}, Args: []Value{fdOut}})
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: resReg, Src1: Imm{Value: ir_constants_and_types.RETURN_SUCCESS}})
	block = append(block, Instruction{Op: ir_constants_and_types.LABEL, Dst: errLabel})
	return block
}

// expandUserAdd реализует добавление пользователя в файл /etc/passwd.
func (e *MacroExpander) expandUserAdd(resReg, userReg, passReg Value) []Instruction {
	var block []Instruction
	passwdPathReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.LOAD_STR, Dst: passwdPathReg, Src1: Str{Value: "/etc/passwd"}})
	fdOut := e.ctx.NextVReg()
	errorLabel := e.ctx.NextLabel()
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: resReg, Src1: Imm{Value: ir_constants_and_types.RETURN_ERROR}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: fdOut, Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, Args: []Value{passwdPathReg, Imm{Value: 1025}, Imm{Value: 420}}})
	block = append(block, Instruction{Op: ir_constants_and_types.CMP, Dst: fdOut, Src1: Imm{Value: 0}})
	block = append(block, Instruction{Op: ir_constants_and_types.JL, Dst: errorLabel})
	userSizeReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.MACRO_STRLEN, Dst: userSizeReg, Src1: userReg})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_WRITE}, Args: []Value{fdOut, userReg, userSizeReg}})
	tailReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.LOAD_STR, Dst: tailReg, Src1: Str{Value: ":x:0:0::/root:/bin/bash\n"}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_WRITE}, Args: []Value{fdOut, tailReg, Imm{Value: 24}}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE}, Args: []Value{fdOut}})
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: resReg, Src1: Imm{Value: ir_constants_and_types.RETURN_SUCCESS}})
	block = append(block, Instruction{Op: ir_constants_and_types.LABEL, Dst: errorLabel})
	return block
}

// expandGetFileSize определяет размер файла с помощью системного вызова lseek.
func (e *MacroExpander) expandGetFileSize(resultReg, filepathReg Value) []Instruction {
	var block []Instruction
	fdReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: resultReg, Src1: Imm{Value: -1}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: fdReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, Args: []Value{filepathReg, Imm{Value: 0}, Imm{Value: 0}}})
	block = append(block, Instruction{Op: ir_constants_and_types.CMP, Dst: fdReg, Src1: Imm{Value: 0}})
	errLaber := e.ctx.NextLabel()
	block = append(block, Instruction{Op: ir_constants_and_types.JL, Dst: errLaber})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_LSEEK}, Args: []Value{fdReg, Imm{Value: 0}, Imm{Value: 2}}})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: e.ctx.NextVReg(), Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE}, Args: []Value{fdReg}})
	block = append(block, Instruction{Op: ir_constants_and_types.LABEL, Dst: errLaber})
	return block
}

// expandSleep реализует приостановку выполнения программы через системный вызов nanosleep.
func (e *MacroExpander) expandSleep(resultReg, sleepTime Value) []Instruction {
	var block []Instruction
	tvSecReg := e.ctx.NextVReg()
	tvNsecReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: tvSecReg, Src1: Imm{Value: 0}})
	block = append(block, Instruction{Op: ir_constants_and_types.MOV, Dst: tvNsecReg, Src1: sleepTime})
	structPtrReg := e.ctx.NextVReg()
	block = append(block, Instruction{Op: ir_constants_and_types.LEA, Dst: structPtrReg, Src1: tvSecReg})
	block = append(block, Instruction{Op: ir_constants_and_types.SYSCALL, Dst: resultReg, Src1: Imm{Value: ir_constants_and_types.SYSCALL_NANOSLEEP}, Args: []Value{structPtrReg, Imm{Value: 0}}})
	return block
}

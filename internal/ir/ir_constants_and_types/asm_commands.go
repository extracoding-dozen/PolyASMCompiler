package ir_constants_and_types

const (
	MOV      Opcode = "MOV"
	ADD      Opcode = "ADD"
	SUB      Opcode = "SUB"
	CMP      Opcode = "CMP"
	JMP      Opcode = "JMP"
	JNE      Opcode = "JNE" // Jump if Not Equal
	JE       Opcode = "JE"  // Jump if Equal
	LABEL    Opcode = "LABEL"
	SYSCALL  Opcode = "SYSCALL"
	LOAD_STR Opcode = "LOAD_STR" // Специальная инструкция для загрузки строк на стек
	JL       Opcode = "JL"
)

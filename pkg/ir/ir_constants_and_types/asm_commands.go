// Package ir_constants_and_types содержит определения номеров системных вызовов,
// архитектурных регистров и типов данных для x86_64.
package ir_constants_and_types

// Команды языка ams .intel noprefix
const (
	MOV      Opcode = "MOV"
	ADD      Opcode = "ADD"
	SUB      Opcode = "SUB"
	CMP      Opcode = "CMP"
	JMP      Opcode = "JMP"
	JNE      Opcode = "JNE"
	JE       Opcode = "JE"
	LABEL    Opcode = "LABEL"
	SYSCALL  Opcode = "SYSCALL"
	LOAD_STR Opcode = "LOAD_STR"
	JL       Opcode = "JL"
	LEA      Opcode = "LEA"
)

// Package ir_constants_and_types содержит определения номеров системных вызовов,
// архитектурных регистров и типов данных для x86_64.
package ir_constants_and_types

// Макросы команд языка, которые не соответствуют вызову одного (syscall)
const (
	MACRO_COPY          Opcode = "MACRO_COPY"
	MACRO_USERADD       Opcode = "MACRO_USERADD"
	MACRO_WRITE         Opcode = "MACRO_WRITE"
	MACRO_STRLEN        Opcode = "MACRO_STRLEN"
	MACRO_GET_FILE_SIZE Opcode = "MACRO_GET_FILE_SIZE"
	MACRO_SLEEP         Opcode = "MACRO_SLEEP"
	MACRO_DECRYPT_STR   Opcode = "MACRO_DECRYPT_STR"
)

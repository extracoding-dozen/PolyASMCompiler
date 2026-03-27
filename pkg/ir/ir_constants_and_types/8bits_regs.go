// Package ir_constants_and_types содержит определения номеров системных вызовов,
// архитектурных регистров и типов данных для x86_64.
package ir_constants_and_types

// Map8BitRegs - словарь соответствий полных регистров и их однобайтовых представлений
var Map8BitRegs = map[string]string{
	RAX: "al",
	RBX: "bl",
	RCX: "cl",
	RDX: "dl",
	R8:  "r8b",
	R9:  "r9b",
	R10: "r10b",
	R11: "r11b",
	R12: "r12b",
	R13: "r13b",
	R14: "r14b",
	R15: "r15b",
}

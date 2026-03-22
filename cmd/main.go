package main

import (
	"fmt"

	"go.mod/internal/ir"
	"go.mod/internal/lexer"
	"go.mod/internal/parser"
)

func main() {
	// 1. Исходный код на твоем языке
	// Простейший скрипт: объявляем пути и копируем файл
	input := `
	string source = "/etc/passwd"
	string target = "/tmp/passwd.bak"

	copy(source, target)
	
	exit(0)
	`

	fmt.Println("[1/4] Лексический и синтаксический анализ (AST)...")
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("[!] Ошибки парсинга:")
		for _, err := range p.Errors() {
			fmt.Println("  -", err)
		}
		return
	}
	fmt.Println("  + AST успешно построено.")

	// 2. Инициализируем общий контекст для счетчиков виртуальных регистров
	ctx := &ir.IRContext{}

	// 3. Генерация High-Level IR
	fmt.Println("[2/4] Генерация абстрактного кода (HLIR)...")
	generator := ir.NewGenerator(ctx)
	generator.Generate(program) // Метод заполняет generator.Instructions

	// Забираем сгенерированные инструкции
	hlir := generator.Instructions

	fmt.Println("\n--- HIGH LEVEL IR (с макросом copy) ---")
	for _, inst := range hlir {
		fmt.Println(inst.String())
	}

	// 4. Раскрутка макросов (Lowering)
	fmt.Println("\n[3/4] Распутывание сложных команд (LLIR)...")
	expander := ir.NewMacroExpander(ctx)

	// Метод Expand принимает массив инструкций и возвращает новый, распутанный массив
	llir := expander.Expand(hlir)

	fmt.Println("\n--- LOW LEVEL IR (чистые системные вызовы) ---")
	for _, inst := range llir {
		fmt.Println(inst.String())
	}

	// 5. Аллокация физических регистров и генерация ASM
	fmt.Println("\n[4/4] Аллокация регистров (x86_64) и генерация NASM...")
	allocator := ir.NewRegisterAllocator()

	// Передаем распутанный IR и общее количество использованных регистров (для выделения места на стеке)
	finalAsm := allocator.Allocate(llir, ctx.GetVRegCount())

	fmt.Println("\n--- ФИНАЛЬНЫЙ КОД (ASM) ---")
	fmt.Println(finalAsm)
}

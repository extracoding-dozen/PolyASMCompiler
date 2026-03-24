/*
input := `
	string log_file = "/tmp/hack.log"
	string backup = "/tmp/hack.bak"
	string hacker_name = "phantom"

		qword size = get_file_size(log_file)

		if size > 0 {
				write(log_file, size, "System compromised\n")

				rename(log_file, backup)
	}

			chmod("/tmp/hack.bak", 0777)
	sleep(10)

	exit(0)
	`
*/

package main

import (
	"fmt"
	"strings"

	"go.mod/external/ui/control_flow_graph_visualizer"
	"go.mod/external/ui/main_page_ui"
	"go.mod/pkg/ir"
	"go.mod/pkg/lexer"
	"go.mod/pkg/obfuscator"
	"go.mod/pkg/parser"
)

func main() {

	orchestrator := func(req main_page_ui.CompileRequest) (string, string, error) {
		var logs strings.Builder

		logs.WriteString("[INFO] Начало компиляции...\n")

		logs.WriteString("[1/4] Лексический и синтаксический анализ (AST)...\n")
		l := lexer.New(req.Code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				logs.WriteString(fmt.Sprintf("  - Ошибка: %s\n", err))
			}
			return "", logs.String(), fmt.Errorf("ошибка парсинга")
		}
		logs.WriteString("  + AST успешно построено.\n")

		ctx := &ir.IRContext{}

		logs.WriteString("[2/4] Генерация абстрактного кода (HLIR)...\n")
		generator := ir.NewGenerator(ctx)
		generator.Generate(program)
		hlir := generator.Instructions

		logs.WriteString("[3/4] Распутывание сложных команд (LLIR)...\n")
		expander := ir.NewMacroExpander(ctx)
		llir := expander.Expand(hlir)

		vis := control_flow_graph_visualizer.NewControlFlowGraphVisualizerImpl()
		vis.GenerateHTMLFile(llir, "clear.html")
		logs.WriteString("  + Граф ДО обфускации сохранен (clear.html)\n")

		if req.EnableObfuscation {
			logs.WriteString("[INFO] Применение обфускации...\n")

			obfsConfig := obfuscator.ObfuscatorConfig{
				EnableSandboxNoise: req.EnableSandboxNoise,
				EnableStringCrypt:  req.EnableStringCrypt,
				EnableOpaquePreds:  req.EnableOpaquePreds,
				NoiseFrequency:     req.NoiseFrequency,
				OpaqueFrequency:    req.OpaqueFrequency,
			}

			obfs := obfuscator.NewObfuscator(expander.GetCtx(), obfsConfig)

			llir = obfs.Obfuscate(llir)

			vis.GenerateHTMLFile(llir, "obfuscated.html")
			logs.WriteString("  + Граф ПОСЛЕ обфускации сохранен (obfuscated.html)\n")
		} else {
			logs.WriteString("[INFO] Обфускация ОТКЛЮЧЕНА пользователем.\n")

			vis.GenerateHTMLFile(llir, "obfuscated.html")
		}

		logs.WriteString("[4/4] Аллокация регистров...\n")
		allocator := ir.NewRegisterAllocator()
		finalAsm := allocator.Allocate(llir, ctx.GetVRegCount())

		logs.WriteString("[SUCCESS] Сборка успешно завершена!\n")

		return finalAsm, logs.String(), nil
	}

	um := main_page_ui.NewApplicationUI("localhost", 8080, orchestrator)

	um.Start()
}

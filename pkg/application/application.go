package application

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

// Application хранит только настройки сервера или глобальные сервисы.
type Application struct {
	host string
	port int
}

func NewApplication() *Application {
	return &Application{
		host: "localhost",
		port: 8080,
	}
}

// makeRequest создает запрос на обработку программой
func (app *Application) makeRequest(req main_page_ui.CompileRequest) (string, string, error) {
	var logs strings.Builder
	lex := lexer.New()
	lex.SetInput(req.Code)
	pars := parser.New(lex)
	ctx := &ir.IRContext{}

	generator := ir.NewGenerator(ctx)
	expander := ir.NewMacroExpander(ctx)
	allocator := ir.NewRegisterAllocator()
	vis := control_flow_graph_visualizer.NewControlFlowGraphVisualizerImpl()

	logs.WriteString("[INFO] Начало компиляции...\n")
	logs.WriteString("[1/4] Лексический и синтаксический анализ (AST)...\n")

	program := pars.ParseProgram()
	if len(pars.Errors()) != 0 {
		for _, err := range pars.Errors() {
			logs.WriteString(fmt.Sprintf("  - Ошибка: %s\n", err))
		}
		return "", logs.String(), fmt.Errorf("ошибка парсинга")
	}
	logs.WriteString("  + AST успешно построено.\n")

	logs.WriteString("[2/4] Генерация абстрактного кода (HLIR)...\n")
	generator.Generate(program)
	hlir := generator.Instructions

	logs.WriteString("[3/4] Распутывание сложных команд (LLIR)...\n")
	llir := expander.Expand(hlir)

	vis.GenerateHTMLFile(llir, "clear.html")
	logs.WriteString("  + Граф ДО обфускации сохранен\n")

	if req.EnableObfuscation {
		logs.WriteString("[INFO] Применение обфускации...\n")

		obfsConfig := obfuscator.ObfuscatorConfig{
			EnableSandboxNoise: req.EnableSandboxNoise,
			EnableStringCrypt:  req.EnableStringCrypt,
			EnableOpaquePreds:  req.EnableOpaquePreds,
			NoiseFrequency:     req.NoiseFrequency,
			OpaqueFrequency:    req.OpaqueFrequency,
			ObfuscateRepeat:    req.RepeatObfuscator,
		}

		obfs := obfuscator.NewObfuscator(expander.GetCtx(), obfsConfig)
		llir = obfs.Obfuscate(llir)

		vis.GenerateHTMLFile(llir, "obfuscated.html")
		logs.WriteString("  + Граф ПОСЛЕ обфускации сохранен\n")
	} else {
		logs.WriteString("[INFO] Обфускация ОТКЛЮЧЕНА пользователем.\n")
	}

	logs.WriteString("[4/4] Аллокация регистров...\n")
	finalAsm := allocator.Allocate(llir, ctx.GetVRegCount())

	logs.WriteString("[SUCCESS] Сборка успешно завершена!\n")

	return finalAsm, logs.String(), nil
}

// Run запускает приложение
func (app *Application) Run() {
	um := main_page_ui.NewApplicationUI(app.host, app.port, app.makeRequest)
	um.Start()
}

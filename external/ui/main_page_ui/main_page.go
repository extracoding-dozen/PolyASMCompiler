package main_page_ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type CompileRequest struct {
	Code              string `json:"code"`
	OptimizationLevel string `json:"optimizationLevel"`
	DebugMode         bool   `json:"debugMode"`

	EnableObfuscation  bool `json:"enableObfuscation"`
	EnableSandboxNoise bool `json:"enableSandboxNoise"`
	EnableStringCrypt  bool `json:"enableStringCrypt"`
	EnableOpaquePreds  bool `json:"enableOpaquePreds"`
	NoiseFrequency     int  `json:"noiseFrequency"`
	OpaqueFrequency    int  `json:"opaqueFrequency"`
}
type CompileResponse struct {
	CompiledCode string `json:"compiledCode"`
	Logs         string `json:"logs"`
}

type CompileHandlerFunc func(req CompileRequest) (asm string, logs string, err error)

type ApplicationUI struct {
	host         string
	port         int
	orchestrator CompileHandlerFunc
}

//go:embed static
var staticFiles embed.FS

func NewApplicationUI(url string, port int, orch CompileHandlerFunc) *ApplicationUI {
	return &ApplicationUI{host: url, port: port, orchestrator: orch}
}

func (ap *ApplicationUI) Start() {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Ошибка при чтении встроенных файлов:", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// 2. Явные ручки для отдачи графов (Безопасно и надежно)
	mux.HandleFunc("/graph/clear", ap.serveClearGraph)
	mux.HandleFunc("/graph/obfuscated", ap.serveObfuscatedGraph)

	// 3. Ручка компиляции
	mux.HandleFunc("/api/compile", ap.compileEndpoint)

	addr := fmt.Sprintf(":%d", ap.port)
	url := "http://" + ap.host + addr

	fmt.Printf("Сервер UI запущен на %s\n", url)
	go ap.openBrowser(url)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// --- НОВЫЕ ФУНКЦИИ ДЛЯ РАЗДАЧИ ГРАФОВ ---

func (ap *ApplicationUI) serveClearGraph(w http.ResponseWriter, r *http.Request) {
	// Читаем файл прямо с диска
	content, err := os.ReadFile("clear.html")
	if err != nil {
		// Если файла нет (еще не скомпилировали) - отдаем красивую ошибку
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h2 style='color:red; font-family:sans-serif;'>Граф еще не сгенерирован!</h2><p style='font-family:sans-serif;'>Сначала нажмите кнопку <b>RUN COMPILATION</b> в интерфейсе.</p>")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

func (ap *ApplicationUI) serveObfuscatedGraph(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("obfuscated.html")
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h2 style='color:red; font-family:sans-serif;'>Граф еще не сгенерирован!</h2><p style='font-family:sans-serif;'>Сначала нажмите кнопку <b>RUN COMPILATION</b> в интерфейсе.</p>")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// ----------------------------------------

func (ap *ApplicationUI) compileEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req CompileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка парсинга", http.StatusBadRequest)
		return
	}

	asm, logs, err := ap.orchestrator(req)
	if err != nil {
		logs += fmt.Sprintf("\n[ОШИБКА]: %v", err)
	}

	res := CompileResponse{CompiledCode: asm, Logs: logs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (ap *ApplicationUI) openBrowser(url string) {
	time.Sleep(100 * time.Millisecond)
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		log.Printf("Не удалось открыть браузер: %v", err)
	}
}

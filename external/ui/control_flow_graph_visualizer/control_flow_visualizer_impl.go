package control_flow_graph_visualizer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.mod/pkg/ir"
	"go.mod/pkg/ir/ir_constants_and_types"
)

// basicBlock представляет собой один узел в графе (окно с кодом)
type basicBlock struct {
	ID           string
	Instructions []ir.Instruction
	TrueTarget   string
	FalseTarget  string
}

type cyElement struct {
	Data    map[string]string `json:"data"`
	Classes string            `json:"classes,omitempty"`
}

type ControlFlowGraphVisualizerImpl struct{}

func NewControlFlowGraphVisualizerImpl() ControlFlowGraphVisualizer {
	return &ControlFlowGraphVisualizerImpl{}
}

// buildCFG разбивает плоский IR-код на блоки
func (cfv *ControlFlowGraphVisualizerImpl) buildCFG(instructions []ir.Instruction) map[string]*basicBlock {
	blocks := make(map[string]*basicBlock)
	var currentBlock *basicBlock
	blockCounter := 0

	newBlock := func(label string) *basicBlock {
		id := label
		if id == "" {
			id = fmt.Sprintf("block_%d", blockCounter)
			blockCounter++
		}
		b := &basicBlock{ID: id, Instructions: []ir.Instruction{}}
		blocks[id] = b
		return b
	}

	currentBlock = newBlock("ENTRY")

	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]

		if inst.Op == ir_constants_and_types.LABEL {
			labelName := inst.Dst.(ir.Lbl).Name

			if len(currentBlock.Instructions) == 0 {
				delete(blocks, currentBlock.ID)
				currentBlock.ID = labelName
				blocks[labelName] = currentBlock
			} else {
				currentBlock.FalseTarget = labelName
				currentBlock = newBlock(labelName)
			}

			// Добавляем саму инструкцию метки для наглядности
			currentBlock.Instructions = append(currentBlock.Instructions, inst)
			continue
		}

		currentBlock.Instructions = append(currentBlock.Instructions, inst)

		if inst.Op == ir_constants_and_types.JMP || inst.Op == ir_constants_and_types.JE || inst.Op == ir_constants_and_types.JNE {
			targetLabel := inst.Dst.(ir.Lbl).Name
			currentBlock.TrueTarget = targetLabel

			if inst.Op != ir_constants_and_types.JMP && i+1 < len(instructions) {
				nextInst := instructions[i+1]
				nextLabel := ""
				if nextInst.Op == ir_constants_and_types.LABEL {
					nextLabel = nextInst.Dst.(ir.Lbl).Name
				} else {
					nextLabel = fmt.Sprintf("block_%d", blockCounter)
				}
				currentBlock.FalseTarget = nextLabel
			}

			currentBlock = newBlock("")
		}
	}

	// ФИНАЛЬНАЯ ЗАЧИСТКА: Удаляем абсолютно пустые блоки из памяти перед рендером
	cleanBlocks := make(map[string]*basicBlock)
	for id, block := range blocks {
		if len(block.Instructions) > 0 {
			cleanBlocks[id] = block
		}
	}

	return cleanBlocks
}

// colorize подсвечивает синтаксис для HTML (Цвета в стиле Dark Theme)
func colorize(inst ir.Instruction) string {
	str := inst.String()

	// Подсветка опкодов
	str = strings.Replace(str, string(inst.Op), fmt.Sprintf("<span class='op'>%s</span>", inst.Op), 1)

	// Подсветка строк и макросов
	if inst.Op == ir_constants_and_types.LOAD_STR {
		str = strings.Replace(str, inst.Src1.String(), fmt.Sprintf("<span class='str'>%s</span>", inst.Src1.String()), 1)
	}
	if inst.Op == ir_constants_and_types.SYSCALL || strings.HasPrefix(string(inst.Op), "MACRO_") {
		str = fmt.Sprintf("<span class='macro'>%s</span>", str)
	}

	return str
}

// GenerateHTMLFile создает независимый HTML файл с графом
func (cfv *ControlFlowGraphVisualizerImpl) GenerateHTMLFile(instructions []ir.Instruction, filename string) error {
	blocks := cfv.buildCFG(instructions)

	var elements []cyElement

	// 1. Создаем Узлы (Nodes)
	for id, block := range blocks {
		if len(block.Instructions) == 0 && id != "ENTRY" {
			continue
		} // Пропускаем пустые

		var htmlContent strings.Builder
		htmlContent.WriteString(fmt.Sprintf("<div class='header'>%s</div><div class='code'>", id))
		for _, inst := range block.Instructions {
			htmlContent.WriteString(colorize(inst) + "<br/>")
		}
		htmlContent.WriteString("</div>")

		elements = append(elements, cyElement{
			Data: map[string]string{"id": id, "label": htmlContent.String()},
		})
	}

	// 2. Создаем Связи (Edges)
	for id, block := range blocks {
		if block.TrueTarget != "" {
			// Безусловный JMP - синий. Условный True (JE/JNE) - зеленый.
			edgeClass := "edge-true"
			lastInst := block.Instructions[len(block.Instructions)-1]
			if lastInst.Op == ir_constants_and_types.JMP {
				edgeClass = "edge-jmp"
			}

			elements = append(elements, cyElement{
				Data:    map[string]string{"source": id, "target": block.TrueTarget},
				Classes: edgeClass,
			})
		}
		if block.FalseTarget != "" {
			// Ветка False (Fallthrough) - красная. Обычный переход - серый.
			edgeClass := "edge-false"
			lastInst := block.Instructions[len(block.Instructions)-1]
			if lastInst.Op != ir_constants_and_types.JE && lastInst.Op != ir_constants_and_types.JNE {
				edgeClass = "edge-next"
			}

			elements = append(elements, cyElement{
				Data:    map[string]string{"source": id, "target": block.FalseTarget},
				Classes: edgeClass,
			})
		}
	}

	jsonData, _ := json.Marshal(elements)

	// 3. Шаблон HTML страницы
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>PolyASM Control Flow Graph</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.23.0/cytoscape.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/dagre/0.8.5/dagre.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/cytoscape-dagre@2.5.0/cytoscape-dagre.min.js"></script>
    <style>
        body { background-color: #1e1e1e; margin: 0; font-family: monospace; }
        #cy { width: 100vw; height: 100vh; display: block; }
        
        /* Стили для генерации SVG-картинок внутри узлов */
        .op { color: #569cd6; font-weight: bold; }
        .str { color: #ce9178; }
        .macro { color: #c586c0; font-style: italic; }
    </style>
</head>
<body>
    <div id="cy"></div>
    <script>
        // Регистрация плагина Dagre (идеальное расположение сверху-вниз)
        cytoscape.use( cytoscapeDagre );

        const elements = ` + string(jsonData) + `;

        const cy = cytoscape({
            container: document.getElementById('cy'),
            elements: elements,
            style: [
                {
                    selector: 'node',
                    style: {
                        'shape': 'round-rectangle',
                        'background-color': '#252526',
                        'border-width': 1,
                        'border-color': '#333',
                        // Магия рендера HTML-таблиц внутри узлов графа через SVG-оболочку
                        'background-image': function(node) {
                            const html = node.data('label');
                            const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200">' +
                                '<foreignObject x="0" y="0" width="100%" height="100%">' +
                                '<div xmlns="http://www.w3.org/1999/xhtml" style="font-family: Consolas, monospace; font-size: 12px; color: #d4d4d4; padding: 10px;">' +
                                html +
                                '</div></foreignObject></svg>';
                            return 'data:image/svg+xml;utf8,' + encodeURIComponent(svg);
                        },
                        'background-fit': 'none',
                        'background-position-x': '0px',
                        'background-position-y': '0px',
                        'width': 350,
                        'height': function(node) {
                            return (node.data('label').match(/<br\/>/g) || []).length * 16 + 40;
                        }
                    }
                },
                {
                    selector: 'edge',
                    style: {
                        'curve-style': 'bezier',
                        'target-arrow-shape': 'triangle',
                        'width': 2
                    }
                },
                { selector: '.edge-jmp', style: { 'line-color': '#007acc', 'target-arrow-color': '#007acc' } },
                { selector: '.edge-true', style: { 'line-color': '#4CAF50', 'target-arrow-color': '#4CAF50' } },
                { selector: '.edge-false', style: { 'line-color': '#F44336', 'target-arrow-color': '#F44336' } },
                { selector: '.edge-next', style: { 'line-color': '#808080', 'target-arrow-color': '#808080' } }
            ],
            layout: {
                name: 'dagre',
                nodeSep: 50,
                rankSep: 100,
                rankDir: 'TB' // Сверху вниз
            }
        });

        // Интерактивность: Зум и панорамирование включены по умолчанию
    </script>
</body>
</html>`

	return os.WriteFile(filename, []byte(htmlTemplate), 0644)
}

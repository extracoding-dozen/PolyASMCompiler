package control_flow_graph_visualizer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.mod/pkg/ir"
	"go.mod/pkg/ir/ir_constants_and_types"
)

// basicBlock представляет собой один узел в графе
type basicBlock struct {
	ID           string
	Instructions []ir.Instruction
	TrueTarget   string
	FalseTarget  string
}

type cyElement struct {
	Data    map[string]interface{} `json:"data"`
	Classes string                 `json:"classes,omitempty"`
}

type ControlFlowGraphVisualizerImpl struct{}

func NewControlFlowGraphVisualizerImpl() ControlFlowGraphVisualizer {
	return &ControlFlowGraphVisualizerImpl{}
}

// Вспомогательная функция для проброса связей в обход пустых блоков
func resolveTarget(blocks map[string]*basicBlock, target string) string {
	visited := make(map[string]bool)
	curr := target

	for {
		if visited[curr] {
			break // Защита от бесконечного цикла
		}
		visited[curr] = true

		b, exists := blocks[curr]
		if !exists {
			break
		}
		// Если блок пустой и у него есть путь дальше (fallthrough) - идем по нему
		if len(b.Instructions) == 0 && b.FalseTarget != "" {
			curr = b.FalseTarget
			continue
		}
		break
	}
	return curr
}

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

			currentBlock.Instructions = append(currentBlock.Instructions, inst)
			continue
		}

		currentBlock.Instructions = append(currentBlock.Instructions, inst)

		if inst.Op == ir_constants_and_types.JMP || inst.Op == ir_constants_and_types.JE || inst.Op == ir_constants_and_types.JNE {
			currentBlock.TrueTarget = inst.Dst.(ir.Lbl).Name

			if inst.Op != ir_constants_and_types.JMP && i+1 < len(instructions) {
				nextLabel := fmt.Sprintf("block_%d", blockCounter)
				if instructions[i+1].Op == ir_constants_and_types.LABEL {
					nextLabel = instructions[i+1].Dst.(ir.Lbl).Name
				}
				currentBlock.FalseTarget = nextLabel
			}

			currentBlock = newBlock("")
		}
	}

	// 1. Убираем пустые блоки из памяти
	cleanBlocks := make(map[string]*basicBlock)
	for id, block := range blocks {
		if len(block.Instructions) > 0 {
			cleanBlocks[id] = block
		}
	}

	// 2. Перенаправляем (rewiring) связи в обход удаленных пустых блоков
	for _, block := range cleanBlocks {
		if block.TrueTarget != "" {
			block.TrueTarget = resolveTarget(blocks, block.TrueTarget)
		}
		if block.FalseTarget != "" {
			block.FalseTarget = resolveTarget(blocks, block.FalseTarget)
		}
	}

	return cleanBlocks
}

// colorize подсвечивает синтаксис (Используем INLINE стили для безопасного SVG рендера)
func colorize(inst ir.Instruction) string {
	str := inst.String()

	// Цвета в стиле VS Code Dark
	colorOp := "#569cd6"  // Синий для инструкций
	colorStr := "#ce9178" // Оранжевый для строк
	colorMac := "#c586c0" // Фиолетовый для макросов

	str = strings.Replace(str, string(inst.Op), fmt.Sprintf("<span style='color:%s; font-weight:bold;'>%s</span>", colorOp, inst.Op), 1)

	if inst.Op == ir_constants_and_types.LOAD_STR {
		str = strings.Replace(str, inst.Src1.String(), fmt.Sprintf("<span style='color:%s;'>%s</span>", colorStr, inst.Src1.String()), 1)
	}
	if inst.Op == ir_constants_and_types.SYSCALL || strings.HasPrefix(string(inst.Op), "MACRO_") {
		str = fmt.Sprintf("<span style='color:%s; font-style:italic;'>%s</span>", colorMac, str)
	}

	return str
}

func (cfv *ControlFlowGraphVisualizerImpl) GenerateHTMLFile(instructions []ir.Instruction, filename string) error {
	blocks := cfv.buildCFG(instructions)
	var elements []cyElement

	// 1. Создаем Узлы
	for id, block := range blocks {
		var htmlContent strings.Builder

		linesCount := len(block.Instructions)

		// ФИКС 1: Если блок пустой, добавляем заглушку NOP
		if linesCount == 0 {
			htmlContent.WriteString("<span style='color:#569cd6; font-weight:bold;'>NOP</span> <span style='color:#6a9955;'>// opaque / empty block</span><br/>")
			linesCount = 1 // Чтобы высота узла считалась минимум для 1 строки
		} else {
			for _, inst := range block.Instructions {
				htmlContent.WriteString(colorize(inst) + "<br/>")
			}
		}

		// Высчитываем высоту узла в зависимости от кол-ва строк кода
		nodeHeight := (linesCount * 18) + 40 // 18px на строку + 40px на отступы и хедер

		elements = append(elements, cyElement{
			Data: map[string]interface{}{
				"id":     id,
				"label":  htmlContent.String(),
				"height": nodeHeight,
			},
		})
	}

	// 2. Создаем Связи (Edges)
	for id, block := range blocks {
		if block.TrueTarget != "" {
			edgeClass := "edge-true"
			if len(block.Instructions) > 0 {
				lastInst := block.Instructions[len(block.Instructions)-1]
				if lastInst.Op == ir_constants_and_types.JMP {
					edgeClass = "edge-jmp"
				}
			}

			// Проверка, что цель существует (защита от битых связей)
			if _, exists := blocks[block.TrueTarget]; exists {
				elements = append(elements, cyElement{
					Data:    map[string]interface{}{"source": id, "target": block.TrueTarget},
					Classes: edgeClass,
				})
			}
		}
		if block.FalseTarget != "" {
			edgeClass := "edge-false"
			if len(block.Instructions) > 0 {
				lastInst := block.Instructions[len(block.Instructions)-1]
				if lastInst.Op != ir_constants_and_types.JE && lastInst.Op != ir_constants_and_types.JNE {
					edgeClass = "edge-next"
				}
			}

			if _, exists := blocks[block.FalseTarget]; exists {
				elements = append(elements, cyElement{
					Data:    map[string]interface{}{"source": id, "target": block.FalseTarget},
					Classes: edgeClass,
				})
			}
		}
	}

	jsonData, _ := json.Marshal(elements)

	// ФИКС 2: Добавлены кнопки ЗУМА и стили в HTML шаблон
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>PolyASM Visualizer</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.23.0/cytoscape.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/dagre/0.8.5/dagre.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/cytoscape-dagre@2.5.0/cytoscape-dagre.min.js"></script>
    <style>
        body { background-color: #1e1e1e; margin: 0; padding: 0; font-family: monospace; overflow: hidden; }
        #cy { width: 100vw; height: 100vh; position: absolute; top: 0; left: 0; }
        #controls { position: absolute; bottom: 20px; right: 20px; background: rgba(0,0,0,0.7); padding: 10px; border-radius: 8px; color: #fff; font-family: sans-serif; font-size: 12px; z-index: 10; border: 1px solid #444; }
        
        /* Стили для кнопок управления зумом */
        .zoom-controls { position: absolute; top: 20px; right: 20px; display: flex; flex-direction: column; gap: 8px; z-index: 1000; }
        .btn-action { width: 40px; height: 40px; background-color: #2d2d30; color: #cccccc; border: 1px solid #444; border-radius: 6px; font-size: 24px; font-weight: bold; cursor: pointer; display: flex; justify-content: center; align-items: center; box-shadow: 0 4px 6px rgba(0,0,0,0.3); transition: 0.1s ease-in-out; }
        .btn-action:hover { background-color: #3e3e42; color: #ffffff; border-color: #666; }
        .btn-action:active { transform: scale(0.95); }
    </style>
</head>
<body>
    <div id="cy"></div>
    <div id="controls">Mouse Wheel: Zoom | Drag: Pan</div>

    <!-- Кнопки зума -->
    <div class="zoom-controls">
        <button class="btn-action" id="zoom-in" title="Zoom In">+</button>
        <button class="btn-action" id="zoom-out" title="Zoom Out">−</button>
        <button class="btn-action" id="zoom-fit" title="Center Graph" style="font-size: 18px;">⛶</button>
    </div>

    <script>
        cytoscape.use( cytoscapeDagre );
        const elements = ` + string(jsonData) + `;

        const cy = cytoscape({
            container: document.getElementById('cy'),
            elements: elements,
            wheelSensitivity: 0.2, /* Делаем зум колесиком более плавным */
            style: [
                {
                    selector: 'node',
                    style: {
                        'shape': 'round-rectangle',
                        'background-color': 'transparent', 
                        'width': 350,
                        'height': 'data(height)',
                        'background-image': function(node) {
                            const html = node.data('label');
                            const id = node.data('id');
                            const h = node.data('height');
                            
                            const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="350" height="' + h + '">' +
                                '<rect width="350" height="' + h + '" rx="8" ry="8" fill="#1e1e1e" stroke="#555" stroke-width="2"/>' +
                                '<rect width="350" height="26" rx="8" ry="8" fill="#2d2d30"/>' +
                                '<path d="M0 20 L0 26 L350 26 L350 20 Z" fill="#2d2d30"/>' +
                                '<text x="12" y="18" fill="#9cdcfe" font-family="sans-serif" font-size="13" font-weight="bold">' + id + '</text>' +
                                '<foreignObject x="0" y="26" width="100%" height="100%">' +
                                '<div xmlns="http://www.w3.org/1999/xhtml" style="font-family: Consolas, \'Courier New\', monospace; font-size: 13px; color: #d4d4d4; padding: 10px; line-height: 1.4; white-space: nowrap;">' +
                                html +
                                '</div></foreignObject></svg>';
                                
                            return 'data:image/svg+xml;utf8,' + encodeURIComponent(svg);
                        },
                        'background-fit': 'none',
                        'background-position-x': '0px',
                        'background-position-y': '0px',
                    }
                },
                {
                    selector: 'edge',
                    style: {
                        'curve-style': 'bezier',
                        'target-arrow-shape': 'triangle',
                        'width': 2,
                        'arrow-scale': 1.5
                    }
                },
                { selector: '.edge-jmp', style: { 'line-color': '#007acc', 'target-arrow-color': '#007acc' } },
                { selector: '.edge-true', style: { 'line-color': '#4CAF50', 'target-arrow-color': '#4CAF50' } },
                { selector: '.edge-false', style: { 'line-color': '#F44336', 'target-arrow-color': '#F44336' } },
                { selector: '.edge-next', style: { 'line-color': '#666666', 'target-arrow-color': '#666666' } }
            ],
            layout: {
                name: 'dagre',
                nodeSep: 60,
                rankSep: 120,
                rankDir: 'TB',
                fit: false
            }
        });

        // ЛОГИКА КНОПОК ЗУМА
        document.getElementById('zoom-in').addEventListener('click', () => {
            cy.zoom(cy.zoom() * 1.25); // Увеличиваем на 25%
        });

        document.getElementById('zoom-out').addEventListener('click', () => {
            cy.zoom(cy.zoom() * 0.8);  // Уменьшаем на 20%
        });

        document.getElementById('zoom-fit').addEventListener('click', () => {
            cy.fit(50); // Отцентровать граф (с отступом 50px от краев экрана)
        });

        // ЛОГИКА ФОКУСИРОВКИ
        cy.ready(function() {
            let rootNode = cy.nodes('#ENTRY');
            
            if (rootNode.empty()) {
                rootNode = cy.nodes().roots().first();
            }
            if (rootNode.empty()) {
                rootNode = cy.nodes().first();
            }

            if (!rootNode.empty()) {
                cy.animate({
                    zoom: 1.1,
                    center: { eles: rootNode },
                }, {
                    duration: 500
                });
            }
        });
    </script>
</body>
</html>`

	return os.WriteFile(filename, []byte(htmlTemplate), 0644)
}

package graph

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.mod/internal/ast"
)

// --- Структуры данных Vis.js ---

type GraphNode struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Group string `json:"group"`
}

type EdgeColor struct {
	Color     string `json:"color"`
	Highlight string `json:"highlight"`
}

type EdgeSmooth struct {
	Type      string  `json:"type"`
	Roundness float64 `json:"roundness"`
}

type GraphEdge struct {
	From   int         `json:"from"`
	To     int         `json:"to"`
	Label  string      `json:"label,omitempty"`
	Color  *EdgeColor  `json:"color,omitempty"`
	Smooth *EdgeSmooth `json:"smooth,omitempty"` // Для огибания графа петлями
}

type ExitNode struct {
	ID    int
	Label string
	Color string
}

type GraphVisualizer struct {
	nodes   []GraphNode
	edges   []GraphEdge
	counter int
}

func GenerateASTGraph(program *ast.Program) string {
	v := &GraphVisualizer{
		nodes: []GraphNode{},
		edges: []GraphEdge{},
	}

	v.buildCFG(program)

	nodesJSON, _ := json.Marshal(v.nodes)
	edgesJSON, _ := json.Marshal(v.edges)

	html := strings.Replace(graphHTMLTemplate, "{{NODES}}", string(nodesJSON), 1)
	html = strings.Replace(html, "{{EDGES}}", string(edgesJSON), 1)

	return html
}

func (v *GraphVisualizer) addNode(label, group string) int {
	v.counter++
	id := v.counter
	v.nodes = append(v.nodes, GraphNode{
		ID:    id,
		Label: label,
		Group: group,
	})
	return id
}

func (v *GraphVisualizer) addEdge(from, to int, label string, color string, smooth *EdgeSmooth) {
	if from == 0 || to == 0 {
		return
	}
	edge := GraphEdge{From: from, To: to, Label: label}
	if color != "" {
		edge.Color = &EdgeColor{Color: color, Highlight: color}
	}
	if smooth != nil {
		edge.Smooth = smooth
	}
	v.edges = append(v.edges, edge)
}

func (v *GraphVisualizer) buildCFG(node ast.Node) (int, []ExitNode) {
	if node == nil {
		return 0, nil
	}

	switch n := node.(type) {
	case *ast.Program:
		entryID := v.addNode("ENTRY", "root")
		currentExits := []ExitNode{{ID: entryID}}

		for _, stmt := range n.Statements {
			stmtEntry, stmtExits := v.buildCFG(stmt)
			if stmtEntry != 0 {
				for _, exit := range currentExits {
					v.addEdge(exit.ID, stmtEntry, exit.Label, exit.Color, nil)
				}
				currentExits = stmtExits
			}
		}

		exitID := v.addNode("EXIT", "root")
		for _, exit := range currentExits {
			v.addEdge(exit.ID, exitID, exit.Label, exit.Color, nil)
		}
		return entryID, []ExitNode{{ID: exitID}}

	case *ast.BlockStatement:
		var blockEntry int
		var currentExits []ExitNode

		for _, stmt := range n.Statements {
			stmtEntry, stmtExits := v.buildCFG(stmt)
			if stmtEntry != 0 {
				if blockEntry == 0 {
					blockEntry = stmtEntry
				} else {
					for _, exit := range currentExits {
						v.addEdge(exit.ID, stmtEntry, exit.Label, exit.Color, nil)
					}
				}
				currentExits = stmtExits
			}
		}
		return blockEntry, currentExits

	case *ast.LetStatement:
		id := v.addNode(fmt.Sprintf("%s = %s", n.Name.String(), n.Value.String()), "statement")
		return id, []ExitNode{{ID: id}}

	case *ast.ExpressionStatement:
		id := v.addNode(n.Expression.String(), "expression")
		return id, []ExitNode{{ID: id}}

	case *ast.ImportStatement:
		id := v.addNode(fmt.Sprintf("import %s", n.Path.String()), "statement")
		return id, []ExitNode{{ID: id}}

	case *ast.IfStatement:
		id := v.addNode(fmt.Sprintf("if (%s)", n.Condition.String()), "control")
		var allExits []ExitNode

		conseqEntry, conseqExits := v.buildCFG(n.Consequence)
		if conseqEntry != 0 {
			v.addEdge(id, conseqEntry, "True", "#4CAF50", nil) // Зеленый
			allExits = append(allExits, conseqExits...)
		}

		if n.Alternative != nil {
			altEntry, altExits := v.buildCFG(n.Alternative)
			if altEntry != 0 {
				v.addEdge(id, altEntry, "False", "#F44336", nil) // Красный
				allExits = append(allExits, altExits...)
			}
		} else {
			allExits = append(allExits, ExitNode{ID: id, Label: "False", Color: "#F44336"})
		}

		return id, allExits

	case *ast.WhileStatement:
		id := v.addNode(fmt.Sprintf("while (%s)", n.Condition.String()), "control")

		bodyEntry, bodyExits := v.buildCFG(n.Body)
		if bodyEntry != 0 {
			v.addEdge(id, bodyEntry, "True", "#4CAF50", nil)

			// --- ИСПРАВЛЕНИЕ: ПЕТЛЯ ТЕПЕРЬ ОГИБАЕТ ГРАФ СБОКУ ---
			loopSmooth := &EdgeSmooth{Type: "curvedCW", Roundness: 0.3}
			for _, exit := range bodyExits {
				v.addEdge(exit.ID, id, "Loop Back", "#2196F3", loopSmooth)
			}
		}

		return id, []ExitNode{{ID: id, Label: "False", Color: "#F44336"}}

	default:
		id := v.addNode(fmt.Sprintf("Unknown: %T", node), "error")
		return id, []ExitNode{{ID: id}}
	}
}

// Шаблон страницы с Vis.js
const graphHTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>IDA Style CFG</title>
    <script type="text/javascript" src="https://unpkg.com/vis-network/standalone/umd/vis-network.min.js"></script>
    <style>
        body, html {
            margin: 0; padding: 0; width: 100%; height: 100%;
            background-color: #1e1e1e; font-family: 'Consolas', monospace; overflow: hidden;
        }
        #mynetwork { width: 100%; height: 100%; }
        #overlay {
            position: absolute; top: 10px; left: 10px; color: #d4d4d4;
            background: rgba(30, 30, 30, 0.8); padding: 10px;
            border: 1px solid #404040; border-radius: 5px; pointer-events: none; z-index: 10;
        }
    </style>
</head>
<body>
<div id="overlay">
    <b>Interactive Control Flow Graph</b><br>
    Scroll: Zoom | Drag: Pan nodes
</div>
<div id="mynetwork"></div>

<script type="text/javascript">
    var nodes = new vis.DataSet({{NODES}});
    var edges = new vis.DataSet({{EDGES}});
    var container = document.getElementById('mynetwork');
    
    var options = {
        layout: {
            hierarchical: {
                direction: "UD",
                sortMethod: "directed",
                levelSeparation: 120, // Расстояние по вертикали
                nodeSpacing: 500,     // ИСПРАВЛЕНИЕ: Значительно увеличен отступ по горизонтали (чтобы текст не наезжал)
                treeSpacing: 450      // Расстояние между ветвями (True/False)
            }
        },
        physics: { enabled: false }, // Оставляем выключенной для строгой иерархии
        interaction: { dragNodes: true, hover: true, zoomView: true, dragView: true },
        edges: {
            smooth: { type: 'cubicBezier', forceDirection: 'vertical', roundness: 0.6 },
            color: { color: '#6a6a6a', highlight: '#ffffff' },
            arrows: { to: { enabled: true, scaleFactor: 0.9 } },
            font: { color: '#ffffff', size: 12, face: 'Consolas', background: '#1e1e1e', strokeWidth: 0 }
        },
        nodes: {
            shape: 'box',
            margin: { top: 12, bottom: 12, left: 15, right: 15 },
            font: { face: 'Consolas', size: 14 }
        },
        groups: {
            root: { color: { background: '#d32f2f', border: '#ff6659' }, font: { color: '#ffffff', bold: true } },
            control: { color: { background: '#2d2d2d', border: '#c586c0' }, font: { color: '#c586c0' } }, 
            statement: { color: { background: '#2d2d2d', border: '#569cd6' }, font: { color: '#569cd6' } },
            expression: { color: { background: '#2d2d2d', border: '#d7ba7d' }, font: { color: '#d7ba7d' } }
        }
    };

    var network = new vis.Network(container, {nodes: nodes, edges: edges}, options);
</script>
</body>
</html>
`

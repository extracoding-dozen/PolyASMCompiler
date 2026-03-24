package drawer

import (
	"fmt"
	"html"
	"strings"

	"go.mod/pkg/ast"
)

// GenerateASTHtml принимает корень AST и возвращает готовую HTML-страницу
func GenerateASTHtml(program *ast.Program) string {
	var sb strings.Builder

	sb.WriteString(htmlHeader)

	sb.WriteString(`<ul class="tree">`)
	sb.WriteString(buildNode("Program", programToHTML(program)))
	sb.WriteString(`</ul>`)

	sb.WriteString(htmlFooter)

	return sb.String()
}

func nodeToHTML(node ast.Node) string {
	if node == nil {
		return ""
	}

	switch n := node.(type) {

	case *ast.LetStatement:
		val := nodeToHTML(n.Value)
		name := spanClass("ident", n.Name.Value)
		title := fmt.Sprintf("LetStatement: %s =", name)
		return buildNode(title, val)

	case *ast.ExpressionStatement:
		return buildNode("ExpressionStatement", nodeToHTML(n.Expression))

	case *ast.BlockStatement:
		var children []string
		for _, stmt := range n.Statements {
			children = append(children, nodeToHTML(stmt))
		}
		return buildNode("BlockStatement", children...)

	case *ast.IfStatement:
		cond := buildNode("Condition", nodeToHTML(n.Condition))
		cons := buildNode("Consequence", nodeToHTML(n.Consequence))

		if n.Alternative != nil {
			alt := buildNode("Alternative", nodeToHTML(n.Alternative))
			return buildNode("IfStatement", cond, cons, alt)
		}
		return buildNode("IfStatement", cond, cons)

	case *ast.WhileStatement:
		cond := buildNode("Condition", nodeToHTML(n.Condition))
		body := buildNode("Body", nodeToHTML(n.Body))
		return buildNode("WhileStatement", cond, body)

	case *ast.ImportStatement:
		path := spanClass("string", html.EscapeString(n.Path.String()))
		return buildLeaf(fmt.Sprintf("ImportStatement: %s", path))

	case *ast.Identifier:
		return buildLeaf(fmt.Sprintf("Identifier: %s", spanClass("ident", html.EscapeString(n.Value))))

	case *ast.IntegerLiteral:
		return buildLeaf(fmt.Sprintf("IntegerLiteral: %s", spanClass("number", fmt.Sprintf("%d", n.Value))))

	case *ast.StringLiteral:
		return buildLeaf(fmt.Sprintf("StringLiteral: %s", spanClass("string", html.EscapeString("\""+n.Value+"\""))))

	case *ast.InfixExpression:
		op := spanClass("operator", html.EscapeString(n.Operator))
		left := buildNode("Left", nodeToHTML(n.Left))
		right := buildNode("Right", nodeToHTML(n.Right))
		return buildNode(fmt.Sprintf("InfixExpression (%s)", op), left, right)

	case *ast.CallExpression:
		funcName := buildNode("Function", nodeToHTML(n.Function))
		var args []string
		for _, arg := range n.Arguments {
			args = append(args, nodeToHTML(arg))
		}

		if len(args) > 0 {
			argsNode := buildNode("Arguments", args...)
			return buildNode("CallExpression", funcName, argsNode)
		}
		return buildNode("CallExpression", funcName, buildLeaf("Arguments: (empty)"))

	default:
		return buildLeaf(fmt.Sprintf("Unknown Node: %T", node))
	}
}

func programToHTML(p *ast.Program) string {
	var sb strings.Builder
	for _, stmt := range p.Statements {
		sb.WriteString(nodeToHTML(stmt))
	}
	return sb.String()
}

func buildNode(title string, children ...string) string {

	hasChildren := false
	for _, c := range children {
		if c != "" {
			hasChildren = true
			break
		}
	}

	if !hasChildren {
		return buildLeaf(title)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<li><span class="caret">%s</span><ul class="nested">`, title))
	for _, child := range children {
		if child != "" {
			sb.WriteString(child)
		}
	}
	sb.WriteString(`</ul></li>`)
	return sb.String()
}

func buildLeaf(content string) string {
	return fmt.Sprintf(`<li><span class="leaf">%s</span></li>`, content)
}

func spanClass(class, text string) string {
	return fmt.Sprintf(`<span class="%s">%s</span>`, class, text)
}

// --- Шаблоны страницы (CSS + JS) ---

const htmlHeader = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AST Visualizer</title>
<style>
	body {
		background-color: #1e1e1e;
		color: #d4d4d4;
		font-family: 'Consolas', 'Courier New', monospace;
		padding: 20px;
	}
	ul, #myUL {
		list-style-type: none;
	}
	#myUL {
		margin: 0;
		padding: 0;
	}
	.caret {
		cursor: pointer;
		user-select: none;
		color: #9cdcfe;
		font-weight: bold;
	}
	.caret::before {
		content: "\25B6";
		color: #808080;
		display: inline-block;
		margin-right: 6px;
		font-size: 12px;
		transition: transform 0.2s;
	}
	.caret-down::before {
		transform: rotate(90deg);
	}
	.nested {
		display: none;
		border-left: 1px dashed #404040;
		padding-left: 20px;
		margin-left: 5px;
		margin-top: 5px;
	}
	.active {
		display: block;
	}
	li {
		margin: 5px 0;
	}
	.leaf {
		color: #cccccc;
		margin-left: 18px;
	}
	
	/* Подсветка синтаксиса */
	.ident { color: #9cdcfe; }
	.number { color: #b5cea8; }
	.string { color: #ce9178; }
	.operator { color: #c586c0; font-weight: bold; }
</style>
</head>
<body>
<h2>Abstract Syntax Tree</h2>
<button onclick="expandAll()">Expand All</button>
<button onclick="collapseAll()">Collapse All</button>
<hr style="border-color: #404040; margin-bottom: 20px;">
`

const htmlFooter = `
<script>
	var toggler = document.getElementsByClassName("caret");
	for (var i = 0; i < toggler.length; i++) {
		toggler[i].addEventListener("click", function() {
			this.parentElement.querySelector(".nested").classList.toggle("active");
			this.classList.toggle("caret-down");
		});
	}

	function expandAll() {
		let nested = document.querySelectorAll('.nested');
		let carets = document.querySelectorAll('.caret');
		nested.forEach(n => n.classList.add('active'));
		carets.forEach(c => c.classList.add('caret-down'));
	}

	function collapseAll() {
		let nested = document.querySelectorAll('.nested');
		let carets = document.querySelectorAll('.caret');
		nested.forEach(n => n.classList.remove('active'));
		carets.forEach(c => c.classList.remove('caret-down'));
	}
</script>
</body>
</html>
`

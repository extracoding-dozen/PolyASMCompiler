package control_flow_graph_visualizer

import "go.mod/pkg/ir"

type ControlFlowGraphVisualizer interface {
	GenerateHTMLFile(instructions []ir.Instruction, filename string) error
}

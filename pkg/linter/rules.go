package linter

import (
	"fmt"
	"strings"

	"github.com/skjutare/mermaid-lint/pkg/parser"
)

// Rule defines the interface for a lint rule.
type Rule interface {
	Name() string
	Description() string
	Check(d *parser.Diagram) []Finding
}

// AllRules returns all available lint rules.
func AllRules() []Rule {
	return []Rule{
		&NoUnknownDiagramType{},
		&NoEmptyDiagram{},
		&ValidDirection{},
		&NoDuplicateNodeIDs{},
		&NodeHasLabel{},
		&NoOrphanNodes{},
	}
}

// --- Rule: no-unknown-diagram-type ---

// NoUnknownDiagramType checks that the diagram type is recognized.
type NoUnknownDiagramType struct{}

func (r *NoUnknownDiagramType) Name() string        { return "no-unknown-diagram-type" }
func (r *NoUnknownDiagramType) Description() string  { return "Diagram type must be a recognized Mermaid type" }

func (r *NoUnknownDiagramType) Check(d *parser.Diagram) []Finding {
	if d.Type == parser.DiagramUnknown && d.TypeRaw != "" {
		return []Finding{{
			Rule:    r.Name(),
			Message: fmt.Sprintf("unknown diagram type %q", d.TypeRaw),
			Line:    d.StartLine,
		}}
	}
	if d.TypeRaw == "" {
		return []Finding{{
			Rule:    r.Name(),
			Message: "no diagram type declaration found",
			Line:    d.StartLine,
		}}
	}
	return nil
}

// --- Rule: no-empty-diagram ---

// NoEmptyDiagram checks that the diagram has meaningful content.
type NoEmptyDiagram struct{}

func (r *NoEmptyDiagram) Name() string        { return "no-empty-diagram" }
func (r *NoEmptyDiagram) Description() string  { return "Diagram must contain at least one element" }

func (r *NoEmptyDiagram) Check(d *parser.Diagram) []Finding {
	// Only check diagram types where we can detect emptiness
	if d.Type != parser.DiagramFlowchart && d.Type != parser.DiagramGraph {
		return nil
	}

	if len(d.Nodes) == 0 && len(d.Edges) == 0 {
		return []Finding{{
			Rule:    r.Name(),
			Message: "diagram has no nodes or edges",
			Line:    d.StartLine,
		}}
	}
	return nil
}

// --- Rule: valid-direction ---

// ValidDirection checks that flowchart direction is valid.
type ValidDirection struct{}

func (r *ValidDirection) Name() string        { return "valid-direction" }
func (r *ValidDirection) Description() string  { return "Flowchart direction must be TB, TD, BT, LR, or RL" }

func (r *ValidDirection) Check(d *parser.Diagram) []Finding {
	if d.Type != parser.DiagramFlowchart && d.Type != parser.DiagramGraph {
		return nil
	}
	if d.Direction == "" {
		// graph without direction defaults to TD, which is fine
		return nil
	}
	dir := strings.ToUpper(d.Direction)
	if !parser.ValidFlowchartDirections[dir] {
		return []Finding{{
			Rule:    r.Name(),
			Message: fmt.Sprintf("invalid flowchart direction %q; must be one of TB, TD, BT, LR, RL", d.Direction),
			Line:    d.StartLine,
		}}
	}
	return nil
}

// --- Rule: no-duplicate-node-ids ---

// NoDuplicateNodeIDs checks for duplicate node IDs within a diagram.
type NoDuplicateNodeIDs struct{}

func (r *NoDuplicateNodeIDs) Name() string        { return "no-duplicate-node-ids" }
func (r *NoDuplicateNodeIDs) Description() string  { return "Node IDs must be unique within a diagram" }

func (r *NoDuplicateNodeIDs) Check(d *parser.Diagram) []Finding {
	if d.Type != parser.DiagramFlowchart && d.Type != parser.DiagramGraph {
		return nil
	}

	seen := make(map[string]int) // id -> first line
	var findings []Finding

	for _, node := range d.Nodes {
		if firstLine, exists := seen[node.ID]; exists {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: fmt.Sprintf("duplicate node ID %q (first defined at line %d)", node.ID, firstLine),
				Line:    node.Line,
			})
		} else {
			seen[node.ID] = node.Line
		}
	}
	return findings
}

// --- Rule: node-has-label ---

// NodeHasLabel checks that nodes have descriptive labels (not just IDs).
type NodeHasLabel struct{}

func (r *NodeHasLabel) Name() string        { return "node-has-label" }
func (r *NodeHasLabel) Description() string  { return "Nodes should have descriptive labels" }

func (r *NodeHasLabel) Check(d *parser.Diagram) []Finding {
	if d.Type != parser.DiagramFlowchart && d.Type != parser.DiagramGraph {
		return nil
	}

	var findings []Finding
	for _, node := range d.Nodes {
		if node.Label == "" {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: fmt.Sprintf("node %q has no label", node.ID),
				Line:    node.Line,
			})
		}
	}
	return findings
}

// --- Rule: no-orphan-nodes ---

// NoOrphanNodes checks for nodes that are not connected to any edges.
type NoOrphanNodes struct{}

func (r *NoOrphanNodes) Name() string        { return "no-orphan-nodes" }
func (r *NoOrphanNodes) Description() string  { return "All nodes should be connected to at least one edge" }

func (r *NoOrphanNodes) Check(d *parser.Diagram) []Finding {
	if d.Type != parser.DiagramFlowchart && d.Type != parser.DiagramGraph {
		return nil
	}

	if len(d.Edges) == 0 {
		return nil // Don't flag orphans if there are no edges at all
	}

	connected := make(map[string]bool)
	for _, edge := range d.Edges {
		connected[edge.From] = true
		connected[edge.To] = true
	}

	var findings []Finding
	for _, node := range d.Nodes {
		if !connected[node.ID] {
			findings = append(findings, Finding{
				Rule:    r.Name(),
				Message: fmt.Sprintf("node %q is not connected to any edge", node.ID),
				Line:    node.Line,
			})
		}
	}
	return findings
}

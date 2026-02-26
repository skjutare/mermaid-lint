// Package parser provides Mermaid diagram parsing capabilities.
package parser

import (
	"regexp"
	"strings"
)

// DiagramType represents a known Mermaid diagram type.
type DiagramType string

const (
	DiagramFlowchart            DiagramType = "flowchart"
	DiagramGraph                DiagramType = "graph"
	DiagramSequence             DiagramType = "sequenceDiagram"
	DiagramClass                DiagramType = "classDiagram"
	DiagramState                DiagramType = "stateDiagram"
	DiagramStateV2              DiagramType = "stateDiagram-v2"
	DiagramER                   DiagramType = "erDiagram"
	DiagramGantt                DiagramType = "gantt"
	DiagramPie                  DiagramType = "pie"
	DiagramGitGraph             DiagramType = "gitGraph"
	DiagramMindmap              DiagramType = "mindmap"
	DiagramTimeline             DiagramType = "timeline"
	DiagramQuadrant             DiagramType = "quadrantChart"
	DiagramRequirement          DiagramType = "requirementDiagram"
	DiagramC4Context            DiagramType = "C4Context"
	DiagramC4Container          DiagramType = "C4Container"
	DiagramC4Component          DiagramType = "C4Component"
	DiagramC4Dynamic            DiagramType = "C4Dynamic"
	DiagramC4Deployment         DiagramType = "C4Deployment"
	DiagramSankey               DiagramType = "sankey-beta"
	DiagramBlock                DiagramType = "block-beta"
	DiagramXYChart              DiagramType = "xychart-beta"
	DiagramUnknown              DiagramType = "unknown"
)

// KnownDiagramTypes lists all recognized Mermaid diagram types.
var KnownDiagramTypes = map[string]DiagramType{
	"flowchart":            DiagramFlowchart,
	"graph":                DiagramGraph,
	"sequencediagram":      DiagramSequence,
	"classdiagram":         DiagramClass,
	"statediagram":         DiagramState,
	"statediagram-v2":      DiagramStateV2,
	"erdiagram":            DiagramER,
	"gantt":                DiagramGantt,
	"pie":                  DiagramPie,
	"gitgraph":             DiagramGitGraph,
	"mindmap":              DiagramMindmap,
	"timeline":             DiagramTimeline,
	"quadrantchart":        DiagramQuadrant,
	"requirementdiagram":   DiagramRequirement,
	"c4context":            DiagramC4Context,
	"c4container":          DiagramC4Container,
	"c4component":          DiagramC4Component,
	"c4dynamic":            DiagramC4Dynamic,
	"c4deployment":         DiagramC4Deployment,
	"sankey-beta":          DiagramSankey,
	"block-beta":           DiagramBlock,
	"xychart-beta":         DiagramXYChart,
}

// ValidFlowchartDirections contains all valid flowchart/graph directions.
var ValidFlowchartDirections = map[string]bool{
	"TB": true, "TD": true, "BT": true, "LR": true, "RL": true,
}

// Node represents a node in a Mermaid diagram.
type Node struct {
	ID    string
	Label string
	Line  int
	Shape string // e.g., "round", "stadium", "rect", "rhombus", "circle", etc.
}

// Edge represents a connection between nodes.
type Edge struct {
	From  string
	To    string
	Label string
	Line  int
	Style string // e.g., "-->", "---", "-.->", "==>"
}

// Diagram represents a parsed Mermaid diagram.
type Diagram struct {
	Type      DiagramType
	TypeRaw   string // The raw type string as written
	Direction string // For flowcharts: TB, TD, BT, LR, RL
	Nodes     []Node
	Edges     []Edge
	Lines     []string // Original source lines
	StartLine int      // Starting line in the original file (1-based)
}

// Parse parses a Mermaid diagram source string into a Diagram.
// startLine is the 1-based line number where this diagram starts in the file.
func Parse(source string, startLine int) *Diagram {
	lines := strings.Split(source, "\n")
	d := &Diagram{
		Lines:     lines,
		StartLine: startLine,
	}

	d.parseType(lines)

	if d.Type == DiagramFlowchart || d.Type == DiagramGraph {
		d.parseFlowchart(lines)
	}

	return d
}

func (d *Diagram) parseType(lines []string) {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}

		// Extract the first word (diagram type keyword)
		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			continue
		}

		keyword := parts[0]
		d.TypeRaw = keyword

		normalized := strings.ToLower(keyword)
		if dt, ok := KnownDiagramTypes[normalized]; ok {
			d.Type = dt
		} else {
			d.Type = DiagramUnknown
		}

		// For flowchart/graph, extract direction
		if (d.Type == DiagramFlowchart || d.Type == DiagramGraph) && len(parts) > 1 {
			d.Direction = parts[1]
		}
		return
	}
}

var (
	// Match node definitions: A, A[label], A(label), A{label}, A((label)), A>label], A[/label/], etc.
	nodeDefPattern = regexp.MustCompile(
		`(?:^|[\s;])([A-Za-z_][A-Za-z0-9_]*)` +
			`(?:` +
			`\[([^\]]*)\]` + // [label] rectangle
			`|\(([^)]*)\)` + // (label) rounded
			`|\{([^}]*)\}` + // {label} rhombus
			`|\[\[([^\]]*)\]\]` + // [[label]] subroutine
			`|\[\(([^)]*)\)\]` + // [(label)] cylinder
			`|\(\[([^\]]*)\]\)` + // ([label]) stadium
			`|\(\(([^)]*)\)\)` + // ((label)) circle
			`|>([^\]]*)\]` + // >label] asymmetric
			`|\[/([^\]]*)/\]` + // [/label/] parallelogram
			`|\[\\([^\]]*)\\\]` + // [\label\] parallelogram alt
			`)?`,
	)

	// Match edges between nodes: A --> B, A --- B, A -.-> B, A ==> B, etc.
	edgePattern = regexp.MustCompile(
		`([A-Za-z_][A-Za-z0-9_]*)` +
			`\s*` +
			`(-->|---|-.->|-.-|==>|===|--[>x]|~~>|~~)` +
			`\s*` +
			`(?:\|([^|]*)\|\s*)?` +
			`([A-Za-z_][A-Za-z0-9_]*)`,
	)
)

func (d *Diagram) parseFlowchart(lines []string) {
	seenNodes := make(map[string]bool)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines, comments, and the diagram declaration
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		if i == 0 || (d.isDeclarationLine(trimmed)) {
			continue
		}

		// Skip subgraph/end/style/class directives
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "subgraph") ||
			lower == "end" ||
			strings.HasPrefix(lower, "style ") ||
			strings.HasPrefix(lower, "classDef ") ||
			strings.HasPrefix(lower, "classdef ") ||
			strings.HasPrefix(lower, "click ") ||
			strings.HasPrefix(lower, "linkstyle ") {
			continue
		}

		lineNum := d.StartLine + i

		// Try to parse edges
		edgeMatches := edgePattern.FindAllStringSubmatch(trimmed, -1)
		for _, match := range edgeMatches {
			fromID := match[1]
			style := match[2]
			label := match[3]
			toID := match[4]

			d.Edges = append(d.Edges, Edge{
				From:  fromID,
				To:    toID,
				Label: label,
				Line:  lineNum,
				Style: style,
			})

			// Register nodes referenced in edges
			if !seenNodes[fromID] {
				seenNodes[fromID] = true
				d.Nodes = append(d.Nodes, Node{ID: fromID, Line: lineNum})
			}
			if !seenNodes[toID] {
				seenNodes[toID] = true
				d.Nodes = append(d.Nodes, Node{ID: toID, Line: lineNum})
			}
		}

		// Try to parse standalone node definitions
		nodeMatches := nodeDefPattern.FindAllStringSubmatch(trimmed, -1)
		for _, match := range nodeMatches {
			id := match[1]
			if id == "" {
				continue
			}
			// Determine label from whichever capture group matched
			label := ""
			shape := "default"
			for j := 2; j < len(match); j++ {
				if match[j] != "" {
					label = match[j]
					switch j {
					case 2:
						shape = "rect"
					case 3:
						shape = "round"
					case 4:
						shape = "rhombus"
					case 5:
						shape = "subroutine"
					case 6:
						shape = "cylinder"
					case 7:
						shape = "stadium"
					case 8:
						shape = "circle"
					case 9:
						shape = "asymmetric"
					case 10:
						shape = "parallelogram"
					case 11:
						shape = "parallelogram-alt"
					}
					break
				}
			}

			if !seenNodes[id] {
				seenNodes[id] = true
				d.Nodes = append(d.Nodes, Node{
					ID:    id,
					Label: label,
					Line:  lineNum,
					Shape: shape,
				})
			} else if label != "" {
				// Update existing node with label info
				for k := range d.Nodes {
					if d.Nodes[k].ID == id && d.Nodes[k].Label == "" {
						d.Nodes[k].Label = label
						d.Nodes[k].Shape = shape
						break
					}
				}
			}
		}
	}
}

func (d *Diagram) isDeclarationLine(line string) bool {
	lower := strings.ToLower(strings.TrimSpace(line))
	for keyword := range KnownDiagramTypes {
		if strings.HasPrefix(lower, keyword) {
			return true
		}
	}
	return false
}

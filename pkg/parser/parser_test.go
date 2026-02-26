package parser

import (
	"testing"
)

func TestParse_FlowchartType(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		wantType  DiagramType
		wantDir   string
	}{
		{"flowchart LR", "flowchart LR\n  A --> B", DiagramFlowchart, "LR"},
		{"flowchart TD", "flowchart TD\n  A --> B", DiagramFlowchart, "TD"},
		{"graph TB", "graph TB\n  A --> B", DiagramGraph, "TB"},
		{"graph no direction", "graph\n  A --> B", DiagramGraph, ""},
		{"sequenceDiagram", "sequenceDiagram\n  Alice->>Bob: Hello", DiagramSequence, ""},
		{"classDiagram", "classDiagram\n  class Animal", DiagramClass, ""},
		{"erDiagram", "erDiagram\n  CUSTOMER ||--o{ ORDER : places", DiagramER, ""},
		{"gantt", "gantt\n  title A Gantt", DiagramGantt, ""},
		{"pie", "pie\n  title Pets", DiagramPie, ""},
		{"stateDiagram-v2", "stateDiagram-v2\n  [*] --> Active", DiagramStateV2, ""},
		{"unknown type", "badtype\n  stuff", DiagramUnknown, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Parse(tt.source, 1)
			if d.Type != tt.wantType {
				t.Errorf("got type %q, want %q", d.Type, tt.wantType)
			}
			if d.Direction != tt.wantDir {
				t.Errorf("got direction %q, want %q", d.Direction, tt.wantDir)
			}
		})
	}
}

func TestParse_SkipsComments(t *testing.T) {
	source := "%% This is a comment\nflowchart LR\n  A --> B"
	d := Parse(source, 1)
	if d.Type != DiagramFlowchart {
		t.Errorf("got type %q, want %q", d.Type, DiagramFlowchart)
	}
}

func TestParse_FlowchartNodes(t *testing.T) {
	source := "flowchart LR\n  A[Start] --> B[End]"
	d := Parse(source, 1)

	if len(d.Nodes) < 2 {
		t.Fatalf("expected at least 2 nodes, got %d", len(d.Nodes))
	}

	nodeByID := make(map[string]Node)
	for _, n := range d.Nodes {
		nodeByID[n.ID] = n
	}

	if n, ok := nodeByID["A"]; !ok {
		t.Error("node A not found")
	} else if n.Label != "Start" {
		t.Errorf("node A label = %q, want %q", n.Label, "Start")
	}

	if n, ok := nodeByID["B"]; !ok {
		t.Error("node B not found")
	} else if n.Label != "End" {
		t.Errorf("node B label = %q, want %q", n.Label, "End")
	}
}

func TestParse_FlowchartEdges(t *testing.T) {
	source := "flowchart LR\n  A --> B\n  B --- C\n  C -.-> D"
	d := Parse(source, 1)

	if len(d.Edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(d.Edges))
	}

	expected := []struct {
		from, to, style string
	}{
		{"A", "B", "-->"},
		{"B", "C", "---"},
		{"C", "D", "-.->"},
	}

	for i, e := range expected {
		if d.Edges[i].From != e.from || d.Edges[i].To != e.to {
			t.Errorf("edge %d: got %s->%s, want %s->%s",
				i, d.Edges[i].From, d.Edges[i].To, e.from, e.to)
		}
		if d.Edges[i].Style != e.style {
			t.Errorf("edge %d: got style %q, want %q", i, d.Edges[i].Style, e.style)
		}
	}
}

func TestParse_StartLine(t *testing.T) {
	source := "flowchart LR\n  A --> B"
	d := Parse(source, 10)
	if d.StartLine != 10 {
		t.Errorf("StartLine = %d, want 10", d.StartLine)
	}
	// Node line should be offset
	for _, n := range d.Nodes {
		if n.Line < 10 {
			t.Errorf("node line %d is less than start line 10", n.Line)
		}
	}
}

func TestParse_EmptySource(t *testing.T) {
	d := Parse("", 1)
	if d.TypeRaw != "" {
		t.Errorf("expected empty TypeRaw for empty source, got %q", d.TypeRaw)
	}
}

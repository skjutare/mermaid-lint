package parser

import (
	"strings"
	"testing"
)

func TestExtractMermaidBlocks_Single(t *testing.T) {
	md := "# Title\n\nSome text\n\n```mermaid\nflowchart LR\n  A --> B\n```\n\nMore text\n"
	blocks, err := ExtractMermaidBlocks(strings.NewReader(md))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	if blocks[0].StartLine != 5 {
		t.Errorf("StartLine = %d, want 5", blocks[0].StartLine)
	}
	if blocks[0].EndLine != 8 {
		t.Errorf("EndLine = %d, want 8", blocks[0].EndLine)
	}
	if !strings.Contains(blocks[0].Source, "flowchart LR") {
		t.Errorf("source should contain diagram, got %q", blocks[0].Source)
	}
}

func TestExtractMermaidBlocks_Multiple(t *testing.T) {
	md := "```mermaid\ngraph TD\n  A --> B\n```\n\nText\n\n```mermaid\nsequenceDiagram\n  Alice->>Bob: Hi\n```\n"
	blocks, err := ExtractMermaidBlocks(strings.NewReader(md))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
}

func TestExtractMermaidBlocks_IgnoresNonMermaid(t *testing.T) {
	md := "```python\nprint('hello')\n```\n\n```javascript\nconsole.log('hi')\n```\n"
	blocks, err := ExtractMermaidBlocks(strings.NewReader(md))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 0 {
		t.Fatalf("expected 0 blocks, got %d", len(blocks))
	}
}

func TestExtractMermaidBlocks_TildeFence(t *testing.T) {
	md := "~~~mermaid\nflowchart LR\n  A --> B\n~~~\n"
	blocks, err := ExtractMermaidBlocks(strings.NewReader(md))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
}

func TestExtractMermaidBlocks_CaseInsensitive(t *testing.T) {
	md := "```Mermaid\nflowchart LR\n  A --> B\n```\n"
	blocks, err := ExtractMermaidBlocks(strings.NewReader(md))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
}

func TestExtractMermaidBlocks_Empty(t *testing.T) {
	blocks, err := ExtractMermaidBlocks(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 0 {
		t.Fatalf("expected 0 blocks, got %d", len(blocks))
	}
}

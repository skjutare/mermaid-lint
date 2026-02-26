package parser

import (
	"bufio"
	"io"
	"strings"
)

// MermaidBlock represents a mermaid code block found in a markdown file.
type MermaidBlock struct {
	Source    string // The mermaid source code (without fences)
	StartLine int   // 1-based line number of the opening fence
	EndLine   int   // 1-based line number of the closing fence
}

// ExtractMermaidBlocks extracts all mermaid code blocks from a markdown reader.
// It looks for fenced code blocks with the "mermaid" language identifier:
//
//	```mermaid
//	...
//	```
func ExtractMermaidBlocks(r io.Reader) ([]MermaidBlock, error) {
	scanner := bufio.NewScanner(r)
	var blocks []MermaidBlock
	var current *MermaidBlock
	var lines []string
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if current == nil {
			// Look for opening fence
			if isMermaidFenceOpen(trimmed) {
				current = &MermaidBlock{StartLine: lineNum}
				lines = nil
			}
		} else {
			// Inside a mermaid block, look for closing fence
			if isClosingFence(trimmed) {
				current.EndLine = lineNum
				current.Source = strings.Join(lines, "\n")
				blocks = append(blocks, *current)
				current = nil
				lines = nil
			} else {
				lines = append(lines, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func isMermaidFenceOpen(line string) bool {
	// Match ```mermaid or ~~~mermaid (with optional whitespace)
	if strings.HasPrefix(line, "```") {
		lang := strings.TrimSpace(strings.TrimPrefix(line, "```"))
		return strings.EqualFold(lang, "mermaid")
	}
	if strings.HasPrefix(line, "~~~") {
		lang := strings.TrimSpace(strings.TrimPrefix(line, "~~~"))
		return strings.EqualFold(lang, "mermaid")
	}
	return false
}

func isClosingFence(line string) bool {
	return line == "```" || line == "~~~"
}

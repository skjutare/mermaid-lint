# mermaid-lint

A lint tool for [Mermaid](https://mermaid.js.org/) diagram files (`.mmd`) and Mermaid code blocks embedded in Markdown files (`.md`).

## Installation

### Homebrew

```bash
brew install skjutare/tap/mermaid-lint
```

### Go

```bash
go install github.com/skjutare/mermaid-lint/cmd/mermaid-lint@latest
```

### From source

```bash
git clone https://github.com/skjutare/mermaid-lint.git
cd mermaid-lint
go build -o mermaid-lint ./cmd/mermaid-lint
```

### Binary releases

Download pre-built binaries for Linux, macOS, and Windows from the [releases page](https://github.com/skjutare/mermaid-lint/releases).

## Usage

```bash
# Lint a single file
mermaid-lint diagram.mmd

# Lint Mermaid blocks inside Markdown
mermaid-lint README.md

# Lint all supported files in a directory (recursive)
mermaid-lint docs/

# Multiple files and directories
mermaid-lint diagram.mmd docs/ notes.md

# Only show warnings and errors
mermaid-lint --severity warning .

# JSON output
mermaid-lint --format json diagrams/

# List available rules
mermaid-lint --list-rules
```

## Supported file types

| Extension              | Description                              |
|------------------------|------------------------------------------|
| `.mmd`, `.mermaid`     | Standalone Mermaid diagram files         |
| `.md`, `.markdown`     | Markdown files with `` ```mermaid `` blocks |

## Rules

| Rule                       | Default Severity | Description                                          |
|----------------------------|------------------|------------------------------------------------------|
| `no-unknown-diagram-type`  | error            | Diagram type must be a recognized Mermaid type       |
| `no-empty-diagram`         | warning          | Diagram must contain at least one element            |
| `valid-direction`          | error            | Flowchart direction must be TB, TD, BT, LR, or RL   |
| `no-duplicate-node-ids`    | warning          | Node IDs must be unique within a diagram             |
| `node-has-label`           | info             | Nodes should have descriptive labels                 |
| `no-orphan-nodes`          | warning          | All nodes should be connected to at least one edge   |

## Configuration

Create a `.mermaid-lint.json` file in your project root to customize rules:

```json
{
  "rules": {
    "no-unknown-diagram-type": { "enabled": true, "severity": "error" },
    "no-empty-diagram": { "enabled": true, "severity": "warning" },
    "valid-direction": { "enabled": true, "severity": "error" },
    "no-duplicate-node-ids": { "enabled": true, "severity": "warning" },
    "node-has-label": { "enabled": false, "severity": "info" },
    "no-orphan-nodes": { "enabled": true, "severity": "warning" }
  }
}
```

Use `--config` to specify a custom config path:

```bash
mermaid-lint --config custom-config.json docs/
```

## Exit codes

| Code | Meaning                          |
|------|----------------------------------|
| 0    | No errors found                  |
| 1    | One or more errors found         |

Warnings and info findings are reported but do not cause a non-zero exit code.

## Supported diagram types

flowchart, graph, sequenceDiagram, classDiagram, stateDiagram, stateDiagram-v2, erDiagram, gantt, pie, gitGraph, mindmap, timeline, quadrantChart, requirementDiagram, C4Context, C4Container, C4Component, C4Dynamic, C4Deployment, sankey-beta, block-beta, xychart-beta

## License

[MIT](LICENSE)

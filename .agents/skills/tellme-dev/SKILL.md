---
name: tellme-dev
description: >-
  Use this skill when you need to build, test, modify, or add new tools/stacks to the tellme codebase.
---

# tellme Development Runbook

Follow these workflows when contributing to or testing the `tellme` project.

## Development Workflows

### 1. Running and Testing Locally
- **Run query directly**: `go run . "do I have python?"`
- **Run the test suite**: `go test ./...`

### 2. Registering a New Tool
To add support for a new command-line tool or NPM package:
1. Open [`registry_data.go`](../../../registry_data.go).
2. Insert your tool in the `standardTools` list.
3. Configure its `ToolConfig` mapping:
   ```go
   ToolRegistry["mytool"] = ToolConfig{
       Names:       []string{"mytool-binary"},
       VersionArgs: []string{"--version"},
       Author:      "Publisher Name",
       Description: "A short, 10-word description of what it does",
       Example:     "`mytool-binary execute`",
   }
   ```
4. If there is an alias (e.g. `python3` -> `python`), map it at the bottom:
   ```go
   ToolRegistry["mytool-alias"] = ToolRegistry["mytool"]
   ```

### 3. Modifying Stacks
To add or update environment stack profiles:
1. Open [`stacks.go`](../../../stacks.go).
2. Edit `StackRegistry` to define the stack, its description, and its list of constituent tool keys.

### 4. Clearing the Lookup Cache
To force a fresh tool scan and bypass cached results:
- **macOS**: `rm -rf ~/Library/Caches/tellme/versions.json`
- **Linux**: `rm -rf ~/.cache/tellme/versions.json`

### 5. Releasing a New Version
To trigger the automated release pipeline (runs tests, compiles, and publishes homebrew formula):
1. Create a version tag matching `v*` (e.g., `v1.1.0`):
   ```bash
   git tag v1.1.0
   ```
2. Push the tag to remote:
   ```bash
   git push origin v1.1.0
   ```


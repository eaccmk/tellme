# Contributing to `tellme`

Thank you for your interest in improving `tellme`! We welcome contributions from the community to help make this local CLI assistant even better and more comprehensive.

As a contributor, please follow these guidelines to ensure smooth reviews and integrations.

---

## Strict Constraints (Non-Negotiable)

1. **100% Offline (Air-gapped):** The tool must NEVER interact with the network or make HTTP requests. All checks and resolution logic must run strictly locally.
2. **Zero Sudo/Zero Elevation:** The tool must NEVER require, ask for, or execute using elevated permissions (`sudo`, admin access). If a folder or binary requires root to check, let it fail safely and return the standard security warning:
   `"Sorry, I never access things that need elevated access such as sudo or admin."`

---

## Local Setup

Ensure you have [Go](https://go.dev/) (v1.26 or later) installed on your system.

1. **Clone the repository:**
   ```bash
   git clone https://github.com/eaccmk/tellme.git
   cd tellme
   ```

2. **Run the CLI tool locally:**
   ```bash
   go run . hi
   ```

3. **Compile the binary:**
   ```bash
   go build -o tellme .
   ```

---

## Running Tests

Before committing any changes, run the full test suite locally to verify code correctness and ensure there are no regressions:

```bash
go test -v ./...
```

---

## How to Add a New Tool or Library

All supported packages and tools are registered in the in-memory lookup map defined in [`registry_data.go`](file:///Users/demo/code/tellme/registry_data.go).

### Step-by-Step Instructions:

1. Open [`registry_data.go`](file:///Users/demo/code/tellme/registry_data.go) and locate the `init()` function.
2. Decide whether the new tool is a standard system binary or a library (such as an NPM package):
   - **For Standard System Tools:** Add the package name to the `standardTools` slice. This automatically creates a default config checking `tool_name --version` in PATH.
   - **For Global NPM Packages:** Add the package name to the `npmPackages` slice. This registers the package under the NPM package manager mapping.
3. If the tool has special configuration needs (e.g. requires checking specific directories or running custom arguments like `-version`), add a custom override block at the bottom of `init()`:

```go
	ToolRegistry["my-special-tool"] = ToolConfig{
		Names:          []string{"my-tool-binary"},
		VersionArgs:     []string{"-version"},
		CommonPaths:    []string{"/opt/special-path/bin/my-tool-binary"},
		PackageManager: "", // Leave blank for standard binary, or set to "npm"
	}
```

4. Run `go test -v ./...` to verify your new configs compile and do not break existing lookup assumptions.

---

## Pull Request Process

1. **Branch Naming:** Create a clean feature branch from `main`:
   ```bash
   git checkout -b feature/add-support-for-mytool
   ```
2. **Commit Messages:** Write clear, concise, and descriptive commit messages.
3. **Tests:** Ensure all unit tests pass before raising a PR.
4. **Target:** Raise your Pull Request targeting the `main` branch.
5. **CI Enforcement:** Every PR triggers a GitHub Actions pipeline (`pr-checks.yml`) that compiles the app and runs the test suite. If any check fails, the PR will be blocked from merging.

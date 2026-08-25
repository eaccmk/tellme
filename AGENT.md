# tellme Workspace Rules

Welcome! These are the guidelines and constraints for developers and AI agents working on the `tellme` project.

## Core Tenets

1. **100% Offline & Air-Gapped**
   - Never write code that makes network calls, HTTP requests, or DNS lookups.
   - All tool detection must happen locally on the system.

2. **Zero Elevation (Safe Execution)**
   - Never use `sudo`, request administrator privileges, or prompt for credentials.
   - Failsafe clean output if access permissions are restricted.

3. **Zero Dependencies**
   - Keep the project lightweight and standard.
   - Avoid adding third-party Go packages (`go.mod` should remain minimal/standard library only) unless absolutely necessary and approved.

4. **Performance & Concurrency**
   - Execute all tool inspections concurrently using goroutines.
   - Keep timeouts strict (maximum 1000ms for command execution) to maintain instantaneous lookups.

5. **Cross-Platform Safety**
   - Ensure detection logic handles macOS/Unix systems and fails gracefully on unsupported platforms.

6. **No Absolute/Local Paths**
   - Never use absolute file URLs (`file:///Users/...` or similar local path schemas) in documentation, comments, or configurations.
   - Always use relative markdown links to refer to other files in the project.


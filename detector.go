package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Installation represents a detected installation of a tool.
type Installation struct {
	Path    string
	Version string
	ModTime time.Time
}

// ToolConfig defines how to search for and check the version of a specific tool.
type ToolConfig struct {
	Names          []string // Executable names to look for (e.g., ["python", "python3"])
	VersionArgs    []string // Arguments to get version (e.g., ["--version"])
	CommonPaths    []string // Additional paths to check outside of standard PATH
	PackageManager string   // Package manager name (e.g. "npm")
	Author         string   // Who built/maintains it
	Description    string   // What it does in < 10 words
	Example        string   // How to use it - a common command
}

// DetectTools scans the system for a given tool name and returns all installations.
func DetectTools(toolName string) ([]Installation, error) {
	lowerName := strings.ToLower(toolName)
	config, ok := ToolRegistry[lowerName]
	if !ok {
		return nil, fmt.Errorf("unsupported tool: %s", toolName)
	}

	// Try loading from cache first
	if cached, found := LoadCache(lowerName); found {
		return cached, nil
	}

	var results []Installation
	var err error

	if config.PackageManager == "npm" {
		results, err = detectNPMPackage(toolName)
	} else {
		results, err = detectSystemTool(toolName, config)
	}

	if err != nil {
		return nil, err
	}

	// Save to cache
	_ = SaveCache(lowerName, results)

	return results, nil
}

func detectSystemTool(toolName string, config ToolConfig) ([]Installation, error) {
	// Gather all search directories
	dirs := []string{}
	pathEnv := os.Getenv("PATH")
	if pathEnv != "" {
		dirs = filepath.SplitList(pathEnv)
	}

	// Collect candidate paths
	candidates := make(map[string]bool)

	// 1. Search PATH directories for the executable names
	for _, dir := range dirs {
		for _, name := range config.Names {
			fullPath := filepath.Join(dir, name)
			if isExecutable(fullPath) {
				candidates[fullPath] = true
			}
		}
	}

	// 2. Add and evaluate common installation paths (glob patterns supported)
	for _, pattern := range config.CommonPaths {
		// Expand any tilde ~ or home dir in pattern
		if strings.HasPrefix(pattern, "~") {
			pattern = filepath.Join(os.Getenv("HOME"), pattern[1:])
		}
		matches, err := filepath.Glob(pattern)
		if err == nil {
			for _, match := range matches {
				if isExecutable(match) {
					candidates[match] = true
				}
			}
		}
	}

	var results []Installation

	// 3. Query versions for candidates with a timeout
	for path := range candidates {
		// Resolve symlinks to find real binary path to prevent duplicates
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			realPath = path
		}

		version, err := getVersionWithTimeout(realPath, config.VersionArgs)
		if err != nil {
			// Check if permission denied/elevated check failed
			if isPermissionDenied(err) {
				return nil, fmt.Errorf("Sorry, I never access things that need elevated access such as sudo or admin. (e.g., trying to run a restricted system binary at %s)", realPath)
			}
			// Skip other errors (e.g. broken symlink)
			continue
		}

		var modTime time.Time
		if stat, err := os.Stat(realPath); err == nil {
			modTime = stat.ModTime()
		}
		results = append(results, Installation{
			Path:    realPath,
			Version: version,
			ModTime: modTime,
		})
	}

	return results, nil
}

// Check if a path exists and is a regular executable file
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	// Basic executable permission check (unix-like)
	return info.Mode()&0111 != 0
}

// Get the version of the tool running its version command with a strict timeout.
func getVersionWithTimeout(binaryPath string, args []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, args...)

	// Suppress standard input, capture stdout & stderr (some tools write version to stderr)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// If context timed out
		if ctx.Err() == context.DeadlineExceeded {
			return "", context.DeadlineExceeded
		}
		return "", err
	}

	// Combine outputs and extract the first non-empty line
	output := stdout.String() + "\n" + stderr.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed, nil
		}
	}

	return "unknown version", nil
}

// Helper to determine if an error was due to insufficient permissions.
func isPermissionDenied(err error) bool {
	if err == nil {
		return false
	}
	// Check standard os error
	if os.IsPermission(err) {
		return true
	}
	// Check syscall error code (e.g., EACCES, EPERM)
	if pathErr, ok := err.(*os.PathError); ok {
		if errno, ok := pathErr.Err.(syscall.Errno); ok {
			return errno == syscall.EACCES || errno == syscall.EPERM
		}
	}
	return false
}

// detectNPMPackage scans standard global node_modules and falls back to running `npm ls -g` with 500ms timeout
func detectNPMPackage(packageName string) ([]Installation, error) {
	homeDir := os.Getenv("HOME")
	searchPaths := []string{
		"/usr/local/lib/node_modules/" + packageName + "/package.json",
		"/opt/homebrew/lib/node_modules/" + packageName + "/package.json",
		homeDir + "/.nvm/versions/node/*/lib/node_modules/" + packageName + "/package.json",
		homeDir + "/.npm-global/lib/node_modules/" + packageName + "/package.json",
		homeDir + "/.config/yarn/global/node_modules/" + packageName + "/package.json",
	}

	var results []Installation

	for _, pattern := range searchPaths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			if isPermissionDenied(err) {
				return nil, fmt.Errorf("Sorry, I never access things that need elevated access such as sudo or admin. (e.g., trying to read restricted global node_modules directory at %s)", pattern)
			}
			continue
		}

		for _, match := range matches {
			version, err := readPackageJSONVersion(match)
			if err != nil {
				if isPermissionDenied(err) {
					return nil, fmt.Errorf("Sorry, I never access things that need elevated access such as sudo or admin. (e.g., trying to read restricted package.json at %s)", match)
				}
				continue
			}
			var modTime time.Time
			if stat, err := os.Stat(match); err == nil {
				modTime = stat.ModTime()
			}
			results = append(results, Installation{
				Path:    filepath.Dir(match),
				Version: version,
				ModTime: modTime,
			})
		}
	}

	// Fallback to running npm CLI query
	if len(results) == 0 {
		version, err := getNPMPackageVersionViaCLI(packageName)
		if err != nil {
			if isPermissionDenied(err) {
				return nil, fmt.Errorf("Sorry, I never access things that need elevated access such as sudo or admin. (e.g., executing npm ls requiring elevated rights)")
			}
			if err == context.DeadlineExceeded {
				return nil, fmt.Errorf("timeout: npm version check took too long for package %s", packageName)
			}
			return nil, nil
		}
		if version != "" {
			results = append(results, Installation{
				Path:    "global node_modules",
				Version: version,
				ModTime: time.Now(),
			})
		}
	}

	return results, nil
}

func readPackageJSONVersion(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(file).Decode(&pkg); err != nil {
		return "", err
	}
	return pkg.Version, nil
}

func getNPMPackageVersionViaCLI(packageName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, npmPath, "ls", "-g", packageName, "--depth=0", "--json")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", context.DeadlineExceeded
		}
	}

	var output struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}

	if err := json.Unmarshal([]byte(stdout.String()), &output); err == nil {
		if pkg, ok := output.Dependencies[packageName]; ok {
			return pkg.Version, nil
		}
	}

	return "", nil
}

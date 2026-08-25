package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// detectUnregisteredBrewTool scans standard Homebrew directory installations for macOS and Linux.
// Returns a ToolConfig if the tool is installed via Homebrew.
func detectUnregisteredBrewTool(packageName string) (ToolConfig, bool) {
	lowerName := strings.ToLower(packageName)

	// Define all common Homebrew Cellar and Caskroom directories
	homeDir, _ := os.UserHomeDir()
	brewRoots := []string{
		"/opt/homebrew",              // macOS ARM
		"/usr/local",                 // macOS Intel
		"/home/linuxbrew/.linuxbrew", // Linux default
	}
	if homeDir != "" {
		brewRoots = append(brewRoots, filepath.Join(homeDir, ".linuxbrew")) // Linux home directory brew
	}

	foundPath := ""
	isCask := false

	// Search Cellar (Formulae) or Caskroom (Casks) using read-only filesystem stats (no sudo)
	for _, root := range brewRoots {
		// Check Cellar
		cellarPath := filepath.Join(root, "Cellar", lowerName)
		if info, err := os.Stat(cellarPath); err == nil && info.IsDir() {
			foundPath = cellarPath
			break
		}
		// Check Caskroom
		caskroomPath := filepath.Join(root, "Caskroom", lowerName)
		if info, err := os.Stat(caskroomPath); err == nil && info.IsDir() {
			foundPath = caskroomPath
			isCask = true
			break
		}
	}

	if foundPath == "" {
		return ToolConfig{}, false
	}

	// Try to locate the executable binary.
	// 1. Check if the binary is already in the system PATH
	binPath, err := exec.LookPath(lowerName)
	if err != nil {
		// 2. Fall back to standard Homebrew bin folders if not found in active PATH
		for _, root := range brewRoots {
			testBin := filepath.Join(root, "bin", lowerName)
			if info, err := os.Stat(testBin); err == nil && !info.IsDir() {
				binPath = testBin
				break
			}
		}
	}

	// If we still can't locate any binary, we cannot execute a version check
	if binPath == "" {
		return ToolConfig{}, false
	}

	pkgType := "Homebrew formula"
	if isCask {
		pkgType = "Homebrew cask"
	}

	// Build a dynamic configuration on the fly
	config := ToolConfig{
		Names:       []string{lowerName},
		VersionArgs: []string{"--version"},
		Author:      "Homebrew Package Manager",
		Description: "Dynamically resolved local " + pkgType,
		Example:     "`" + lowerName + " --help`",
	}

	return config, true
}

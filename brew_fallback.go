package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// detectUnregisteredBrewTool scans standard Homebrew directory installations for macOS and Linux.
// Returns the ToolConfig, resolved installations, and true if found.
func detectUnregisteredBrewTool(packageName string) (ToolConfig, []Installation, bool) {
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
		return ToolConfig{}, nil, false
	}

	// Try to locate the executable binary.
	binPath := ""
	if bin, err := exec.LookPath(lowerName); err == nil {
		binPath = bin
	} else {
		// Fall back to standard Homebrew bin folders
		for _, root := range brewRoots {
			testBin := filepath.Join(root, "bin", lowerName)
			if info, err := os.Stat(testBin); err == nil && !info.IsDir() {
				binPath = testBin
				break
			}
		}
	}

	// Scan the folder for installed version subdirectories (e.g. 1.0.0, 3.1.2)
	files, err := os.ReadDir(foundPath)
	if err != nil {
		return ToolConfig{}, nil, false
	}

	var installations []Installation
	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			versionName := f.Name()
			versionPath := filepath.Join(foundPath, versionName)

			var modTime time.Time
			if stat, err := os.Stat(versionPath); err == nil {
				modTime = stat.ModTime()
			} else {
				modTime = time.Now()
			}

			// If we have a binary, use it; otherwise, point directly to the local Cellar path
			displayPath := binPath
			if displayPath == "" {
				displayPath = versionPath
			}

			installations = append(installations, Installation{
				Path:    displayPath,
				Version: versionName,
				ModTime: modTime,
			})
		}
	}

	if len(installations) == 0 {
		return ToolConfig{}, nil, false
	}

	pkgType := "Homebrew formula"
	if isCask {
		pkgType = "Homebrew cask"
	}

	exampleCmd := "`" + lowerName + " --help`"
	if binPath == "" {
		exampleCmd = "Library package (no binary)"
	}

	// Build a dynamic configuration on the fly
	config := ToolConfig{
		Names:          []string{lowerName},
		VersionArgs:    []string{"--version"},
		PackageManager: "brew_fallback",
		Author:         "Homebrew Package Manager",
		Description:    "Dynamically resolved local " + pkgType,
		Example:        exampleCmd,
	}

	return config, installations, true
}

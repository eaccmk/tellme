package main

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// TestIsPermissionDenied verifies that permission errors are correctly classified.
func TestIsPermissionDenied(t *testing.T) {
	if isPermissionDenied(nil) {
		t.Error("Expected false for nil error")
	}

	// Mocking a standard Permission error
	permErr := fs.ErrPermission
	if !isPermissionDenied(permErr) {
		t.Error("Expected true for fs.ErrPermission")
	}

	// Mocking path-specific permission error
	pathErr := &os.PathError{
		Op:   "open",
		Path: "/restricted",
		Err:  syscall.EACCES,
	}
	if !isPermissionDenied(pathErr) {
		t.Error("Expected true for PathError with EACCES")
	}

	// Other random errors should return false
	randomErr := errors.New("something went wrong")
	if isPermissionDenied(randomErr) {
		t.Error("Expected false for standard error")
	}
}

// TestGetVersionWithTimeout verifies the execution and timeout limits.
func TestGetVersionWithTimeout(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tellme-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock executable script that prints version info
	mockScriptPath := filepath.Join(tempDir, "mock-tool")
	scriptContent := "#!/bin/sh\necho \"mock version 1.2.3\"\n"
	err = os.WriteFile(mockScriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write mock script: %v", err)
	}

	// Test successful version retrieval
	version, err := getVersionWithTimeout(mockScriptPath, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if version != "mock version 1.2.3" {
		t.Errorf("Expected version 'mock version 1.2.3', got %q", version)
	}

	// Create a mock script that hangs / sleeps
	hangingScriptPath := filepath.Join(tempDir, "hanging-tool")
	hangingContent := "#!/bin/sh\nexec sleep 2\necho \"done\"\n"
	err = os.WriteFile(hangingScriptPath, []byte(hangingContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write hanging script: %v", err)
	}

	// Test that it times out (timeout is 1s, sleep is 2s)
	startTime := time.Now()
	_, err = getVersionWithTimeout(hangingScriptPath, nil)
	duration := time.Since(startTime)

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
	}
	if duration > 1500*time.Millisecond {
		t.Errorf("Expected timeout to occur around 1s, but took %v", duration)
	}
}

// TestToolRegistry verifies that the registry contains the requested tools and aliases work.
func TestToolRegistry(t *testing.T) {
	requiredTools := []string{"git", "python", "python3", "go", "node", "awscli", "kubernetes-cli"}
	for _, tool := range requiredTools {
		if _, ok := ToolRegistry[tool]; !ok {
			t.Errorf("Expected tool %q to be present in ToolRegistry", tool)
		}
	}
}

// TestNPMPackageTimeout simulates a slow npm executable to verify the 500ms timeout fails safely.
func TestNPMPackageTimeout(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tellme-npm-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write mock npm that sleeps for 2 seconds
	mockNPMPath := filepath.Join(tempDir, "npm")
	scriptContent := "#!/bin/sh\nexec sleep 2\n"
	err = os.WriteFile(mockNPMPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write mock npm: %v", err)
	}

	// Prepend tempDir to PATH so mock npm is resolved first, but standard shell tools can still be found
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", tempDir+string(os.PathListSeparator)+origPath)
	defer os.Setenv("PATH", origPath)

	startTime := time.Now()
	_, err = getNPMPackageVersionViaCLI("lodash")
	duration := time.Since(startTime)

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
	}
	if duration > 1500*time.Millisecond {
		t.Errorf("Expected timeout to occur around 1s, but took %v", duration)
	}
}

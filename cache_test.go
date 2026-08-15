package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestCacheSaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tellme-cache-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set cacheFile to a temporary file
	originalCacheFile := cacheFile
	cacheFile = filepath.Join(tempDir, "versions.json")
	defer func() { cacheFile = originalCacheFile }()

	installs := []Installation{
		{Path: "/usr/bin/mock", Version: "1.0.0"},
	}

	// 1. Assert Load on missing cache returns false
	_, found := LoadCache("mocktool")
	if found {
		t.Error("Expected found = false on empty cache")
	}

	// 2. Save cache
	err = SaveCache("mocktool", installs)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	// 3. Load cache and verify
	loadedInstalls, found := LoadCache("mocktool")
	if !found {
		t.Error("Expected found = true after saving cache")
	}
	if len(loadedInstalls) != 1 || loadedInstalls[0].Path != "/usr/bin/mock" {
		t.Errorf("Loaded cache mismatch: %v", loadedInstalls)
	}

	// 4. Verify TTL Expiration
	// Modify the cache file directly to make the entry timestamp older than 24 hours
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("Failed to read cache file: %v", err)
	}
	var registryCache map[string]ToolCache
	if err := json.Unmarshal(data, &registryCache); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Set timestamp to 25 hours ago
	entry := registryCache["mocktool"]
	entry.Timestamp = time.Now().Add(-25 * time.Hour)
	registryCache["mocktool"] = entry

	newData, _ := json.Marshal(registryCache)
	_ = os.WriteFile(cacheFile, newData, 0644)

	// Try loading again, should fail TTL check
	_, found = LoadCache("mocktool")
	if found {
		t.Error("Expected found = false for expired cache entry (TTL > 24h)")
	}
}

func TestCacheConcurrency(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tellme-cache-concurrency")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalCacheFile := cacheFile
	cacheFile = filepath.Join(tempDir, "versions.json")
	defer func() { cacheFile = originalCacheFile }()

	var wg sync.WaitGroup
	// Run 20 concurrent readers/writers
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			inst := []Installation{{Path: "/path", Version: "1.0.0"}}
			_ = SaveCache("tool", inst)
		}(i)

		go func() {
			defer wg.Done()
			_, _ = LoadCache("tool")
		}()
	}

	wg.Wait()
}

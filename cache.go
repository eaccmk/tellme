package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ToolCache stores the cached installation details and the fetch timestamp.
type ToolCache struct {
	Installations []Installation `json:"installations"`
	Timestamp     time.Time      `json:"timestamp"`
}

var (
	cacheMutex sync.RWMutex
	cacheFile  string
)

func init() {
	// Initialize cache file path in OS cache directory
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		userCacheDir = os.TempDir()
	}
	cacheFile = filepath.Join(userCacheDir, "tellme", "versions.json")
}

// LoadCache retrieves cached installations for a tool if it is valid (under 24h old).
func LoadCache(toolName string) ([]Installation, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, false
	}

	var registryCache map[string]ToolCache
	if err := json.Unmarshal(data, &registryCache); err != nil {
		return nil, false
	}

	entry, exists := registryCache[toolName]
	if !exists {
		return nil, false
	}

	// 24-hour expiration check
	if time.Since(entry.Timestamp) > 24*time.Hour {
		return nil, false
	}

	return entry.Installations, true
}

// SaveCache updates the cache file with the detected installations of a tool.
func SaveCache(toolName string, installations []Installation) error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Read existing cache to avoid overwriting other tools
	registryCache := make(map[string]ToolCache)
	data, err := os.ReadFile(cacheFile)
	if err == nil {
		_ = json.Unmarshal(data, &registryCache)
	}

	// Update entry
	registryCache[toolName] = ToolCache{
		Installations: installations,
		Timestamp:     time.Now(),
	}

	// Create parent directory if missing
	dir := filepath.Dir(cacheFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	newData, err := json.MarshalIndent(registryCache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, newData, 0644)
}

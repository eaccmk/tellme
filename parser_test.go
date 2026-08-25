package main

import (
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      []string
		isTemporal    bool
		expectedStack string
	}{
		{
			name:          "Simple check",
			input:         "python",
			expected:      []string{"python"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Do I have style query",
			input:         "do I have python ?",
			expected:      []string{"python"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Where is my query",
			input:         "where is my rust",
			expected:      []string{"rust"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Is installed style query",
			input:         "is java installed?",
			expected:      []string{"java"},
			isTemporal:    true, // "installed" flags temporal intent
			expectedStack: "",
		},
		{
			name:          "Have I got style query",
			input:         "have I got ruby",
			expected:      []string{"ruby"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Command prefix in conversational text",
			input:         "tell me if I have go installed on this machine",
			expected:      []string{"go"},
			isTemporal:    true,
			expectedStack: "",
		},
		{
			name:          "Unknown query fallback to last word",
			input:         "is claude installed",
			expected:      []string{"claude"},
			isTemporal:    true,
			expectedStack: "",
		},
		{
			name:          "Aggressive punctuation stripping at the end of word",
			input:         "do I have python? and express; and lodash!",
			expected:      []string{"python", "express", "lodash"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Preserve internal punctuation like node.js",
			input:         "do I have node.js?",
			expected:      []string{"node.js"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Empty input",
			input:         "",
			expected:      nil,
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Multi-tool check with and",
			input:         "do I have python and node?",
			expected:      []string{"python", "node"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Multi-tool check with commas and and",
			input:         "where is go, rust, and ruby?",
			expected:      []string{"go", "rust", "ruby"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Temporal query with when",
			input:         "when was python installed?",
			expected:      []string{"python"},
			isTemporal:    true,
			expectedStack: "",
		},
		{
			name:          "Emoji queries",
			input:         "🐍 and 🦀",
			expected:      []string{"🐍", "🦀"},
			isTemporal:    false,
			expectedStack: "",
		},
		{
			name:          "Valid Stack check (web)",
			input:         "do I have the web stack?",
			expected:      []string{"node", "yarn", "docker", "nginx", "git"},
			isTemporal:    false,
			expectedStack: "web",
		},
		{
			name:          "Stack check with typo (stak)",
			input:         "tellme backend stak",
			expected:      []string{"go", "python", "rust", "docker", "redis", "postgresql"},
			isTemporal:    false,
			expectedStack: "backend",
		},
		{
			name:          "Stack check with typo (stac)",
			input:         "mobile stac",
			expected:      []string{"cocoapods", "ios-deploy", "java", "gradle"},
			isTemporal:    false,
			expectedStack: "mobile",
		},
		{
			name:          "Stack check without stack name",
			input:         "tellme stack",
			expected:      nil,
			isTemporal:    false,
			expectedStack: "unknown",
		},
		{
			name:          "Stack check with unknown stack name",
			input:         "cloud stak",
			expected:      nil,
			isTemporal:    false,
			expectedStack: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, isTemporal, stack := ParseQuery(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseQuery(%q) targets = %q; want %q", tt.input, got, tt.expected)
			}
			if isTemporal != tt.isTemporal {
				t.Errorf("ParseQuery(%q) isTemporal = %v; want %v", tt.input, isTemporal, tt.isTemporal)
			}
			if stack != tt.expectedStack {
				t.Errorf("ParseQuery(%q) stack = %q; want %q", tt.input, stack, tt.expectedStack)
			}
		})
	}
}

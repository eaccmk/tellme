package main

import (
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   []string
		isTemporal bool
	}{
		{
			name:       "Simple check",
			input:      "python",
			expected:   []string{"python"},
			isTemporal: false,
		},
		{
			name:       "Do I have style query",
			input:      "do I have python ?",
			expected:   []string{"python"},
			isTemporal: false,
		},
		{
			name:       "Where is my query",
			input:      "where is my rust",
			expected:   []string{"rust"},
			isTemporal: false,
		},
		{
			name:       "Is installed style query",
			input:      "is java installed?",
			expected:   []string{"java"},
			isTemporal: true, // "installed" flags temporal intent
		},
		{
			name:       "Have I got style query",
			input:      "have I got ruby",
			expected:   []string{"ruby"},
			isTemporal: false,
		},
		{
			name:       "Command prefix in conversational text",
			input:      "tell me if I have go installed on this machine",
			expected:   []string{"go"},
			isTemporal: true,
		},
		{
			name:       "Unknown query fallback to last word",
			input:      "is claude installed",
			expected:   []string{"claude"},
			isTemporal: true,
		},
		{
			name:       "Aggressive punctuation stripping at the end of word",
			input:      "do I have python? and express; and lodash!",
			expected:   []string{"python", "express", "lodash"},
			isTemporal: false,
		},
		{
			name:       "Preserve internal punctuation like node.js",
			input:      "do I have node.js?",
			expected:   []string{"node.js"},
			isTemporal: false,
		},
		{
			name:       "Empty input",
			input:      "",
			expected:   nil,
			isTemporal: false,
		},
		{
			name:       "Multi-tool check with and",
			input:      "do I have python and node?",
			expected:   []string{"python", "node"},
			isTemporal: false,
		},
		{
			name:       "Multi-tool check with commas and and",
			input:      "where is go, rust, and ruby?",
			expected:   []string{"go", "rust", "ruby"},
			isTemporal: false,
		},
		{
			name:       "Temporal query with when",
			input:      "when was python installed?",
			expected:   []string{"python"},
			isTemporal: true,
		},
		{
			name:       "Emoji queries",
			input:      "🐍 and 🦀",
			expected:   []string{"🐍", "🦀"},
			isTemporal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, isTemporal := ParseQuery(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseQuery(%q) targets = %q; want %q", tt.input, got, tt.expected)
			}
			if isTemporal != tt.isTemporal {
				t.Errorf("ParseQuery(%q) isTemporal = %v; want %v", tt.input, isTemporal, tt.isTemporal)
			}
		})
	}
}

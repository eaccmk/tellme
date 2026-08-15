package main

import (
	"strings"
	"testing"
	"time"
)

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		contains string
	}{
		{
			name:     "Just now",
			input:    now.Add(-10 * time.Second),
			contains: "just now",
		},
		{
			name:     "Minutes ago",
			input:    now.Add(-5 * time.Minute),
			contains: "5 minutes ago",
		},
		{
			name:     "Hours ago",
			input:    now.Add(-2 * time.Hour),
			contains: "2 hours ago",
		},
		{
			name:     "Days ago",
			input:    now.Add(-3 * 24 * time.Hour),
			contains: "3 days ago",
		},
		{
			name:     "Months ago",
			input:    now.Add(-4 * 30 * 24 * time.Hour),
			contains: "4 months ago",
		},
		{
			name:     "Years ago",
			input:    now.Add(-2 * 365 * 24 * time.Hour),
			contains: "2 years ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatRelativeTime(tt.input)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("FormatRelativeTime() = %q; expected to contain %q", got, tt.contains)
			}
		})
	}
}

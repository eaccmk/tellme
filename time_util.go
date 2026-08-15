package main

import (
	"fmt"
	"math"
	"time"
)

// FormatRelativeTime formats a time.Time into a relative duration (e.g., "3 months ago (May 14, 2026)").
func FormatRelativeTime(t time.Time) string {
	diff := time.Since(t)
	formattedDate := t.Format("Jan 02, 2006")

	if diff < 0 {
		return fmt.Sprintf("in the future (%s)", formattedDate)
	}

	seconds := diff.Seconds()
	if seconds < 60 {
		return fmt.Sprintf("just now (%s)", formattedDate)
	}

	minutes := diff.Minutes()
	if minutes < 60 {
		m := int(math.Round(minutes))
		if m == 1 {
			return fmt.Sprintf("1 minute ago (%s)", formattedDate)
		}
		return fmt.Sprintf("%d minutes ago (%s)", m, formattedDate)
	}

	hours := diff.Hours()
	if hours < 24 {
		h := int(math.Round(hours))
		if h == 1 {
			return fmt.Sprintf("1 hour ago (%s)", formattedDate)
		}
		return fmt.Sprintf("%d hours ago (%s)", h, formattedDate)
	}

	days := hours / 24
	if days < 30 {
		d := int(math.Round(days))
		if d == 1 {
			return fmt.Sprintf("1 day ago (%s)", formattedDate)
		}
		return fmt.Sprintf("%d days ago (%s)", d, formattedDate)
	}

	months := days / 30.44
	if months < 12 {
		m := int(math.Round(months))
		if m == 1 {
			return fmt.Sprintf("1 month ago (%s)", formattedDate)
		}
		return fmt.Sprintf("%d months ago (%s)", m, formattedDate)
	}

	years := days / 365.2425
	y := int(math.Round(years))
	if y == 1 {
		return fmt.Sprintf("1 year ago (%s)", formattedDate)
	}
	return fmt.Sprintf("%d years ago (%s)", y, formattedDate)
}

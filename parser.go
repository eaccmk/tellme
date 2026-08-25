package main

import (
	"strings"
)

// ParseQuery takes a conversational input string and attempts to extract
// a slice of target tool names (nouns) that the user is asking about,
// returns whether there is a temporal query intent (e.g. "when was X installed"),
// and returns the name of a detected stack profile if it's a stack query.
// It handles conjunctions like "and" or "or", and comma delimiters.
func ParseQuery(query string) ([]string, bool, string) {
	// Check if this is a stack query
	hasStackWord := false
	for _, word := range strings.Fields(query) {
		if isStackWord(word) {
			hasStackWord = true
			break
		}
	}

	if hasStackWord {
		stackName := detectStackInQuery(query)
		if stackName != "" {
			return StackRegistry[stackName].Tools, false, stackName
		}
		return nil, false, "unknown"
	}

	// 1. Normalize the query: convert to lowercase, replace commas with " and ", and split by whitespace
	query = strings.ToLower(query)

	// Check for temporal intent (when, date, updated, installed, modified, age)
	temporalWords := map[string]bool{
		"when":      true,
		"date":      true,
		"updated":   true,
		"installed": true,
		"modified":  true,
		"age":       true,
	}

	isTemporal := false
	rawFields := strings.Fields(query)
	for _, f := range rawFields {
		clean := strings.TrimRight(f, "?.!;,:")
		if temporalWords[clean] {
			isTemporal = true
		}
	}

	query = strings.ReplaceAll(query, ",", " and ")
	rawWords := strings.Fields(query)

	var cleanWords []string
	for _, word := range rawWords {
		// Aggressively strip trailing punctuation from the end of each word
		cleanWord := strings.TrimRight(word, "?.!;,:")
		if cleanWord != "" {
			cleanWords = append(cleanWords, cleanWord)
		}
	}

	// Reconstruct the clean query
	cleanQuery := strings.Join(cleanWords, " ")

	// 2. Split query by conjunctions "and" and "or"
	segments := splitConjunctions(cleanQuery)

	// Include temporal intent words in stopwords so they are not treated as target tools
	stopwords := map[string]bool{
		"do":        true,
		"i":         true,
		"have":      true,
		"where":     true,
		"is":        true,
		"my":        true,
		"installed": true,
		"present":   true,
		"got":       true,
		"any":       true,
		"a":         true,
		"an":        true,
		"the":       true,
		"we":        true,
		"you":       true,
		"if":        true,
		"there":     true,
		"on":        true,
		"this":      true,
		"system":    true,
		"machine":   true,
		"tellme":    true,
		"tell":      true,
		"me":        true,
		"find":      true,
		"locate":    true,
		"check":     true,
		"for":       true,
		// Temporal query noise
		"when":     true,
		"date":     true,
		"updated":  true,
		"modified": true,
		"age":      true,
		"was":      true,
		"did":      true,
		"install":  true,
	}

	var results []string
	seen := make(map[string]bool)

	// 3. For each segment, extract the noun target tool
	for _, segment := range segments {
		words := strings.Fields(segment)
		if len(words) == 0 {
			continue
		}

		var extracted []string
		for _, word := range words {
			if !stopwords[word] {
				extracted = append(extracted, word)
			}
		}

		var target string
		if len(extracted) == 0 {
			target = words[len(words)-1]
		} else {
			target = strings.Join(extracted, " ")
		}

		target = strings.TrimSpace(target)
		if target != "" && !seen[target] {
			seen[target] = true
			results = append(results, target)
		}
	}

	return results, isTemporal, ""
}

// Helper to split a string by "and" or "or" word tokens
func splitConjunctions(s string) []string {
	words := strings.Fields(s)
	var current []string
	var segments []string

	for _, w := range words {
		if w == "and" || w == "or" {
			if len(current) > 0 {
				segments = append(segments, strings.Join(current, " "))
				current = []string{}
			}
		} else {
			current = append(current, w)
		}
	}

	if len(current) > 0 {
		segments = append(segments, strings.Join(current, " "))
	}

	return segments
}

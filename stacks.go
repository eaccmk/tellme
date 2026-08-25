package main

import (
	"strings"
)

// StackProfile defines a collection of tools related to a developer stack.
type StackProfile struct {
	Name        string
	Description string
	Tools       []string
}

// StackRegistry maps normalized stack names to their profiles.
var StackRegistry = map[string]StackProfile{
	"web": {
		Name:        "Web Development Stack",
		Description: "Tools for web application development and hosting",
		Tools:       []string{"node", "yarn", "docker", "nginx", "git"},
	},
	"mobile": {
		Name:        "Mobile Development Stack",
		Description: "Tools for iOS and Android app development",
		Tools:       []string{"cocoapods", "ios-deploy", "java", "gradle"},
	},
	"backend": {
		Name:        "Backend Development Stack",
		Description: "Languages and tools for server-side development",
		Tools:       []string{"go", "python", "rust", "docker", "redis", "postgresql"},
	},
	"db": {
		Name:        "Database Systems Stack",
		Description: "Relational and non-relational database management systems",
		Tools:       []string{"postgresql", "redis", "mysql", "sqlite", "mariadb", "mongodb", "elasticsearch"},
	},
	"devops": {
		Name:        "DevOps & Infrastructure Stack",
		Description: "Infrastructure as code, containerization, and cloud management tools",
		Tools:       []string{"terraform", "ansible", "kubernetes-cli", "docker", "awscli", "azure-cli", "kubernetes-helm"},
	},
	"frontend": {
		Name:        "Frontend Development Stack",
		Description: "Tools for building client-side web interfaces",
		Tools:       []string{"node", "yarn", "typescript", "webpack", "bower"},
	},
}

// isStackWord checks if the given word is "stack" or a likely typo of it.
func isStackWord(word string) bool {
	w := strings.ToLower(strings.TrimRight(word, "?.!;,:"))
	if !strings.HasPrefix(w, "st") {
		return false
	}
	return levenshteinDistance(w, "stack") <= 2 || levenshteinDistance(w, "stacks") <= 2
}

// detectStackInQuery scans the words of the query to identify if the user is asking
// for a specific stack profile. Returns the normalized stack name if found.
func detectStackInQuery(query string) string {
	words := strings.Fields(strings.ToLower(query))
	
	// Map common stack name aliases to their canonical name in StackRegistry
	aliases := map[string]string{
		"web":      "web",
		"mobile":   "mobile",
		"backend":  "backend",
		"db":       "db",
		"database": "db",
		"devops":   "devops",
		"ops":      "devops",
		"frontend": "frontend",
	}

	for _, w := range words {
		clean := strings.TrimRight(w, "?.!;,:")
		if canonical, ok := aliases[clean]; ok {
			return canonical
		}
	}
	return ""
}

// levenshteinDistance computes the distance between two strings to check for typos.
func levenshteinDistance(s, t string) int {
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := 0; i <= len(s); i++ {
		d[i][0] = i
	}
	for j := 0; j <= len(t); j++ {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			substitutionCost := 0
			if s[i-1] != t[j-1] {
				substitutionCost = 1
			}
			d[i][j] = minOfThree(
				d[i-1][j]+1,
				d[i][j-1]+1,
				d[i-1][j-1]+substitutionCost,
			)
		}
	}
	return d[len(s)][len(t)]
}

func minOfThree(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

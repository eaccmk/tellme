package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// Parse command line arguments (excluding the program name itself)
	args := os.Args[1:]

	// 1. Interactive REPL Mode (if zero arguments)
	if len(args) == 0 {
		runREPL()
		return
	}

	// 2. Shell initialization command (the shell fix)
	if len(args) == 1 && strings.ToLower(args[0]) == "init" {
		runShellInit()
		return
	}

	// 3. Greeting check
	if len(args) == 1 {
		cmd := strings.ToLower(args[0])
		if cmd == "hi" || cmd == "hello" {
			printGreeting()
			return
		}
	}

	// 4. Normal query execution
	rawQuery := strings.Join(args, " ")
	runDetectionQuery(rawQuery)
}

func runShellInit() {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		fmt.Println("alias tellme='noglob tellme'")
	}
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("tellme> ")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" || line == "quit" {
			break
		}
		if line != "" {
			runDetectionQuery(line)
		}
		fmt.Print("tellme> ")
	}
}

func runDetectionQuery(rawQuery string) {
	// Cap query length at 256 characters
	if len(rawQuery) > 256 {
		fmt.Println("Error: Query payload exceeds maximum limit of 256 characters.")
		return
	}

	targets, _ := ParseQuery(rawQuery)
	if len(targets) == 0 {
		printGreeting()
		return
	}

	type ToolResult struct {
		Target        string
		ProperName    string
		Installations []Installation
		Err           error
	}

	results := make([]ToolResult, len(targets))
	var wg sync.WaitGroup

	// Hardening checks: Spinner and formatting options
	useSpinner := isTTY() && os.Getenv("NO_COLOR") == ""
	useEmojis := isTTY() && os.Getenv("NO_COLOR") == ""

	// Start a spinner in a background goroutine ONLY if outputting to a TTY and NO_COLOR is not set
	var cancelSpinner context.CancelFunc
	spinnerDone := make(chan struct{})
	close(spinnerDone) // Default close

	if useSpinner {
		var spinnerCtx context.Context
		spinnerCtx, cancelSpinner = context.WithCancel(context.Background())
		spinnerDone = make(chan struct{})
		go func() {
			defer close(spinnerDone)
			frames := []string{".  ", ".. ", "...", ".. ", ".  "}
			i := 0
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-spinnerCtx.Done():
					fmt.Print("\r\033[K") // Clear line
					return
				case <-ticker.C:
					fmt.Printf("\rSearching%s", frames[i])
					i = (i + 1) % len(frames)
				}
			}
		}()
	}

	for idx, target := range targets {
		wg.Add(1)
		go func(i int, t string) {
			defer wg.Done()
			_, supported := ToolRegistry[t]
			if !supported {
				results[i] = ToolResult{
					Target: t,
					Err:    fmt.Errorf("unsupported"),
				}
				return
			}

			properName := capitalize(t)
			installations, err := DetectTools(t)
			results[i] = ToolResult{
				Target:        t,
				ProperName:    properName,
				Installations: installations,
				Err:           err,
			}
		}(idx, target)
	}

	wg.Wait()
	if useSpinner {
		cancelSpinner()
		<-spinnerDone
	}

	// Display results in order
	for _, res := range results {
		if res.Err != nil {
			if res.Err.Error() == "unsupported" {
				fmt.Printf("I don't recognize '%s' on this system yet.\n\n", res.Target)
				if useEmojis {
					fmt.Println("🤝 PRs Welcome!")
				} else {
					fmt.Println("PRs Welcome!")
				}
				fmt.Println("Help us expand the registry. Check out our contributing guide to add it yourself:")
				fmt.Println("https://github.com/eaccmk/tellme/blob/main/CONTRIBUTING.md")
			} else {
				fmt.Println(res.Err.Error())
			}
			continue
		}

		if len(res.Installations) == 0 {
			config := ToolRegistry[res.Target]
			if config.PackageManager == "npm" {
				fmt.Printf("No installations of the library '%s' found on this system.\n", res.Target)
			} else {
				fmt.Printf("No installations of %s found on this system.\n", res.ProperName)
			}
			continue
		}

		config := ToolRegistry[res.Target]

		for _, inst := range res.Installations {
			modTime := inst.ModTime
			if (modTime.IsZero() || modTime.Year() <= 1970) && inst.Path != "global node_modules" {
				if stat, err := os.Stat(inst.Path); err == nil {
					modTime = stat.ModTime()
				}
			}

			ageStr := ""
			showWhen := true
			if modTime.IsZero() || modTime.Year() <= 1970 {
				if strings.HasPrefix(inst.Path, "/usr/bin/") || strings.HasPrefix(inst.Path, "/bin/") || strings.HasPrefix(inst.Path, "/System/") {
					ageStr = "Pre-installed macOS system binary"
				} else {
					showWhen = false
				}
			} else {
				ageStr = "Installed/updated " + FormatRelativeTime(modTime)
			}

			author := config.Author
			if author == "" {
				author = "Unknown Publisher"
			}
			desc := config.Description
			if desc == "" {
				desc = "Developer tool or package"
			}
			example := config.Example
			if example == "" {
				example = "`" + res.Target + " --help`"
			}

			// Format lines (strip emojis if useEmojis is false)
			pkgPrefix := "📦"
			whatPrefix := "📌 What:   "
			wherePrefix := "📍 Where:  "
			whenPrefix := "🕒 When:   "
			howPrefix := "🛠️ How:    "

			if !useEmojis {
				pkgPrefix = ""
				whatPrefix = "What:   "
				wherePrefix = "Where:  "
				whenPrefix = "When:   "
				howPrefix = "How:    "
			}

			sep := "--------------------------------------------------"

			if useColors := useSpinner; useColors {
				fmt.Printf("\033[1;34m%s %s (%s) - %s\033[0m\n", pkgPrefix, res.ProperName, inst.Version, author)
				fmt.Println("\033[34m" + sep + "\033[0m")
				fmt.Printf("\033[1;32m%s\033[0m%s\n", whatPrefix, desc)
				fmt.Printf("\033[1;32m%s\033[0m%s\n", wherePrefix, inst.Path)
				if showWhen {
					fmt.Printf("\033[1;32m%s\033[0m%s\n", whenPrefix, ageStr)
				}
				fmt.Printf("\033[1;32m%s\033[0mTry running: %s\n\n", howPrefix, example)
			} else {
				fmt.Printf("%s %s (%s) - %s\n", pkgPrefix, res.ProperName, inst.Version, author)
				fmt.Println(sep)
				fmt.Printf("%s%s\n", whatPrefix, desc)
				fmt.Printf("%s%s\n", wherePrefix, inst.Path)
				if showWhen {
					fmt.Printf("%s%s\n", whenPrefix, ageStr)
				}
				fmt.Printf("%sTry running: %s\n\n", howPrefix, example)
			}
		}
	}
}

func printGreeting() {
	fmt.Println("Hi, you can search like: do I have python ? or do I have claude installed")
}

func capitalize(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

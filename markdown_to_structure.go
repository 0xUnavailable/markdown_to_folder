package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type InputFormat string

const (
	FormatLayered InputFormat = "layered"
	FormatTree    InputFormat = "tree"
)

// Removes all tree-drawing characters and leading whitespace
var treePrefixRegex = regexp.MustCompile(`^[\s│├└─\-\+\|]+`)

// Extract filename starting at first '_' OR first valid filename character
// Supports:
//   _file.py
//   text.py
//   README.md
var nameExtractionRegex = regexp.MustCompile(`^([_A-Za-z0-9][A-Za-z0-9._-]*)$`)

func main() {
	inputFile := flag.String("input", "structure.md", "Path to the markdown file")
	formatType := flag.String("format", "tree", "Input format: 'layered' or 'tree'")
	flag.Parse()

	format := InputFormat(*formatType)
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if format == FormatLayered {
		parseLayeredFormat(file)
	} else {
		parseTreeFormat(file)
	}

	fmt.Println("\nRepository structure created successfully!")
}

func parseTreeFormat(file *os.File) {
	scanner := bufio.NewScanner(file)
	currentPath := []string{}
	rootDir := ""
	previousDepth := -1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Root directory
		if rootDir == "" {
			rootDir = strings.TrimSpace(strings.TrimSuffix(line, "/"))
			os.MkdirAll(rootDir, 0755)
			fmt.Printf("Created directory: %s\n", rootDir)
			continue
		}

		// Determine depth
		depth := 0
		found := false
		for i, r := range line {
			if !strings.ContainsRune("│├└─-+ |", r) && !unicode.IsSpace(r) {
				depth = i / 4
				found = true
				break
			}
		}
		if !found {
			continue
		}

		// --- CRITICAL FIX ---
		// 1. Strip tree prefix
		stripped := treePrefixRegex.ReplaceAllString(line, "")
		stripped = strings.TrimSpace(stripped)

		// 2. Extract filename
		match := nameExtractionRegex.FindStringSubmatch(stripped)
		if len(match) < 2 {
			continue
		}
		name := match[1]

		// Manage nesting
		if depth <= previousDepth {
			backtrack := previousDepth - depth + 1
			if backtrack > len(currentPath) {
				backtrack = len(currentPath)
			}
			currentPath = currentPath[:len(currentPath)-backtrack]
		}
		previousDepth = depth

		// Build path
		pathParts := append([]string{rootDir}, currentPath...)
		pathParts = append(pathParts, name)
		fullPath := filepath.Join(pathParts...)

		createFileOrDirectory(fullPath, name, &currentPath)
	}
}

func parseLayeredFormat(file *os.File) {
	scanner := bufio.NewScanner(file)
	currentPath := []string{}
	rootDir := ""

	for scanner.Scan() {
		originalLine := scanner.Text()
		line := strings.TrimSpace(originalLine)
		if line == "" {
			continue
		}

		indent := 0
		for _, r := range originalLine {
			if r == ' ' {
				indent++
			} else {
				break
			}
		}

		var name string
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			name = strings.TrimSpace(line[2:])
		} else if indent == 0 {
			name = line
		} else {
			continue
		}

		if rootDir == "" && indent == 0 {
			rootDir = name
			os.MkdirAll(rootDir, 0755)
			fmt.Printf("Created directory: %s\n", rootDir)
			continue
		}

		level := indent / 2
		if level > 0 {
			level--
		}

		if len(currentPath) > level {
			currentPath = currentPath[:level]
		}

		pathParts := append([]string{rootDir}, currentPath...)
		pathParts = append(pathParts, name)
		fullPath := filepath.Join(pathParts...)

		createFileOrDirectory(fullPath, name, &currentPath)
	}
}

func createFileOrDirectory(fullPath, name string, currentPath *[]string) {
	if strings.Contains(name, ".") {
		parent := filepath.Dir(fullPath)
		os.MkdirAll(parent, 0755)
		f, _ := os.Create(fullPath)
		f.Close()
		fmt.Printf("Created file: %s\n", name)
	} else {
		os.MkdirAll(fullPath, 0755)
		fmt.Printf("Created directory: %s\n", name)
		*currentPath = append(*currentPath, name)
	}
}

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
	var rootDir string
	// currentPath will store the directory stack
	currentPath := []string{}
	// track the depth of each level to know when to pop from the stack
	depths := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 1. Identify the root directory (first line)
		if rootDir == "" {
			rootDir = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "./"), "/"))
			os.MkdirAll(rootDir, 0755)
			fmt.Printf("Created root: %s\n", rootDir)
			continue
		}

		// 2. Determine visual depth by finding the first alphanumeric character
		contentIndex := strings.IndexFunc(line, func(r rune) bool {
			return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.'
		})
		
		if contentIndex == -1 {
			continue
		}

		// 3. Clean the name
		name := strings.TrimSpace(line[contentIndex:])
		name = strings.TrimSuffix(name, "/")

		// 4. Backtrack the stack: If current depth is <= previous depth, 
		// pop until we are at the parent level
		for len(depths) > 0 && contentIndex <= depths[len(depths)-1] {
			depths = depths[:len(depths)-1]
			currentPath = currentPath[:len(currentPath)-1]
		}

		// 5. Build full path
		fullPath := filepath.Join(rootDir, filepath.Join(currentPath...), name)

		// 6. Create file or directory
		if strings.Contains(name, ".") {
			// It's a file
			os.MkdirAll(filepath.Dir(fullPath), 0755)
			f, _ := os.Create(fullPath)
			f.Close()
			fmt.Printf("Created file: %s\n", fullPath)
		} else {
			// It's a directory
			os.MkdirAll(fullPath, 0755)
			fmt.Printf("Created folder: %s\n", fullPath)
			// Add to stack for children
			currentPath = append(currentPath, name)
			depths = append(depths, contentIndex)
		}
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

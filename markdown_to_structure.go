package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define command-line flag for markdown file path
	inputFile := flag.String("input", "structure.md", "Path to the markdown file containing the repository structure")
	flag.Parse()

	// Open the markdown file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening markdown file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentPath := make([]string, 0)
	rootDir := ""

	for scanner.Scan() {
		originalLine := scanner.Text()
		line := strings.TrimSpace(originalLine)

		if line == "" {
			continue
		}

		// Calculate indentation level from original line (count leading spaces)
		indentLevel := 0
		for _, char := range originalLine {
			if char == ' ' {
				indentLevel++
			} else {
				break
			}
		}

		// Parse the line (expecting markdown list item starting with - or *)
		var name string
		if strings.HasPrefix(line, "- ") {
			name = strings.TrimSpace(line[2:])
		} else if strings.HasPrefix(line, "* ") {
			name = strings.TrimSpace(line[2:])
		} else if indentLevel == 0 {
			name = line // Root directory
		} else {
			continue // Skip lines that don't match expected format
		}

		// Handle root directory (first item with indent 0)
		if rootDir == "" && indentLevel == 0 {
			rootDir = name
			err = os.MkdirAll(rootDir, 0755)
			if err != nil {
				fmt.Printf("Error creating root directory %s: %v\n", rootDir, err)
				os.Exit(1)
			}
			continue
		}

		// Calculate target level (each 2 spaces = 1 level)
		targetLevel := indentLevel / 2

		// Adjust current path based on indentation level
		// We need to be at targetLevel - 1 (since we're about to add this item)
		if targetLevel > 0 {
			targetLevel-- // Adjust because currentPath represents parent directories
		}

		// Trim currentPath to match the target level
		if len(currentPath) > targetLevel {
			currentPath = currentPath[:targetLevel]
		}

		// Construct the full path
		pathParts := append([]string{rootDir}, currentPath...)
		pathParts = append(pathParts, name)
		fullPath := filepath.Join(pathParts...)

		// Check if it's a file (contains an extension) or directory
		if strings.Contains(name, ".") {
			// Create parent directory if it doesn't exist
			parentDir := filepath.Dir(fullPath)
			err := os.MkdirAll(parentDir, 0755)
			if err != nil {
				fmt.Printf("Error creating parent directory %s: %v\n", parentDir, err)
				continue
			}

			// Create file
			f, err := os.Create(fullPath)
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", fullPath, err)
				continue
			}
			f.Close()
			fmt.Printf("Created file: %s\n", fullPath)
		} else {
			// Create directory
			err := os.MkdirAll(fullPath, 0755)
			if err != nil {
				fmt.Printf("Error creating directory %s: %v\n", fullPath, err)
				continue
			}
			fmt.Printf("Created directory: %s\n", fullPath)
			
			// Add to current path for nested items
			currentPath = append(currentPath, name)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Println("\nRepository structure created successfully!")
}
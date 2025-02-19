package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// returns the root directory of fortune files
func getFortuneDirectory() (string, error) {
	cmd := exec.Command("fortune", "-f")
	pipe, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(pipe)
	if scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "/"); idx >= 0 {
			return line[idx:], nil
		}
	}
	return "", fmt.Errorf("fortune directory not found")
}

// returns all valid fortune files
func findFortuneFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() ||
			filepath.Ext(path) == ".dat" ||
			strings.Contains(path, "/off/") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

// selects a random quote from a random file
func getRandomQuote(files []string, rng *rand.Rand) (string, error) {
	randomFile := files[rng.Intn(len(files))]
	content, err := os.ReadFile(randomFile)
	if err != nil {
		return "", err
	}
	quotes := strings.Split(string(content), "%")
	if len(quotes) > 1 {
		return quotes[rng.Intn(len(quotes)-1)+1], nil
	}
	return "", fmt.Errorf("no quotes found in file")
}

func main() {
	// Initialize random number generator with a source
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Get fortune directory path
	root, err := getFortuneDirectory()
	if err != nil {
		log.Fatal("Failed to get fortune directory:", err)
	}

	// Find all valid fortune files
	files, err := findFortuneFiles(root)
	if err != nil {
		log.Fatal("Failed to find fortune files:", err)
	}

	// Select and print a random quote
	if len(files) > 0 {
		quote, err := getRandomQuote(files, rng)
		if err != nil {
			log.Fatal("Failed to get random quote:", err)
		}
		fmt.Print(quote)
	}
}

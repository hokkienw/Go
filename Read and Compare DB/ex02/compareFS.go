package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	oldFile := flag.String("old", "", "Old filesystem snapshot file")
	newFile := flag.String("new", "", "New filesystem snapshot file")
	flag.Parse()

	if *oldFile == "" || *newFile == "" {
		fmt.Println("Print both old and new files using -old and -new")
		os.Exit(1)
	}

	oldSnapshot, err := readSnapshot(*oldFile)
	if err != nil {
		fmt.Println("Error reading old file:", err)
		os.Exit(1)
	}

	newSnapshot, err := readSnapshot(*newFile)
	if err != nil {
		fmt.Println("Error reading new file:", err)
		os.Exit(1)
	}

	compareSnapshots(oldSnapshot, newSnapshot)
}

func readSnapshot(filename string) (map[string]bool, error) {
	snapshot := make(map[string]bool)
	file, err := os.Open(filename)
	if err != nil {
		return snapshot, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		snapshot[path] = true
	}

	if err := scanner.Err(); err != nil {
		return snapshot, err
	}

	return snapshot, nil
}

func compareSnapshots(oldSnapshot, newSnapshot map[string]bool) {
	for path := range newSnapshot {
		if _, exists := oldSnapshot[path]; !exists {
			fmt.Printf("ADDED %s\n", path)
		}
	}

	for path := range oldSnapshot {
		if _, exists := newSnapshot[path]; !exists {
			fmt.Printf("REMOVED %s\n", path)
		}
	}
}

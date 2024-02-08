package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	showFiles    bool
	showDirs     bool
	showSymlinks bool
)

func init() {
	flag.BoolVar(&showFiles, "f", false, "Print files")
	flag.BoolVar(&showDirs, "d", false, "Print directories")
	flag.BoolVar(&showSymlinks, "sl", false, "Print symlinks")
	flag.Parse()
}

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if (showFiles && !info.IsDir() && !info.Mode().IsRegular()) ||
		(showDirs && info.IsDir()) ||
		(showSymlinks && info.Mode()&os.ModeSymlink != 0) {
		fmt.Println(path)
	}

	return nil
}

func main() {
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Usage: ./myFind [-f] [-d] [-sl] /path/to/dir")
		return
	}

	root := args[0]

	err := filepath.Walk(root, visit)
	if err != nil {
		fmt.Println(err)
	}
}

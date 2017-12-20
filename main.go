package main

import (
	"flag"
	"errors"
	"fmt"
	"path/filepath"
	"os"
)


var keysSlice []string

func parseArgs() (string, error) {
	dir := flag.String("dir", "", "Path to directory that should be scanned")
	flag.Parse()
	if *dir == "" {
		return "", errors.New("no directory specified")
	}
	return *dir, nil
}

func main() {
	directory, err := parseArgs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	err2 := filepath.Walk(directory, visit)

	fmt.Printf("len=%d cap=%d %v\n", len(keysSlice), cap(keysSlice), keysSlice)

	fmt.Printf("filepath.Walk() returned %v\n", err2)
}
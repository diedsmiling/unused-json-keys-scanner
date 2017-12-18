package main

import (
	"flag"
	"os"
	"errors"
	"fmt"
)

func parseArgs() (string, error) {
	dir := flag.String("dir", "", "Path to directory that should be scanned")
	flag.Parse()
	if *dir == "" {
		return "", errors.New("no directory specified")
	}
	return *dir, nil
}

func (d dir) scan() bool {
	fmt.Printf("Scaning directory  %s", d.name)
	return true
}

func main() {
	directory, err := parseArgs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	var rootDir = dir{name: directory}
	rootDir.scan()
}
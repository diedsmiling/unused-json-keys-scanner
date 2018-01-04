package main

import (
	"flag"
	"errors"
	"strings"
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/buger/jsonparser"
	"bytes"
	"log"
)

type Key struct {
	Key string
	File string
	Used bool
}

var keysSlice []Key

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		os.Exit(1)
	}
}

func isDir(f os.FileInfo) (r bool) {
	switch mode := f.Mode(); {
	case mode.IsDir():
		r = true
	case mode.IsRegular():
		r = false
	}
	return

}

func collectJsons(path string, f os.FileInfo, err error) error {
	if isDir(f) {
		return nil
	}
	body, err := ioutil.ReadFile(path)
	failOnError(err, "Could not parse")
	isJ := isJson(body)
	if isJ {
		gatherKeys("", path, body)
	}
	return nil
}

func scan(path string, f os.FileInfo, err error) error {
	body, err := ioutil.ReadFile(path)
	for i, element := range keysSlice {
		if strings.Contains(string(body), element.Key) {
			keysSlice[i].Used = true
		}
	}
	return nil
}


func gatherKeys(parentKey string, path string, json []byte)  {
	jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		var newKey Key
		if dataType == jsonparser.Object {
			if parentKey != "" {
				gatherKeys(concatTwoStrings(parentKey, string(key)), path, value)
			} else {
				gatherKeys(string(key), path, value)
			}
		} else {
			if parentKey != "" {
				newKey = Key{Key: concatTwoStrings(parentKey, string(key)), File: path, Used: false}
			} else {
				newKey = Key{Key: string(key), File: path, Used: false}
			}
			keysSlice = append(keysSlice, newKey)
		}
		return nil
	})
}

func concatTwoStrings(a string, b string) string{
	list := []string{a, ".", b}
	var str bytes.Buffer

	for _, l := range list {
		str.WriteString(l)
	}
	return str.String()
}

func isJson(buf []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(buf, &js) == nil
}

func parseArgs() (string, bool, error) {
	dir := flag.String("dir", "", "Path to directory that should be scanned")
	toDelete := flag.Bool("delete", false, "Delete unused properties")
	flag.Parse()
	if *dir == "" {
		return "", false, errors.New("no directory specified")
	}
	return *dir, *toDelete, nil
}

func Filter(vs []Key, f func(Key) bool) []Key {
	vsf := make([]Key, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func main() {
	directory, toDelete, err := parseArgs()
	failOnError(err, "Could not parse args")

	err2 := filepath.Walk(directory, collectJsons)
	failOnError(err2, "Could not parse args")
	fmt.Printf("len=%d cap=%d %v\n", len(keysSlice), cap(keysSlice), keysSlice)

	err3 := filepath.Walk(directory, scan)
	keysSlice = Filter(keysSlice, func(v Key) bool {
		return !v.Used
	})

    for _, element := range keysSlice {
		fmt.Printf("Used key %v in file: %w\n", element.Key, element.File)
	}

	if toDelete {
		fmt.Printf("Deleting ")
	}

	fmt.Printf("length ", len(keysSlice))

	failOnError(err3, "Could not parse args")
}
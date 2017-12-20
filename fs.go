package main

import (
	"os"
	"fmt"
	"github.com/tidwall/gjson"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

func visit(path string, f os.FileInfo, err error) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not parse: ", err)
	}
	isJ := isJson(bytes)
	if isJ {
		var f interface{}
		err := json.Unmarshal(bytes, &f)
		if err != nil {
			fmt.Printf("er")
		}

		result := gjson.Parse(string(bytes))
		gatherKeys("", result)
	}
	return nil
}

func gatherKeys(parentKey string, json gjson.Result)  {
	json.ForEach(func(key, value gjson.Result) bool {
		if parentKey != "" {
			keysSlice = append(keysSlice, concatTwoStrings(parentKey, key.String()))
		} else {
			keysSlice = append(keysSlice, key.String())
		}

		fmt.Println("parent key", parentKey)
		fmt.Println("key", key)
		if value.Type == gjson.JSON {
			gatherKeys(key.String(), value)
		}
		return true
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


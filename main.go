package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Mapping struct {
	hasher map[string]string
}

func main() {

	// Open our jsonFile
	jsonFile, err := os.Open("properties.json") //
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened properties.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	//read json file
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var mapper Mapping

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(byteValue, &mapper.hasher)

	err = filepath.Walk("./", mapper.visit)
	if err != nil {
		panic(err)
	}
}

func (mapper *Mapping) visit(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !!fi.IsDir() {
		return nil //
	}

	fileExt := []string{"*.html", "*.js", "CNAME"}

	for _, val := range fileExt {
		matched, err := filepath.Match(val, fi.Name())
		if err != nil {
			return err
		}
		if matched {
			read, err := ioutil.ReadFile(path) //read file with ext
			if err != nil {
				panic(err)
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			var newContents string
			newContents = string(read)
			for key, value := range mapper.hasher {
				newContents = strings.Replace(newContents, key, value, -1)
			}

			err = ioutil.WriteFile(path, []byte(newContents), 0)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//LoadDataFile loads a json data file.
func LoadDataFile(fileName string) []interface{} {
	// Load services file
	dataFile, err := ioutil.ReadFile("services.json")
	if err != nil {
		fmt.Println("Can't open services.json")
		os.Exit(1)
	}
	var v []interface{}
	if err := json.Unmarshal(dataFile, &v); err != nil {
		fmt.Println("Can't open config file")
		os.Exit(1)
	}
	return v
}

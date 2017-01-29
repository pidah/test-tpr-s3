package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

//LoadDataFile loads a json data file.
func LoadDataFile(fileName string) []Service {
	// Load services file
	dataFile, err := ioutil.ReadFile("services.json")
	if err != nil {
		fmt.Println("Can't open services.json")
		os.Exit(1)
	}

	services := []Service{}
	if err := json.Unmarshal(dataFile, &services); err != nil {
		fmt.Println("Can't open config file")
		os.Exit(1)
	}
          for _, service := range services {
    //fmt.Printf( "The service '%s' is available at  '%s'\n", k, v[k] );
		fmt.Println(service)
  }
	return services
}

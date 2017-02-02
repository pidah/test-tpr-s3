package main

import (
//	"encoding/json"
        "gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
	"os"
)

type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

//LoadDataFile loads a yaml data file.
func LoadDataFile(fileName string) []Service {
	// Load services file
	dataFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Can't open", fileName)
		os.Exit(1)
	}

	services := []Service{}
	if err := yaml.Unmarshal(dataFile, &services); err != nil {
		fmt.Println("Can't open config file")
		os.Exit(1)
	}
	for _, service := range services {
		//fmt.Printf( "The service '%s' is available at  '%s'\n", k, v[k] );
		fmt.Println(service)
	}
	return services
}

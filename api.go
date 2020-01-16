package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// API holds the list of endpoints from the definition file
type API struct {
	Endpoints []Endpoint `json:"api"`
}

// ReadAPI reads and parses the API defuiniion from a JSON encoded file
func ReadAPI(file string) (*API, error) {

	jsonFile, err := os.Open(file)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var api = new(API)

	json.Unmarshal([]byte(byteValue), api)

	fmt.Println("Found", len(api.Endpoints), "Endpoints")

	return api, nil
}

// Process generates the binding and handler files for the API
func (api *API) Process(bindingDir string, handlerDir string, appName string) error {
	for _, endpoint := range api.Endpoints {
		err := endpoint.Generate(bindingDir, handlerDir, appName)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

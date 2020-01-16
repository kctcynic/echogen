package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {

	var fileNamePtr = flag.String("file", "api.json", "name of the api definition file")
	var appNamePtr = flag.String("name", "app", "name of the app directory")

	flag.Parse()

	apiFileName := *fileNamePtr
	appName := *appNamePtr
	bindingDirName := "bindings"
	handlerDirName := "handlers"

	fmt.Println("Processing", apiFileName)

	//
	// Make sure the bindings directory exists
	//
	err := checkDirectory(bindingDirName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = checkDirectory(handlerDirName)
	if err != nil {
		fmt.Println(err)
		return
	}

	//
	// Make sure the handlers directory exists
	//
	api, err := ReadAPI(apiFileName)
	api.Process(bindingDirName, handlerDirName, appName)
}

// checkDirectory validates that a directory exists or creates it if necessary
func checkDirectory(name string) error {
	mode, err := os.Stat(name)
	if os.IsNotExist(err) {
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			return err
		}
	} else if !mode.IsDir() {
		// It's not a directory, bail...
		return errors.New("cannot create directory")
	}

	return nil
}

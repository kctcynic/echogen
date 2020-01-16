package main

import (
	"bufio"
	"fmt"
	"os"
)

// Endpoint holds the definition of the parameters and results for and API endpoint
type Endpoint struct {
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
	Method  string `json:"method,omitempty"`
	URL     string `json:"url,omitempty"`
	Params  []Attr `json:"params,omitempty"`
	Results []Attr `json:"results,omitempty"`
}

// Attr holds the definition of one of the endpoint arguments
type Attr struct {
	Tag   string `json:"tag"`
	Field string `json:"field"`
	Type  string `json:"type"`
}

// Generate creates the files to bind the api parameters and results for an endpoint
func (endpoint *Endpoint) Generate(bindingDir string, handlerDir string, appName string) error {

	println("Processing", endpoint.Name)

	if len(endpoint.Params) > 0 {
		paramFile := endpoint.Prefix + "_parameters.go"
		println("Creating", paramFile)

		file, err := os.Create(bindingDir + "/" + paramFile)
		if err != nil {
			return err
		}
		defer file.Close()
		endpoint.generateParameters(file)
	}

	if len(endpoint.Results) > 0 {
		resultFile := endpoint.Prefix + "_results.go"
		println("Creating", resultFile)

		file, err := os.Create(bindingDir + "/" + resultFile)
		if err != nil {
			return err
		}
		defer file.Close()
		endpoint.generateResults(file)
	}

	if len(endpoint.URL) > 0 {
		handlerFile := endpoint.Prefix + ".go"
		println("Creating", handlerFile)

		file, err := os.Create(handlerDir + "/" + handlerFile)
		if err != nil {
			return err
		}
		defer file.Close()
		endpoint.generateHandler(file, appName)

		implFile := endpoint.Prefix + "Impl.go"
		_, err = os.Stat(handlerDir + "/" + implFile)
		if !os.IsNotExist(err) {
			return nil
		}

		//
		// Impl file does not exist, create it
		//
		impl, err := os.Create(handlerDir + "/" + implFile)
		if err != nil {
			return err
		}
		defer impl.Close()
		endpoint.generateImpl(impl, appName)

	}

	return nil
}

// generateParameters generates the binding class for the endpoint parameters
func (endpoint *Endpoint) generateParameters(file *os.File) error {

	w := bufio.NewWriter(file)

	fmt.Fprintf(w, "package bindings\n\n")

	fmt.Fprintf(w, "// %sParameters is used for binding to JSON parameters for %s\n", endpoint.Name, endpoint.URL)
	fmt.Fprintf(w, "type %sParameters struct {\n", endpoint.Name)

	for _, param := range endpoint.Params {
		fmt.Fprintf(w, "\t%s %s `json:\"%s\"`\n", param.Field, param.Type, param.Tag)
	}

	fmt.Fprintf(w, "}\n")

	w.Flush()

	return nil
}

// generateResults generates the binding class for the endpoint results
func (endpoint *Endpoint) generateResults(file *os.File) error {
	w := bufio.NewWriter(file)

	fmt.Fprintf(w, "package bindings\n\n")

	if len(endpoint.URL) > 0 {
		fmt.Fprintf(w, "// %sResults is used for binding to JSON results for %s\n", endpoint.Name, endpoint.URL)
	} else {
		fmt.Fprintf(w, "// %sResults is used for binding to JSON results\n", endpoint.Name)
	}

	fmt.Fprintf(w, "type %sResults struct {\n", endpoint.Name)

	for _, result := range endpoint.Results {
		fmt.Fprintf(w, "\t%s %s `json:\"%s\"`\n", result.Field, result.Type, result.Tag)
	}

	fmt.Fprintf(w, "}\n")

	w.Flush()

	return nil
}

// generateHandler generates the handler class for the endpoint
func (endpoint *Endpoint) generateHandler(file *os.File, appName string) error {
	w := bufio.NewWriter(file)

	fmt.Fprintf(w, "package handlers\n\n")

	fmt.Fprint(w, "import (\n\t\"net/http\"\n")
	fmt.Fprintf(w, "\n\t\"%s/bindings\"\n\n", appName)
	fmt.Fprint(w, "\t\"github.com/labstack/echo/v4\"\n)\n\n")

	fmt.Fprintf(w, "// %s is used for handling requests on %s\n", endpoint.Name, endpoint.URL)

	fmt.Fprintf(w, "func %s(c echo.Context) error {\n", endpoint.Name)

	if endpoint.Params != nil {
		fmt.Fprintf(w, "\tvar params = new(bindings.%sParameters)\n", endpoint.Name)
		fmt.Fprint(w, "\tif err := c.Bind(params); err != nil {\n")

		fmt.Fprint(w, "\t\tmsg := new(bindings.ErrorResults)\n")
		fmt.Fprint(w, "\t\tmsg.ErrorMessage = \"Invalid Parameters\"\n")
		fmt.Fprint(w, "\t\treturn c.JSON(http.StatusOK, msg)\n")

		fmt.Fprint(w, "\t}\n")
	} else {
		fmt.Fprintf(w, "\tvar params *bindings.%sParameters = nil\n", endpoint.Name)
	}

	fmt.Fprintf(w, "\tresult, err := %sImpl(c, params)\n", endpoint.Name)

	fmt.Fprint(w, "\tif err != nil {\n")
	fmt.Fprint(w, "\t\tmsg := new(bindings.ErrorResults)\n")
	fmt.Fprint(w, "\t\tmsg.ErrorMessage = err.Error()\n")
	fmt.Fprint(w, "\t\treturn c.JSON(http.StatusOK, msg)\n")

	fmt.Fprint(w, "\t}\n")
	fmt.Fprint(w, "\treturn c.JSON(http.StatusOK, result)\n")

	fmt.Fprintf(w, "}\n")

	w.Flush()

	return nil
}

// generateImpl generates the implementation class for the endpoint if it does not yet exist
func (endpoint *Endpoint) generateImpl(file *os.File, appName string) error {
	w := bufio.NewWriter(file)

	fmt.Fprintf(w, "package handlers\n\n")

	fmt.Fprint(w, "import (\n")
	fmt.Fprintf(w, "\t\"%s/bindings\"\n\n", appName)
	fmt.Fprint(w, "\t\"github.com/labstack/echo/v4\"\n)\n")

	fmt.Fprintf(w, "// %sImpl is used for handling requests on %s\n", endpoint.Name, endpoint.URL)

	fmt.Fprintf(w, "func %sImpl(c echo.Context, params *bindings.%sParameters) (result *bindings.%sResults, err error) {\n",
		endpoint.Name, endpoint.Name, endpoint.Name)

	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "\t//\n")
	fmt.Fprint(w, "\t// IMPLEMENTATION GOES HERE\n")
	fmt.Fprint(w, "\t//\n")
	fmt.Fprint(w, "\n")

	fmt.Fprint(w, "\treturn nil, nil\n")
	fmt.Fprintf(w, "}\n")

	w.Flush()

	return nil
}

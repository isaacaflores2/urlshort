package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/isaacaflores2/urlshort"
)

func main() {
	const defaultYamlFilePath = "defaultPaths.yml"
	const defaultJsonFilePath = "defaultPaths.json"

	yamlFilePathPrt := flag.String("f", defaultYamlFilePath, "Optional flag for yaml file")
	jsonFlagPrt := flag.Bool("useJson", false, "Optional flag to read jsonFile")
	flag.Parse()

	fmt.Println("JsonFlag: ", *jsonFlagPrt)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	redirectFilePath := *yamlFilePathPrt
	if *jsonFlagPrt == true {
		redirectFilePath = defaultJsonFilePath
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	redirectFile, err := ioutil.ReadFile(redirectFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Using redirect file: ", redirectFilePath)
	handler, err := createHandlerForFile([]byte(redirectFile), mapHandler, *jsonFlagPrt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func createHandlerForFile(redirectFile []byte, fallback http.HandlerFunc, isJsonFile bool) (http.HandlerFunc, error) {
	if isJsonFile {
		return urlshort.JSONHandler(redirectFile, fallback)
	}

	return urlshort.YAMLHandler(redirectFile, fallback)
}

package urlshort

import (
	"fmt"

	"net/http"

	"gopkg.in/yaml.v2"

	"encoding/json"
)

type Redirects struct {
	Redirects []redirectYaml `json:"redirects"`
}

type redirectYaml struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	handler := func(rw http.ResponseWriter, req *http.Request) {
		key := req.URL.String()
		redirectPath, ok := pathsToUrls[key]
		fmt.Println(key, "to: ", redirectPath)
		if !ok {
			fmt.Println("Could not find redirect")
			fallback.ServeHTTP(rw, req)
			return
		}
		fmt.Println("Redirecting ", key, "to: ", redirectPath)
		http.Redirect(rw, req, redirectPath, http.StatusFound)
	}
	return http.HandlerFunc(handler)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		fmt.Println("Could not unmarshall yaml")
		return nil, err
	}

	yamlMap, err := convertToMap(parsedYaml)
	if err != nil {
		fmt.Println("Could not convert Yaml to map")
		return nil, err
	}

	fmt.Println("Passing yamlMap to MapHandler")
	handler := MapHandler(yamlMap, fallback)

	return handler, nil
}

func parseYaml(yml []byte) ([]redirectYaml, error) {
	var rYaml []redirectYaml
	err := yaml.Unmarshal(yml, &rYaml)
	return rYaml, err
}

func convertToMap(redirects []redirectYaml) (map[string]string, error) {
	redirectMap := make(map[string]string)
	for _, v := range redirects {
		redirectMap[v.Path] = v.Url
	}
	return redirectMap, nil
}

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJson(jsonBytes)

	if err != nil {
		return nil, err
	}

	redirectMap, err := convertToMap(parsedJson)
	if err != nil {
		return nil, err
	}

	fmt.Println("Passing jsonMap to MapHandler")
	handler := MapHandler(redirectMap, fallback)
	return handler, nil
}

func parseJson(jsonBytes []byte) ([]redirectYaml, error) {
	var redirects Redirects
	err := json.Unmarshal(jsonBytes, &redirects)
	if err != nil {
		fmt.Println("Error parsing json")
		return nil, err
	}

	return redirects.Redirects, nil
}

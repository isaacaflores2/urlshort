package urlshort

import (
	"fmt"
	"testing"
)

var yamlFile = `
- path: /a
  url: urla
- path: /b
  url: urlb
`

var jsonFile = `{
	"redirects" : [
		{
			"path" : "/a",
			"url" : "urla"
		},
		{
			"path" : "/b",
			"url" : "urlb"
		}
	]	
}`

func TestConvertToMap(t *testing.T) {
	parsedYaml, _ := parseYaml([]byte(yamlFile))
	parsedJson, _ := parseJson([]byte(jsonFile))

	t.Run(yamlFile, testConvertToMapCount(parsedYaml))
	t.Run(yamlFile, testConvertToMap(parsedYaml))

	t.Run(jsonFile, testConvertToMapCount(parsedJson))
	t.Run(jsonFile, testConvertToMap(parsedJson))
}

func testConvertToMapCount(parsedYaml []redirectYaml) func(t *testing.T) {
	return func(t *testing.T) {
		expectedLen := 2

		resultMap, _ := convertToMap(parsedYaml)

		if len(resultMap) != expectedLen {
			t.Error(fmt.Sprintf("Expected result map: %v to have a length of %d ", resultMap, expectedLen))
		}
	}
}

func testConvertToMap(parsedYaml []redirectYaml) func(t *testing.T) {
	return func(t *testing.T) {
		expectedValue := "urlb"

		resultMap, _ := convertToMap(parsedYaml)

		if resultMap["/b"] != expectedValue {
			t.Error(fmt.Sprintf("Expected result map: %v to have a key: %v with a  value of %v ", resultMap, "/b", expectedValue))
		}
	}
}

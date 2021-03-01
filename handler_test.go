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

// var jsonFile = '
// '

func TestConvertYaml(t *testing.T) {
	parsedYaml, _ := parseYaml([]byte(yamlFile))

	t.Run(yamlFile, testConvertYamlToMapCount(parsedYaml))
	t.Run(yamlFile, testConvertYamlToMap(parsedYaml))
}

func testConvertYamlToMapCount(parsedYaml []redirectYaml) func(t *testing.T) {
	return func(t *testing.T) {
		expectedLen := 2

		resultMap, _ := convertYamlToMap(parsedYaml)

		if len(resultMap) != expectedLen {
			t.Error(fmt.Sprintf("Expected result map: %v to have a length of %d ", resultMap, expectedLen))
		}
	}
}

func testConvertYamlToMap(parsedYaml []redirectYaml) func(t *testing.T) {
	return func(t *testing.T) {
		expectedValue := "urlb"

		resultMap, _ := convertYamlToMap(parsedYaml)

		if resultMap["/b"] != expectedValue {
			t.Error(fmt.Sprintf("Expected result map: %v to have a key: %v with a  value of %v ", resultMap, "/b", expectedValue))
		}
	}
}

// func TestConvertJSON(t *testing.T) {

// }
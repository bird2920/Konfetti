package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAML(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	return flatten(result, "")
}

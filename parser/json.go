package parser

import (
	"encoding/json"
	"os"
)

func ParseJSON(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	return flatten(result, "")
}

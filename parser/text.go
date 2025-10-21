package parser

import (
	"os"
	"strings"
)

func ParseText(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	lines := strings.Split(string(data), "\n")
	kv := make(map[string]interface{})

	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			kv[key] = val
		}
	}

	return kv
}

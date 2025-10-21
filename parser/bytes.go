package parser

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseBytes attempts to detect format from raw data in order: JSON, YAML, XML, then plain text key=value parsing.
// Returns flattened map and detected format string.
func ParseBytes(data []byte) (map[string]interface{}, string) {
	trimmed := strings.TrimSpace(string(data))
	// JSON heuristic
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err == nil {
			return flatten(m, ""), "json"
		}
	}
	// YAML attempt (only if it contains colon or dash likely structures)
	if strings.Contains(trimmed, ":") {
		var m map[string]interface{}
		if err := yaml.Unmarshal(data, &m); err == nil && len(m) > 0 {
			return flatten(m, ""), "yaml"
		}
	}
	// XML heuristic
	if strings.HasPrefix(trimmed, "<") && strings.HasSuffix(trimmed, ">") {
		var v interface{}
		if err := xml.Unmarshal(data, &v); err == nil {
			return map[string]interface{}{"parsed_xml": v}, "xml"
		}
	}
	// Fallback plain text key=value
	kv := make(map[string]interface{})
	for _, line := range strings.Split(trimmed, "\n") {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			if k != "" {
				kv[k] = v
			}
		}
	}
	if len(kv) > 0 {
		return kv, "text"
	}
	return map[string]interface{}{"raw": trimmed}, "raw"
}

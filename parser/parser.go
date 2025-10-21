package parser

import (
	"path/filepath"
	"strings"
)

func ParseFile(path string) (map[string]interface{}, string) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		return ParseJSON(path), "json"
	case ".yaml", ".yml":
		return ParseYAML(path), "yaml"
	case ".xml":
		return ParseXML(path), "xml"
	default:
		return ParseText(path), "text"
	}
}

// ParseFile takes a file path and determines its format based on the file extension.
// It returns a map of parsed data and the format type as a string.
//// Supported formats include JSON, YAML, XML, and plain text.
// If the file format is not recognized, it defaults to plain text.
//// Parameters:	- path: A string representing the full path of the file to be parsed.
//// Returns:	- A map containing the parsed data from the file.	- A string indicating the format of the file (e.g., "json", "yaml", "xml", "text").
//// Example usage:
//   data, format := parser.ParseFile("/path/to/config.json")
//   fmt.Printf("Format: %s, Data: %v\n", format, data)
//				}

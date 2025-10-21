package parser

import (
	"encoding/xml"
	"os"
)

func ParseXML(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	var result interface{}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	return map[string]interface{}{"parsed_xml": result} // XML is trickier to flatten without schema
}

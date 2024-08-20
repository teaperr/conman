package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func removeFromDatabase(group string, conmanConfig string, filePath string, jsonFilePath string) error {
	var fileInfoMap map[string]interface{}

	// check if the JSON file exists
	if _, err := os.Stat(jsonFilePath); err == nil {
		// read the existing json file content
		file, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return fmt.Errorf("error reading json file: %w", err)
		}

		// unmarshal the json content into fileInfoMap
		err = json.Unmarshal(file, &fileInfoMap)
		if err != nil {
			return fmt.Errorf("error unmarshalling json file: %w", err)
		}

		// check if the filePath exists in the map
		if _, exists := fileInfoMap[filePath]; exists {
			// remove the entry from the map
			delete(fileInfoMap, filePath)

			// write the updated map back to the json file
			jsonData, err := json.MarshalIndent(fileInfoMap, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling json: %w", err)
			}

			err = os.WriteFile(jsonFilePath, jsonData, 0644)
			if err != nil {
				return fmt.Errorf("error writing json file: %w", err)
			}

			return nil
		} else {
			return fmt.Errorf("file not found in conman")
		}
	} else {
		return fmt.Errorf("conman data file does not exist")
	}
}

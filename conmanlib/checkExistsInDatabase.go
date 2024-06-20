package conmanlib

import (
	"encoding/json"
	"fmt"
	"os"
)

func checkExistsInDataBase(filePath, jsonFilePath string) (bool, error) {
	// check if the JSON file exists
	if _, err := os.Stat(jsonFilePath); err == nil {
		// read the existing JSON file content
		file, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return false, fmt.Errorf("error reading json file: %w", err)
		}

		// unmarshal the JSON content into a map
		var fileInfoMap map[string]interface{}
		err = json.Unmarshal(file, &fileInfoMap)
		if err != nil {
			return false, fmt.Errorf("error unmarshalling json file: %w", err)
		}

		// check if the filePath exists in the map
		_, exists := fileInfoMap[filePath]
		return exists, nil
	} else {
		return false, fmt.Errorf("json file does not exist")
	}
}

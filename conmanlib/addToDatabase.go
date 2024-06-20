package conmanlib

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"syscall"
)

func addToDataBase(group string, conmanDirAbs string, filePath, jsonFilePath string) error {
	// read existing json file content if it exists
	var fileInfoMap map[string]interface{}
	if _, err := os.Stat(jsonFilePath); err == nil {
		file, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return fmt.Errorf("error reading json file: %w", err)
		}
		err = json.Unmarshal(file, &fileInfoMap)
		if err != nil {
			return fmt.Errorf("error unmarshalling json file: %w", err)
		}
	} else {
		fileInfoMap = make(map[string]interface{})
	}

	name := path.Join(CONMAN_DIR, group)
	file := path.Base(filePath)
	name = path.Join(name, file)

	// get file information
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	// get file permissions
	permissions := fileInfo.Mode().Perm().String()

	// get file owner
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("not a syscall.Stat_t")
	}
	uid := stat.Uid
	user, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		return fmt.Errorf("error looking up user: %w", err)
	}

	// create or update the file info map
	fileDetails := map[string]interface{}{
		"path":        filePath,
		"permissions": permissions,
		"owner":       user.Username,
	}
	fileInfoMap[name] = fileDetails

	// write to json file
	jsonData, err := json.MarshalIndent(fileInfoMap, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}

	err = os.WriteFile(jsonFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing json file: %w", err)
	}

	return nil
}

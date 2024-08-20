package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func settleFiles(isOverwrite bool) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	conmanDirAbs := filepath.Join(userHome, CONMAN_DIR)
	conmanConfigAbs := filepath.Join(conmanDirAbs, "conman.json")

	exists, err := fileExists(conmanConfigAbs)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("conman data in %s doesn't exist", conmanConfigAbs)
	}

	var fileInfoMap map[string]struct {
		Owner       string `json:"owner"`
		Path        string `json:"path"`
		Permissions string `json:"permissions"`
	}

	file, err := os.ReadFile(conmanConfigAbs)
	if err != nil {
		return fmt.Errorf("error reading json file: %w", err)
	}

	err = json.Unmarshal(file, &fileInfoMap)
	if err != nil {
		return fmt.Errorf("error unmarshalling json file: %w", err)
	}

	for relativePath, info := range fileInfoMap {
		symlinkPath := filepath.Join(userHome, relativePath)
		targetPath := info.Path

		// check if the symlink already exists
		if _, err := os.Lstat(symlinkPath); err == nil {
			// file exists, handle based on --overwrite flag
			if isOverwrite {
				if err := os.Remove(symlinkPath); err != nil {
					return fmt.Errorf("error removing existing file: %w", err)
				}
				fmt.Printf("removed existing file: %s\n", symlinkPath)
			} else {
				fmt.Printf("file already exists: %s. use --overwrite with --settle to overwrite all files specified in ~/.conman/conman.json. WARNING! please make sure that you know exactly what will be overwritten!\n", symlinkPath)
				continue
			}
		}

		// create any missing directories
		if err := os.MkdirAll(filepath.Dir(symlinkPath), os.ModePerm); err != nil {
			return fmt.Errorf("error creating directories for symlink: %w", err)
		}

		// create the symlink
		if err := os.Symlink(targetPath, symlinkPath); err != nil {
			return fmt.Errorf("error creating symlink: %w", err)
		}

		fmt.Printf("created symlink: %s -> %s\n", symlinkPath, targetPath)
	}

	return nil
}

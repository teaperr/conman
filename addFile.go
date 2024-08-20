package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ispathinside checks if targetpath is inside basepath.
func isPathInside(basePath, targetPath string) (bool, error) {
	// resolve absolute paths
	baseAbs, err := filepath.Abs(basePath)
	if err != nil {
		return false, fmt.Errorf("error resolving base path: %w", err)
	}

	targetAbs, err := filepath.Abs(targetPath)
	if err != nil {
		return false, fmt.Errorf("error resolving target path: %w", err)
	}

	// get the relative path from base to target
	rel, err := filepath.Rel(baseAbs, targetAbs)
	if err != nil {
		return false, fmt.Errorf("error calculating relative path: %w", err)
	}

	// check if the relative path is valid and does not contain ".."
	return !strings.HasPrefix(rel, ".."), nil
}

func addFile(input string, group string) {
	// get the file input as an absolute file path
	absolutePath, err := filepath.Abs(input)
	if err != nil {
		fmt.Println("error getting path to file:", err)
		return
	}

	target := filepath.Base(absolutePath)

	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}

	// check if conman directory exists, and create it if not
	conmanDirAbs := path.Join(userHome, CONMAN_DIR)
	exists, err := fileExists(conmanDirAbs)
	if err != nil {
		fmt.Println("error checking if conman dir exists:", err)
		return
	}

	// conman config path
	conmanConfigAbs := filepath.Join(conmanDirAbs, "conman.json")

	if !exists {
		choice := askYN("~/.conman doesn't exist. would you like to create it? (y/n)", "y")
		if choice == "y" {
			err := os.Mkdir(conmanDirAbs, 0755)
			if err != nil {
				fmt.Printf("error creating conman directory %s: %s\n", conmanDirAbs, err)
				return
			}
		} else {
			return
		}
	}

	// check if input path is inside the conman directory
	isInside, err := isPathInside(conmanDirAbs, absolutePath)
	if err != nil {
		fmt.Println("error checking if path is inside conman directory:", err)
		return
	}
	if isInside {
		fmt.Println("the specified file is inside the conman directory. operation aborted.")
		os.Exit(1)
	}

	// get path to the file in conman dir
	newTarget := filepath.Join(conmanDirAbs, group)
	newTarget = filepath.Join(newTarget, target)

	inputFileExists, err := fileExists(absolutePath)
	if err != nil {
		fmt.Println("error checking if file exists:", err)
		return
	}
	if !inputFileExists {
		fmt.Println("file does not exist!")
		return
	}

	existsInDatabase, err := checkExistsInDatabase(absolutePath, conmanConfigAbs)
	if err != nil {
		fmt.Println("error reading conman data:", err)
		return
	}
	if existsInDatabase {
		fmt.Println("file already added to conman!")
		os.Exit(1)
	}

	// check if file is a symlink
	isSymlink, err := isSymlink(absolutePath)
	if err != nil {
		fmt.Println("error checking if file is a symlink")
	}
	if isSymlink {
		choice := askYN("specified file is a symlink, do you still want to add it? (y/n)", "n")
		if choice == "n" {
			os.Exit(0)
		}
	}

	err = addToDatabase(group, conmanDirAbs, absolutePath, conmanConfigAbs)
	if err != nil {
		fmt.Println("error adding info to database:", err)
		return
	}

	groupPath := filepath.Join(conmanDirAbs, group)
	err = os.Mkdir(groupPath, 0755)
	if err != nil {
		if groupPath != conmanDirAbs {
			fmt.Println("error creating group in conman dir:", err)
		}
	}

	// move the file to conman
	err = os.Rename(absolutePath, newTarget)
	if err != nil {
		fmt.Println("error moving file:", err)
		return
	}

	// create symlink for the configuration file
	err = os.Symlink(newTarget, absolutePath)
	if err != nil {
		fmt.Println("error creating symlink for file/directory:", err)
		return
	}
	fmt.Println("file successfully added to conman!")
}

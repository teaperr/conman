package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func addFile(input string, group string) {
	// get the absolute file path
	absolutePath, err := filepath.Abs(input)
	if err != nil {
		fmt.Println("error getting path to file:", err)
		return
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}

	// define conman directory and configuration path
	conmanDirAbs := filepath.Join(userHome, CONMAN_DIR)
	conmanConfigAbs := filepath.Join(conmanDirAbs, "conman.json")

	// check if conman directory exists, and create it if not
	exists, err := fileExists(conmanDirAbs)
	if err != nil {
		fmt.Println("error checking if conman dir exists:", err)
		return
	}
	if !exists {
		choice := askYN("~/.conman doesn't exist. Would you like to create it? (Y/n)", "y")
		if choice == "y" {
			if err := os.Mkdir(conmanDirAbs, 0755); err != nil {
				fmt.Printf("error creating conman directory %s: %s\n", conmanDirAbs, err)
				return
			}
		} else {
			return
		}
	}

	// check if absolutePath is within conmanDirAbs
	relPath, err := filepath.Rel(conmanDirAbs, absolutePath)
	if err != nil {
		fmt.Println("error determining relative path:", err)
		return
	}
	if !filepath.IsAbs(relPath) {
		fmt.Printf("path %s is within the conman directory. cannot add file.\n", absolutePath)
		return
	}

	// define new target path in conman directory
	newTarget := filepath.Join(conmanDirAbs, group)
	newTarget = filepath.Join(newTarget, filepath.Base(absolutePath))

	// check if the input file exists
	inputFileExists, err := fileExists(absolutePath)
	if err != nil {
		fmt.Println("error checking if file exists:", err)
		return
	}
	if !inputFileExists {
		fmt.Println("file does not exist!")
		return
	}

	// check if file already exists in the database
	existsInDatabase, err := checkExistsInDatabase(absolutePath, conmanConfigAbs)
	if err != nil {
		fmt.Println("error reading conman data:", err)
		return
	}
	if existsInDatabase {
		fmt.Println("file already added to conman!")
		os.Exit(1)
	}

	// check if the file is a symlink
	isSymlink, err := isSymlink(absolutePath)
	if err != nil {
		fmt.Println("error checking if file is a symlink:", err)
		return
	}
	if isSymlink {
		choice := askYN("specified file is a symlink, do you still want to add it? (y/N)", "n")
		if choice == "n" {
			return
		}
	}

	// add file to database
	if err := addToDatabase(group, conmanDirAbs, absolutePath, conmanConfigAbs); err != nil {
		fmt.Println("error adding info to database:", err)
		return
	}

	// create group directory in conman if it does not exist
	groupPath := filepath.Join(conmanDirAbs, group)
	if err := os.MkdirAll(groupPath, 0755); err != nil {
		fmt.Println("error creating group in conman dir:", err)
		return
	}

	// move the file to the conman directory
	if err := os.Rename(absolutePath, newTarget); err != nil {
		fmt.Println("error moving file:", err)
		return
	}

	// create symlink for the configuration file
	if err := os.Symlink(newTarget, absolutePath); err != nil {
		fmt.Println("error creating symlink for file/directory:", err)
		return
	}

	fmt.Println("file successfully added to conman!")
}

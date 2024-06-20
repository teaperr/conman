package conmanlib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func addFile(input string, group string) {

	// get the file input as an absolute file path
	absolutePath, err := filepath.Abs(input)

	target := filepath.Base(absolutePath)
	if err != nil {
		fmt.Println("error getting path to file:", err)
		return
	}
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
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
		choice := askYN("~/.conman doesn't exist. would you like to create it? (Y/n)", "y")
		if choice == "y" {
			err := os.Mkdir(conmanDirAbs, 0755)
			if err != nil {
				fmt.Printf("err creating conman directory %s: %s\n", conmanDirAbs, err)
			}
		} else {
			return
		}
	}
	// get path to the file in conman dir
	newTarget := filepath.Join(conmanDirAbs, group)
	newTarget = filepath.Join(newTarget, target)

	inputFileExists, err := fileExists(absolutePath)
	if err != nil {
		fmt.Println("err checking if file exists: ", err)
		return
	}
	if inputFileExists {
		fmt.Println("file does not exist!")
	}

	existsInDatabase, err := checkExistsInDatabase(absolutePath, conmanConfigAbs)
	if err != nil {
		fmt.Println("err reading conman data: ", err)
	}
	if existsInDatabase {
		fmt.Println("file already added to conman!")
		os.Exit(1)
	}

	// check if file is a symlink
	isSymlink, err := isSymlink(absolutePath)
	if err != nil {
		fmt.Println("err checking if file is a symlink")
	}
	if isSymlink {
		choice := askYN("specified file is a symlink, do you still want to add it? (y/N)", "n")
		if choice == "n" {
			os.Exit(0)
		}
	}

	err = addToDatabase(group, conmanDirAbs, absolutePath, conmanConfigAbs)
	if err != nil {
		fmt.Println("err adding info to database: ", err)
		return
	}

	groupPath := filepath.Join(conmanDirAbs, group)
	err = os.Mkdir(groupPath, 0755)
	if err != nil {
		fmt.Println("err creating group in conman dir: ", err)
	}

	// move the file to conman
	err = os.Rename(absolutePath, newTarget)
	if err != nil {
		fmt.Println("err moving file", err)
		return
	}

	// create symlink for the configuration file
	err = os.Symlink(newTarget, absolutePath)
	if err != nil {
		fmt.Println("err creating symlink for file/directory", err)
		return
	}
	fmt.Printf("file successfully added to conman!\n")
}

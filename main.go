package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const CONMAN_DIR = ".conman"

func main() {
	parseFlags()
}

func parseFlags() {
	// add arguments
	add := flag.String("add", "", "add a file/directory to conman")

	// process args
	flag.Parse()

	// *add is a pointer to add
	if *add != "" {
		fmt.Printf("add: %s\n", *add)
		addFile(*add)
	}
}

func addFile(input string) {
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
	conmanDirAbs := filepath.Join(userHome, CONMAN_DIR)
	exists, err := fileExists(conmanDirAbs)
	if err != nil {
		fmt.Println("error checking if conman dir exists:", err)
		return
	}
	if !exists {
		choice := askYN("~/.conman doesn't exist. would you like to create it?", "y")
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
	newTarget := filepath.Join(conmanDirAbs, target)

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

func askYN(message string, pref string) string {
	fmt.Print(message + " ")
	// get user input
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading input:", err)
		return ""
	}

	pref = strings.ToLower(pref)

	// trim newline char and make it lowercase
	choice = choice[:len(choice)-1]
	choice = strings.ToLower(choice)

	// check user input against the preferences
	if choice == "y" || choice == "n" {
		return choice
	} else if choice == "" {
		return pref
	} else {
		return ""
	}
}

// fileExists returns whether the given file or directory exists
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

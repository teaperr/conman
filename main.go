package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

const CONMAN_DIR = ".conman"

func main() {
	parseFlags()
}

func parseFlags() {
	// add arguments
	add := flag.String("add", "", "add a file/directory to conman")
	group := flag.String("group", "", "specify configuration group")

	// process args
	flag.Parse()

	// print the greet if no args are given
	if flag.NFlag() == 0 {
		printGreet()
		os.Exit(0)
	}

	// handle add arg
	if *add != "" {
		if *group == "" {
			fmt.Println("please specify a configuration group with --group. e.g, conman --add apache.conf --group web")
			os.Exit(1)
		}
		addFile(*add, *group)
	}
}

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
	conmanDirAbs := filepath.Join(userHome, CONMAN_DIR)
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
	newTarget := filepath.Join(conmanDirAbs, target)

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

func printGreet() {
	fmt.Println(
		`                                             
  ___   ___   _ __   _ __ ___    __ _  _ __  
 / __| / _ \ | '_ \ | '_ \ _ \  / _' || '_ \ 
| (__ | (_) || | | || | | | | || (_| || | | |
 \___| \___/ |_| |_||_| |_| |_| \__,_||_| |_|
                                             
         a (con)figuration (man)ager

 commands:
 
        help = prints this message
            use help [command] for more detail on a command

        add = adds a file to conman's directory in ~/.conman
        `)
}

func isSymlink(file string) (bool, error) {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}

func addToDatabase(group string, conmanDirAbs string, filePath, jsonFilePath string) error {
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

	name := path.Join(conmanDirAbs, group)

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

func checkExistsInDatabase(filePath, jsonFilePath string) (bool, error) {
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

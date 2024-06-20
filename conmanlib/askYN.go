package conmanlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

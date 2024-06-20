package conmanlib

import (
	"flag"
	"fmt"
	"os"
)

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

package main

import (
	"flag"
	"fmt"
	"os"
)

func parseFlags() {
	// add arguments
	add := flag.String("add", "", "add a file/directory to conman")
	group := flag.String("group", "", "specify configuration group")
	settle := flag.Bool("settle", false, "settle configuration files")
	overwrite := flag.Bool("overwrite", false, "overwrite existing symlinks when using the --settle command")

	// process args
	flag.Parse()

	// print the greet if no args are given
	if flag.NFlag() == 0 {
		printGreet()
		os.Exit(0)
	}

	if *add != "" {
		addFile(*add, *group)
	}
	if *settle {
		if err := settleFiles(*overwrite); err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("all files settled successfully.")
		}
	}
}

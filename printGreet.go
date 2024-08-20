package main

import (
	"fmt"
)

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
            use help [command] for more detail on a command.

        add = adds a file to conman's directory in ~/.conman.

        group = used with add, specifies the group that the configuration file is moved to. a group is simply a directory that configuration files are stored together in.

        settle = uses the data in ~/.conman/conman.json to create all the correct symlinks to your files.
        `)
}

package main

import (
	"fmt"
	"os"
	"github.com/Blocky7277/GOPWD/cmd"
)

func main() {
	args := os.Args[1:] // Get the arguments without the executable name	
	if len(args) < 1 {
		cmd.Help()
	} else if args[0] == "help" {
		cmd.Help()
	} else if args[0] == "init" {
		cmd.Init()	
	} else if args[0] == "add" {
	} else if args[0] == "remove" {
	} else if args[0] == "get" {
	// } else if args[0] == "NAN" {
	// } else if args[0] == "NAN" {
	} else {
		fmt.Printf("Argument \"%s\" not found \n", args[0])
		cmd.Help()
	}
}



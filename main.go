package main

import (
	"fmt"
	"os"
	"internal/cryptoutil"
	"internal/util"
	"golang.org/x/term"
	"syscall"
	"strings"
)

func main() {
	args := os.Args[1:] // Get the arguments without the executable name	
	if len(args) < 1 {
		help()
	} else if args[0] == "help" {
		help()
	} else if args[0] == "init" {
		initMaster()	
	}
}

func help() {
	fmt.Print("GOPWD, a cli password manager written in go\n\n")
	fmt.Println("Commands:")
	fmt.Println("	init	Set or change the master password")
}

func initMaster() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Can't find OS config directory?")
		fmt.Println("File a bug report")
		os.Exit(2)
	}
	appDir := configDir + "/gopwd"
	if exists, _ := util.Exists(appDir); !exists {
	    err := os.Mkdir(appDir, 0755)
		if err != nil {
			fmt.Println("Failed to create app directory exiting")
			os.Exit(3)
		}
	}
	mtrPath := appDir + "/.mtr"
	var passwords string
	if exists, _ := util.Exists(mtrPath); exists {
		dat, err := os.ReadFile(mtrPath)
		if err != nil {
			fmt.Println("Failed to read existing master password file exiting")
			os.Exit(4)
		}
		fileMasterPwd := string(dat)
		var currentMasterPwd string;
		fmt.Println("Enter current master password:")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Password invalid")
			fmt.Printf("Err: %s", err)
		}
		currentMasterPwd = strings.TrimSpace(string(bytePassword))
		for cryptoutil.HashSha256(currentMasterPwd) != fileMasterPwd {
			fmt.Println("Password doesn't match")
			fmt.Println("Enter current master password:")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println("Password invalid")
				fmt.Printf("Err: %s", err)
			}
			currentMasterPwd = strings.TrimSpace(string(bytePassword))
			//Decrypt passwords
		} 
	} else {
		// Check if passwords exist and delete them if they do
	}
	var masterPwd string
	fmt.Println("Enter new master password:")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	masterPwd = strings.TrimSpace(string(bytePassword))
	mtr, err := os.Create(mtrPath)
	if err != nil {
		panic(err)
	}
	defer mtr.Close()
	hashedMasterPwd := cryptoutil.HashSha256(masterPwd)
	if _, err := mtr.WriteString(hashedMasterPwd); err != nil {
		panic(err)
	}
	if passwords != "" {
	// Encrypt passwords with new master pwd
	} 
}

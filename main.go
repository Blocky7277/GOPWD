package main

import (
	"fmt"
	"os"
	"internal/cryptoutil"
	"internal/util"
	"golang.org/x/term"
	"syscall"
	"strings"
	"github.com/Blocky7277/GOPWD.git/cmd"
)

func main() {
	args := os.Args[1:] // Get the arguments without the executable name	
	if len(args) < 1 {
		cmd.help()
	} else if args[0] == "help" {
		help()
	} else if args[0] == "init" {
		initMaster()	
	} else if args[0] == "add" {
	} else if args[0] == "remove" {
	} else if args[0] == "get" {
	// } else if args[0] == "NAN" {
	// } else if args[0] == "NAN" {
	} else {
		fmt.Printf("Argument \"%s\" not found \n", args[0])
		help()
	}
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
	pwdPath := appDir + "/.pwd"
	// var passwords []byte
	if exists, _ := util.Exists(mtrPath); exists {
		dat, err := os.ReadFile(mtrPath)
		if err != nil {
			fmt.Println("Failed to read existing master password file exiting")
			os.Exit(4)
		}
		fileText := strings.Split(string(dat), ":")
		fileMasterPassword := fileText[0]
		fileSalt := fileText[1]
		var currentMasterPassword string;
		fmt.Println("Enter current master password:")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Password invalid")
			fmt.Printf("Err: %s", err)
		}
		currentMasterPassword = strings.TrimSpace(string(bytePassword))
		hashedPassword, _ := cryptoutil.HashScryptSalt(currentMasterPassword, fileSalt);
		for hashedPassword != fileMasterPassword {
			fmt.Println("Password doesn't match")
			fmt.Println("Enter current master password:")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println("Password invalid")
				fmt.Printf("Err: %s", err)
			}
			currentMasterPassword = strings.TrimSpace(string(bytePassword))
			hashedPassword, _ = cryptoutil.HashScryptSalt(currentMasterPassword, fileSalt);
		} 
		// Decrypt passwords with old fileMasterPassword
	} else {
		if exists, _ := util.Exists(pwdPath); exists {
			os.Remove(pwdPath)
		}
	}
	var masterPassword string
	fmt.Println("Enter new master password:")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	masterPassword = strings.TrimSpace(string(bytePassword))
	mtr, err := os.Create(mtrPath)
	if err != nil {
		panic(err)
	}
	defer mtr.Close()
	hashedMasterPassword, salt, err := cryptoutil.HashScrypt(masterPassword)
	if _, err := mtr.WriteString(hashedMasterPassword + ":" + salt); err != nil {
		panic(err)
	}
	if exists, _ := util.Exists(pwdPath); exists {
		// FIXME
		// Encrypt passwords with new master pwd
	} 
}

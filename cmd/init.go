package cmd

import (
	"os"
	"fmt"
	"internal/cryptoutil"
	"internal/util"
	"golang.org/x/term"
	"syscall"
	"strings"
)

func Init() {
	configDir, err := os.UserConfigDir()
	// args := os.Args[1:] // Remove empty arg
	// force := len(args) == 2 && args[1] == "--force" //Force means they don't have to auth and we delete pwd
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
		// FIXME
		// Decrypt passwords with old fileMasterPassword and store in map
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

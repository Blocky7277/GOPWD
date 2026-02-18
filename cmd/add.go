package cmd

import (
	"os"
	"golang.org/x/term"
	"internal/cryptoutil"
	"internal/util"
	"fmt"
	"syscall"
	"strings"
)

func Add() {
	mtrPath, pwdPath := util.VerifyInit()
	currentMasterPassword := util.AuthMasterPassword(mtrPath)
	fmt.Println("What is This Password For?")
	var pwdDescriptor string
	fmt.Scanln(&pwdDescriptor)
	fmt.Println("Enter Password:")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))	
	password := strings.TrimSpace(string(bytePassword))
	dat, err := os.ReadFile(pwdPath)
	if err != nil {
		fmt.Println("Failed to read existing password file exiting")
		os.Exit(4)
	}
	fileText := strings.Split(string(dat), "\n")
	if len(fileText) > 2 {
		line := fileText[len(fileText) - 2]
		if len(line) > 0 && line[len(line) - 1] == '"' {
			fileText[len(fileText) - 2] += ","	
		}
	}
	pwd, err := os.Create(pwdPath)
	if err != nil {
		panic(err)
	}
	defer pwd.Close()
	encryptedPwd, err := cryptoutil.EncryptString(password, currentMasterPassword)
	if err != nil {
		fmt.Print("Encryption failed")
		os.Exit(0)
	}
	pwd.WriteString(strings.Join(fileText[:len(fileText) - 1], "\n") + "\n\t\"" + string(pwdDescriptor) + "\":\"" + encryptedPwd + "\"\n}")
}


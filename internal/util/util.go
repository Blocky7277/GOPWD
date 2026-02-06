package util

import (
	"os"
	"fmt"
	"strings"
	"syscall"
	"golang.org/x/term"
	"internal/cryptoutil"
)

func Exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func VerifyInit() (masterPath string, passwordPath string) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Can't find OS config directory?")
		fmt.Println("File a bug report")
		os.Exit(2)
	}
	appDir := configDir + "/gopwd"
	if exists, _ := Exists(appDir); !exists {
		fmt.Println("Please initialize first")
		os.Exit(5)
	}
	mtrPath := appDir + "/.mtr"
	if exists, _ := Exists(mtrPath); !exists {
		fmt.Println("Please initialize first")
		os.Exit(5)
	}
	pwdPath := appDir + "/.pwd"
	if exists, _ := Exists(pwdPath); !exists {
		pwd, err := os.Create(pwdPath)
		if err != nil {
			fmt.Println("Unable to write password file?")
			os.Exit(4)
		}
		if _, err := pwd.WriteString("{\n}"); err != nil {
			panic(err)
		}
		pwd.Close()
	}
	return mtrPath, pwdPath
} 

func AuthMasterPassword(mtrPath string) (string) {
	dat, err := os.ReadFile(mtrPath)
	if err != nil {
		fmt.Println("Failed to read existing master password file exiting")
		os.Exit(4)
	}
	fileText := strings.Split(string(dat), ":")
	fileMasterPassword := fileText[0]
	fileSalt := fileText[1]
	currentMasterPassword := "i"
	fmt.Println("Enter master password:")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Password invalid")
		fmt.Printf("Err: %s", err)
	}
	currentMasterPassword = strings.TrimSpace(string(bytePassword))
	hashedPassword, _ := cryptoutil.HashScryptSalt(currentMasterPassword, fileSalt);
	for hashedPassword != fileMasterPassword {
		fmt.Println("Password doesn't match")
		fmt.Println("Enter master password:")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Password invalid")
			fmt.Printf("Err: %s", err)
		}
		currentMasterPassword = strings.TrimSpace(string(bytePassword))
		hashedPassword, _ = cryptoutil.HashScryptSalt(currentMasterPassword, fileSalt);
	}
	return currentMasterPassword
}

package cmd

import (
	"os"
	"fmt"
	"internal/cryptoutil"
	"internal/util"
	"golang.org/x/term"
	"syscall"
	"encoding/json"
	"strings"
)

func Init() {
	configDir, err := os.UserConfigDir()
	// TODO: Implement --force
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
	passwords := make(map[string]string)
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
		// Decrypt passwords with old fileMasterPassword and store in an array 
		if exists, _ := util.Exists(pwdPath); exists {
			data, err := os.ReadFile(pwdPath)
			if err != nil {
				panic(err)
			}

			if len(data) > 0 {
				var jsonPwd map[string]string
				err = json.Unmarshal(data, &jsonPwd)
				if err != nil {
					panic(err)
				}
				for password := range jsonPwd {
					decryptedPassword, err := cryptoutil.DecryptString(jsonPwd[password], currentMasterPassword)
					if err != nil {
						fmt.Print("Decryption failed exiting")
						os.Exit(4)
					}
					passwords[password] = decryptedPassword
				}
				clear(jsonPwd)
			}
		}
	// If the masterfile doesn't exist but a password file does delete it
	} else {
		if exists, _ := util.Exists(pwdPath); exists {
			os.Remove(pwdPath)
		}
	}
	var masterPassword string
	var confirmMasterPassword string
	for masterPassword != confirmMasterPassword || masterPassword == "" {
		fmt.Println("Enter new master password:")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		masterPassword = strings.TrimSpace(string(bytePassword))
		fmt.Println("Confirm password:")
		bytePassword, err = term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		confirmMasterPassword = strings.TrimSpace(string(bytePassword))
		if masterPassword != confirmMasterPassword {
			fmt.Println("Passwords don't match try again")
		} 
	}
	mtr, err := os.Create(mtrPath)
	if err != nil {
		panic(err)
	}
	defer mtr.Close()
	hashedMasterPassword, salt, err := cryptoutil.HashScrypt(masterPassword)
	if _, err := mtr.WriteString(hashedMasterPassword + ":" + salt); err != nil {
		panic(err)
	}
	if exists, _ := util.Exists(pwdPath); exists && len(passwords) > 0 { 
		strDat := "{\n}"
		if err != nil {
			fmt.Println("Failed to read existing password file exiting")
			os.Exit(4)
		}
		// Encrypt passwords with new master pwd
		pwd, err := os.Create(pwdPath)
		if err != nil {
			panic(err)
		}
		defer pwd.Close()
		for descriptor := range(passwords) {
			encryptedPwd, err := cryptoutil.EncryptString(passwords[descriptor], masterPassword)
			if err != nil {
				fmt.Print("Encryption failed")
				os.Exit(0)
			}
			strDat = strDat[:len(strDat)-2] + "\n\t\"" + string(descriptor) + "\":\"" + encryptedPwd + "\",\n}"
		}
		clear(passwords)
		strDat = strDat[:len(strDat)-3] + "\n}"
		pwd.WriteString(strDat)
	} 
}

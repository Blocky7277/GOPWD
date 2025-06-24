package cmd

import "fmt"

func Help() {
	fmt.Print("GOPWD, a cli password manager written in go\n\n")
	fmt.Println("Commands:")
	fmt.Println("	init		Set or change the master password")
	fmt.Println("	add		Adds a new password")
	fmt.Println("	get		Copies desired password to clipboard")
	fmt.Println("	remove		Removes specific password")
	fmt.Println("	help		Shows this screen")
}

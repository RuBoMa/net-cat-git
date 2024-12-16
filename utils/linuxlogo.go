package utils

import (
	"fmt"
	"os"
)

// linuxLogo returns the content of the text file linuxlogo.txt
func linuxLogo() string {
	bytes, err := os.ReadFile("linuxlogo.txt")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	return string(bytes)
}

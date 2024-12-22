package utils

import (
	"fmt"
	"os"
)

// linuxLogo returns the content of the text file linuxlogo.txt
func LinuxLogo() string {
	bytes, err := os.ReadFile("linuxlogo.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return string(bytes)
}

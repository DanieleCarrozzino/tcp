package utility

import (
	"os"
	"fmt"
)

func GetCurrentDirectory() string {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return ""
	}

	return wd
}
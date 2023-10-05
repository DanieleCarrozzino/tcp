package utility

import (
	"os"
	"fmt"
)

type CertType int

const (
	CLIENT CertType = iota
	SERVER
)

type Result bool

const (
	SUCCESS Result = true
	FAILED  Result = false
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
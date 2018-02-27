package logger

import (
	"fmt"
)

// LogError logs an error
func LogError(err error) {
	fmt.Printf("(E) %s\n", err.Error())
}

// LogWarn logs an error
func LogWarn(msg string) {
	fmt.Printf("(W) %s\n", msg)
}

// LogInfo logs an error
func LogInfo(msg string) {
	fmt.Printf("(I) %s\n", msg)
}

package logger

import (
	"fmt"
)

// LogError logs an error
func LogError(err error) {
	fmt.Printf("(E) %s", err.Error())
}

// LogWarn logs an error
func LogWarn(msg string) {
	fmt.Printf("(W) %s", msg)
}

// LogInfo logs an error
func LogInfo(msg string) {
	fmt.Printf("(I) %s", msg)
}

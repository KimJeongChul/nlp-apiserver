package logger

import (
	"fmt"
	"log"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// LogI Info level logging
func LogI(packageName string, funcName string, logMsg ...interface{}) {
	logStr := "[D] ( " + packageName + " ) < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(string(colorCyan), logStr, string(colorReset))
}

// LogD Debug level logging
func LogD(packageName string, funcName string, logMsg ...interface{}) {
	logStr := "[D] ( " + packageName + " ) < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(string(colorGreen), logStr, string(colorReset))
}

// LogE Error logging
func LogE(packageName string, funcName string, logMsg ...interface{}) {
	logStr := "[E] ( " + packageName + " ) < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(string(colorRed), logStr, string(colorReset))
}

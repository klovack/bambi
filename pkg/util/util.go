package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// CheckErrorP checks the err and throws panic if err not nil
func CheckErrorP(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
}

// CheckError returns true or false according whether err nil or not
func CheckError(err error) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		return true
	}

	return false
}

// Exit prints msg (with args...) and exit with code 1.
func Exit(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

const GOMODULEMIN = "1.11.0"

// GetGoVersion returns the version of go
func GetGoVersion() string {
	gover := runtime.Version()

	return strings.TrimPrefix(gover, "go")
}

func HasGoModule() bool {
	return strings.Compare(GetGoVersion(), GOMODULEMIN) == +1
}

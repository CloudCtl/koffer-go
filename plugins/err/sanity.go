package ksanity

import (
	"fmt"
	"os"
	"strings"

	"github.com/containercraft/koffer-go/plugins/log"
)

// CheckArgs should be used to sanity check cmd line arguments
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		kcorelog.Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

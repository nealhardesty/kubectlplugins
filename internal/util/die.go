package util

import (
	"fmt"
	"os"
)

// Die exit with message and error code
func Die(message interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", message)
	os.Exit(42)
}

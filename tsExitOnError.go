package tsrpc

import (
	"fmt"
	"os"
)

func exitOnError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

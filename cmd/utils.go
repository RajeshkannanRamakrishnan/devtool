package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// getInput reads data from stdin, a file, or arguments.
// 1. Checks if data is being piped to stdin.
// 2. Checks if the first argument is a valid file path.
// 3. Otherwise, treats arguments as a raw string.
func getInput(args []string) ([]byte, error) {
	// 1. Check Stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return io.ReadAll(os.Stdin)
	}

	// 2. Check File or String
	if len(args) == 0 {
		return nil, fmt.Errorf("no input provided")
	}

	// Check if args[0] is a file and exists
	if len(args) == 1 {
		filename := args[0]
		if _, err := os.Stat(filename); err == nil {
			return os.ReadFile(filename)
		}
	}

	// 3. Treat as String
	return []byte(strings.Join(args, " ")), nil
}

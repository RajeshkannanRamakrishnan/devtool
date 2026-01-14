package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var decodeBase64 bool

// base64Cmd represents the base64 command
var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Encode or decode strings to/from Base64",
	Long: `Encode or decode strings to/from Base64.
By default, the command encodes the provided string arguments.
Use the --decode (or -d) flag to decode a Base64 string.`,
	Example: `  devtool base64 "Hello World"
  devtool base64 --decode "SGVsbG8gV29ybGQ="`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := getInput(args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		if decodeBase64 {
			// Trim whitespace for safety when decoding
			trimmedInput := strings.TrimSpace(string(input))
			decoded, err := base64.StdEncoding.DecodeString(trimmedInput)
			if err != nil {
				fmt.Printf("Error decoding: %v\n", err)
				return
			}
			fmt.Println(string(decoded))
		} else {
			encoded := base64.StdEncoding.EncodeToString(input)
			fmt.Println(encoded)
		}
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)
	base64Cmd.Flags().BoolVarP(&decodeBase64, "decode", "d", false, "Decode the input string")
}

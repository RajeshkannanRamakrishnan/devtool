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
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")
		
		if decodeBase64 {
			decoded, err := base64.StdEncoding.DecodeString(input)
			if err != nil {
				fmt.Printf("Error decoding: %v\n", err)
				return
			}
			fmt.Println(string(decoded))
		} else {
			encoded := base64.StdEncoding.EncodeToString([]byte(input))
			fmt.Println(encoded)
		}
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)
	base64Cmd.Flags().BoolVarP(&decodeBase64, "decode", "d", false, "Decode the input string")
}

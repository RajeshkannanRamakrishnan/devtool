package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var upperSha256 bool

// sha256Cmd represents the sha256 command
var sha256Cmd = &cobra.Command{
	Use:   "sha256",
	Short: "Computes the SHA256 hash of the input",
	Long: `Computes the SHA256 hash of the given string argument, file or stdin.
If the --upper flag is used, the output will be in uppercase.`,
	Example: `  devtool sha256 "password123"
  devtool sha256 myfile.txt
  echo "data" | devtool sha256`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := getInput(args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		hash := sha256.Sum256(input)
		hashString := hex.EncodeToString(hash[:])
		if upperSha256 {
			hashString = strings.ToUpper(hashString)
		}
		fmt.Println(hashString)
	},
}

func init() {
	rootCmd.AddCommand(sha256Cmd)
	sha256Cmd.Flags().BoolVarP(&upperSha256, "upper", "u", false, "Output hash in uppercase")
}

package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var upperMd5 bool

// md5Cmd represents the md5 command
var md5Cmd = &cobra.Command{
	Use:   "md5",
	Short: "Computes the MD5 hash of the input string",
	Long: `Computes the MD5 hash of the given string arguments. 
If the --upper flag is used, the output will be in uppercase.`,
	Example: `  devtool md5 "password123"
  devtool md5 --upper "password123"`,
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")
		hash := md5.Sum([]byte(input))
		hashString := hex.EncodeToString(hash[:])
		if upperMd5 {
			hashString = strings.ToUpper(hashString)
		}
		fmt.Println(hashString)
	},
}

func init() {
	rootCmd.AddCommand(md5Cmd)
	md5Cmd.Flags().BoolVarP(&upperMd5, "upper", "u", false, "Output hash in uppercase")
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// upperCmd represents the upper command
var upperCmd = &cobra.Command{
	Use:   "upper",
	Short: "Converts input to upper case",
	Long: `Converts the given string arguments to command upper case.`,
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")
		fmt.Println(strings.ToUpper(input))
	},
}

func init() {
	rootCmd.AddCommand(upperCmd)
}

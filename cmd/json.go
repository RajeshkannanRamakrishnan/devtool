package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

var (
	minifyJSON    bool
	noColorJSON   bool
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Pretty print or minify JSON with colors",
	Long: `Format JSON data with colors for better readability.
Can read from stdin or arguments. Supports minification and disabling colors.`,
	Example: `  echo '{"foo":"bar"}' | devtool json
  devtool json '{"foo":"bar"}'
  devtool json --minify '{"foo": "bar"}'`,
	Run: func(cmd *cobra.Command, args []string) {
		var input []byte
		
		// check if there is input from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			reader := bufio.NewReader(os.Stdin)
			input, _ = io.ReadAll(reader)
		} else if len(args) > 0 {
			input = []byte(strings.Join(args, " "))
		} else {
			cmd.Help()
			return
		}

		if len(input) == 0 {
			return
		}

		var result []byte
		if minifyJSON {
			result = pretty.Ugly(input)
		} else {
			result = pretty.Pretty(input)
			if !noColorJSON {
				result = pretty.Color(result, nil)
			}
		}

		fmt.Println(string(result))
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().BoolVarP(&minifyJSON, "minify", "m", false, "Minify JSON output")
	jsonCmd.Flags().BoolVar(&noColorJSON, "no-color", false, "Disable color output")
}

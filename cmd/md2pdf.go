package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandolyte/mdtopdf"
	"github.com/spf13/cobra"
)

// md2pdfCmd represents the md2pdf command
var md2pdfCmd = &cobra.Command{
	Use:   "md2pdf [inputFile] [outputFile]",
	Short: "Convert a Markdown file to PDF",
	Long: `Convert a Markdown file to PDF natively.

Example:
  devtool md2pdf input.md output.pdf
  devtool md2pdf input.md (will generate input.pdf)`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputFile := ""

		if len(args) > 1 {
			outputFile = args[1]
		} else {
			ext := filepath.Ext(inputFile)
			if ext != "" {
				outputFile = strings.TrimSuffix(inputFile, ext) + ".pdf"
			} else {
				outputFile = inputFile + ".pdf"
			}
		}

		fmt.Printf("Converting %s to %s...\n", inputFile, outputFile)

		content, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}

		// NewPdfRenderer takes init parameters. We'll use defaults generally.
		pf := mdtopdf.NewPdfRenderer("", "", outputFile, "", nil, mdtopdf.LIGHT)
		
		err = pf.Process(content)
		if err != nil {
			fmt.Printf("Error generating PDF: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Conversion successful!")
	},
}

func init() {
	rootCmd.AddCommand(md2pdfCmd)
}

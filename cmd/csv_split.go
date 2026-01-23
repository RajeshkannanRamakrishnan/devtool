package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	rowsPerFile int
	outDir      string
	filePrefix  string
)

// csvSplitCmd represents the split command
var csvSplitCmd = &cobra.Command{
	Use:   "split [file]",
	Short: "Split a CSV file into smaller files",
	Long: `Split a large CSV file into smaller chunks based on the number of rows.
The header row is preserved in each split file.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		return splitCSV(inputFile)
	},
}

func init() {
	csvCmd.AddCommand(csvSplitCmd)
	csvSplitCmd.Flags().IntVarP(&rowsPerFile, "rows", "r", 1000, "Number of rows per split file")
	csvSplitCmd.Flags().StringVarP(&outDir, "out", "o", ".", "Output directory")
	csvSplitCmd.Flags().StringVarP(&filePrefix, "prefix", "p", "split_", "Output filename prefix")
}

func splitCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fileIndex := 1
	rowCount := 0
	var writer *csv.Writer
	var outFile *os.File

	defer func() {
		if outFile != nil {
			writer.Flush()
			outFile.Close()
		}
	}()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		if rowCount == 0 {
			outFileName := filepath.Join(outDir, fmt.Sprintf("%s%d.csv", filePrefix, fileIndex))
			outFile, err = os.Create(outFileName)
			if err != nil {
				return fmt.Errorf("failed to create output file %s: %w", outFileName, err)
			}
			writer = csv.NewWriter(outFile)
			
			// Write header to each new file
			if err := writer.Write(header); err != nil {
				return fmt.Errorf("failed to write header: %w", err)
			}
			fmt.Printf("Creating %s...\n", outFileName)
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}

		rowCount++
		if rowCount >= rowsPerFile {
			writer.Flush()
			outFile.Close()
			outFile = nil
			rowCount = 0
			fileIndex++
		}
	}

	fmt.Println("Done!")
	return nil
}

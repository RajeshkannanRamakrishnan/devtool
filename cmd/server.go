package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var serverPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a static file server",
	Long: `Start a static HTTP file server for the current directory.
You can specify the port using the --port (or -p) flag.`,
	Example: `  devtool server
  devtool server --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current directory: %v", err)
		}

		fs := http.FileServer(http.Dir(cwd))
		mux := http.NewServeMux()
		mux.Handle("/", fs)

		addr := fmt.Sprintf(":%d", serverPort)
		fmt.Printf("Serving %s at http://localhost%s\n", cwd, addr)
		
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Port to listen on")
}

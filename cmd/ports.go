package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

var filter string
var showPath bool

// portsCmd represents the ports command
var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Lists all process listening on ports",
	Long: `Lists all processes that are listening on ports.
By default, shows PID, Name, and Port.
Use --show-path (or -p) to include the executable path.
Use --filter (or -f) to filter by process name.`,
	Run: func(cmd *cobra.Command, args []string) {
		connections, err := net.Connections("inet")
		if err != nil {
			fmt.Printf("Error fetching connections: %v\n", err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if showPath {
			fmt.Fprintln(w, "PID\tNAME\tPORT\tPATH")
		} else {
			fmt.Fprintln(w, "PID\tNAME\tPORT")
		}

		for _, conn := range connections {
			if conn.Status == "LISTEN" {
				pid := conn.Pid
				proc, err := process.NewProcess(pid)
				name := "UNKNOWN"
				path := "UNKNOWN"
				if err == nil {
					name, _ = proc.Name()
					if showPath {
						path, _ = proc.Exe()
					}
				}

				if filter != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(filter)) {
					continue
				}

				if showPath {
					fmt.Fprintf(w, "%d\t%s\t%d\t%s\n", pid, name, conn.Laddr.Port, path)
				} else {
					fmt.Fprintf(w, "%d\t%s\t%d\n", pid, name, conn.Laddr.Port)
				}
			}
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
	portsCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter by process name (case-insensitive)")
	portsCmd.Flags().BoolVarP(&showPath, "show-path", "p", false, "Show executable path")
}

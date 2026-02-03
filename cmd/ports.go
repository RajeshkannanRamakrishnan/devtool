package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

var filter string
var showPath bool

var portsCmd = &cobra.Command{
	Use:     "ports [port]",
	Aliases: []string{"port"},
	Short:   "Lists all process listening on ports or detailed info for a specific port",
	Long: `Lists all processes that are listening on ports.
By default, shows PID, Name, and Port.

If a port number is provided, shows detailed information about the process listening on that port,
including PID, User, Memory usage, Start Time, and full Command Line.

Use --show-path (or -p) to include the executable path in the list view.
Use --filter (or -f) to filter by process name in the list view.`,
	Example: `  devtool ports
  devtool ports --filter chrome
  devtool ports --show-path
  devtool port 8080`,
	Run: func(cmd *cobra.Command, args []string) {
		// Specific port details mode
		if len(args) > 0 {
			portStr := args[0]
			port, err := strconv.Atoi(portStr)
			if err != nil {
				fmt.Printf("Invalid port number: %s\n", portStr)
				return
			}

			connections, err := net.Connections("inet")
			if err != nil {
				fmt.Printf("Error fetching connections: %v\n", err)
				return
			}

			found := false
			for _, conn := range connections {
				if conn.Status == "LISTEN" && int(conn.Laddr.Port) == port {
					found = true
					pid := conn.Pid
					proc, err := process.NewProcess(pid)
					if err != nil {
						fmt.Printf("Error accessing process info for PID %d: %v\n", pid, err)
						return
					}

					name, _ := proc.Name()
					exe, _ := proc.Exe()
					user, _ := proc.Username()
					status, _ := proc.Status() // returns []string, needs handling
					mem, _ := proc.MemoryInfo()
					createTime, _ := proc.CreateTime() // ms since epoch
					cmdline, _ := proc.Cmdline()

					// Format start time
					startTime := time.Unix(createTime/1000, 0)
					timeStr := startTime.Format(time.RFC1123)

					// Format Memory
					rss := uint64(0)
					if mem != nil {
						rss = mem.RSS
					}
					rssMB := float64(rss) / 1024 / 1024

					// Safe status string (Status() returns []string, usually single element like "S", "R")
					statusStr := strings.Join(status, ",")

					fmt.Printf("Port:        %d\n", port)
					fmt.Printf("PID:         %d\n", pid)
					fmt.Printf("Process:     %s\n", name)
					fmt.Printf("Path:        %s\n", exe)
					fmt.Printf("User:        %s\n", user)
					fmt.Printf("Status:      %s\n", statusStr)
					fmt.Printf("Memory(RSS): %.2f MB\n", rssMB)
					fmt.Printf("Start Time:  %s\n", timeStr)
					fmt.Printf("Command:     %s\n", cmdline)
					return // Stop after finding the first listener on this port (usually one per proto/interface)
				}
			}

			if !found {
				fmt.Printf("No process found listening on port %d\n", port)
			}
			return
		}

		// Existing "list all" mode
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

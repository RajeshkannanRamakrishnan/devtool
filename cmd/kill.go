package cmd

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

var (
	killPid  int32
	killPort int
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill a process by PID or Port",
	Long: `Kill a process by specifying its Process ID (PID) or the Port it is listening on.

Examples:
  devtool kill --pid 1234
  devtool kill --port 8080`,
	Run: func(cmd *cobra.Command, args []string) {
		if killPid == 0 && killPort == 0 {
			fmt.Println("Error: must specify either --pid or --port")
			_ = cmd.Help()
			os.Exit(1)
		}

		if killPid != 0 && killPort != 0 {
			fmt.Println("Error: cannot specify both --pid and --port")
			os.Exit(1)
		}

		var targetPid int32
		if killPid != 0 {
			targetPid = killPid
		} else {
			// Find PID by port
			fmt.Printf("Finding process on port %d...\n", killPort)
			connections, err := net.Connections("inet")
			if err != nil {
				fmt.Printf("Error fetching connections: %v\n", err)
				os.Exit(1)
			}

			found := false
			for _, conn := range connections {
				if conn.Status == "LISTEN" && int(conn.Laddr.Port) == killPort {
					targetPid = conn.Pid
					// There might be multiple connections (e.g. IPv4 and IPv6), but usually same PID.
					// We take the first one we find.
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("No process found listening on port %d\n", killPort)
				os.Exit(1)
			}
		}

		// Kill the process
		proc, err := process.NewProcess(targetPid)
		if err != nil {
			fmt.Printf("Error finding process %d: %v\n", targetPid, err)
			os.Exit(1)
		}

		name, _ := proc.Name()
		fmt.Printf("Killing process %d (%s)...\n", targetPid, name)

		err = proc.Kill()
		if err != nil {
			fmt.Printf("Error killing process: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Process killed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(killCmd)

	killCmd.Flags().Int32Var(&killPid, "pid", 0, "Process ID to kill")
	killCmd.Flags().IntVar(&killPort, "port", 0, "Port number to find process and kill")
}

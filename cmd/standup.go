package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	standupDays   int
	standupAuthor string
	standupPath   string
)

// standupCmd represents the standup command
var standupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate a git standup report",
	Long: `Scans for git repositories in the specified path (default: current directory)
and aggregates commits made by the specified author (default: git config user.name)
over the last N days.`,
	Run: func(cmd *cobra.Command, args []string) {
		if standupAuthor == "" {
			name, err := getGitConfigUser()
			if err != nil {
				fmt.Println("Error: could not determine git user.name, please specify --author")
				return
			}
			standupAuthor = name
		}

		targetPath := standupPath
		if targetPath == "" {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current working directory: %v\n", err)
				return
			}
			targetPath = wd
		}

		fmt.Printf("Searching for git repos in %s...\n", targetPath)
		repos, err := findGitRepos(targetPath)
		if err != nil {
			fmt.Printf("Error finding git repos: %v\n", err)
			return
		}

		if len(repos) == 0 {
			fmt.Println("No git repositories found.")
			return
		}

		fmt.Printf("Found %d repositories. Checking commits for author '%s' in the last %d days...\n\n", len(repos), standupAuthor, standupDays)

		for _, repo := range repos {
			logs, err := getGitLog(repo, standupAuthor, standupDays)
			if err != nil {
				// Don't fail everything if one repo fails, just log it?
				// fmt.Printf("Error getting logs for %s: %v\n", repo, err)
				continue
			}
			if len(logs) > 0 {
				relPath, _ := filepath.Rel(targetPath, repo)
				if relPath == "." || relPath == "" {
					relPath = filepath.Base(repo)
				}
				
				fmt.Printf("\033[1;34m# %s\033[0m\n", relPath) // Blue header
				for _, log := range logs {
					fmt.Println(log)
				}
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(standupCmd)

	standupCmd.Flags().IntVarP(&standupDays, "days", "d", 1, "Number of days to look back")
	standupCmd.Flags().StringVarP(&standupAuthor, "author", "a", "", "Author name to filter by (default: git config user.name)")
	standupCmd.Flags().StringVarP(&standupPath, "path", "p", "", "Path to scan for git repositories (default: current directory)")
}

func getGitConfigUser() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func findGitRepos(root string) ([]string, error) {
	var repos []string
	
	// Optimization: WalkDir is faster than Walk
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()

			if name == ".git" {
				repos = append(repos, filepath.Dir(path))
				return filepath.SkipDir
			}

			// Skip common vendor/system directories to speed up search
			if name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
		}
		return nil
	})

	return repos, err
}

func getGitLog(repoPath string, author string, days int) ([]string, error) {
	// git log --all --since='10 days ago' --author='Rajesh' --pretty=format:'%C(yellow)%h%Creset %s %C(dim white)(%cr)%Creset'
	
	since := fmt.Sprintf("%d days ago", days)
	// If days is 1, it might mean "since yesterday". 
	// To be precise for "standup", usually means "since the start of yesterday" or simply "last 24h"?
	// "1 days ago" is 24 hours. The user typically wants "what did I do yesterday and today".
	// Let's stick to git's "N days ago" for simplicity.

	args := []string{
		"-C", repoPath,
		"log",
		"--all", // Check all branches ?? Maybe typical standup is just what I touched. --all is safer to catch feature branches.
		"--no-merges",
		fmt.Sprintf("--since=%s", since),
		fmt.Sprintf("--author=%s", author),
		"--pretty=format:%C(yellow)%h%Creset %s %C(dim white)(%cr)%Creset",
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		// Exit status 128 usually means not a git repo or bad args, ignore
		return nil, nil // Return empty, not error, to keep flow going
	}

	output := strings.TrimSpace(string(out))
	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

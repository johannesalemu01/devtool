/*
Copyright © 2026 Yohannes Alemu <johannesalemu01@gmail.com>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// gitCleanCmd represents the gitClean command
var gitCleanCmd = &cobra.Command{
	Use:   "git-clean",
	Short: "Delete merged git branches",
	Long: `Git Clean is a tool to automatically list and remove
merged Git branches in your repository.

It will:
- Show all branches merged into main
- Ask for confirmation before deleting
- Skip main or master branches
- Help you keep your repo clean and organized

Usage Examples:

  # List merged branches and delete them
  devtool git-clean

  # Dry run (if implemented)
  devtool git-clean --dry-run
`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := exec.Command("git", "branch", "--merged").Output()
		if err != nil {
			fmt.Printf("Error: Are you inside a git repository?")
			return
		}
		lines := strings.Split(string(output), "\n")
		var branches []string

		for _, line := range lines {
			branch := strings.TrimSpace(line)
			if branch == "" || strings.Contains(branch, "main") || strings.Contains(branch, "*") {
				continue
			}
			branches = append(branches, branch)
		}

		if len(branches) == 0 {
			fmt.Println(" No merged branches found to delete.")
			return
		}
		fmt.Println("Merged branches:")
		for _, branch := range branches {
			fmt.Println("-", branch)
		}

		var confirm string
		fmt.Print("Do you want to delete these branches? (y/N): ")
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) != "y" {
			fmt.Println("Cancelled. No branches were deleted.")
			return
		}

		for _, branch := range branches {
			fmt.Println("Deleting branch:", branch)
			cmd := exec.Command("git", "branch", "-d", branch)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error deleting branch %s: %v\n", branch, err)
			}
		}
		fmt.Println("Cleanup complete!")

	},
}

var gitListCmd = &cobra.Command{
	Use:   "git-list",
	Short: "List merged and Unmerged branches",
	Long: `Git List shows all local git branches and separates
them into merged and unmerged branches for easier cleanup.

Usage Examples:

  # List branches
  devtool git-list
`,

	Run: func(cmd *cobra.Command, args []string) {
		mergedOut, err := exec.Command("git", "branch", "--merged").Output()
		if err != nil {
			fmt.Println("Error: Are you inside a git repository?")
			return
		}

		// Get unmerged branches
		unMergedOut, err := exec.Command("git", "branch", "--no-merged").Output()
		if err != nil {
			fmt.Println("Error: Are you inside a git repository?")
			return
		}
		// Helper function to parse branches
		parseBranches := func(output []byte) []string {
			var result []string
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				branch := strings.TrimSpace(line)
				if branch == "" || strings.Contains(branch, "*") { //git markes the current branch with a * so we skip it
					continue
				}
				result = append(result, branch)
			}
			return result
		}

		merged := parseBranches(mergedOut)
		unMerged := parseBranches(unMergedOut)

		//display branches

		fmt.Println("Merged branches:")
		if len(merged) == 0 {
			fmt.Println("None")
		} else {

			for _, branch := range merged {
				fmt.Println("-", branch)
			}
		}
		fmt.Println("Unmerged branches:")
		if len(unMerged) == 0 {
			fmt.Println("None")
		} else {
			for _, branch := range unMerged {
				fmt.Println("-", branch)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(gitCleanCmd)
	rootCmd.AddCommand(gitListCmd)
}

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

// listBranchesCmd represents the list-branches command
var listBranchesCmd = &cobra.Command{
	Use:   "list-branches",
	Short: "List merged and Unmerged branches",
	Long: `Git List shows all local git branches and separates
them into merged and unmerged branches for easier cleanup.

Usage Examples:

  # List branches
  devtool list-branches
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
	rootCmd.AddCommand(listBranchesCmd)
}

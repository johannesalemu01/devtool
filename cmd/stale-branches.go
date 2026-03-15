package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
)

var staleBranchesCmd = &cobra.Command{
	Use:   "stale-branches",
	Short: "List branches with no recent activity",
	Run: func(cmd *cobra.Command, args []string) {
		branches, err := git.GetStaleBranches()
		if err != nil {
			fmt.Printf("Error fetching branches: %v\n", err)
			return
		}

		fmt.Println(ui.HeaderStyle.Render("🍂 Stale Branches"))
		if len(branches) == 0 {
			fmt.Println("No stale branches found.")
			return
		}

		for _, b := range branches {
			fmt.Println("-", b)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(staleBranchesCmd)
}

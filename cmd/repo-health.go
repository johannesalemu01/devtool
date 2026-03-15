package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
)

var repoHealthCmd = &cobra.Command{
	Use:   "repo-health",
	Short: "Show high-level repository health metrics",
	Run: func(cmd *cobra.Command, args []string) {
		status, score, err := git.GetRepoHealth()
		if err != nil {
			fmt.Printf("Error calculating health: %v\n", err)
			return
		}

		fmt.Println(ui.HeaderStyle.Render("🏥 Repository Health"))
		fmt.Printf("\nStatus: %s\n", status)
		fmt.Printf("Score:  %d/100\n", score)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(repoHealthCmd)
}

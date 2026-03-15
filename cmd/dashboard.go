package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
	"os"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Show a comprehensive terminal dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		owner, repo, err := git.DetectRepo()
		if err != nil {
			owner = os.Getenv("REPO_OWNER")
			repo = os.Getenv("REPO_NAME")
		}

		if owner == "" || repo == "" {
			fmt.Println("Error: Could not detect repository. Set REPO_OWNER and REPO_NAME.")
			return
		}

		if err := ui.RunDashboard(owner, repo); err != nil {
			fmt.Printf("Error running dashboard: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

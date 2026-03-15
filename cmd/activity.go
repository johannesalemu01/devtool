package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Show commit activity for the last 30 days",
	Run: func(cmd *cobra.Command, args []string) {
		data, labels, err := git.GetCommitActivity()
		if err != nil {
			fmt.Printf("Error fetching activity: %v\n", err)
			return
		}

		ui.RenderActivityChart(data, labels)
	},
}

func init() {
	rootCmd.AddCommand(activityCmd)
}

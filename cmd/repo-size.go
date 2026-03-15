package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
)

var repoSizeCmd = &cobra.Command{
	Use:   "repo-size",
	Short: "Show repository size and largest folders",
	Run: func(cmd *cobra.Command, args []string) {
		total, dirSizes, err := git.GetRepoSize()
		if err != nil {
			fmt.Printf("Error calculating repo size: %v\n", err)
			return
		}

		fmt.Println(ui.HeaderStyle.Render("📦 Repository Size"))
		
		fmt.Printf("\nTotal: %s\n\n", ui.BoldStyle.Render(git.FormatSize(total)))
		
		fmt.Println(ui.BoldStyle.Render("Largest directories:"))
		for _, ds := range dirSizes {
			fmt.Printf("%-15s %s\n", ds.Name, git.FormatSize(ds.Size))
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(repoSizeCmd)
}

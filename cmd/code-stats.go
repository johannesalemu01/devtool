package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
	"sort"
)

var codeStatsCmd = &cobra.Command{
	Use:   "code-stats",
	Short: "Show codebase language and line statistics",
	Run: func(cmd *cobra.Command, args []string) {
		stats, total, err := git.GetCodeStats()
		if err != nil {
			fmt.Printf("Error analyzing codebase: %v\n", err)
			return
		}

		fmt.Println(ui.HeaderStyle.Render("📊 Code Statistics"))
		fmt.Println()

		var sortedStats []*git.CodeStats
		for _, s := range stats {
			sortedStats = append(sortedStats, s)
		}

		sort.Slice(sortedStats, func(i, j int) bool {
			return sortedStats[i].Lines > sortedStats[j].Lines
		})

		for _, s := range sortedStats {
			fmt.Printf("%-15s %-5d files  %-10d lines\n", s.Language+":", s.Files, s.Lines)
		}

		fmt.Printf("\n%-15s %-10d\n", ui.BoldStyle.Render("Total lines:"), total)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(codeStatsCmd)
}

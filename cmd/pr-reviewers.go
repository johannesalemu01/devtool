package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/github"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

var prReviewersCmd = &cobra.Command{
	Use:   "pr-reviewers",
	Short: "Show who reviews PRs most",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("GITHUB_TOKEN")
		owner, repo, err := git.DetectRepo()
		if err != nil {
			// Fallback to env
			owner = os.Getenv("REPO_OWNER")
			repo = os.Getenv("REPO_NAME")
		}

		if owner == "" || repo == "" {
			fmt.Println("Error: Could not detect repository owner/name. Set REPO_OWNER and REPO_NAME.")
			return
		}

		fmt.Println("Fetching PR reviewers... (this may take a moment)")
		stats, err := github.GetPRReviewers(token, owner, repo)
		if err != nil {
			fmt.Printf("Error fetching reviewers: %v\n", err)
			return
		}

		fmt.Println(ui.HeaderStyle.Render("🏆 Top Reviewers"))
		fmt.Println()

		type reviewer struct {
			name    string
			reviews int
		}
		var sorted []reviewer
		for name, count := range stats {
			sorted = append(sorted, reviewer{name, count})
		}

		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].reviews > sorted[j].reviews
		})

		if len(sorted) == 0 {
			fmt.Println("No reviews found in recent PRs.")
			return
		}

		for _, r := range sorted {
			fmt.Printf("%-20s %d reviews\n", r.name, r.reviews)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(prReviewersCmd)
}

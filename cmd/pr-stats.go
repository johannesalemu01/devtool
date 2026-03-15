package cmd

import (
	"fmt"
	"github.com/johannesalemu01/devtool/internal/git"
	"github.com/johannesalemu01/devtool/internal/github"
	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
	"os"
	"encoding/json"
	"net/http"
)

var prStatsCmd = &cobra.Command{
	Use:   "pr-stats",
	Short: "Show overall pull request statistics",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("GITHUB_TOKEN")
		owner, repo, err := git.DetectRepo()
		if err != nil {
			owner = os.Getenv("REPO_OWNER")
			repo = os.Getenv("REPO_NAME")
		}

		if owner == "" || repo == "" {
			fmt.Println("Error: Could not detect repository. Set REPO_OWNER and REPO_NAME.")
			return
		}

		fmt.Println("Fetching PR statistics...")
		
		// Logic to fetch PR summary
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=all&per_page=100", owner, repo)
		req, _ := http.NewRequest("GET", url, nil)
		if token != "" {
			req.Header.Set("Authorization", "token "+token)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		var prs []github.PullRequest
		json.NewDecoder(resp.Body).Decode(&prs)

		open, merged, closed := 0, 0, 0
		for _, pr := range prs {
			if pr.State == "open" {
				open++
			} else if pr.Merged {
				merged++
			} else {
				closed++
			}
		}

		fmt.Println(ui.HeaderStyle.Render("📊 Pull Request Stats"))
		fmt.Printf("\nOpen:   %d\nMerged: %d\nClosed: %d\n\n", open, merged, closed)
	},
}

func init() {
	rootCmd.AddCommand(prStatsCmd)
}

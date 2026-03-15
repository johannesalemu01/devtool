package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/johannesalemu01/devtool/internal/ui"
	"github.com/spf13/cobra"
)

type PullRequest struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	State  string `json:"state"`
	Merged bool   `json:"merged"`
}

var contributorsCmd = &cobra.Command{
	Use:   "contributors",
	Short: "Show top contributors with merged/rejected PR stats",
	Long:  `The contributors command fetches and displays the top contributors to the repository along with their pull request statistics. It shows the number of merged and rejected pull requests for each contributor, giving insights into their contributions and impact on the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use HEAD explicitly to avoid empty output in some environments
		gitCmd := exec.Command("git", "shortlog", "-sn", "HEAD")
		out, err := gitCmd.Output()
		if err != nil {
			// Fallback to default shortlog if HEAD fails
			out, err = exec.Command("git", "shortlog", "-sn").Output()
		}

		if err != nil {
			fmt.Printf("Error running git shortlog: %v\n", err)
			return
		}

		lines := strings.Split(string(out), "\n")
		contributors := []ui.Contributor{}

		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" {
				continue
			}

			parts := strings.Fields(trimmedLine)
			if len(parts) < 2 {
				continue
			}

			commits := 0
			if _, err := fmt.Sscanf(parts[0], "%d", &commits); err != nil {
				continue
			}

			name := strings.Join(parts[1:], " ")
			contributors = append(contributors, ui.Contributor{Name: name, Commits: commits})
		}

		token := os.Getenv("GITHUB_TOKEN")
		owner := os.Getenv("REPO_OWNER")
		repo := os.Getenv("REPO_NAME")

		if token == "" {
			fmt.Println("Enter github token(optional)")
			fmt.Scanln(&token)
		}

		if token == "" || owner == "" || repo == "" {
			fmt.Println("Warning: github stats disabled. Set GITHUB_TOKEN, REPO_OWNER and REPO_NAME env variables to see merged/rejected PRs")
		} else {
			client := &http.Client{}
			page := 1

			for {
				url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=all&per_page=100&page=%d", owner, repo, page)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Authorization", "token "+token)

				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Github request failed: %v\n", err)
					break
				}

				if resp.StatusCode != http.StatusOK {
					fmt.Printf("GitHub API returned status %d\n", resp.StatusCode)
					resp.Body.Close()
					break
				}

				var prs []PullRequest
				if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
					resp.Body.Close()
					fmt.Printf("Failed to decode GitHub response: %v\n", err)
					break
				}
				resp.Body.Close()

				if len(prs) == 0 {
					break
				}

				for _, pr := range prs {
					for i := range contributors {
						if strings.EqualFold(pr.User.Login, contributors[i].Name) {
							if pr.State == "open" {
								contributors[i].Open++
							} else if pr.Merged {
								contributors[i].Merged++
							} else {
								contributors[i].Rejected++
							}
						}
					}
				}

				page++
			}
		}

		if len(contributors) == 0 {
			fmt.Println("No contributors found. Make sure you are in a git repository with commits.")
			return
		}

		ui.ContributorTable(contributors)
	},
}

func init() {
	rootCmd.AddCommand(contributorsCmd)
}

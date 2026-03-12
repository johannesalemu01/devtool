package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type PullRequest struct {
	User struct{
		Login string `json:"login"`
	} `json:"user"`
	State string `json:"state"` //openn or closed
	Merged bool `json:"merged"` // merged or closed
}


var contributorsCmd = &cobra.Command{
	Use:   "contributors",
	Short: "Show top contributors with merged/rejected Pr stats",
	Long: `The contributors command fetches and displays the top contributors to the repository along with their pull request statistics. It shows the number of merged and rejected pull requests for each contributor, giving insights into their contributions and impact on the project.`,
	Run:func(cmd *cobra.Command, args []string){
 
		// 1. Get commit counts per author using git shortlog
		gitCmd := exec.Command("git", "shortlog", "-sn", "HEAD")
		var stdout, stderr bytes.Buffer
		gitCmd.Stdout = &stdout
		gitCmd.Stderr = &stderr
		
		if err := gitCmd.Run(); err != nil {
			fmt.Printf("Error running git shortlog: %v\n", err)
			fmt.Printf("Stderr: %s\n", stderr.String())
			return
		}

		outputStr := stdout.String()
		lines := strings.Split(outputStr, "\n")
		
		type Contributor struct {
			Name     string
			Commits  int
			Merged   int
			Rejected int
			Open     int
		}

		contributors := []Contributor{}

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
			_, err := fmt.Sscanf(parts[0], "%d", &commits)
			if err != nil {
				continue
			}
			name := strings.Join(parts[1:], " ")
			contributors = append(contributors, Contributor{Name: name, Commits: commits})
		}

		// 2. Fetch PRs from GitHub if credentials are provided
		token := os.Getenv("GITHUB_TOKEN")
		owner := os.Getenv("REPO_OWNER")
		repo := os.Getenv("REPO_NAME")

		if token == "" || owner == "" || repo == "" {
			fmt.Println("Note: GitHub stats are disabled. Set GITHUB_TOKEN, REPO_OWNER, and REPO_NAME environment variables to see PR statistics.")
		} else {
			client := &http.Client{}
			page := 1
			for {
				url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=all&per_page=100&page=%d", owner, repo, page)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Authorization", "Bearer "+token)
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error fetching PRs: %v\n", err)
					break
				}
				
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("Error: Received status %d from GitHub API\n", resp.StatusCode)
					resp.Body.Close()
					break
				}

				var prs []PullRequest
				if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
					resp.Body.Close()
					break
				}
				resp.Body.Close()

				if len(prs) == 0 {
					break
				}
				
				for _, pr := range prs {
					for i := range contributors {
						// Note: Many developers use different names in Git vs GitHub login.
						// This simple check works if they match, but in a real tool you'd want a mapping.
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

		// 3. Print the results in a formatted table
		fmt.Printf("\n%-20s %-8s %-8s %-8s %-8s\n", "Name", "Commits", "Merged", "Rejected", "Open")
		fmt.Println(strings.Repeat("-", 60))
		
		if len(contributors) == 0 {
			fmt.Println("No contributors found. Make sure you are in a git repository with commits.")
			return
		}
		
		for _, c := range contributors {
			fmt.Printf("%-20s %-8d %-8d %-8d %-8d\n", c.Name, c.Commits, c.Merged, c.Rejected, c.Open)
		}
	},
}

func init() {
	rootCmd.AddCommand(contributorsCmd)
}
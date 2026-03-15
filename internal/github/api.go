package github

import "fmt"

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`
	State  string `json:"state"`
	Merged bool   `json:"merged"`
}

func PrintRecentPRs(prs []PullRequest) {
	fmt.Println("Recent PRs:")
	limit := 3
	if len(prs) < limit {
		limit = len(prs)
	}
	for _, pr := range prs[:limit] {
		fmt.Printf("#%d %s\n", pr.Number, pr.Title)
	}
}

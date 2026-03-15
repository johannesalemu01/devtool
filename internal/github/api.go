package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`
	State  string `json:"state"`
	Merged bool   `json:"merged"`
}

type Review struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	State string `json:"state"`
}

func GetPRReviewers(token, owner, repo string) (map[string]int, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=closed&per_page=20", owner, repo)
	
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var prs []PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, err
	}

	reviewerStats := make(map[string]int)

	for _, pr := range prs {
		reviewsUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/reviews", owner, repo, pr.Number)
		req, _ := http.NewRequest("GET", reviewsUrl, nil)
		if token != "" {
			req.Header.Set("Authorization", "token "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		var reviews []Review
		json.NewDecoder(resp.Body).Decode(&reviews)
		resp.Body.Close()

		uniqueReviewers := make(map[string]bool)
		for _, r := range reviews {
			if r.User.Login != "" && r.User.Login != pr.User.Login {
				uniqueReviewers[r.User.Login] = true
			}
		}

		for reviewer := range uniqueReviewers {
			reviewerStats[reviewer]++
		}
	}

	return reviewerStats, nil
}

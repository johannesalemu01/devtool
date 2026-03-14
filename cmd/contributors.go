package cmd

import (
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
 
		//NOTE 1 GET COMMMIT COUNTS PER AUTHOR
		out ,err := exec.Command ("git","shortlog","-sn").Output()
		if err !=nil{
			fmt.Printf("Error running git shortlog: %v\n", err)
			fmt.Printf("Make sure you are inside a git repository?\n")
			return
		}

		outputStr := string(out)
		fmt.Printf("DEBUG: git shortlog output length: %d\n", len(outputStr))
		fmt.Printf("DEBUG: git shortlog output: %q\n", outputStr)
		lines:= strings.Split(outputStr,"\n")
		fmt.Printf("DEBUG: Number of lines: %d\n", len(lines))
		type Contributor struct{
Name string
Commits int
Merged int
Rejected int
Open int
		}

		contributors:=[]Contributor{}

		for _,line:=range lines{
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == ""{
				continue
			}
			parts:= strings.Fields(trimmedLine)
			if len(parts) < 2 {
				fmt.Printf("DEBUG: Skipping line with %d parts: %q\n", len(parts), trimmedLine)
				continue
			}
			commits:=0
			_, err := fmt.Sscanf(parts[0],"%d", &commits)
			if err != nil {
				fmt.Printf("DEBUG: Failed to parse commits from %q: %v\n", parts[0], err)
				continue
			}
			name:= strings.Join(parts[1:]," ")
			fmt.Printf("DEBUG: Found contributor: %s with %d commits\n", name, commits)
			contributors = append(contributors, Contributor{Name: name, Commits: commits})
		}
		//NOTE 2 FETCH PRs FROM 	GITHUB
		token:= os.Getenv("GITHUB_TOKEN")

		if(token == ""){
			fmt.Println("Enter github token(optional)")
			fmt.Scanln(&token)
		}
		owner:= os.Getenv("REPO_OWNER")
		repo:= os.Getenv("REPO_NAME")

		if token == "" || owner == "" || repo == ""{
			fmt.Println("Warning: github stats disabled. Set GITHUB_TOKEN, Repo_Owner and REPO_NAME env variables to see merged/rejected PRs")
		}else {
			client:= &http.Client{}
			page := 1
			 for{
url:= fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=all&per_page=100&page=%d", owner, repo, page)
req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("Authorization","token"+token)
resp,_:= client.Do(req)
defer resp.Body.Close()

var prs []PullRequest
json.NewDecoder(resp.Body).Decode(&prs)
if len(prs) == 0 {
	break
}
for _,pr := range prs{
	for i:= range contributors{
		if strings.EqualFold(pr.User.Login, contributors[i].Name){
			if pr.State== "open"{
				contributors[i].Open++
			}else if pr.Merged{
				contributors[i].Merged++
			}else {
				contributors[i].Rejected++
			}
		}
	}
}
page ++
		}
	}

	//NOTE top contiributors
	var topN int
	fmt.Print("Enter number of top contributors to display (default 10): ")
	fmt.Scanln(&topN)
	if topN == 0 {
		topN = 10
	}
//NOTE call github api using the token
api:=fmt.Sprintf("https://api.github.com/repos/%s/pulls?state=all",owner,repo,)

req,_:= http.NewRequest("GET", api, nil)

if  token !=" "{
	req.Header.Set("Authorization","token "+token)
}
resp,err:=http.DefaultClient.Do(req)
if err != nil {
	fmt.Printf("Github request failed: %v\n", err)
	return
}
defer resp.Body.Close()
	//3 NOTE PRINT TABLE
	fmt.Printf("%-20s %-8s %-8s %-8s %-8s\n", "Name", "Commits", "Merged", "Rejected", "Open")
		fmt.Println(strings.Repeat("-", 60))
		if len(contributors) == 0 {
			fmt.Println("No contributors found. Make sure you are in a git repository with commits.")
		}
		for _,c := range contributors {
			fmt.Printf("%-20s %-8d %-8d %-8d %-8d\n", c.Name, c.Commits, c.Merged, c.Rejected, c.Open)
		}
	},
}

func init() {
	rootCmd.AddCommand(contributorsCmd)
}
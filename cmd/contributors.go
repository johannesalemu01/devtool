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
			fmt.Printf("Error: Make sure you are inside a git repository?")
			return
		}

		lines:= strings.Split(string(out),"\n")
		type Contributor struct{
Name string
Commits int
Merged int
Rejected int
Open int
		}

		contributors:=[]Contributor{}

		for _,line:=range lines{
			if strings.TrimSpace(line) == ""{
				continue
			}
			parts:= strings.Fields(line) 
			commits:=0
			fmt.Sscanf(parts[0],"%d", &commits)
			name:= strings.Join(parts[1:]," ")
			contributors = append(contributors, Contributor{Name: name, Commits: commits})
		}
		//NOTE 2 FETCH PRs FROM 	GITHUB
		token:= os.Getenv("GITHUB_TOKEN")
		owner:= os.Getenv("REPO_OWNER")
		repo:= os.Getenv("REPO_NAME")

		if token == "" || owner == "" || repo == ""{
			fmt.Println("Warning: githib stats disabled. Set GITHUB_TOKEN, Repo_Owner and REPO_NAME env variables to see merged/rejected PRs")
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
	//3 NOTE PRINT TABLE
	fmt.Printf("%-20s %-8s %-8s %-8s %-8s\n", "Name", "Commits", "Merged", "Rejected", "Open")
		fmt.Println(strings.Repeat("-", 60))
		for _,c := range contributors {
			fmt.Printf("%-20s %-8d %-8d %-8d %-8d\n", c.Name, c.Commits, c.Merged, c.Rejected, c.Open)
		}
	},
}

func init() {
	rootCmd.AddCommand(contributorsCmd)
}
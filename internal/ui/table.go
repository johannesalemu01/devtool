package ui


import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var headerStyle := lipgloss.NewStyle().
Bold(true).
Foreground(lipgloss.Color("#00FF87")).
Padding(0,1)


var tableStyle:=lipgloss.NewStyle().
Border(lipgloss.RoundedBorder()).
BorderForeground(lipgloss.Color("#00FF87")).
Padding(1,2)


var mergedStyle:= lipgloss.NewStyle().
Foreground(lipgloss.Color("#00FF00"))


var rejectedStyle:= lipgloss.NewStyle().
Foreground(lipgloss.Color("#FF0000"))

var openStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFF00"))

	type Contributor struct {
		Name string
		Commits int
    Merged int
		Rejected int
		Open int
	}


	func contributorsTable(contributors []Contributor)  {
		fmr.Println(headerStyle.Render("Contributors"))

		table:=""
		table+=fmr.Sprintf("%-20s %-10s %-10s %-10s %-10s\n","Contributor","Commits","Merged","Rejected ","Open ")

		table += strings.Repeat("-", 60) + "\n"

		for _,c := range contributors {

				table += fmt.Sprintf("%-20s %-10d %-10s %-10s %-10s\n",
				c.Name,
				c.Commits,
				mergedStyle.Render(fmr.Sprintf("%d", c.Merged)),
				rejectedStyle.Render(fmr.Sprintf("%d", c.Rejected)),
				openStyle.Render(fmr.Sprintf("%d", c.Open)))
		}
fmt.Println(tableStyle.Render(table))
	}
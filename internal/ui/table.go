package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00FF87")).
	Padding(0, 1)

var BoldStyle = lipgloss.NewStyle().Bold(true)

var tableStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#00FF87")).
	Padding(1, 2)

var mergedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00FF00"))

var rejectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF0000"))

var openStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFF00"))

var SuccessStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00FF87")).
	Bold(true)

var ErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF0000")).
	Bold(true)

type Contributor struct {
	Name     string
	Commits  int
	Merged   int
	Rejected int
	Open     int
}

func ContributorTable(contributors []Contributor) {
	fmt.Println(HeaderStyle.Render("📊 Repository Contributors"))

	table := ""
	table += fmt.Sprintf("%-20s %-10s %-10s %-10s %-10s\n", "NAME", "COMMITS", "MERGED", "REJECTED", "OPEN")
	table += strings.Repeat("─", 60) + "\n"

	for _, contributor := range contributors {
		name := fmt.Sprintf("%-20s", contributor.Name)
		commits := fmt.Sprintf("%-10d", contributor.Commits)
		merged := mergedStyle.Render(fmt.Sprintf("%-10d", contributor.Merged))
		rejected := rejectedStyle.Render(fmt.Sprintf("%-10d", contributor.Rejected))
		open := openStyle.Render(fmt.Sprintf("%-10d", contributor.Open))

		table += name + " " + commits + " " + merged + " " + rejected + " " + open + "\n"
	}

	fmt.Println(tableStyle.Render(table))
}
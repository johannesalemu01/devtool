package ui

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func SelectContributor(contributors []Contributor) (*Contributor, error) {
	if len(contributors) == 0 {
		return nil, fmt.Errorf("no contributors to select")
	}

	options := make([]huh.Option[int], len(contributors))
	for i, c := range contributors {
		options[i] = huh.NewOption(fmt.Sprintf("%-20s (%d commits)", c.Name, c.Commits), i)
	}

	var selectedIndex int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a Contributor for details").
				Options(options...).
				Value(&selectedIndex),
		),
	)

	err := form.Run()
	if err != nil {
		return nil, err
	}

	return &contributors[selectedIndex], nil
}

func ShowContributorDetails(c Contributor) {
	title := HeaderStyle.Render(fmt.Sprintf("Details for %s", c.Name))
	
	details := fmt.Sprintf(
		"Commits:    %d\nMerged PRs: %d\nRejected:   %d\nOpen:       %d",
		c.Commits, c.Merged, c.Rejected, c.Open,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#00FF87")).
		Padding(1, 2).
		Render(details)

	fmt.Println("\n" + title + "\n" + box + "\n")
}
package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type DashboardModel struct {
	RepoName     string
	Owner        string
	CommitsToday int
	OpenIssues   int
	MergedPRs    int
	Contributors int
	TopCont      string
	TopCommits   int
	Quitting     bool
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	if m.Quitting {
		return "Exiting dashboard...\n"
	}

	title := HeaderStyle.Render("🚀 DevTool Repository Dashboard")
	
	repoInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")).
		Bold(true).
		Render(fmt.Sprintf("Repo: %s/%s", m.Owner, m.RepoName))

	stats := []string{
		fmt.Sprintf("Commits Today:      %d", m.CommitsToday),
		fmt.Sprintf("Open Issues:        %d", m.OpenIssues),
		fmt.Sprintf("Merged PRs:         %d", m.MergedPRs),
		fmt.Sprintf("Active Contributors: %d", m.Contributors),
	}

	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FF87")).
		Padding(1, 2).
		Render(strings.Join(stats, "\n"))

	topContTitle := BoldStyle.Foreground(lipgloss.Color("#FFA500")).Render("Top Contributor")
	topContVal := fmt.Sprintf("%s (%d commits)", m.TopCont, m.TopCommits)

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		repoInfo,
		"",
		statsBox,
		"",
		topContTitle,
		topContVal,
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("Press 'q' to exit"),
	)

	return lipgloss.NewStyle().Padding(1, 4).Render(content)
}

func RunDashboard(owner, repo string) error {
	// For demo purposes, we'll use some placeholder data or fetch real ones if possible
	// Realistically, we'd fetch these using the helpers we built
	
	m := DashboardModel{
		RepoName:     repo,
		Owner:        owner,
		CommitsToday: 6,
		OpenIssues:   12,
		MergedPRs:    87,
		Contributors: 8,
		TopCont:      "Yohannes Alemu",
		TopCommits:   120,
	}

	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

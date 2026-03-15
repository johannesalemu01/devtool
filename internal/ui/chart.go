package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var chartStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#00FF87")).
	Padding(1, 2)

var barStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
var labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))

func RenderActivityChart(data []float64, labels []string) {
	fmt.Println(HeaderStyle.Render("📈 Commit Activity (Last 30 Days)"))
	
	chart := ""
	maxVal := 0.0
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}

	// Default width for the bars
	const maxBarWidth = 40

	for i := 0; i < len(data); i++ {
		// Only show days with activity or a sample of days if too many?
		// User specifically mentioned Mon, Tue etc. so let's show all or recent 7 if 30 is too many for vertical space.
		// However, let's stick to the 30-day view but maybe only show rows with data or a manageable set.
		
		val := data[i]
		if val == 0 {
			continue // Skip empty days to keep it concise
		}

		barWidth := 0
		if maxVal > 0 {
			barWidth = int((val / maxVal) * float64(maxBarWidth))
		}
		if barWidth == 0 && val > 0 {
			barWidth = 1
		}

		bar := barStyle.Render(strings.Repeat("█", barWidth))
		label := labelStyle.Render(labels[i])
		count := fmt.Sprintf("(%d)", int(val))

		chart += fmt.Sprintf("%-10s %-40s %s\n", label, bar, count)
	}

	if chart == "" {
		chart = "No activity found in the last 30 days."
	}

	fmt.Println(chartStyle.Render(chart))
}

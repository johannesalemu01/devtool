package ui

import (
	"fmt"

	"github.com/guptarohit/asciigraph"
	"github.com/charmbracelet/lipgloss"
)

var chartStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#00FF87")).
	Padding(1, 2)

func RenderActivityChart(data []float64, labels []string) {
	fmt.Println(HeaderStyle.Render("📈 Commit Activity (Last 30 Days)"))
	
	graph := asciigraph.Plot(data, 
		asciigraph.Height(10), 
		asciigraph.Width(60),
		asciigraph.Caption("Commits per day"),
	)

	fmt.Println(chartStyle.Render(graph))
}

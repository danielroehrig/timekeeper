package themes

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Background lipgloss.Color
	Foreground lipgloss.Color
	Accent     lipgloss.Color
	AltAccent  lipgloss.Color
}

var (
	InputStyle     = lipgloss.NewStyle()
	SubtextStyle   = lipgloss.NewStyle()
	BorderedWidget = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
)

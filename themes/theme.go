package themes

import "github.com/charmbracelet/lipgloss"

type Theme interface {
	Background() lipgloss.Color
	Foreground() lipgloss.Color
	Accent() lipgloss.Color
	AltAccent() lipgloss.Color
	Subtext() lipgloss.Color
	NormalStyle() lipgloss.Style
	SubtextStyle() lipgloss.Style
	WidgetStyle() lipgloss.Style
	ActiveWidgetStyle() lipgloss.Style
}

package themes

import (
	"github.com/charmbracelet/lipgloss"
)

type TokyoNight struct {
	background lipgloss.Color
	foreground lipgloss.Color
	accent     lipgloss.Color
	altAccent  lipgloss.Color
	subtext    lipgloss.Color
}

func NewTokyoNight() TokyoNight {
	return TokyoNight{
		background: lipgloss.Color("#1f2335"),
		foreground: lipgloss.Color("#a9b1d6"),
		accent:     lipgloss.Color("#bb9af7"),
		altAccent:  lipgloss.Color("#ff007c"),
		subtext:    lipgloss.Color("#737aa2"),
	}
}

func (t TokyoNight) Background() lipgloss.Color {
	return t.background
}

func (t TokyoNight) Foreground() lipgloss.Color {
	return t.foreground
}

func (t TokyoNight) Accent() lipgloss.Color {
	return t.accent
}

func (t TokyoNight) AltAccent() lipgloss.Color {
	return t.altAccent
}

func (t TokyoNight) Subtext() lipgloss.Color {
	return t.subtext
}

func (t TokyoNight) NormalStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.foreground).Bold(false)
}

func (t TokyoNight) SubtextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.subtext).Bold(false)
}

func (t TokyoNight) WidgetStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(t.foreground)
}

func (t TokyoNight) ActiveWidgetStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(t.accent)
}

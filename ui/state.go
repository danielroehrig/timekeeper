package ui

import tea "github.com/charmbracelet/bubbletea"

type State interface {
	KeyPressed(key string) tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	View() string
	Init() tea.Cmd
}

type StateChangeMsg struct {
	NextState *State
}

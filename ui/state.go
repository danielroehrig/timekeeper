package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type State interface {
	KeyPressed(msg tea.KeyMsg) tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	Init() tea.Cmd
}

type StateChangeMsg struct {
	NextState State
}

func ChangeState(s State) tea.Cmd {
	log.Println("change state called")
	return func() tea.Msg {
		return StateChangeMsg{NextState: s}
	}
}

package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/models"
)

type GetRunningTaskMsg struct{}

type RunningTaskModel struct {
	Entry *models.Entry
}

func (m *RunningTaskModel) KeyPressed(key tea.KeyMsg) tea.Cmd {
	return nil
}

func (m *RunningTaskModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (m *RunningTaskModel) View() string {
	return m.Entry.Name
}

func (m *RunningTaskModel) Init() tea.Cmd {
	return nil
}

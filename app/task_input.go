package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/models"
	"time"
)

type StartRunningMsg struct {
	runningTask *models.Entry
}

func (m model) handleKeypressTaskInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "enter":
		return m, func() tea.Msg {
			runningTask := &models.Entry{
				Start: time.Now(),
				End:   nil,
				Name:  m.taskEntry.Value(),
			}
			return StartRunningMsg{runningTask: runningTask}
		}
	default:
		v, cmd := m.taskEntry.Update(msg)
		m.taskEntry = v
		return m, cmd
	}
}

func (m model) taskInputView() string {
	width := (m.width / 2) - 2
	if m.focused == TaskInput {
		return borderedWidget.Width(width).BorderForeground(m.theme.Accent).Render(m.taskEntry.View())
	} else {
		return borderedWidget.Width(width).Render(m.taskEntry.View())
	}
}

package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielroehrig/timekeeper/log"
	"time"
)

type StopRunningTaskMsg struct{}

func (m model) handleKeypressTaskRunning(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.Type
	switch key {
	case tea.KeySpace:
		log.Debugf("Is space triggered?")
		return m, func() tea.Msg {
			return StopRunningTaskMsg{}
		}
	}
	return m, nil
}

func (m model) runningTaskView() string {
	elapsed := time.Since(m.runningTask.Start)
	inner := lipgloss.JoinVertical(lipgloss.Left, inputStyle.Render(m.runningTask.Name), subtextStyle.Render(elapsed.Round(time.Second).String()))
	if m.focused == TaskRunning {
		return borderedWidget.BorderForeground(m.theme.Accent).Render(inner)
	}
	return borderedWidget.Render(inner)
}

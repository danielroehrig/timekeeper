package app

import tea "github.com/charmbracelet/bubbletea"

func (m model) handleKeypressTaskList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	v, cmd := m.entryList.Update(msg)
	m.entryList = v
	return m, cmd
}

func (m model) TaskListView() string {
	width := (m.width / 2) - 2
	if m.focused == EntryList {
		return borderedWidget.BorderForeground(m.theme.Accent).Width(width).Render(m.entryList.View())
	}
	return borderedWidget.Width(width).Render(m.entryList.View())
}

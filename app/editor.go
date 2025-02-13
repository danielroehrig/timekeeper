package app

import tea "github.com/charmbracelet/bubbletea"

func (m model) handleKeypressEditor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	v, cmd := m.description.Update(msg)
	m.description = v
	return m, cmd
}

func (m model) EditorView() string {
	height := m.height - 4
	width := (m.width / 2) - 1
	if m.focused == Editor {
		return borderedWidget.Width(width).Height(height).BorderForeground(m.theme.Accent).Render(m.description.View())
	}
	return borderedWidget.Width(width).Height(height).Render(m.description.View())
}

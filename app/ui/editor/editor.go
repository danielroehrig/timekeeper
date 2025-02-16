package editor

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateEditorContentsMessage struct {
	value string
}

type Model struct {
	content textarea.Model
}

func New() Model {
	t := textarea.New()
	t.Focus()
	return Model{
		content: t,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypressEditor(msg)
	}
	return m, nil
}
func (m Model) View() string {
	return m.content.View()
}
func (m Model) handleKeypressEditor(msg tea.KeyMsg) (Model, tea.Cmd) {
	v, cmd := m.content.Update(msg)
	m.content = v
	updateCmd := func() tea.Msg {
		return UpdateEditorContentsMessage{
			value: v.Value(),
		}
	}
	return m, tea.Batch(updateCmd, cmd)
}

func (m Model) EditorView() string {
	return m.content.View()
}

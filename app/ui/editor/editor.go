package editor

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/log"
	"github.com/danielroehrig/timekeeper/models"
)

type EntryEditedMsg struct {
	Entry *models.Entry
}
type EntryListSelectedMsg struct {
	Entry *models.Entry
}

type Model struct {
	content textarea.Model
	entry   *models.Entry
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
	case EntryListSelectedMsg:
		m.entry = msg.Entry
		log.Debugf("Update the editor with %s", msg.Entry.Content)
		m.content.SetValue(msg.Entry.Content)
		return m, nil

	}
	return m, nil
}
func (m Model) View() string {
	return m.content.View()
}

func (m Model) StatusBar() string {
	return "<tab> current task"
}

func (m Model) handleKeypressEditor(msg tea.KeyMsg) (Model, tea.Cmd) {
	v, _ := m.content.Update(msg)
	m.content = v
	m.entry.Content = v.Value()
	return m, func() tea.Msg {
		return EntryEditedMsg{
			Entry: m.entry,
		}
	}
}

func (m Model) EditorView() string {
	return m.content.View()
}

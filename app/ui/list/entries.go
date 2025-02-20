package list

import (
	bl "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/danielroehrig/timekeeper/themes"
)

type Model struct {
	list  bl.Model
	theme themes.Theme
}

type EntriesLoadedMsg struct {
	Entries []*models.Entry
}

type EntrySelectedMsg struct{}

type AddEntryMsg struct {
	Entry *models.Entry
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New(theme themes.Theme) Model {
	delegates := NewEntryListDelegate(theme)
	entryList := bl.New(nil, delegates, 40, 10)
	return Model{
		list:  entryList,
		theme: theme,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypressTaskList(msg)
	case EntriesLoadedMsg:
		m.list = convertEntriesToList(msg.Entries, m.theme)
		return m, nil
	case AddEntryMsg:
		return m, m.list.InsertItem(0, msg.Entry)
	case bl.FilterMatchesMsg:
		m.list, _ = m.list.Update(msg)
	}
	return m, nil
}
func (m Model) View() string {
	return m.list.View()
}

func (m Model) handleKeypressTaskList(msg tea.KeyMsg) (Model, tea.Cmd) {
	v, cmd := m.list.Update(msg)
	m.list = v
	if msg.Type == tea.KeyEnter {
		return m, tea.Batch(cmd, func() tea.Msg {
			return EntrySelectedMsg{}
		})
	}
	return m, cmd
}

func convertEntriesToList(entries []*models.Entry, theme themes.Theme) bl.Model {
	listEntries := make([]bl.Item, 0, len(entries))
	for _, entry := range entries {
		listEntries = append(listEntries, entry)
	}
	m := bl.New(listEntries, NewEntryListDelegate(theme), 40, 20)
	m.SetShowStatusBar(false)
	m.SetShowTitle(false)
	return m
}

func (m Model) StatusBar() string {
	return "see list"
}

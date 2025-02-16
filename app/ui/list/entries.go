package list

import (
	bl "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/models"
)

type Model struct {
	list bl.Model
}

type EntriesLoadedMsg struct {
	Entries []*models.Entry
}

type AddEntryMsg struct {
	Entry *models.Entry
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() Model {
	entries := make([]bl.Item, 0)
	entryList := bl.New(entries, EntryListDelegate{}, 40, 10)
	return Model{
		list: entryList,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypressTaskList(msg)
	case EntriesLoadedMsg:
		m.list = convertEntriesToList(msg.Entries)
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
	return m, cmd
}

func convertEntriesToList(entries []*models.Entry) bl.Model {
	listEntries := make([]bl.Item, 0, len(entries))
	for _, entry := range entries {
		listEntries = append(listEntries, entry)
	}
	return bl.New(listEntries, EntryListDelegate{}, 40, 10)
}

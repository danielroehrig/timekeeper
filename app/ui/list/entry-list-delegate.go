package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/danielroehrig/timekeeper/themes"
	"io"
	"time"
)

type EntryListDelegate struct {
	theme themes.Theme
}

func NewEntryListDelegate(theme themes.Theme) EntryListDelegate {
	return EntryListDelegate{
		theme: theme,
	}
}

type EntryChangedMsg struct {
	SelectedEntry *models.Entry
}

func (d EntryListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	e, ok := item.(*models.Entry)
	if !ok {
		return
	}
	dur := 0 * time.Second
	if e.End != nil {
		dur = e.End.Sub(e.Start)
	}
	now := time.Now()
	var dateString string

	if e.Start.YearDay() == now.YearDay() && e.Start.Year() == now.Year() {
		dateString = "today"
	} else if e.Start.YearDay() == now.YearDay()-1 && e.Start.Year() == now.Year() {
		dateString = "yesterday"
	} else {
		dateString = e.Start.Format("2006-01-02")
	}
	var endTime string
	if e.End != nil {
		endTime = e.End.Format("15:04")
	} else {
		endTime = "unknown"
	}
	dateString = d.theme.SubtextStyle().Render(fmt.Sprintf("%s %s - %s: %s", dateString, e.Start.Format("15:04"), endTime, dur.Round(time.Minute)))
	var taskString string
	if m.Index() == index {
		taskString = d.theme.AccentStyle().Render(e.Name)
	} else {
		taskString = d.theme.NormalStyle().Render(e.Name)
	}

	fmt.Fprintf(w, "%s\n%s", dateString, taskString)
}

func (d EntryListDelegate) Height() int {
	return 2
}

func (d EntryListDelegate) Spacing() int {
	return 1
}

func (d EntryListDelegate) Update(_ tea.Msg, m *list.Model) tea.Cmd {
	return func() tea.Msg {
		return EntryChangedMsg{
			SelectedEntry: m.SelectedItem().(*models.Entry),
		}
	}
}

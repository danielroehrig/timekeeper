package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/entries"
	"io"
)

type EntryListDelegate struct{}

func (d EntryListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	e, ok := item.(*entries.Entry)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, e.Name)
	fmt.Fprint(w, str)
}

func (d EntryListDelegate) Height() int {
	return 1
}

func (d EntryListDelegate) Spacing() int {
	return 0
}

func (d EntryListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

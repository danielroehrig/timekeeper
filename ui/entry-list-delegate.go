package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/log"
	"github.com/danielroehrig/timekeeper/models"
	"io"
	"time"
)

type EntryListDelegate struct{}

func (d EntryListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	e, ok := item.(*models.Entry)
	if !ok {
		return
	}
	dur := 0 * time.Second
	if e.End != nil {
		dur = e.End.Sub(e.Start)
	}
	str := fmt.Sprintf("%d. %s - %s", index+1, e.Name, dur.Round(time.Second))
	fmt.Fprint(w, str)
}

func (d EntryListDelegate) Height() int {
	return 1
}

func (d EntryListDelegate) Spacing() int {
	return 0
}

func (d EntryListDelegate) Update(_ tea.Msg, m *list.Model) tea.Cmd {
	log.Debugf("list update %d", m.Index())
	return nil
}

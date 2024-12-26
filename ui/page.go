package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/log"
)

type Page interface {
	KeyPressed(msg tea.KeyMsg) tea.Cmd
	tea.Model
}

type PageChangeMsg struct {
	NextPage Page
}

func ChangePage(s Page) tea.Cmd {
	log.Debugf("change page called")
	return func() tea.Msg {
		return PageChangeMsg{NextPage: s}
	}
}

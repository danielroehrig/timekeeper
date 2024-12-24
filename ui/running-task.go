package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/log"
	"github.com/danielroehrig/timekeeper/models"
	"time"
)

type TickMsg struct{}

type RunningTaskModel struct {
	Entry     *models.Entry
	Running   bool
	MainModel tea.Model
}

func (m *RunningTaskModel) KeyPressed(key tea.KeyMsg) tea.Cmd {
	log.Debugf("KeyPressed: %s", key)
	switch key.String() {
	case "ctrl+c":
		return tea.Quit
	}
	return nil
}

func (m *RunningTaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debugf("Updated received: %s", msg)
	return m.MainModel.Update(msg)
}

func (m *RunningTaskModel) View() string {
	log.Debugf("Running Task View Called")
	now := time.Now()
	return fmt.Sprintf("%s\n%f", m.Entry.Name, now.Sub(m.Entry.Start).Seconds())
}

func (m *RunningTaskModel) Init() tea.Cmd {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				log.Debugf("Tick")
				m.Update(TickMsg{})
			}
		}
	}()
	return nil
}

func (m *RunningTaskModel) String() string {
	return "Running Task"
}

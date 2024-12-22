package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/models"
	"log"
	"time"
)

type TickMsg struct{}

type RunningTaskModel struct {
	Entry   *models.Entry
	Running bool
}

func (m *RunningTaskModel) KeyPressed(key tea.KeyMsg) tea.Cmd {
	return nil
}

func (m *RunningTaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *RunningTaskModel) View() string {
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
				log.Println("herrruu")
				m.Update(TickMsg{})
			}
		}
	}()
	return nil
}

func (m *RunningTaskModel) String() string {
	return "Running Task"
}

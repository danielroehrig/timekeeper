package task

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/danielroehrig/timekeeper/themes"
	"time"
)

type StartRunningMsg struct {
	RunningTask *models.Entry
}
type StopRunningTaskMsg struct{}
type EditRunningTaskMsg struct{}

type state byte

const (
	input state = iota
	running
)

type Model struct {
	state       state
	task        textinput.Model
	runningTask *models.Entry
	focused     bool
	width       int
	theme       themes.Theme
}

func New(theme themes.Theme) Model {
	m := Model{
		state:       input,
		task:        textinput.New(),
		runningTask: nil,
		focused:     true,
		width:       10,    // might be needed to tweak max input characters or placeholder message
		theme:       theme, // might be needed to style inner components
	}
	m.task.Placeholder = "Tell me what you are doing"
	m.task.Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypressTaskInput(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case cursor.BlinkMsg:
		m.task, cmd = m.task.Update(msg)
		return m, cmd
	case StartRunningMsg:
		m.state = running
		m.runningTask = msg.RunningTask
		return m, nil
	case StopRunningTaskMsg:
		m.state = input
		m.task.Reset()
		m.runningTask = nil
		return m, nil
	}
	return m, nil
}

func (m Model) handleKeypressTaskInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.Type
	if m.state == input {
		switch key {
		case tea.KeyEnter:
			return m, func() tea.Msg {
				runningTask := &models.Entry{
					Start: time.Now(),
					End:   nil,
					Name:  m.task.Value(),
				}
				return StartRunningMsg{RunningTask: runningTask}
			}
		default:
			v, cmd := m.task.Update(msg)
			m.task = v
			return m, cmd
		}

	} else {
		switch key {
		case tea.KeySpace:
			return m, func() tea.Msg {
				return StopRunningTaskMsg{}
			}
		case tea.KeyEnter:
			return m, func() tea.Msg {
				return EditRunningTaskMsg{}
			}

		}
		return m, nil
	}
}

func (m Model) View() string {
	if m.state == input {
		return m.task.View()
	} else {
		elapsed := time.Since(m.runningTask.Start).Round(time.Second).String()
		return lipgloss.JoinVertical(lipgloss.Left, themes.InputStyle.Render(m.runningTask.Name), themes.SubtextStyle.Render(elapsed))
	}
}

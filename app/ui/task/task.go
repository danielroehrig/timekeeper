package task

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/danielroehrig/timekeeper/themes"
	"time"
)

type StartRunningMsg struct {
	RunningTask *models.Entry
}

type Model struct {
	task    textinput.Model
	focused bool
	width   int
	theme   themes.Theme
}

func New(theme themes.Theme) Model {
	m := Model{
		task:    textinput.New(),
		focused: true,
		width:   10,    // might be needed to tweak max input characters or placeholder message
		theme:   theme, // might be needed to style inner components
	}
	m.task.Placeholder = "Tell me what you are doing"
	m.task.Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypressTaskInput(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.task, cmd = m.task.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) handleKeypressTaskInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "enter":
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
}

func (m Model) View() string {
	return m.task.View()
}

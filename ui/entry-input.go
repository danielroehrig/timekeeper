package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	input textinput.Model
}

type AddEntryMsg struct {
	Description string
}

func NewInputModel() *InputModel {
	entryText := textinput.New()
	entryText.Placeholder = "What are you doing right now?"
	entryText.Focus()

	return &InputModel{
		input: entryText,
	}
}

func (i *InputModel) KeyPressed(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()
	switch key {
	case "ctrl+c":
		return tea.Quit
	case "enter":
		return func() tea.Msg {
			return AddEntryMsg{Description: i.input.Value()}
		}
	}
	m, cmd := i.input.Update(msg)
	i.input = m
	return cmd
}

func (i *InputModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	i.input, cmd = i.input.Update(msg)
	return cmd
}

func (i *InputModel) View() string {
	return i.input.View()
}

func (i *InputModel) Init() tea.Cmd {
	return textinput.Blink
}

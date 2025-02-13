package app

import (
	"github.com/danielroehrig/timekeeper/app/ui"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dbaccess "github.com/danielroehrig/timekeeper/db"
	"github.com/danielroehrig/timekeeper/log"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/danielroehrig/timekeeper/themes"
	"github.com/ostafen/clover/v2"
)

type Focused byte

const (
	TaskInput Focused = iota
	TaskRunning
	EntryList
	Editor
)

type model struct {
	db          *clover.DB
	focused     Focused
	runningTask *models.Entry
	stopwatch   stopwatch.Model
	taskEntry   textinput.Model
	entryList   list.Model
	description textarea.Model
	theme       themes.Theme
	width       int
	height      int
}

type EntriesLoadedMsg struct {
	entries []*models.Entry
}

type AddEntryMsg struct{}
type EntryAddedMsg struct{}

type (
	NextFocusMsg struct{}
	PrevFocusMsg struct{}
)

var (
	inputStyle     = lipgloss.NewStyle()
	subtextStyle   = lipgloss.NewStyle()
	borderedWidget = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
)

func initialModel(db *clover.DB) model {
	entryText := textinput.New()
	entryText.Placeholder = "What are you doing right now?"
	entryText.Focus()

	entries := make([]list.Item, 0)
	entryList := list.New(entries, ui.EntryListDelegate{}, 40, 10)
	theme := themes.TokyoNight
	inputStyle = inputStyle.Bold(false).Foreground(theme.Accent)

	return model{
		db:          db,
		focused:     TaskInput,
		taskEntry:   entryText,
		stopwatch:   stopwatch.New(),
		entryList:   entryList,
		description: textarea.New(),
		theme:       themes.TokyoNight,
		width:       10,
		height:      10,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(loadEntries(m.db), m.stopwatch.Init(), m.stopwatch.Start(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Debugf("main update %v", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)
	case EntriesLoadedMsg:
		log.Debugf("Received entries from database")
		m.entryList = convertEntriesToList(msg.entries)
		return m, nil
	case StartRunningMsg:
		m.runningTask = msg.runningTask
		m.description.Focus()
		m.focused = Editor
	case EntryAddedMsg:
		m.runningTask = nil
		m.taskEntry.Reset()
		m.taskEntry.Focus()
		m.focused = TaskInput
	case stopwatch.TickMsg:
		log.Debugf("Tick Message received")
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	case stopwatch.StartStopMsg:
		log.Debugf("Start Stop Message")
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	case NextFocusMsg:
		log.Debugf("Next Focus Message received: %d", m.focused)
		switch m.focused {
		case TaskInput:
			m.entryList.FilterInput.Focus()
			m.focused = EntryList
		case TaskRunning:
			m.entryList.FilterInput.Focus()
			m.focused = EntryList
		case EntryList:
			m.description.Focus()
			m.focused = Editor
		case Editor:
			if m.runningTask == nil {
				m.taskEntry.Focus()
				m.focused = TaskInput
			} else {
				m.taskEntry.Focus()
				m.focused = TaskRunning
			}
		}
	case StopRunningTaskMsg:
		log.Debugf("Stop Running Task Message")
		taskEnd := time.Now()
		m.runningTask.End = &taskEnd
		return m, func() tea.Msg {
			return AddEntryMsg{}
		}
	case AddEntryMsg:
		err := dbaccess.AddEntry(m.db, m.runningTask)
		if err != nil {
			log.Errorf("Error adding entry: %v", err)
		}
		m.entryList.InsertItem(0, m.runningTask)
		return m, func() tea.Msg {
			return EntryAddedMsg{}
		}
	case tea.WindowSizeMsg:
		log.Debugf("Window Size Changed")
		m.width, m.height = msg.Width, msg.Height
	case list.FilterMatchesMsg:
		log.Debugf("Filter Matches Message")
		m.entryList, _ = m.entryList.Update(msg)
	}
	return m, cmd
}

func (m model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	log.Debugf("new key msg %s and focused was %v", msg.String(), m.focused)
	key := msg.String()
	switch key {
	case "ctrl+c":
		return m, tea.Quit
	case "tab":
		return m, func() tea.Msg {
			return NextFocusMsg{}
		}
	}
	switch m.focused {
	case TaskInput:
		return m.handleKeypressTaskInput(msg)
	case TaskRunning:
		return m.handleKeypressTaskRunning(msg)
	case Editor:
		return m.handleKeypressEditor(msg)
	case EntryList:
		return m.handleKeypressTaskList(msg)
	default:
		log.Debugf("no handle for focus: %v", m.focused)
		return m, nil
	}
}

func (m model) View() string {
	log.Debugf("main view called")
	headline := lipgloss.NewStyle().Bold(true).Foreground(m.theme.AltAccent).PaddingLeft(2).PaddingTop(1).MarginBottom(1)

	var t string
	if m.runningTask == nil {
		t = m.taskInputView()
	} else {
		t = m.runningTaskView()
	}

	s := lipgloss.JoinVertical(lipgloss.Top, headline.Render("Timekeeper"),
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.JoinVertical(lipgloss.Top, t, m.TaskListView()),
			m.EditorView()))
	return s
}

func Run(db *clover.DB) error {
	p := tea.NewProgram(initialModel(db), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func loadEntries(db *clover.DB) tea.Cmd {
	log.Infof("Loading entries...")
	return func() tea.Msg {
		loadedEntries := dbaccess.LoadEntries(db)
		return EntriesLoadedMsg{entries: loadedEntries}
	}
}

func convertEntriesToList(entries []*models.Entry) list.Model {
	listEntries := make([]list.Item, 0, len(entries))
	for _, entry := range entries {
		listEntries = append(listEntries, entry)
	}
	return list.New(listEntries, ui.EntryListDelegate{}, 40, 10)
}

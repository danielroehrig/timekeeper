package app

import (
	"github.com/danielroehrig/timekeeper/app/ui/editor"
	l "github.com/danielroehrig/timekeeper/app/ui/list"
	"github.com/danielroehrig/timekeeper/app/ui/task"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
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
	Task Focused = iota
	EntryList
	Editor
)

type model struct {
	db          *clover.DB
	focused     Focused
	runningTask *models.Entry
	dirtyTask   *models.Entry
	stopwatch   stopwatch.Model
	task        task.Model
	entryList   l.Model
	editor      editor.Model
	theme       themes.Theme
	width       int
	height      int
}

type AddEntryMsg struct {
	Entry *models.Entry
}
type EntryAddedMsg struct{}
type NextFocusMsg struct{}

func initialModel(db *clover.DB) model {
	theme := themes.TokyoNight
	return model{
		db:        db,
		focused:   Task,
		task:      task.New(theme),
		stopwatch: stopwatch.New(),
		entryList: l.New(),
		editor:    editor.New(),
		theme:     themes.TokyoNight,
		width:     10,
		height:    10,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(loadEntries(m.db), m.stopwatch.Init(), m.stopwatch.Start(), cursor.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	//log.Debugf("main update %v", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)
	case l.EntriesLoadedMsg:
		log.Debugf("Received entries from database")
		m.entryList, cmd = m.entryList.Update(msg)
		return m, cmd
	case task.StartRunningMsg:
		log.Debugf("Starting running task: %v", msg)
		m.runningTask = msg.RunningTask
		m.editor, _ = m.editor.Update(editor.EntryListSelectedMsg{Entry: msg.RunningTask})
		m.focused = Editor
		m.task, cmd = m.task.Update(msg)
		return m, cmd
	case EntryAddedMsg:
		m.runningTask = nil
		m.focused = Task
	case stopwatch.TickMsg:
		//log.Debugf("Tick Message received")
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	case stopwatch.StartStopMsg:
		log.Debugf("Start Stop Message")
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		// todo only the active input should have the blink animation
		var dc, te tea.Cmd
		m.editor, dc = m.editor.Update(msg)
		m.task, te = m.task.Update(msg)
		return m, tea.Batch(dc, te)
	case NextFocusMsg:
		log.Debugf("Next Focus Message received: %d", m.focused)
		m.saveChanges()
		switch m.focused {
		case Task:
			m.focused = EntryList
		case EntryList:
			if m.runningTask != nil {
				m.editor, cmd = m.editor.Update(editor.EntryListSelectedMsg{Entry: m.runningTask})
			}
			m.focused = Task
		case Editor:
			if m.runningTask != nil {
				m.editor, cmd = m.editor.Update(editor.EntryListSelectedMsg{Entry: m.runningTask})
			}
			m.focused = Task
		}
	case task.StopRunningTaskMsg:
		log.Debugf("Stop Running Task Message")
		taskEnd := time.Now()
		m.runningTask.End = &taskEnd
		m.task, cmd = m.task.Update(msg)
		return m, tea.Batch(cmd, func() tea.Msg {
			return AddEntryMsg{
				Entry: m.runningTask,
			}
		})
	case AddEntryMsg:
		log.Debugf("Add Entry Message: %v", msg)
		err := dbaccess.AddEntry(m.db, m.runningTask)
		if err != nil {
			log.Errorf("Error adding entry: %v", err)
		}
		m.entryList, _ = m.entryList.Update(l.AddEntryMsg{Entry: msg.Entry})
		return m, func() tea.Msg {
			return EntryAddedMsg{}
		}
	case tea.WindowSizeMsg:
		log.Debugf("Window Size Changed")
		m.width, m.height = msg.Width, msg.Height
		m.task, _ = m.task.Update(msg)
	case list.FilterMatchesMsg:
		log.Debugf("Filter Matches Message")
		m.entryList, _ = m.entryList.Update(msg)
	case editor.EntryEditedMsg:
		if msg.Entry != m.runningTask {
			log.Debugf("replacing entry: %v", msg.Entry)
			m.dirtyTask = msg.Entry
		}
	case l.EntryChangedMsg:
		log.Debugf("Select Entry Message")
		m.saveChanges()
		m.editor, _ = m.editor.Update(editor.EntryListSelectedMsg{Entry: msg.SelectedEntry})
		return m, nil
	case l.EntrySelectedMsg:
		log.Debugf("Edit Entry Message")
		m.saveChanges()
		m.focused = Editor
	case task.EditRunningTaskMsg:
		m.saveChanges()
		if m.runningTask != nil {
			m.editor, cmd = m.editor.Update(editor.EntryListSelectedMsg{Entry: m.runningTask})
		}
		m.focused = Editor
	}
	return m, cmd
}

// todo make async
func (m *model) saveChanges() tea.Cmd {
	// Todo error handling, messages, everything
	if m.dirtyTask != nil {
		log.Debugf("Saving changes to database")
		dbaccess.UpdateEntry(m.db, m.dirtyTask)
		m.dirtyTask = nil
	}
	return nil
}

func (m model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	//log.Debugf("new key msg %s and focused was %v", msg.String(), m.focused)
	key := msg.String()
	switch key {
	case "ctrl+c":
		m.saveChanges()
		return m, tea.Quit
	case "tab":
		return m, func() tea.Msg {
			return NextFocusMsg{}
		}
	}
	switch m.focused {
	case Task:
		tm, cmd := m.task.Update(msg)
		m.task = tm
		return m, cmd
	case Editor:
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		return m, cmd
	case EntryList:
		el, cmd := m.entryList.Update(msg)
		m.entryList = el
		return m, cmd
	default:
		log.Debugf("no handle for focus: %v", m.focused)
		return m, nil
	}
}

func (m model) View() string {
	//log.Debugf("main view called")
	headline := lipgloss.NewStyle().Bold(true).Foreground(m.theme.AltAccent).PaddingLeft(2).PaddingTop(1).MarginBottom(1)
	leftWidth := (m.width / 2) - 2

	var t, li, e string
	t = themes.BorderedWidget.Width(leftWidth).Render(m.task.View())
	li = themes.BorderedWidget.Width(leftWidth).Render(m.entryList.View())
	e = themes.BorderedWidget.Width(leftWidth).Render(m.editor.View())

	switch m.focused {
	case Task:
		t = themes.BorderedWidget.BorderForeground(m.theme.Accent).Width(leftWidth).Render(m.task.View())
	case EntryList:
		li = themes.BorderedWidget.BorderForeground(m.theme.Accent).Width(leftWidth).Render(m.entryList.View())
	case Editor:
		e = themes.BorderedWidget.BorderForeground(m.theme.Accent).Width(leftWidth).Render(m.editor.View())
	}

	s := lipgloss.JoinVertical(lipgloss.Top, headline.Render("Timekeeper"),
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.JoinVertical(lipgloss.Top, t, li),
			e))
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
		return l.EntriesLoadedMsg{Entries: loadedEntries}
	}
}

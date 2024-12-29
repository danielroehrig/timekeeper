package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	dbaccess "github.com/danielroehrig/timekeeper/db"
	"github.com/danielroehrig/timekeeper/log"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/danielroehrig/timekeeper/themes"
	"github.com/danielroehrig/timekeeper/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ostafen/clover"
	"github.com/spf13/viper"
)

type Focused byte

const (
	TaskInput Focused = iota
	EntryList
	Editor
)

type model struct {
	focused     Focused
	runningTask *models.Entry
	stopwatch   stopwatch.Model
	taskEntry   textinput.Model
	entryList   list.Model
	description textarea.Model
	theme       themes.Theme
}

type EntriesLoadedMsg struct {
	entries []*models.Entry
}

type AddEntryMsg struct {
	description string
}
type EntryAddedMsg struct {
	entry *models.Entry
}

var entryText = textinput.New()

func initialModel() model {
	entryText.Placeholder = "What are you doing right now?"
	entryText.Focus()

	entries := make([]list.Item, 0)

	return model{
		focused:     TaskInput,
		taskEntry:   entryText,
		stopwatch:   stopwatch.New(),
		entryList:   list.New(entries, ui.EntryListDelegate{}, 40, 10),
		description: textarea.New(),
		theme:       themes.TokyoNight,
	}
}

var db *clover.DB

func (m model) Init() tea.Cmd {
	return tea.Sequence(loadEntries(db), m.stopwatch.Init(), m.stopwatch.Start(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Debugf("main update %v", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Debugf("new key msg %s and focused was %v", msg.String(), m.focused)
		key := msg.String()
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.focused {
			case TaskInput:
				return m, func() tea.Msg {
					return AddEntryMsg{description: m.taskEntry.Value()}
				}
			default:
				log.Debugf("Currently Focused: %v", m.focused)
				return m, nil
			}
		}
		switch m.focused {
		case TaskInput:
			v, cmd := m.taskEntry.Update(msg)
			m.taskEntry = v
			return m, cmd
		}
	case AddEntryMsg:
		return m, addNewEntryToDatabase(db, msg.description)
	case EntriesLoadedMsg:
		log.Debugf("Received entries from database")
		m.entryList = convertEntriesToList(msg.entries)
		return m, nil
	case EntryAddedMsg:
		m.runningTask = msg.entry
		m.entryList.InsertItem(0, msg.entry)
		m.description.Focus()
		m.focused = Editor
	case stopwatch.TickMsg:
		log.Debugf("Tick Message received")
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	case stopwatch.StartStopMsg:
		log.Debugf("Start Stop Message")
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	}
	return m, cmd
}

func (m model) View() string {
	log.Debugf("main view called")
	headline := lipgloss.NewStyle().Bold(true).Foreground(m.theme.AltAccent).PaddingLeft(2).PaddingTop(1).MarginBottom(1)
	inputStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent).MarginBottom(0).Border(lipgloss.RoundedBorder())

	var t string
	if m.runningTask != nil {
		t = inputStyle.Render("Some Running Task")
	} else {
		t = inputStyle.Render(m.taskEntry.View())
	}

	s := lipgloss.JoinVertical(lipgloss.Top, headline.Render("Timekeeper"),
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.JoinVertical(lipgloss.Top, t, inputStyle.Render(m.entryList.View())),
			inputStyle.Render(m.description.View())))
	s += m.stopwatch.View()
	return s
}

func main() {
	switch strings.ToLower(os.Getenv("LOGLEVEL")) {
	case "debug":
		log.SetLogLevel(log.LevelDebug)
	case "info":
		log.SetLogLevel(log.LevelInfo)
	case "warn":
		log.SetLogLevel(log.LevelWarn)
	case "error":
		log.SetLogLevel(log.LevelError)
	}
	f, err := tea.LogToFile(path.Join(os.TempDir(), "timekeeper.log"), "")
	if err != nil {
		log.Errorf("Failed to open log file: %v", err)
	}
	defer f.Close()

	loadConfig()
	db = dbaccess.OpenDatabase()

	defer dbaccess.CloseDatabase(db)

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Errorf("Error running program: %v", err)
	}
}

func loadConfig() {
	log.Infof("Loading configuration...")
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Errorf("could not find config dir. Aborting. %s", err)
	}
	configFile := filepath.Join(configDir, "timekeeper", "config.yml")
	err = os.Mkdir(filepath.Dir(configFile), 0755)
	if err != nil && !os.IsExist(err) {
		log.Errorf("could not create config folder %v", err)
	}
	viper.SetConfigFile(configFile)
	viper.SetDefault("someValue", "foobar")
	viper.Set("foo", "bar")
	err = viper.WriteConfig()
	if err != nil {
		log.Errorf("could not write to config %v", err)
	}
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

func addNewEntryToDatabase(db *clover.DB, description string) tea.Cmd {
	log.Debugf("Adding new entry...")
	return func() tea.Msg {
		e := &models.Entry{
			Start: time.Now(),
			End:   nil,
			Name:  description,
		}
		err := dbaccess.AddEntry(db, e)
		if err != nil {
			log.Errorf("Could not write to database: %v", err)
		}
		return EntryAddedMsg{entry: e}
	}
}

package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	dbaccess "github.com/danielroehrig/timekeeper/db"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/danielroehrig/timekeeper/themes"
	"github.com/danielroehrig/timekeeper/ui"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ostafen/clover"
	"github.com/spf13/viper"
)

type model struct {
	currentState  ui.State
	taskIsRunning bool
	entryList     list.Model
	theme         themes.Theme
}

type EntriesLoadedMsg struct {
	entries []*models.Entry
}

type EntryAddedMsg struct {
	entry *models.Entry
}

func initialModel() model {
	entryText := ui.NewInputModel()

	entries := make([]list.Item, 0)
	return model{
		currentState: entryText,
		entryList:    list.New(entries, ui.EntryListDelegate{}, 40, 10),
		theme:        themes.TokyoNight,
	}
}

var style = lipgloss.
	NewStyle().
	Bold(true).
	PaddingTop(2).
	PaddingLeft(4)

var db *clover.DB

func (m model) Init() tea.Cmd {
	return tea.Sequence(loadEntries(db), m.currentState.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.currentState.KeyPressed(msg)
	case ui.AddEntryMsg:
		return m, addNewEntryToDatabase(db, msg.Description)
	case EntriesLoadedMsg:
		log.Println("Received entries from database")
		m.entryList = convertEntriesToList(msg.entries)
		return m, nil
	case EntryAddedMsg:
		m.entryList.InsertItem(0, msg.entry)
		s := &ui.RunningTaskModel{Entry: msg.entry}
		return m, ui.ChangeState(s)
	case ui.StateChangeMsg:
		log.Println("Received state change")
		m.currentState = msg.NextState
		return m, nil
	}
	m.currentState.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var headline = lipgloss.NewStyle().Bold(true).Foreground(m.theme.AltAccent).PaddingLeft(2).PaddingTop(1).MarginBottom(1)
	var inputStyle = lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent).MarginBottom(2)

	s := headline.Render("Timekeeper")
	s += fmt.Sprintf("\n%s", inputStyle.Render(m.currentState.View()))
	s += fmt.Sprintf("\n%s", inputStyle.Render(m.entryList.View()))
	return s
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	loadConfig()
	db = dbaccess.OpenDatabase()

	defer dbaccess.CloseDatabase(db)
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func loadConfig() {
	log.Println("Loading configuration...")
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("could not find config dir. Aborting. %s", err)
	}
	configFile := filepath.Join(configDir, "timekeeper", "config.yml")
	err = os.Mkdir(filepath.Dir(configFile), 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("could not create config folder %v", err)
	}
	viper.SetConfigFile(configFile)
	viper.SetDefault("someValue", "foobar")
	viper.Set("foo", "bar")
	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("could not write to config %v", err)
	}
}

func loadEntries(db *clover.DB) tea.Cmd {
	log.Println("Loading entries...")
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
	log.Println("Adding new entry...")
	return func() tea.Msg {
		e := &models.Entry{
			Start: time.Now(),
			End:   nil,
			Name:  description,
		}
		err := dbaccess.AddEntry(db, e)
		if err != nil {
			log.Fatalf("Could not write to database: %v", err)
		}
		return EntryAddedMsg{entry: e}
	}

}

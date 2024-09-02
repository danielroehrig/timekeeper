package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
	entryInput    textinput.Model
	taskIsRunning bool
	entryList     list.Model
	entries       []*models.Entry
	theme         themes.Theme
}

func initialModel() model {
	entryText := textinput.New()
	entryText.Placeholder = "What are you doing right now?"
	entryText.Focus()
	loadedEntries := dbaccess.LoadEntries(db)
	listEntries := make([]list.Item, 0, len(loadedEntries))
	for _, entry := range loadedEntries {
		listEntries = append(listEntries, entry)
	}

	return model{
		entryInput: entryText,
		entryList:  list.New(listEntries, ui.EntryListDelegate{}, 40, 10),
		theme:      themes.TokyoNight,
	}
}

var style = lipgloss.
	NewStyle().
	Bold(true).
	PaddingTop(2).
	PaddingLeft(4)

var db *clover.DB

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			e := &models.Entry{
				Start: time.Now(),
				End:   nil,
				Name:  m.entryInput.Value(),
			}
			err := dbaccess.AddEntry(db, e)
			if err != nil {
				log.Fatalf("Could not write to database: %v", err)
			}
			m.entries = append(m.entries, e)
			m.entryList.InsertItem(0, e)
			m.entryInput.Reset()
			return m, nil
		}
	}
	m.entryInput, cmd = m.entryInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var headline = lipgloss.NewStyle().Bold(true).Foreground(m.theme.AltAccent).PaddingLeft(2).PaddingTop(1).MarginBottom(1)
	var inputStyle = lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent).MarginBottom(2)

	s := headline.Render("Timekeeper")
	s += fmt.Sprintf("\n%s", inputStyle.Render(m.entryInput.View()))
	s += fmt.Sprintf("\n%s", inputStyle.Render(m.entryList.View()))
	return s
}

func main() {
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

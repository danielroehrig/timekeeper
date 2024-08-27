package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/danielroehrig/timekeeper/entries"
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
	entryInput textinput.Model
	entryList  list.Model
	entries    []*entries.Entry
	theme      themes.Theme
}

func initialModel() model {
	entryText := textinput.New()
	entryText.Placeholder = "What are you doing right now?"
	entryText.Focus()
	entrs := loadEntries(db)
	return model{
		entryInput: entryText,
		entryList:  list.New(entrs, ui.EntryListDelegate{}, 40, 10),
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
			doc := clover.NewDocument()
			doc.Set("name", m.entryInput.Value())
			doc.Set("start", time.Now())
			doc.Set("end", nil)
			_, err := db.InsertOne("entries", doc)
			if err != nil {
				log.Fatalf("Could not write to database: %v", err)
			}
		}
	}
	m.entryInput, cmd = m.entryInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var headline = lipgloss.NewStyle().Bold(true).Foreground(m.theme.AltAccent).Border(lipgloss.RoundedBorder())
	var inputStyle = lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent).Border(lipgloss.RoundedBorder())

	s := headline.Render("Timekeeper")
	s += fmt.Sprintf("\n%s", inputStyle.Render(m.entryInput.View()))
	s += fmt.Sprintf("\n%s", inputStyle.Render(m.entryList.View()))
	return s
}

func main() {
	loadConfig()
	db = openDatabase()

	defer closeDatabase(db)
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func openDatabase() *clover.DB {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("could not find config dir. Aborting. %s", err)
	}
	db, err := clover.Open(filepath.Join(configDir, "timekeeper"))
	if err != nil {
		log.Fatalf("could not open database. Aborting. %s", err)
	}

	hasEntriesCollection, err := db.HasCollection("entries")
	if err != nil {
		log.Fatalf("could not check if there are entries collection. Aborting. %s", err)
	}
	if !hasEntriesCollection {
		err = db.CreateCollection("entries")
		if err != nil {
			log.Fatalf("could not create collection. Aborting. %s", err)
		}
	}
	return db
}

func loadEntries(db *clover.DB) []list.Item {
	docs, err := db.Query("entries").FindAll()
	if err != nil {
		log.Fatalf("could not list entries. Aborting. %s", err)
	}
	entry := &struct {
		Name  string    `clover:"name"`
		End   time.Time `clover:"end"`
		Start time.Time `clover:"start"`
	}{}
	items := make([]list.Item, 0, len(docs))
	for _, doc := range docs {
		doc.Unmarshal(entry)
		items = append(items, &entries.Entry{
			Start: entry.Start,
			End:   entry.End,
			Name:  entry.Name,
		})
	}
	return items
}

func closeDatabase(db *clover.DB) {
	db.ExportCollection("entries", "entries.json")
	log.Printf("closing database file")
	err := db.Close()
	if err != nil {
		log.Fatalf("could not close db. %s", err)
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

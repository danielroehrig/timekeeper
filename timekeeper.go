package main

import (
	"fmt"
	"github.com/danielroehrig/timekeeper/themes"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ostafen/clover/v2"
	"github.com/spf13/viper"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	theme    themes.Theme
}

func initialModel() model {
	return model{
		choices:  []string{"test 1", "test2"},
		selected: make(map[int]struct{}),
		theme:    themes.TokyoNight,
	}
}

var style = lipgloss.
	NewStyle().
	Bold(true).
	PaddingTop(2).
	PaddingLeft(4)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var newStyle = lipgloss.NewStyle().Foreground(m.theme.Foreground).Inherit(style)
	var headline = lipgloss.NewStyle().Bold(true).Foreground(m.theme.AltAccent).Border(lipgloss.RoundedBorder())

	s := headline.Render("Timekeeper")
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += newStyle.Render(fmt.Sprintf("\n%s [%s] %s\n", cursor, checked, choice))
	}

	s += "\nPress q\n"
	return s
}

func main() {
	loadConfig()
	db := openDatabase()
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
	return db
}

func closeDatabase(db *clover.DB) {
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

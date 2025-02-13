package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielroehrig/timekeeper/app"
	dbaccess "github.com/danielroehrig/timekeeper/db"
	"github.com/danielroehrig/timekeeper/log"
	"github.com/ostafen/clover/v2"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var db *clover.DB

func main() {
	// set up loggin
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

	// log configs
	loadConfig()

	// set up database access
	db = dbaccess.OpenDatabase()
	defer dbaccess.CloseDatabase(db)

	// run the app
	if err := app.Run(db); err != nil {
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

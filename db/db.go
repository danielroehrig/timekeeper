package db

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/ostafen/clover"
	"log"
	"os"
	"path/filepath"
	"time"
)

func OpenDatabase() *clover.DB {
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

func LoadEntries(db *clover.DB) []list.Item {
	docs, err := db.Query("entries").FindAll()
	if err != nil {
		log.Fatalf("could not list entries. Aborting. %s", err)
	}
	entry := &struct {
		Name  string     `clover:"name"`
		End   *time.Time `clover:"end"`
		Start time.Time  `clover:"start"`
	}{}
	items := make([]list.Item, 0, len(docs))
	for _, doc := range docs {
		doc.Unmarshal(entry)
		items = append(items, &models.Entry{
			Start: entry.Start,
			End:   entry.End,
			Name:  entry.Name,
		})
	}
	return items
}

func CloseDatabase(db *clover.DB) {
	db.ExportCollection("entries", "entries.json")
	log.Printf("closing database file")
	err := db.Close()
	if err != nil {
		log.Fatalf("could not close db. %s", err)
	}
}

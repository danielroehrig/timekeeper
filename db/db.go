package db

import (
	"fmt"
	"github.com/danielroehrig/timekeeper/log"
	"github.com/danielroehrig/timekeeper/models"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	"os"
	"path"
	"path/filepath"
	"time"
)

const collectionName = "entries"

type entry struct {
	ObjectId string
	Name     string     `clover:"name"`
	End      *time.Time `clover:"end"`
	Start    time.Time  `clover:"start"`
	Content  string     `clover:"content"`
}

func OpenDatabase() *clover.DB {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Errorf("could not find config dir. Aborting. %s", err)
	}
	db, err := clover.Open(filepath.Join(configDir, "timekeeper"))
	if err != nil {
		log.Errorf("could not open database. Aborting. %s", err)
	}

	hasEntriesCollection, err := db.HasCollection(collectionName)
	if err != nil {
		log.Errorf("could not check if there are entries collection. Aborting. %s", err)
	}
	if !hasEntriesCollection {
		err = db.CreateCollection(collectionName)
		if err != nil {
			log.Errorf("could not create collection. Aborting. %s", err)
		}
	}
	return db
}

func LoadEntries(db *clover.DB) []*models.Entry {
	docs, err := db.FindAll(query.NewQuery(collectionName))
	if err != nil {
		log.Errorf("could not list entries. Aborting. %s", err)
	}
	items := make([]*models.Entry, 0, len(docs))
	for _, doc := range docs {
		entry, err := unmarshallDoc(doc)
		if err != nil {
			log.Errorf("loading entries failed: %s", err)
		}
		items = append(items, entry)
	}
	return items
}

func AddEntry(db *clover.DB, e *models.Entry) error {
	doc := document.NewDocument()
	doc.Set("name", e.Name)
	doc.Set("start", e.Start)
	doc.Set("end", nil)
	doc.Set("content", e.Content)
	id, err := db.InsertOne("entries", doc)
	if err != nil {
		return fmt.Errorf("could not write to database: %w", err)
	}
	e.ObjectId = id
	return nil
}

func CloseDatabase(db *clover.DB) {
	log.Infof("closing database file")
	db.ExportCollection(collectionName, path.Join(os.TempDir(), "timekeeper.json"))
	err := db.Close()
	if err != nil {
		log.Errorf("could not close db. %s", err)
	}
}

func GetRunning(db *clover.DB) (*models.Entry, error) {
	entries, err := db.FindAll(query.NewQuery(collectionName).Where(query.Field("end").IsNil()))
	if err != nil {
		log.Errorf("could not list entries. Aborting. %s", err)
	}
	if len(entries) > 1 {
		return nil, fmt.Errorf("more than one running tasks")
	}
	return unmarshallDoc(entries[0])
}

func ExportEntries(db *clover.DB) error {
	return db.ExportCollection(collectionName, "entries.json")
}

func UpdateEntry(db *clover.DB, e *models.Entry) error {
	return db.UpdateById(collectionName, e.ObjectId, func(doc *document.Document) *document.Document {
		doc.Set("name", e.Name)
		doc.Set("start", e.Start)
		doc.Set("end", e.End)
		doc.Set("content", e.Content)
		return doc
	})
}

func unmarshallDoc(doc *document.Document) (*models.Entry, error) {
	entry := &entry{}
	err := doc.Unmarshal(entry)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal document. %s", err)
	}
	doc.ObjectId()
	return &models.Entry{
		ObjectId: doc.ObjectId(),
		Start:    entry.Start,
		End:      entry.End,
		Name:     entry.Name,
		Content:  entry.Content,
	}, nil
}

package models

import (
	"time"
)

type Entry struct {
	ObjectId string
	Start    time.Time
	End      *time.Time
	Name     string
	Content  string
}

func (e *Entry) FilterValue() string {
	return e.Name
}

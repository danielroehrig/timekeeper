package entries

import "time"

type Entry struct {
	Start time.Time
	End   time.Time
	Name  string
}

func (e *Entry) FilterValue() string {
	return e.Name
}

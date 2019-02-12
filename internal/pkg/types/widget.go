package types

import (
	"fmt"
	"time"
)

// Widget concre
type Widget struct {
	ID        string    // 16
	Source    string    // 16 byte memory allloc on 64 addr space
	CreatedAt time.Time // 24
	Broken    bool      // 1 byte memory allloc on 64 addr space
}

// NewWidget public Widget constructor
func NewWidget(source string, broken bool) *Widget {
	return &Widget{
		ID:        fmt.Sprintf("%v", time.Now().UnixNano()),
		Source:    source,
		CreatedAt: time.Now(),
		Broken:    broken,
	}
}

// String
func (w *Widget) String() string {
	return fmt.Sprintf("id=%v source=%s time=%s broken=%v", w.ID, w.Source, w.CreatedAt.Format("11:06:39.12340"), w.Broken)
}

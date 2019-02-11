package types

import (
	"fmt"
	"time"
)

// Widget concre
type Widget struct {
	id        string    // 16
	source    string    // 16 byte memory allloc on 64 addr space
	createdAt time.Time // 24
	broken    bool      // 1 byte memory allloc on 64 addr space
}

// NewWidget public Widget constructor
func NewWidget(source string, broken bool) Widget {
	return Widget{
		id:        fmt.Sprintf("%v", time.Now().UnixNano()),
		source:    source,
		createdAt: time.Now(),
		broken:    broken,
	}
}

// String
func (w Widget) String() string {
	return fmt.Sprintf("id=%v source=%s time=%s broken=%v", w.id, w.source, w.createdAt.Format("11:06:39.12340"), w.broken)
}

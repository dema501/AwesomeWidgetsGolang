package types

import (
	"fmt"
	"math/rand"
	"time"
)

// Producer can produce
type Producer struct {
	ID        string
	CreatedAt time.Time
}

// NewProducer constructor func
func NewProducer(id int) *Producer {
	return &Producer{
		ID:        fmt.Sprintf("producer-%d", id),
		CreatedAt: time.Now(),
	}
}

// Produce implement interface
func (p *Producer) Produce(produceCount, brokenCount int, dataCh chan<- Widget) error {
	var brokenPositions map[int]bool

	if brokenCount > 0 {
		brokenPositions = p.getBrokenWidgetsPosMap(produceCount, brokenCount)
	}

	for i := 0; i < produceCount; i++ {
		broken, ok := brokenPositions[i]
		if !ok {
			broken = false
		}
		dataCh <- NewWidget(p.ID, broken)
	}
	return nil
}

func (p *Producer) getBrokenWidgetsPosMap(produceCount, brokenCount int) map[int]bool {
	brokenPositions := make(map[int]bool)

	for i := 0; i < brokenCount; i++ {
		if produceCount > brokenCount {
			for {
				rnd := rand.Intn(produceCount)
				if _, ok := brokenPositions[rnd]; !ok {
					brokenPositions[rnd] = true
					break
				}
			}
		} else {
			brokenPositions[i] = true
		}
	}

	return brokenPositions
}

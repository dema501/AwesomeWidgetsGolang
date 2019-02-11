package types

import (
	"fmt"
	"sync"
	"time"
)

// Consumer knows how to consume
type Consumer struct {
	ID        string
	CreatedAt time.Time
}

// NewConsumer constructor func
func NewConsumer(id int) *Consumer {
	return &Consumer{
		ID:        fmt.Sprintf("consumer-%d", id),
		CreatedAt: time.Now(),
	}
}

// Consume method
func (c *Consumer) Consume(dataCh <-chan Widget, stopCh chan bool, toStopCh chan string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	for {
		// operation here is to try to exit the consumer
		// goroutine as early as possible.
		select {
		case <-stopCh:
			return
		default:
		}

		select {
		case <-stopCh:
			return
		// ok will be false if dataCh has been closed
		case w, ok := <-dataCh:
			if ok {
				fmt.Printf("%v consumes [%v] in %v time\n", c.ID, w, time.Since(w.createdAt))
				toStopCh <- c.ID
				continue
			} else {
				return
			}
		}
	}
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/dema501/AwesomeWidgetsGolang/internal/pkg/consumer"
	fl "github.com/dema501/AwesomeWidgetsGolang/internal/pkg/flag"
	"github.com/dema501/AwesomeWidgetsGolang/internal/pkg/producer"
	"github.com/dema501/AwesomeWidgetsGolang/internal/pkg/types"
	"github.com/pkg/errors"
)

func string2int(p *int) func(value string) error {
	return func(value string) error {
		parsedVal, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "Error occur during ParseInt of value: %v", value)
		}

		*p = int(parsedVal)
		return nil
	}
}

func moderator(total int, stopCh chan bool, toStop chan string) {
	for {
		// Same as the sender goroutine, the try-receive
		// operation here is to try to exit the receiver
		// goroutine as early as possible.
		select {
		case <-stopCh:
			return
		default:
		}

		select {
		case stBy, ok := <-toStop:
			if ok {
				total--

				if total == 0 {
					fmt.Println("last consumer was:", stBy)
					close(stopCh)

					return
				}
			}
		}
	}
}

func main() {
	// Default values
	cfg := types.Config{
		WidgetsCount:       10,
		ProducersCount:     1,
		ConsumersCount:     1,
		BrokenWidgetsCount: -1,
	}

	// Custom flag module inspired by flag https://golang.org/pkg/flag/
	// just to meet requirements, please read README.md
	var f fl.FlagSet

	f.Val("n", "Sets the number of widgets created (integer)", string2int(&cfg.WidgetsCount))
	f.Val("p", "Sets the number of producers created (integer)", string2int(&cfg.ProducersCount))
	f.Val("c", "Sets the number of consumers created (integer)", string2int(&cfg.ConsumersCount))
	f.Val("k", "Sets the `k`th widget to be broken (integer)", string2int(&cfg.BrokenWidgetsCount))

	// Parse Arguments and show an error if it's exists
	// ignore ErrFlagNotDefined errors
	if err := f.Parse(os.Args[1:]); err != nil {
		switch errors.Cause(err) {
		case fl.ErrFlagNotDefined:
			{
				fmt.Printf("[WARN] %s", err)
			}
		case fl.ErrHelp:
			{
				os.Exit(0) // success
			}
		default:
			{
				fmt.Printf("[ERROR] '%s'\n", err)
				f.DefaultUsage()
				os.Exit(2)
			}
		}
	}

	rand.Seed(time.Now().UnixNano())

	// dataCh is widget channel
	dataCh := make(chan *types.Widget, cfg.ProducersCount)

	// stopCh is an additional signal channel.
	// Its sender is the moderator goroutine shown below.
	// Its receivers are all producers and consumers of dataCh.
	stopCh := make(chan bool)

	// The channel toStop is used to notify the moderator
	// to close the additional signal channel (stopCh).
	toStop := make(chan string, 1)

	// moderator knows how to stop consumers
	go moderator(cfg.ProducersCount*cfg.WidgetsCount, stopCh, toStop)

	// senders
	for i := 0; i < cfg.ProducersCount; i++ {
		p := producer.NewProducer(i)
		go p.Produce(cfg.WidgetsCount, cfg.BrokenWidgetsCount, dataCh)
	}

	// receivers
	wgConsumers := sync.WaitGroup{}

	// We know in advance how many concurrent consumers to run
	wgConsumers.Add(cfg.ConsumersCount)

	for i := 0; i < cfg.ConsumersCount; i++ {
		c := consumer.NewConsumer(i)
		//c.Consume will decrease wgConsumers on defer
		go c.Consume(dataCh, stopCh, toStop, &wgConsumers)
	}

	// wait when all concurrent consumers will finish their job
	wgConsumers.Wait()
}

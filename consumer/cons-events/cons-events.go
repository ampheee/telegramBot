package cons_events

import (
	"github.com/ampheee/telegramBot/v2/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() {
	for {
		Events, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("ERR consumer %s", err.Error())
			continue
		}
		if len(Events) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

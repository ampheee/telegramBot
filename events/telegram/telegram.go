package telegram

import (
	"github.com/ampheee/telegramBot/v2/clients/telegram"
	"github.com/ampheee/telegramBot/v2/lib/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func newProcessor(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

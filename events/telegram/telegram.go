package telegram

import (
	"errors"
	"github.com/ampheee/telegramBot/v2/clients/telegram"
	"github.com/ampheee/telegramBot/v2/events"
	"github.com/ampheee/telegramBot/v2/lib/errs"
	"github.com/ampheee/telegramBot/v2/lib/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

func NewProcessor(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

var (
	ErrUnknownMetaType  = errors.New("unknown meta")
	ErrUnknownEventType = errors.New("unknown process")
)

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errs.Wrap("cant get event", err)
	}
	if len(updates) == 0 {
		return nil, errs.Wrap("updates are empty", err)
	}
	res := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		res = append(res, event(update))
	}
	p.offset = updates[len(updates)-1].UpdateId + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errs.Wrap("cant process event", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(e events.Event) error {
	meta, err := meta(e)
	if err != nil {
		return errs.Wrap("cant process message", err)
	}
	if err := p.doCmd(e.Text, meta.ChatID, meta.UserName); err != nil {
		return errs.Wrap("Cant process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errs.Wrap("cant get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)
	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}
	if uType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.Id,
			UserName: u.Message.From.UserName,
		}
	}
	return res
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}

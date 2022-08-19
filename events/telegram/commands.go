package telegram

import (
	"errors"
	"github.com/ampheee/telegramBot/v2/lib/errs"
	"github.com/ampheee/telegramBot/v2/lib/storage"
	"log"
	"net/url"
	"strings"
)

const (
	Random = "/random"
	Help   = "/help"
	Start  = "/start"
)

func (p *Processor) doCmd(text string, chatId int, userName string) error {
	text = strings.TrimSpace(text)

	if isAdd(text) {
		return p.savePage(chatId, text, userName)
	}

	log.Printf("new command '%s' from %s in chat %d", text, userName, chatId)
	switch text {
	case Random:
		return p.sendRandom(chatId, userName)
	case Help:
		return p.Help(chatId)
	case Start:
		return p.Hello(chatId)
	default:
		return p.tg.SendMessages(chatId, msUnknown)
	}
}

func (p *Processor) savePage(chatId int, pageUrl, username string) (err error) {
	defer func() { err = errs.Wrap("cant do savePage", err) }()
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}
	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}
	if isExist {
		return p.tg.SendMessages(chatId, msAlExist)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := p.tg.SendMessages(chatId, msSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatId int, userName string) (err error) {
	defer func() { err = errs.Wrap("cant do sendRandom", err) }()
	page, err := p.storage.PickRandom(userName)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessages(chatId, msNoSavedPages)
	}
	if err := p.tg.SendMessages(chatId, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) Help(chatId int) error {
	return p.tg.SendMessages(chatId, msHelp)
}

func (p *Processor) Hello(chatId int) error {
	return p.tg.SendMessages(chatId, msHello)
}

func isAdd(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

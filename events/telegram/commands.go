package telegram

import (
	"errors"
	"github.com/ampheee/telegramBot/v2/lib/errs"
	"github.com/ampheee/telegramBot/v2/storage"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	random = "/random"
	help   = "/help"
	start  = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case random:
		return p.sendRandom(chatID, username)
	case help:
		return p.tg.SendMessages(chatID, msHelp)
	case start:
		return p.tg.SendMessages(chatID, msHello)
	default:
		return p.tg.SendMessages(chatID, msUnknown)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = errs.Wrap("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessages(chatID, msAlExist)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessages(chatID, msSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = errs.Wrap("can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		if os.IsNotExist(err) {
			return p.tg.SendMessages(chatID, "You hadnt save link")
		} else {
			return err
		}
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessages(chatID, msNoSavedPages)
	}

	if err := p.tg.SendMessages(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

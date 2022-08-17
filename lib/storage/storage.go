package storage

import (
	"crypto/sha1"
	"fmt"
	"github.com/ampheee/telegramBot/v2/lib/errors"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userN string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", errors.Wrap("Error! Can`t calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", errors.Wrap("Error! Can`t calculate hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

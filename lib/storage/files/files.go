package files

import (
	"encoding/gob"
	"github.com/ampheee/telegramBot/v2/lib/errors"
	"github.com/ampheee/telegramBot/v2/lib/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	permission = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}
func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = errors.Wrap("cant init storage", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(filePath, permission); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}
	filePath = filepath.Join(filePath, fName)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	err = gob.NewEncoder(file).Encode(page)
	if err != nil {
		return err
	}
	return nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) PickRandom(username string) (page *storage.Page, err error) {
	filePath := filepath.Join(s.basePath, username)
	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.Wrap("NO FILES IN STORAGE", err)
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]
	return s.decode(filepath.Join(filePath, file.Name()))
}

func (s Storage) decode(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap("cant decode storage page", err)
	}
	defer func() { _ = f.Close() }()
	var p storage.Page
	if err != gob.NewDecoder(f).Decode(&p) {
		return nil, errors.Wrap("cant decode storage page", err)
	}
	return &p, err
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return errors.Wrap("cant remove page", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)
	if err := os.Remove(path); err != nil {
		return errors.Wrap("cant remove page", err)
	}
	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, errors.Wrap("cant check page exist", err)
	}
	filePath := filepath.Join(s.basePath, p.UserName, fName)
	switch _, err = os.Stat(filePath); {
	case err == os.ErrNotExist:
		return false, nil
	case err != nil:
		return false, errors.Wrap("cant check page exist", err)
	}
	return true, nil
}

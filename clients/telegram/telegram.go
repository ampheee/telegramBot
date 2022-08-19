package telegram

import (
	"encoding/json"
	"github.com/ampheee/telegramBot/v2/lib/errs"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdateMethod   = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func NewClient(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) SendMessages(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	_, err := c.doReq(sendMessageMethod, q)
	if err != nil {
		return errs.Wrap("error", err)
	}
	return err
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = errs.Wrap("ERROR: can`t do request", err) }()
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	data, err := c.doReq(getUpdateMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) doReq(method string, q url.Values) (data []byte, err error) {
	defer func() { err = errs.Wrap("ERROR: can`t do request", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

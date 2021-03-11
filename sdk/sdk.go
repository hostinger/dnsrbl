package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Request struct {
	IP      string
	Action  string
	Author  string
	Comment string
}

type Client interface {
	Allow(ctx context.Context, ip, author, comment string) error
	Block(ctx context.Context, ip, author, comment string) error
	Unblock(ctx context.Context, ip, author, comment string) error
}

type client struct {
	http *http.Client
	url  string
	key  string
}

func NewClient(key, url string) Client {
	return &client{
		http: &http.Client{
			Timeout: time.Second * 5,
		},
		url: url,
		key: key,
	}
}

func (c *client) Call(ctx context.Context, method, url string, data io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.url, url), data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating new request object")
	}
	req.Header.Add("X-API-Key", c.key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed executing request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read response body")
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("Unknown response from API: %s", string(body))
	}

	return body, nil
}

func (c *client) ExecuteAction(ctx context.Context, ip, action, author, comment string) error {
	if net.ParseIP(ip) == nil {
		return &net.ParseError{
			Type: "IPv4 Address",
			Text: ip,
		}
	}
	b := &Request{
		IP:      ip,
		Action:  action,
		Author:  author,
		Comment: comment,
	}
	body, err := json.Marshal(b)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal request into JSON")
	}
	_, err = c.Call(ctx, "POST", "addresses", bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "Failed to execute POST request")
	}
	return nil
}

func (c *client) Allow(ctx context.Context, ip, author, comment string) error {
	return c.ExecuteAction(ctx, ip, "Allow", author, comment)
}

func (c *client) Block(ctx context.Context, ip, author, comment string) error {
	return c.ExecuteAction(ctx, ip, "Block", author, comment)
}

func (c *client) Unblock(ctx context.Context, ip, author, comment string) error {
	return c.ExecuteAction(ctx, ip, "Unblock", author, comment)
}

package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hostinger/dnsrbl/pkg/dnsutils"
)

type Client struct {
	Client  *http.Client
	BaseURL string
	Key     string
}

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

type RRSet struct {
	Records    []Record `json:"records"`
	ChangeType string   `json:"changetype"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	TTL        int      `json:"ttl"`
}

type Zone struct {
	RRSets []RRSet `json:"rrsets"`
}

func NewClient(scheme, host, port, key string) (*Client, error) {
	baseURL := fmt.Sprintf("%s://%s:%s/api/v1/servers/localhost", scheme, host, port)
	return &Client{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Key:     key,
		BaseURL: baseURL,
	}, nil
}

func (c *Client) Call(ctx context.Context, uri, method string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", c.Key)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("error: %s", string(body))
	}

	return nil
}

func (c *Client) PatchZone(ctx context.Context, ip, zone, action string) error {
	if net.ParseIP(ip) == nil {
		return errors.New("argument 'ip' must be a valid IP address")
	}
	reverseIP := dnsutils.ReverseAddress(strings.Split(ip, "."))
	data := Zone{
		RRSets: []RRSet{
			{
				Name:       fmt.Sprintf("%s.%s.", reverseIP, zone),
				Type:       "A",
				TTL:        3600,
				ChangeType: action,
				Records: []Record{
					{
						Content:  "127.0.0.1",
						Disabled: false,
					},
				},
			},
		},
	}
	uri := fmt.Sprintf("%s/zones/%s", c.BaseURL, zone)
	if err := c.Call(ctx, uri, "PATCH", data); err != nil {
		return err
	}
	return nil
}

func (c *Client) Block(ctx context.Context, ip, zone string) error {
	return c.PatchZone(ctx, ip, zone, "REPLACE")
}

func (c *Client) Unblock(ctx context.Context, ip, zone string) error {
	return c.PatchZone(ctx, ip, zone, "DELETE")
}

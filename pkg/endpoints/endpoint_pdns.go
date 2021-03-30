package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hostinger/hbl/pkg/logger"
	"github.com/hostinger/hbl/pkg/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type pdnsEndpoint struct {
	l       logger.Logger
	client  *http.Client
	baseURL string
	scheme  string
	zone    string
	host    string
	port    string
	key     string
}

func NewPDNSEndpoint(l logger.Logger) Endpoint {
	l.Info("Starting execution of NewPDNSEndpoint", zap.String("endpoint", "PowerDNS"))
	c := &pdnsEndpoint{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		scheme: os.Getenv("PDNS_API_SCHEME"),
		zone:   os.Getenv("PDNS_API_ZONE"),
		host:   os.Getenv("PDNS_API_HOST"),
		port:   os.Getenv("PDNS_API_PORT"),
		key:    os.Getenv("PDNS_API_KEY"),
		l:      l,
	}
	c.baseURL = fmt.Sprintf("%s://%s:%s/api/v1/servers/localhost", c.scheme, c.host, c.port)
	l.Info("Finished execution of NewPDNSEndpoint", zap.String("endpoint", "PowerDNS"))
	return c
}

func (c *pdnsEndpoint) Call(ctx context.Context, uri, method string, code int, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		c.l.Error(
			"Failed to marshal body into JSON",
			zap.String("endpoint", "PowerDNS"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed to marshal body into JSON")
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(body))
	if err != nil {
		c.l.Error(
			"Failed to create new request object",
			zap.String("endpoint", "PowerDNS"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed creating new request object")
	}
	req.Header.Set("X-API-Key", c.key)

	resp, err := c.client.Do(req)
	if err != nil {
		c.l.Error(
			"Failed to execute request",
			zap.String("endpoint", "PowerDNS"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed executing request")
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		c.l.Error(
			"Failed to read response body",
			zap.String("endpoint", "PowerDNS"),
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "Failed to read response body")
	}

	if resp.StatusCode != code {
		c.l.Error(
			"Uknown response from PowerDNS API",
			zap.String("endpoint", "PowerDNS"),
			zap.String("error", string(body)),
		)
		return nil, errors.Errorf("Unknown response from API: %s", string(body))
	}

	return body, nil
}

func (c *pdnsEndpoint) PatchZone(ctx context.Context, ip, action string) error {
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
	reverseIP := utils.ReverseAddress(strings.Split(ip, "."))
	data := Zone{
		RRSets: []RRSet{
			{
				Name:       fmt.Sprintf("%s.%s.", reverseIP, c.zone),
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
	uri := fmt.Sprintf("%s/zones/%s", c.baseURL, c.zone)
	if _, err := c.Call(ctx, uri, "PATCH", 204, data); err != nil {
		return err
	}
	return nil
}

func (c *pdnsEndpoint) SearchZone(ctx context.Context, ip string) error {
	type SearchResult struct {
		Content    string `json:"content"`
		Disabled   bool   `json:"disabled"`
		Name       string `json:"name"`
		ObjectType string `json:"object_type"`
		TTL        int    `json:"ttl"`
		Type       string `json:"type"`
		Zone       string `json:"zone"`
		ZoneID     string `json:"zone_id"`
	}
	if net.ParseIP(ip) == nil {
		return errors.New("Argument 'IP' must be a valid IP address")
	}
	reverseIP := utils.ReverseAddress(strings.Split(ip, "."))
	uri := fmt.Sprintf("%s/search-data?q=%s.%s&object_type=%s&max=1", c.baseURL, reverseIP, c.zone, "record")

	resp, err := c.Call(ctx, uri, "GET", 200, nil)
	if err != nil {
		return err
	}

	var results []SearchResult
	if err := json.Unmarshal(resp, &results); err != nil {
		return errors.Wrap(err, "Failed to unmarshal JSON")
	}

	if len(results) == 0 {
		return errors.New("Address doesn't exist")
	}

	return nil
}

func (c *pdnsEndpoint) Name() string {
	return "PowerDNS"
}

func (c *pdnsEndpoint) Block(ctx context.Context, ip string) error {
	return c.PatchZone(ctx, ip, "REPLACE")
}

func (c *pdnsEndpoint) Unblock(ctx context.Context, ip string) error {
	return c.PatchZone(ctx, ip, "DELETE")
}

func (c *pdnsEndpoint) Exists(ctx context.Context, ip string) error {
	return c.SearchZone(ctx, ip)
}

func (c *pdnsEndpoint) Sync(ctx context.Context, ip string) error {
	if err := c.Exists(ctx, ip); err != nil {
		return c.PatchZone(ctx, ip, "REPLACE")
	}
	return nil
}

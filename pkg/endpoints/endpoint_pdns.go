package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/hostinger/dnsrbl/pkg/utils"
)

type PDNSEndpoint struct {
	Client  *http.Client
	BaseURL string
	Scheme  string `env:"PDNS_API_SCHEME,required"`
	Zone    string `env:"PDNS_API_ZONE,required"`
	Host    string `env:"PDNS_API_HOST,required"`
	Port    string `env:"PDNS_API_PORT,required"`
	Key     string `env:"PDNS_API_KEY,required"`
}

func NewPDNSEndpoint() Endpoint {
	c := &PDNSEndpoint{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	if err := env.Parse(c); err != nil {
		log.Fatalf("Endpoints: PDNSEndpoint: %s", err)
	}
	c.BaseURL = fmt.Sprintf("%s://%s:%s/api/v1/servers/localhost", c.Scheme, c.Host, c.Port)
	return c
}

func (c *PDNSEndpoint) Call(ctx context.Context, uri, method string, data interface{}) error {
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

func (c *PDNSEndpoint) PatchZone(ctx context.Context, ip, action string) error {
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
	if net.ParseIP(ip) == nil {
		return errors.New("argument 'ip' must be a valid IP address")
	}
	reverseIP := utils.ReverseAddress(strings.Split(ip, "."))
	data := Zone{
		RRSets: []RRSet{
			{
				Name:       fmt.Sprintf("%s.%s.", reverseIP, c.Zone),
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
	uri := fmt.Sprintf("%s/zones/%s", c.BaseURL, c.Zone)
	if err := c.Call(ctx, uri, "PATCH", data); err != nil {
		return err
	}
	return nil
}

func (c *PDNSEndpoint) Name() string {
	return "PowerDNS"
}

func (c *PDNSEndpoint) Block(ctx context.Context, ip string) error {
	return c.PatchZone(ctx, ip, "REPLACE")
}

func (c *PDNSEndpoint) Unblock(ctx context.Context, ip string) error {
	return c.PatchZone(ctx, ip, "DELETE")
}

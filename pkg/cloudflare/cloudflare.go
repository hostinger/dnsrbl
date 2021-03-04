package cloudflare

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudflare/cloudflare-go"
)

type Client struct {
	Client    *cloudflare.API
	AccountID string
	Email     string
	Key       string
}

func NewClient(accountID string, email string, key string) (*Client, error) {
	api, err := cloudflare.New(key, email)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:    api,
		AccountID: accountID,
		Email:     email,
		Key:       key,
	}, nil
}

func (c *Client) Block(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("Address '%s' is not a valid IP address.", ip)
	}
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	response, err := c.Client.CreateAccountAccessRule(context.Background(), c.AccountID, rule)
	if err != nil || response.Success == false {
		return err
	}
	return nil
}

func (c *Client) Unblock(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("Address '%s' is not a valid IP address.", ip)
	}
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	rules, err := c.Client.ListAccountAccessRules(context.Background(), c.AccountID, rule, 1)
	if err != nil {
		return err
	}
	if rules.Count <= 0 || rules.Count > 1 {
		return fmt.Errorf("AccessRule for IP address '%s' was not found. ", ip)
	}
	response, err := c.Client.DeleteAccountAccessRule(context.Background(), c.AccountID, rules.Result[0].ID)
	if err != nil || response.Success == false {
		return err
	}
	return nil
}

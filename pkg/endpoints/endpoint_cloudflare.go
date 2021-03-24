package endpoints

import (
	"context"
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/cloudflare/cloudflare-go"
)

type CloudflareEndpoint struct {
	Client    *cloudflare.API
	AccountID string `env:"CF_API_ACCOUNT,required"`
	Email     string `env:"CF_API_EMAIL,required"`
	Key       string `env:"CF_API_KEY,required"`
}

func NewCloudflareEndpoint() Endpoint {
	e := &CloudflareEndpoint{}
	if err := env.Parse(e); err != nil {
		log.Fatalf("Endpoints: CloudflareEndpoint: %s", err)
	}
	api, err := cloudflare.New(e.Key, e.Email)
	if err != nil {
		log.Fatalf("Endpoints: CloudflareEndpoint: %s", err)
	}
	e.Client = api
	return e
}

func (c *CloudflareEndpoint) Name() string {
	return "Cloudflare"
}

func (c *CloudflareEndpoint) Block(ctx context.Context, ip string) error {
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	response, err := c.Client.CreateAccountAccessRule(ctx, c.AccountID, rule)
	if err != nil || !response.Success {
		return err
	}
	return nil
}

func (c *CloudflareEndpoint) Unblock(ctx context.Context, ip string) error {
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	rules, err := c.Client.ListAccountAccessRules(ctx, c.AccountID, rule, 1)
	if err != nil {
		return err
	}
	if rules.Count <= 0 || rules.Count > 1 {
		return fmt.Errorf("AccessRule for IP address '%s' was not found. ", ip)
	}
	response, err := c.Client.DeleteAccountAccessRule(ctx, c.AccountID, rules.Result[0].ID)
	if err != nil || !response.Success {
		return err
	}
	return nil
}

func (c *CloudflareEndpoint) Exists(ctx context.Context, ip string) error {
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	rules, err := c.Client.ListAccountAccessRules(ctx, c.AccountID, rule, 1)
	if err != nil {
		return err
	}
	if rules.Count <= 0 || rules.Count > 1 {
		return fmt.Errorf("Address '%s' doesn't exist", ip)
	}
	return nil
}

func (c *CloudflareEndpoint) Sync(ctx context.Context, ip string) error {
	if err := c.Exists(ctx, ip); err != nil {
		return c.Block(ctx, ip)
	}
	return nil
}

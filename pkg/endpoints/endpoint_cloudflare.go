package endpoints

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"go.uber.org/zap"
)

type cloudflareEndpoint struct {
	l       *zap.Logger
	client  *cloudflare.API
	account string
	email   string
	key     string
}

func NewCloudflareEndpoint(l *zap.Logger) Endpoint {
	e := &cloudflareEndpoint{
		account: os.Getenv("CF_API_ACCOUNT"),
		email:   os.Getenv("CF_API_EMAIL"),
		key:     os.Getenv("CF_API_KEY"),
		l:       l,
	}
	api, err := cloudflare.New(e.key, e.email)
	if err != nil {
		log.Fatalf("Endpoints: CloudflareEndpoint: %s", err)
	}
	e.client = api
	return e
}

func (c *cloudflareEndpoint) Name() string {
	return "Cloudflare"
}

func (c *cloudflareEndpoint) Block(ctx context.Context, ip string) error {
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	response, err := c.client.CreateAccountAccessRule(ctx, c.account, rule)
	if err != nil || !response.Success {
		return err
	}
	return nil
}

func (c *cloudflareEndpoint) Unblock(ctx context.Context, ip string) error {
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	rules, err := c.client.ListAccountAccessRules(ctx, c.account, rule, 1)
	if err != nil {
		return err
	}
	if rules.Count <= 0 || rules.Count > 1 {
		return fmt.Errorf("AccessRule for IP address '%s' was not found. ", ip)
	}
	response, err := c.client.DeleteAccountAccessRule(ctx, c.account, rules.Result[0].ID)
	if err != nil || !response.Success {
		return err
	}
	return nil
}

func (c *cloudflareEndpoint) Exists(ctx context.Context, ip string) error {
	rule := cloudflare.AccessRule{
		Mode: "block",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip,
		},
		Notes: "Created automatically by HBL API.",
	}
	rules, err := c.client.ListAccountAccessRules(ctx, c.account, rule, 1)
	if err != nil {
		return err
	}
	if rules.Count <= 0 || rules.Count > 1 {
		return fmt.Errorf("Address '%s' doesn't exist", ip)
	}
	return nil
}

func (c *cloudflareEndpoint) Sync(ctx context.Context, ip string) error {
	if err := c.Exists(ctx, ip); err != nil {
		return c.Block(ctx, ip)
	}
	return nil
}

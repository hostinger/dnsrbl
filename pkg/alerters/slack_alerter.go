package alerters

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/slack-go/slack"
)

type SlackAlerter struct {
	URL      string `env:"SLACK_WEBHOOK_URL,required"`
	Username string `env:"SLACK_WEBHOOK_USERNAME,required"`
	Channel  string `env:"SLACK_WEBHOOK_CHANNEL,required"`
}

func NewSlackAlerter() Alerter {
	c := &SlackAlerter{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("Alerters: SlackAlerter: %s", err)
	}
	return c
}

func (s *SlackAlerter) Name() string {
	return "Slack"
}

func (s *SlackAlerter) Alert(ctx context.Context, alert *Alert) {
	message := &slack.WebhookMessage{
		Username: s.Username,
		Channel:  s.Channel,
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.SectionBlock{
					Type: "section",
					Text: &slack.TextBlockObject{
						Type: "mrkdwn",
						Text: fmt.Sprintf("Received a new *%s* action for the following address.", strings.ToTitle(alert.Action)),
					},
					Fields: []*slack.TextBlockObject{
						{
							Type: "mrkdwn",
							Text: fmt.Sprintf("*IP*\n%s", alert.IP),
						},
						{
							Type: "mrkdwn",
							Text: fmt.Sprintf("*Reason*\n%s", alert.Comment),
						},
					},
				},
			},
		},
	}
	if err := slack.PostWebhookContext(ctx, s.URL, message); err != nil {
		log.Printf("SlackAlerter: Alert: %s", err)
	}
}

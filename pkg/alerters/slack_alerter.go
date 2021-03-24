package alerters

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type slackAlerter struct {
	l        *zap.Logger
	url      string
	username string
	channel  string
}

func NewSlackAlerter(l *zap.Logger) Alerter {
	c := &slackAlerter{
		l:        l,
		url:      os.Getenv("SLACK_WEBHOOK_URL"),
		channel:  os.Getenv("SLACK_WEBHOOK_CHANNEL"),
		username: os.Getenv("SLACK_WEBHOOK_USERNAME"),
	}
	return c
}

func (s *slackAlerter) Name() string {
	return "Slack"
}

func (s *slackAlerter) Alert(ctx context.Context, alert *Alert) {
	message := &slack.WebhookMessage{
		Username: s.username,
		Channel:  s.channel,
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
	if err := slack.PostWebhookContext(ctx, s.url, message); err != nil {
		log.Printf("SlackAlerter: Alert: %s", err)
	}
}

package alerters

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hostinger/hbl/pkg/logger"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type slackAlerter struct {
	l        logger.Logger
	url      string
	username string
	channel  string
}

func NewSlackAlerter(l logger.Logger) Alerter {
	l.Info("Starting execution of NewSlackAlerter", zap.String("alerter", "Slack"))
	c := &slackAlerter{
		l:        l,
		url:      os.Getenv("SLACK_WEBHOOK_URL"),
		channel:  os.Getenv("SLACK_WEBHOOK_CHANNEL"),
		username: os.Getenv("SLACK_WEBHOOK_USERNAME"),
	}
	l.Info("Finished execution of NewSlackAlerter", zap.String("alerter", "Slack"))
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
		s.l.Error(
			"Failed to execute PostWebhookContext",
			zap.String("alerter", "Slack"),
			zap.Error(err),
		)
	}
}

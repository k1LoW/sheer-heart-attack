package cmd

import (
	"fmt"
	"time"

	slack "github.com/monochromegane/slack-incoming-webhooks"
	"github.com/spf13/cast"
	"go.uber.org/zap/zapcore"
)

func notifySlack(webhookURL string, slackChannel string, trackFields []trackField) func(zapcore.Entry) error {
	return func(e zapcore.Entry) error {
		name := "Sheer Heart Attack"
		emoji := ":bomb:"
		color := "#DBA6CC"
		var (
			prefix   string
			hostname string
		)
		switch e.Message {
		case executeMessage:
			prefix = ":boom:"
			color = "#B61972"
		case timeoutMessage:
			prefix = ":hourglass:"
		}
		payload := slack.Payload{
			Channel:   slackChannel,
			IconEmoji: emoji,
			Username:  name,
		}
		attachment := slack.Attachment{
			Title:      fmt.Sprintf("%s %s", prefix, e.Message),
			Text:       "_\"Hey, look over <!here>.\"_",
			Fallback:   e.Message,
			Color:      color,
			Timestamp:  time.Now().Unix(),
			MarkdownIn: []string{"fields"},
		}
		for _, f := range trackFields {
			switch f.key {
			case "command":
				attachment.AddField(&slack.Field{
					Title: fmt.Sprintf("--%s", f.key),
					Value: fmt.Sprintf("```%s```", cast.ToString(f.value)),
					Short: false,
				})
			case "hostname":
				hostname = cast.ToString(f.value)
			case "slack-channel":
				continue
			default:
				attachment.AddField(&slack.Field{
					Title: fmt.Sprintf("--%s", f.key),
					Value: cast.ToString(f.value),
					Short: true,
				})
			}
		}
		attachment.Footer = hostname
		payload.AddAttachment(&attachment)
		slack.Client{
			WebhookURL: webhookURL,
		}.Post(&payload)
		return nil
	}
}

package cmd

import (
	"fmt"
	"time"

	slack "github.com/monochromegane/slack-incoming-webhooks"
	"github.com/spf13/cast"
	"go.uber.org/zap/zapcore"
)

func notifySlack(webhookURL string, slackChannel string, slackMention string, trackFields []trackField) func(zapcore.Entry) error {
	text := "_\"Hey, look over here.\"_"
	if slackMention == "@here" {
		text = "_\"Hey, look over <!here>.\"_"
	} else if slackMention != "" {
		text = fmt.Sprintf("_\"Hey <%s>, look over here.\"_", slackMention)
	}
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
		case errorMessage, executeFailedMessage:
			prefix = ":bangbang:"
			color = "#C7243A"
		}
		payload := slack.Payload{
			Channel:   slackChannel,
			IconEmoji: emoji,
			Username:  name,
		}
		attachment := slack.Attachment{
			Title:      fmt.Sprintf("%s %s", prefix, e.Message),
			Text:       text,
			Fallback:   e.Message,
			Color:      color,
			Timestamp:  time.Now().Unix(),
			MarkdownIn: []string{"fields"},
		}
		for _, f := range trackFields {
			switch f.key {
			case "command":
				c := cast.ToString(f.value)
				v := fmt.Sprintf("```%s```", c)
				if c == "" {
					v = "no command"
				}
				attachment.AddField(&slack.Field{
					Title: fmt.Sprintf("--%s", f.key),
					Value: v,
					Short: false,
				})
			case "hostname":
				hostname = cast.ToString(f.value)
			case "slack-channel":
				continue
			case "pid":
				v := cast.ToString(f.value)
				if v == "0" {
					continue
				}
				attachment.AddField(&slack.Field{
					Title: fmt.Sprintf("--%s", f.key),
					Value: v,
					Short: true,
				})
			case "name":
				v := cast.ToString(f.value)
				if v == "" {
					continue
				}
				attachment.AddField(&slack.Field{
					Title: fmt.Sprintf("--%s", f.key),
					Value: v,
					Short: true,
				})
			case "log-path":
				l := cast.ToString(f.value)
				v := fmt.Sprintf("```%s```", l)
				attachment.AddField(&slack.Field{
					Title: "log path",
					Value: v,
					Short: false,
				})
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
		err := slack.Client{
			WebhookURL: webhookURL,
		}.Post(&payload)
		if err != nil {
			return err
		}
		return nil
	}
}

package options

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Songmu/prompter"
	"github.com/gobuffalo/packr/v2"
	"github.com/k1LoW/metr/metrics"
	"github.com/labstack/gommon/color"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const CollectInterval = time.Duration(500) * time.Millisecond

type Options struct {
	nonInteractive bool
	options        []string
	localizer      *i18n.Localizer
}

var langs = []language.Tag{
	language.English,
	language.Japanese,
}
var matcher = language.NewMatcher(langs)

// NewOptions ...
func NewOptions(
	pid int32,
	name string,
	threshold string,
	interval string,
	attempts int,
	commands []string,
	times int,
	timeout string,
	slackChannel string,
	slackMention string,
	nonInteractive bool,
	lang string,
) (*Options, error) {
	box := packr.New("i18n", "../i18n")
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, l := range langs {
		path := fmt.Sprintf("%s.toml", l.String())
		d, err := box.Find(path)
		if err != nil {
			return nil, err
		}
		bundle.MustParseMessageFileBytes(d, path)
	}
	matched, _, _ := matcher.Match(language.Make(lang))

	fmt.Printf("Detected language: %s\n", color.Magenta(display.English.Tags().Name(matched)))
	fmt.Println("")

	o := &Options{
		nonInteractive: nonInteractive,
		options:        []string{},
		localizer:      i18n.NewLocalizer(bundle, matched.String()),
	}
	pid, name, err := o.process(pid, name)
	if err != nil {
		return o, err
	}
	err = o.threshold(threshold, pid, name)
	if err != nil {
		return o, err
	}
	err = o.interval(interval)
	if err != nil {
		return o, err
	}
	err = o.attempts(attempts)
	if err != nil {
		return o, err
	}
	err = o.command(commands, true)
	if err != nil {
		return o, err
	}
	err = o.times(times)
	if err != nil {
		return o, err
	}
	err = o.timeout(timeout)
	if err != nil {
		return o, err
	}
	slackChannel, err = o.slackChannel(slackChannel)
	if err != nil {
		return o, err
	}
	if slackChannel != "" {
		err := o.slackMention(slackMention)
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

// Get ...
func (o *Options) Get() []string {
	return o.options
}

func (o *Options) process(pid int32, name string) (int32, string, error) {
	if pid > 0 && name != "" {
		return pid, name, errors.New("you can only use either --pid or --name")
	}
	if o.nonInteractive {
		if pid > 0 {
			o.options = append(o.options, []string{"--pid", strconv.Itoa(int(pid))}...)
		}
		if name != "" {
			o.options = append(o.options, []string{"--name", name}...)
		}
		return pid, name, nil
	}

	processStr := strconv.Itoa(int(pid))
	if processStr == "0" {
		processStr = name
	}

	fmt.Printf("%s ... %s\n", color.Magenta("--pid", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionPID"}))
	fmt.Printf("%s ... %s\n", color.Magenta("--name", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionName"}))
	fmt.Println("")

	processStr = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "processPromptMessage"}), processStr)

	pidInt64, err := strconv.ParseInt(processStr, 10, 32)
	if err == nil {
		pid = int32(pidInt64)
		p, err := process.NewProcess(pid)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			return o.process(pid, name)
		}
		name, err = p.Name()
		if err != nil || name == "" {
			_, _ = fmt.Fprintf(os.Stderr, "No process found: %d\n", pid)
			return o.process(pid, name)
		}

		fmt.Printf("Target process name: %s\n", color.Magenta(name))
		fmt.Println("")
		o.options = append(o.options, []string{"--pid", processStr}...)
		return pid, "", nil
	}

	if processStr != "" {
		fmt.Printf("Target process name: %s\n", color.Magenta(processStr))
		fmt.Println("")
		o.options = append(o.options, []string{"--pid", processStr}...)
		return 0, processStr, nil
	}

	fmt.Println(color.Magenta("Track only host metrics"))
	fmt.Println("")
	return 0, "", nil
}

func (o *Options) threshold(threshold string, pid int32, name string) error {
	if o.nonInteractive {
		o.options = append(o.options, []string{"--threshold", threshold}...)
		return nil
	}
	var (
		m   *metrics.Metrics
		err error
	)

	switch {
	case pid > 0:
		m, err = metrics.GetMetrics(CollectInterval, pid)
		if err != nil {
			return err
		}
	case name != "":
		m, err = metrics.GetMetricsByName(CollectInterval, name)
		if err != nil {
			return err
		}
	default:
		m, err = metrics.GetMetrics(CollectInterval)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s ... %s\n", color.Magenta("--threshold", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionThreshold"}))
	fmt.Println("")
	fmt.Printf("%s\n", color.Magenta(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "availableMetrics"}), color.B))

	m.Each(func(metric metrics.Metric, value interface{}) {
		fmt.Printf("  %s (now:%s %s): %s\n", color.White(metric.Name), color.Magenta(fmt.Sprintf(metric.Format, value)), metric.Unit, metric.Description)
	})

	fmt.Printf("%s\n", color.Magenta(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "availableOperators"}), color.B))
	fmt.Printf("  %s\n", "+, -, *, /, ==, !=, <, >, <=, >=, not, and, or, !, &&, ||")
	fmt.Println("")
	threshold = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "thresholdPromptMessage"}), threshold)
	fmt.Println("")
	o.options = append(o.options, []string{"--threshold", threshold}...)
	return nil
}

func (o *Options) interval(interval string) error {
	intervalStr := interval
	if o.nonInteractive {
		o.options = append(o.options, []string{"--interval", intervalStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--interval", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionInterval"}))
	fmt.Println("")
	intervalStr = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "intervalPromptMessage"}), intervalStr)
	fmt.Println("")

	o.options = append(o.options, []string{"--interval", intervalStr}...)
	return nil
}

func (o *Options) attempts(attempts int) error {
	attemptsStr := strconv.Itoa(attempts)
	if o.nonInteractive {
		o.options = append(o.options, []string{"--attempts", attemptsStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--attempts", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionAttempts"}))
	fmt.Println("")
	attemptsStr = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "attemptsPromptMessage"}), attemptsStr)
	fmt.Println("")
	o.options = append(o.options, []string{"--attempts", attemptsStr}...)
	return nil
}

func (o *Options) command(commands []string, first bool) error {
	if o.nonInteractive {
		if len(commands) > 0 {
			for _, c := range commands {
				o.options = append(o.options, []string{"--command", c}...)
			}
		}
		return nil
	}
	if first {
		fmt.Printf("%s ... %s\n", color.Magenta("--command", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionCommand"}))
		fmt.Println("")
		fmt.Printf("%s\n", color.White(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "additionalEnvironmentVariables"}), color.B))
		fmt.Printf("  %s: %s\n", color.White("$PID", color.B), "PID of the process.")
		fmt.Println("")
	}
	if len(commands) > 0 {
		fmt.Println("")
		fmt.Printf("%s\n", color.Magenta("Execution commands"))
		for _, c := range commands {
			fmt.Printf("%s %s\n", color.Magenta(">"), color.White(c))
		}
		fmt.Println("")
	}

	command := prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandPromptMessage"}), "")
	fmt.Println("")
	if command != "" {
		commands = append(commands, command)
	}
	if len(commands) == 0 {
		return nil
	}
	yn := prompter.YN(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandYNMessage"}), false)
	if yn {
		return o.command(commands, false)
	}
	for _, c := range commands {
		o.options = append(o.options, []string{"--command", c}...)
	}
	return nil
}

func (o *Options) times(times int) error {
	timesStr := strconv.Itoa(times)
	if o.nonInteractive {
		o.options = append(o.options, []string{"--times", timesStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--times", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionTimes"}))
	fmt.Println("")
	timesStr = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "timesPromptMessage"}), strconv.Itoa(times))
	fmt.Println("")
	o.options = append(o.options, []string{"--times", timesStr}...)
	return nil
}

func (o *Options) timeout(timeout string) error {
	timeoutStr := timeout
	if o.nonInteractive {
		o.options = append(o.options, []string{"--timeout", timeoutStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--timeout", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionTimeout"}))
	fmt.Println("")
	timeoutStr = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "timeoutPromptMessage"}), timeoutStr)
	fmt.Println("")
	o.options = append(o.options, []string{"--timeout", timeoutStr}...)
	return nil
}

func (o *Options) slackChannel(slackChannel string) (string, error) {
	if o.nonInteractive {
		if slackChannel != "" {
			o.options = append(o.options, []string{"--slack-channel", slackChannel}...)
		}
		return slackChannel, nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--slack-channel", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionSlackChannel"}))
	fmt.Println("")
	url, urlErr := GetEnvSlackIncommingWebhook()
	if urlErr == nil {
		fmt.Printf("%s: %s\n", "Slack Incomming Webhook URL", color.Magenta(url))
		fmt.Println("")
	}
	yn := prompter.YN(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "slackChannelYNMessage"}), true)
	if !yn {
		fmt.Println("")
		return "", nil
	}
	if urlErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", urlErr)
		url = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "slackWebhookPromptMessage"}), "")
		if url == "" {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", errors.New("invalid URL"))
			return o.slackChannel(slackChannel)
		}
		err := os.Setenv("SLACK_INCOMMING_WEBHOOK_URL", url)
		if err != nil {
			return "", err
		}
	}
	slackChannel = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "slackChannelPromptMessage"}), slackChannel)
	fmt.Println("")
	o.options = append(o.options, []string{"--slack-channel", slackChannel}...)
	return slackChannel, nil
}

func (o *Options) slackMention(slackMention string) error {
	if o.nonInteractive {
		if slackMention == "" {
			return nil
		} else {
			o.options = append(o.options, []string{"--slack-mention", slackMention}...)
			return nil
		}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--slack-mention", color.B), o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "optionSlackMention"}))
	fmt.Println("")
	yn := prompter.YN(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "slackMentionYNMessage"}), true)
	if !yn {
		fmt.Println("")
		return nil
	}
	slackMention = prompter.Prompt(o.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "slackMentionPromptMessage"}), slackMention)
	fmt.Println("")
	o.options = append(o.options, []string{"--slack-mention", slackMention}...)
	return nil
}

// GetEnvSlackIncommingWebhook return slack incomming webhook URL via os.Envirion
func GetEnvSlackIncommingWebhook() (string, error) {
	envKeys := []string{
		"SLACK_INCOMMING_WEBHOOK_URL",
		"SLACK_WEBHOOK_URL",
		"SLACK_URL",
	}
	for _, key := range envKeys {
		if url := os.Getenv(key); url != "" {
			return url, nil
		}
	}
	return "", fmt.Errorf("slack incomming webhook url environment variables are not found %s", envKeys)
}

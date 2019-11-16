package options

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Songmu/prompter"
	"github.com/k1LoW/metr/metrics"
	"github.com/labstack/gommon/color"
	"github.com/shirou/gopsutil/process"
)

const CollectInterval = time.Duration(500) * time.Millisecond

type Options struct {
	nonInteractive bool
	options        []string
}

// NewOptions ...
func NewOptions(
	pid int32,
	name string,
	threshold string,
	interval int,
	attempts int,
	command string,
	times int,
	timeout int,
	slackChannel string,
	slackMention string,
	nonInteractive bool,
) (*Options, error) {
	o := &Options{
		nonInteractive: nonInteractive,
		options:        []string{},
	}
	pid, name, err := o.Process(pid, name)
	if err != nil {
		return o, err
	}
	err = o.Threshold(threshold, pid, name)
	if err != nil {
		return o, err
	}
	err = o.Interval(interval)
	if err != nil {
		return o, err
	}
	err = o.Attempts(attempts)
	if err != nil {
		return o, err
	}
	err = o.Command(command)
	if err != nil {
		return o, err
	}
	err = o.Times(times)
	if err != nil {
		return o, err
	}
	err = o.Timeout(timeout)
	if err != nil {
		return o, err
	}
	slackChannel, err = o.SlackChannel(slackChannel)
	if err != nil {
		return o, err
	}
	if slackChannel != "" {
		err := o.SlackMention(slackMention)
		if err != nil {
			return o, err
		}
	}
	return o, nil

}

func (o *Options) Get() []string {
	return o.options
}

func (o *Options) Process(pid int32, name string) (int32, string, error) {
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

	fmt.Printf("%s ... %s\n", color.Magenta("--pid", color.B), "PID of the process.")
	fmt.Printf("%s ... %s\n", color.Magenta("--name", color.B), "name of the process.")
	fmt.Println("")

	processStr = prompter.Prompt("Enter PID or name of the process (If empty, sheer-heart-atack track only host metrics)", processStr)

	pidInt64, err := strconv.ParseInt(processStr, 10, 32)
	if err == nil {
		pid = int32(pidInt64)
		p, err := process.NewProcess(pid)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			return o.Process(pid, name)
		}
		name, err = p.Name()
		if err != nil || name == "" {
			_, _ = fmt.Fprintf(os.Stderr, "No process found: %d\n", pid)
			return o.Process(pid, name)
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

func (o *Options) Threshold(threshold string, pid int32, name string) error {
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

	fmt.Printf("%s ... %s\n", color.Magenta("--threshold", color.B), "Threshold conditions.")
	fmt.Println("")
	fmt.Printf("%s\n", color.Magenta("Available Metrics", color.B))

	m.Each(func(metric metrics.Metric, value interface{}) {
		fmt.Printf("  %s (now:%s %s): %s\n", color.White(metric.Name), color.Magenta(fmt.Sprintf(metric.Format, value)), metric.Unit, metric.Description)
	})

	fmt.Printf("%s\n", color.Magenta("Available Operators", color.B))
	fmt.Printf("  %s\n", "+, -, *, /, ==, !=, <, >, <=, >=, not, and, or, !, &&, ||")
	fmt.Println("")
	threshold = prompter.Prompt("Enter threshold", threshold)
	fmt.Println("")
	o.options = append(o.options, []string{"--threshold", threshold}...)
	return nil
}

func (o *Options) Interval(interval int) error {
	intervalStr := strconv.Itoa(interval)
	if o.nonInteractive {
		o.options = append(o.options, []string{"--interval", intervalStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--interval", color.B), "Interval of checking if the threshold exceeded (seconds).")
	fmt.Println("")
	intervalStr = prompter.Prompt("Enter interval", intervalStr)
	fmt.Println("")
	o.options = append(o.options, []string{"--interval", intervalStr}...)
	return nil
}

func (o *Options) Attempts(attempts int) error {
	attemptsStr := strconv.Itoa(attempts)
	if o.nonInteractive {
		o.options = append(o.options, []string{"--attempts", attemptsStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--attempts", color.B), "Maximum number of attempts continuously exceeding the threshold.")
	fmt.Println("")
	attemptsStr = prompter.Prompt("Enter attempts", attemptsStr)
	fmt.Println("")
	o.options = append(o.options, []string{"--attempts", attemptsStr}...)
	return nil
}

func (o *Options) Command(command string) error {
	if o.nonInteractive {
		if command != "" {
			o.options = append(o.options, []string{"--command", command}...)
		}
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--command", color.B), "Command to execute when the maximum number of attempts is exceeded.execution. If command execution time > 'interval' * 3, kill command.")
	fmt.Println("")
	fmt.Printf("%s\n", color.White("Additional Environment Variables", color.B))
	fmt.Printf("  %s: %s\n", color.White("$PID", color.B), "PID of the process.")
	fmt.Println("")
	command = prompter.Prompt("Enter command", command)
	fmt.Println("")
	if command != "" {
		o.options = append(o.options, []string{"--command", command}...)
	}
	return nil
}

func (o *Options) Times(times int) error {
	timesStr := strconv.Itoa(times)
	if o.nonInteractive {
		o.options = append(o.options, []string{"--times", timesStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--times", color.B), "Maximum number of command executions. If times < 1, track and execute until timeout.")
	fmt.Println("")
	timesStr = prompter.Prompt("Enter times", strconv.Itoa(times))
	fmt.Println("")
	o.options = append(o.options, []string{"--times", timesStr}...)
	return nil
}

func (o *Options) Timeout(timeout int) error {
	timeoutStr := strconv.Itoa(timeout)
	if o.nonInteractive {
		o.options = append(o.options, []string{"--timeout", timeoutStr}...)
		return nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--timeout", color.B), "Timeout of tracking (seconds).")
	fmt.Println("")
	timeoutStr = prompter.Prompt("Enter timeout", timeoutStr)
	fmt.Println("")
	o.options = append(o.options, []string{"--timeout", timeoutStr}...)
	return nil
}

func (o *Options) SlackChannel(slackChannel string) (string, error) {
	if o.nonInteractive {
		if slackChannel != "" {
			o.options = append(o.options, []string{"--slack-channel", slackChannel}...)
		}
		return slackChannel, nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--slack-channel", color.B), "Slack channel to notify.")
	fmt.Println("")
	url, urlErr := GetEnvSlackIncommingWebhook()
	if urlErr == nil {
		fmt.Printf("%s: %s\n", "Slack Incomming Webhook URL", color.Magenta(url))
		fmt.Println("")
	}
	yn := prompter.YN("Do you want to notify slack channel?", true)
	if !yn {
		fmt.Println("")
		return "", nil
	}
	if urlErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", urlErr)
		url = prompter.Prompt("Enter slack incoming webhook URL", "")
		if url == "" {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", errors.New("invalid URL"))
			return o.SlackChannel(slackChannel)
		}
		err := os.Setenv("SLACK_INCOMMING_WEBHOOK_URL", url)
		if err != nil {
			return "", err
		}
	}
	slackChannel = prompter.Prompt("Enter slack channel", slackChannel)
	fmt.Println("")
	o.options = append(o.options, []string{"--slack-channel", slackChannel}...)
	return slackChannel, nil
}

func (o *Options) SlackMention(slackMention string) error {
	if o.nonInteractive {
		if slackMention == "" {
			return nil
		} else {
			o.options = append(o.options, []string{"--slack-mention", slackMention}...)
			return nil
		}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--slack-mention", color.B), "Slack mention.")
	fmt.Println("")
	yn := prompter.YN("Do you want to mention?", true)
	if !yn {
		fmt.Println("")
		return nil
	}
	slackMention = prompter.Prompt("Enter mention [@here or user_id (ex. @UXXXXXXXX)]", slackMention)
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

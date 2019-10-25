package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Songmu/prompter"
	"github.com/k1LoW/metr/metrics"
	"github.com/labstack/gommon/color"
	"github.com/shirou/gopsutil/process"
)

type option []string

const collectInterval = 500

// optionProcess ...
func optionProcess(pid int32, name string, nonInteractive bool) (int32, string, option, error) {
	if pid > 0 && name != "" {
		return pid, name, option{}, errors.New("you can only use either --pid or --name")
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
			return optionProcess(pid, name, nonInteractive)
		}
		name, err = p.Name()
		if err != nil || name == "" {
			_, _ = fmt.Fprintf(os.Stderr, "No process found: %d\n", pid)
			return optionProcess(pid, name, nonInteractive)
		}

		fmt.Printf("Target process name: %s\n", color.Magenta(name))
		fmt.Println("")
		return pid, "", option{"--pid", processStr}, nil
	}

	if processStr != "" {
		fmt.Printf("Target process name: %s\n", color.Magenta(processStr))
		fmt.Println("")
		return 0, processStr, option{"--name", processStr}, nil
	}

	fmt.Println(color.Magenta("Track only host metrics"))
	fmt.Println("")
	return 0, "", option{}, nil
}

// optionThreshold ...
func optionThreshold(threshold string, pid int32, name string, nonInteractive bool) (option, error) {
	if nonInteractive {
		return option{"--threshold", threshold}, nil
	}
	var (
		m   *metrics.Metrics
		err error
	)

	switch {
	case pid > 0:
		m, err = metrics.GetMetrics(collectInterval, pid)
		if err != nil {
			return option{}, err
		}
	case name != "":
		m, err = metrics.GetMetricsByName(collectInterval, name)
		if err != nil {
			return option{}, err
		}
	default:
		m, err = metrics.GetMetrics(collectInterval)
		if err != nil {
			return option{}, err
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
	return option{"--threshold", threshold}, nil
}

// optionInterval ...
func optionInterval(interval int, nonInteractive bool) (option, error) {
	intervalStr := strconv.Itoa(interval)
	if nonInteractive {
		return option{"--interval", intervalStr}, nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--interval", color.B), "Interval of checking if the threshold exceeded (seconds).")
	fmt.Println("")
	intervalStr = prompter.Prompt("Enter interval", intervalStr)
	fmt.Println("")
	return option{"--interval", intervalStr}, nil
}

// optionAttempts ...
func optionAttempts(attempts int, nonInteractive bool) (option, error) {
	attemptsStr := strconv.Itoa(attempts)
	if nonInteractive {
		return option{"--attempts", attemptsStr}, nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--attempts", color.B), "Maximum number of attempts continuously exceeding the threshold.")
	fmt.Println("")
	attemptsStr = prompter.Prompt("Enter attempts", attemptsStr)
	fmt.Println("")
	return option{"--attempts", attemptsStr}, nil
}

// optionCommand ...
func optionCommand(command string, nonInteractive bool) (option, error) {
	if nonInteractive {
		return option{"--command", command}, nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--command", color.B), "Command to execute when the maximum number of attempts is exceeded.execution. If command execution time > 'interval' * 3, kill command.")
	fmt.Println("")
	fmt.Printf("%s\n", color.White("Additional Environment Variables", color.B))
	fmt.Printf("  %s: %s\n", color.White("$PID", color.B), "PID of the process.")
	fmt.Println("")
	command = prompter.Prompt("Enter command", command)
	fmt.Println("")
	return option{"--command", command}, nil
}

// optionTimes ...
func optionTimes(times int, nonInteractive bool) (option, error) {
	timesStr := strconv.Itoa(times)
	if nonInteractive {
		return option{"--times", timesStr}, nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--times", color.B), "Maximum number of command executions. If times < 1, track and execute until timeout.")
	fmt.Println("")
	timesStr = prompter.Prompt("Enter times", strconv.Itoa(times))
	fmt.Println("")
	return option{"--times", timesStr}, nil
}

// optionTimeout ...
func optionTimeout(timeout int, nonInteractive bool) (option, error) {
	timeoutStr := strconv.Itoa(timeout)
	if nonInteractive {
		return option{"--timeout", timeoutStr}, nil
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--tineout", color.B), "Timeout of tracking (seconds).")
	fmt.Println("")
	timeoutStr = prompter.Prompt("Enter timeout", timeoutStr)
	fmt.Println("")
	return option{"--timeout", timeoutStr}, nil
}

// optionSlackChannel ...
func optionSlackChannel(slackChannel string, nonInteractive bool) (option, error) {
	if nonInteractive {
		if slackChannel == "" {
			return option{}, nil
		} else {
			return option{"--slack-channel", slackChannel}, nil
		}
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
		return option{}, nil
	}
	if urlErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", urlErr)
		url = prompter.Prompt("Enter slack incoming webhook URL", "")
		if url == "" {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", errors.New("invalid URL"))
			return optionSlackChannel(slackChannel, nonInteractive)
		}
		err := os.Setenv("SLACK_INCOMMING_WEBHOOK_URL", url)
		if err != nil {
			return option{}, err
		}
	}
	slackChannel = prompter.Prompt("Enter slack channel", slackChannel)
	fmt.Println("")
	return option{"--slack-channel", slackChannel}, nil
}

// optionSlackMention ...
func optionSlackMention(slackMention string, nonInteractive bool) (option, error) {
	if nonInteractive {
		if slackMention == "" {
			return option{}, nil
		} else {
			return option{"--slack-mention", slackMention}, nil
		}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--slack-mention", color.B), "Slack mention.")
	fmt.Println("")
	yn := prompter.YN("Do you want to mention?", true)
	if !yn {
		fmt.Println("")
		return option{}, nil
	}
	slackMention = prompter.Prompt("Enter mention [@here or user_id (ex. @UXXXXXXXX)]", slackMention)
	fmt.Println("")
	return option{"--slack-mention", slackMention}, nil
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

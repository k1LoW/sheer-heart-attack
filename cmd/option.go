package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Songmu/prompter"
	"github.com/k1LoW/sheer-heart-attack/metrics"
	"github.com/labstack/gommon/color"
	"github.com/shirou/gopsutil/process"
)

type option []string

// optionPID ...
func optionPID(pid int32, nonInteractive bool) (option, error) {
	fmt.Printf("%s ... %s\n", color.Magenta("--pid", color.B), "PID of the process.")
	pidStr := strconv.Itoa(int(pid))
	if pidStr == "0" {
		pidStr = ""
	}
	fmt.Println("")
	pidStr = prompter.Prompt("Enter pid", pidStr)
	pidInt32, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return optionPID(pid, nonInteractive)
	}
	p, err := process.NewProcess(int32(pidInt32))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return optionPID(pid, nonInteractive)
	}
	name, err := p.Name()
	if err != nil || name == "" {
		_, _ = fmt.Fprintf(os.Stderr, "No process found: %s\n", pidStr)
		return optionPID(pid, nonInteractive)
	}
	fmt.Printf("Target process: %s\n", name)

	fmt.Println("")
	return option{"--pid", pidStr}, nil
}

// optionThreshold ...
func optionThreshold(threshold string, pid int32, nonInteractive bool) (option, error) {
	if nonInteractive {
		return option{"--threshold", threshold}, nil
	}
	m, err := metrics.Get(pid)
	if err != nil {
		return option{}, err
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--threshold", color.B), "Threshold conditions.")
	fmt.Println("")
	fmt.Printf("%s\n", color.Magenta("Available Metrics", color.B))
	mlist := metrics.List()
	for _, metric := range mlist {
		if of, ok := m[metric.Name]; ok {
			fmt.Printf("  %s (now:%s): %s\n", color.White(metric.Name), color.Magenta(fmt.Sprintf(metric.Format, of)), metric.Description)
		}
	}
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
	fmt.Printf("%s ... %s\n", color.Magenta("--command", color.B), "Command to execute when the maximum number of attempts is exceeded.")
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

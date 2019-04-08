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
func optionPID(pid int32, nonInteractive bool) option {
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
	return []string{"--pid", pidStr}
}

// optionThreshold ...
func optionThreshold(threshold string, pid int32, nonInteractive bool) option {
	if nonInteractive {
		return []string{"--threshold", threshold}
	}
	m, err := metrics.Get(pid)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return optionThreshold(threshold, pid, nonInteractive)
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--threshold", color.B), "Threshold conditions.")
	fmt.Println("")
	fmt.Printf("%s\n", color.Magenta("Available Metrics", color.B))
	fmt.Printf("  %s (now:%s): %s\n", color.White("cpu"), color.Magenta(fmt.Sprintf("%f", m["cpu"])), "The percentage of the CPU time the process uses (percent).")
	fmt.Printf("  %s (now:%s): %s\n", color.White("mem"), color.Magenta(fmt.Sprintf("%f", m["mem"])), "The percentage of the total RAM the process uses (percent).")
	fmt.Printf("  %s (now:%s): %s\n", color.White("rss"), color.Magenta(fmt.Sprintf("%d", m["rss"])), "The non-swapped physical memory the process uses (bytes).")
	fmt.Printf("  %s (now:%s): %s\n", color.White("vms"), color.Magenta(fmt.Sprintf("%d", m["vms"])), "The amount of virtual memory the process uses (bytes).")
	fmt.Printf("  %s (now:%s): %s\n", color.White("swap"), color.Magenta(fmt.Sprintf("%d", m["swap"])), "The amount of memory that has been swapped out to disk the process uses (bytes).")
	fmt.Printf("  %s (now:%s): %s\n", color.White("connections"), color.Magenta(fmt.Sprintf("%d", m["connections"])), "The amount of connections(TCP, UDP or UNIX) the process uses.")
	if of, ok := m["open_files"]; ok {
		fmt.Printf("  %s (now:%d): %s\n", color.White("open_files"), color.Magenta(fmt.Sprintf("%d", of)), "The amount of files and file discripters opend by the process.")
	}
	fmt.Printf("  %s (now:%s): %s\n", color.White("host_cpu"), color.Magenta(fmt.Sprintf("%f", m["host_cpu"])), "The percentage of cpu used.")
	fmt.Printf("  %s (now:%s): %s\n", color.White("host_mem"), color.Magenta(fmt.Sprintf("%f", m["host_mem"])), "The percentage of RAM used.")
	fmt.Printf("  %s (now:%s): %s\n", color.White("host_swap"), color.Magenta(fmt.Sprintf("%d", m["host_swap"])), "The amount of memory that has been swapped out to disk (bytes).")
	fmt.Printf("  %s (now:%s): %s\n", color.White("load1"), color.Magenta(fmt.Sprintf("%f", m["load1"])), "Load avarage for 1 minute.")
	fmt.Printf("  %s (now:%s): %s\n", color.White("load5"), color.Magenta(fmt.Sprintf("%f", m["load5"])), "Load avarage for 5 minutes.")
	fmt.Printf("  %s (now:%s): %s\n", color.White("load15"), color.Magenta(fmt.Sprintf("%f", m["load15"])), "Load avarage for 15 minutes.")
	fmt.Printf("%s\n", color.Magenta("Available Operators", color.B))
	fmt.Printf("  %s\n", "+, -, *, /, ==, !=, <, >, <=, >=, not, and, or, !, &&, ||")
	fmt.Println("")
	threshold = prompter.Prompt("Enter threshold", threshold)
	fmt.Println("")
	return []string{"--threshold", threshold}
}

// optionInterval ...
func optionInterval(interval int, nonInteractive bool) option {
	intervalStr := strconv.Itoa(interval)
	if nonInteractive {
		return []string{"--interval", intervalStr}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--interval", color.B), "Interval of checking if the threshold exceeded (seconds).")
	fmt.Println("")
	intervalStr = prompter.Prompt("Enter interval", intervalStr)
	fmt.Println("")
	return []string{"--interval", intervalStr}
}

// optionAttempts ...
func optionAttempts(attempts int, nonInteractive bool) option {
	attemptsStr := strconv.Itoa(attempts)
	if nonInteractive {
		return []string{"--attempts", attemptsStr}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--attempts", color.B), "Maximum number of attempts continuously exceeding the threshold.")
	fmt.Println("")
	attemptsStr = prompter.Prompt("Enter attempts", attemptsStr)
	fmt.Println("")
	return []string{"--attempts", attemptsStr}
}

// optionCommand ...
func optionCommand(command string, nonInteractive bool) option {
	if nonInteractive {
		return []string{"--command", command}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--command", color.B), "Command to execute when the maximum number of attempts is exceeded.")
	fmt.Println("")
	fmt.Printf("%s\n", color.White("Additional Environment Variables", color.B))
	fmt.Printf("  %s: %s\n", color.White("$PID", color.B), "PID of the process.")
	fmt.Println("")
	command = prompter.Prompt("Enter command", command)
	fmt.Println("")
	return []string{"--command", command}
}

// optionTimes ...
func optionTimes(times int, nonInteractive bool) option {
	timesStr := strconv.Itoa(times)
	if nonInteractive {
		return []string{"--times", timesStr}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--times", color.B), "Maximum number of command executions. If times < 1, track and execute until timeout.")
	fmt.Println("")
	timesStr = prompter.Prompt("Enter times", strconv.Itoa(times))
	fmt.Println("")
	return []string{"--times", timesStr}
}

// optionTimeout ...
func optionTimeout(timeout int, nonInteractive bool) option {
	timeoutStr := strconv.Itoa(timeout)
	if nonInteractive {
		return []string{"--timeout", timeoutStr}
	}
	fmt.Printf("%s ... %s\n", color.Magenta("--tineout", color.B), "Timeout of tracking (seconds).")
	fmt.Println("")
	timeoutStr = prompter.Prompt("Enter timeout", timeoutStr)
	fmt.Println("")
	return []string{"--timeout", timeoutStr}
}

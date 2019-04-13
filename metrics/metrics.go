package metrics

import (
	"os/exec"
	"strconv"
)

type Metric struct {
	Name        string
	Description string
	Format      string
}

// List ...
func List() []Metric {
	return []Metric{
		{"cpu", "Percentage of the CPU time the process uses.", "%f"},
		{"mem", "Percentage of the total RAM the process uses.", "%f"},
		{"rss", "Non-swapped physical memory the process uses (bytes).", "%d"},
		{"vms", "Amount of virtual memory the process uses (bytes).", "%d"},
		{"swap", "Amount of memory that has been swapped out to disk the process uses (bytes).", "%d"},
		{"connections", "Amount of connections(TCP, UDP or UNIX) the process uses.", "%d"},
		{"open_files", "Amount of files and file discripters opend by the process.", "%d"},
		{"host_cpu", "Percentage of cpu used.", "%f"},
		{"host_mem", "Percentage of RAM used.", "%f"},
		{"host_swap", "Amount of memory that has been swapped out to disk (bytes).", "%d"},

		{"user", "Percentage of CPU utilization that occurred while executing at the user level.", "%f"},
		{"system", "Percentage of CPU utilization that occurred while executing at the system level.", "%f"},
		{"idle", "Percentage of time that the CPU or CPUs were idle and the system did not have an outstanding disk I/O request.", "%f"},
		{"nice", "Percentage of CPU utilization that occurred while executing at the user level with nice priority.", "%f"},
		{"iowait", "Percentage of time that CPUs were idle during which the system had an outstanding disk I/O request.", "%f"},
		{"irq", "Percentage of time spent by CPUs to service hardware interrupts.", "%f"},
		{"softirq", "Percentage of time spent by CPUs to service software interrupts.", "%f"},
		{"steal", "Percentage of time spent in involuntary wait by the virtual CPUs while the hypervisor was servicing another virtual processor.", "%f"},
		{"guest", "Percentage of time spent by CPUs to run a virtual processor.", "%f"},
		{"guest_nice", "Percentage of time spent by CPUs to run a virtual processor with nice priority.", "%f"},
		// {"stolen", "", "%f"},

		{"load1", "Load avarage for 1 minute.", "%f"},
		{"load5", "Load avarage for 5 minutes.", "%f"},
		{"load15", "Load avarage for 15 minutes.", "%f"},
	}
}

func ClkTck() float64 {
	tck := float64(128)
	out, err := exec.Command("/usr/bin/getconf", "CLK_TCK").Output()
	if err == nil {
		i, err := strconv.ParseFloat(string(out), 64)
		if err == nil {
			tck = float64(i)
		}
	}
	return tck
}

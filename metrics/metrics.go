package metrics

type Metric struct {
	Name        string
	Description string
	Format      string
}

// List ...
func List() []Metric {
	return []Metric{
		{"cpu", "The percentage of the CPU time the process uses (percent).", "%f"},
		{"mem", "The percentage of the total RAM the process uses (percent).", "%f"},
		{"rss", "The non-swapped physical memory the process uses (bytes).", "%d"},
		{"vms", "The amount of virtual memory the process uses (bytes).", "%d"},
		{"swap", "The amount of memory that has been swapped out to disk the process uses (bytes).", "%d"},
		{"connections", "The amount of connections(TCP, UDP or UNIX) the process uses.", "%d"},
		{"open_files", "The amount of files and file discripters opend by the process.", "%d"},
		{"host_cpu", "The percentage of cpu used.", "%f"},
		{"host_mem", "The percentage of RAM used.", "%f"},
		{"host_swap", "The amount of memory that has been swapped out to disk (bytes).", "%d"},
		{"load1", "Load avarage for 1 minute.", "%f"},
		{"load5", "Load avarage for 5 minutes.", "%f"},
		{"load15", "Load avarage for 15 minutes.", "%f"},
	}
}

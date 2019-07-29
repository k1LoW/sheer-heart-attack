// +build linux

package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

func Get(pid int32) (map[string]interface{}, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return map[string]interface{}{}, err
	}
	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return map[string]interface{}{}, err
	}
	memPercent, err := p.MemoryPercent()
	if err != nil {
		return map[string]interface{}{}, err
	}
	memInfo, err := p.MemoryInfo()
	if err != nil {
		return map[string]interface{}{}, err
	}

	memoryMaps, err := p.MemoryMaps(true)
	if err != nil {
		return map[string]interface{}{}, err
	}
	maps := *memoryMaps
	swap := maps[0].Swap

	connections, err := p.Connections()
	if err != nil {
		return map[string]interface{}{}, err
	}

	openFiles, err := p.OpenFiles()
	if err != nil {
		return map[string]interface{}{}, err
	}

	hostCpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return map[string]interface{}{}, err
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return map[string]interface{}{}, err
	}

	sw, err := mem.SwapMemory()
	if err != nil {
		return map[string]interface{}{}, err
	}

	ts, err := cpu.Times(false)
	if err != nil {
		return map[string]interface{}{}, err
	}
	total := ts[0].Total()

	l, err := load.Avg()
	if err != nil {
		return map[string]interface{}{}, err
	}

	m := map[string]interface{}{
		"cpu":         cpuPercent,
		"mem":         memPercent,
		"rss":         memInfo.RSS,
		"vms":         memInfo.VMS,
		"swap":        swap,
		"connections": len(connections),
		"open_files":  len(openFiles),
		"host_cpu":    hostCpuPercent[0],
		"host_mem":    vm.UsedPercent,
		"host_swap":   sw.Used,
		"user":        ts[0].User / total * 100,
		"system":      ts[0].System / total * 100,
		"idle":        ts[0].Idle / total * 100,
		"nice":        ts[0].Nice / total * 100,
		"iowait":      ts[0].Iowait / total * 100,
		"irq":         ts[0].Irq / total * 100,
		"softirq":     ts[0].Softirq / total * 100,
		"steal":       ts[0].Steal / total * 100,
		"guest":       ts[0].Guest / total * 100,
		"guest_nice":  ts[0].GuestNice / total * 100,
		// "stolen":      ts[0].Stolen / total * 100,
		"load1":  l.Load1,
		"load5":  l.Load5,
		"load15": l.Load15,
	}
	return m, nil
}

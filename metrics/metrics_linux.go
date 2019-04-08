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
		"load1":       l.Load1,
		"load5":       l.Load5,
		"load15":      l.Load15,
	}
	return m, nil
}

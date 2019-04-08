// +build linux
package metrics

import (
	"github.com/shirou/gopsutil/process"
)

func Get(pid int32) (map[string]interface{}, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return map[string]interface{}{}, err
	}
	cpu, err := p.CPUPercent()
	if err != nil {
		return map[string]interface{}{}, err
	}
	mem, err := p.MemoryPercent()
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

	stat := map[string]interface{}{
		"cpu":         cpu,
		"mem":         mem,
		"rss":         memInfo.RSS,
		"vms":         memInfo.VMS,
		"swap":        swap,
		"connections": len(connections),
	}
	return stat, nil
}

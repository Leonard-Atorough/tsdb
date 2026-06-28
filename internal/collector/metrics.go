package collector

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemMetrics struct {
	CPUDuration time.Duration
}

func NewSystemMetrics(cpuDuration time.Duration) *SystemMetrics {
	return &SystemMetrics{
		CPUDuration: cpuDuration,
	}
}

func (sm *SystemMetrics) CPUUsage() (float64, error) {
	cpuPercentages, err := cpu.Percent(sm.CPUDuration, false)
	if err != nil {
		return 0, err
	}
	return cpuPercentages[0], nil
}

func (sm *SystemMetrics) MemoryUsage() (float64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.UsedPercent, nil
}

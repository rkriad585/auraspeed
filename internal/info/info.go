package info

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfo holds system information collected from the host.
type SystemInfo struct {
	OS       string
	CPU      string
	Memory   string
	Disk     string
	Hostname string
}

// GetSystemInfo collects and returns system information.
// It gathers OS, CPU, memory, disk, and hostname data.
func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{}

	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}
	info.OS = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	info.Hostname = hostInfo.Hostname

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}
	if len(cpuInfo) > 0 {
		cores, _ := cpu.Counts(true)
		info.CPU = fmt.Sprintf("%s (%d cores)", cpuInfo[0].ModelName, cores)
	} else {
		info.CPU = fmt.Sprintf("%d cores", runtime.NumCPU())
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}
	memTotalGB := float64(memInfo.Total) / (1024 * 1024 * 1024)
	memUsedGB := float64(memInfo.Used) / (1024 * 1024 * 1024)
	memPercent := memInfo.UsedPercent
	info.Memory = fmt.Sprintf("%.0f GB total, %.0f GB used (%.0f%%)", memTotalGB, memUsedGB, memPercent)

	diskInfo, err := disk.Usage("/")
	if err != nil {
		diskInfo, err = disk.Usage("C:")
		if err != nil {
			return nil, fmt.Errorf("failed to get disk info: %w", err)
		}
	}
	diskTotalGB := float64(diskInfo.Total) / (1024 * 1024 * 1024)
	diskUsedGB := float64(diskInfo.Used) / (1024 * 1024 * 1024)
	diskPercent := diskInfo.UsedPercent
	info.Disk = fmt.Sprintf("%.0f GB total, %.0f GB used (%.0f%%)", diskTotalGB, diskUsedGB, diskPercent)

	if info.Hostname == "" {
		hostname, _ := os.Hostname()
		info.Hostname = hostname
	}

	return info, nil
}

package cpu

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// darwinCPUReader macOS平台的CPU读取器
type darwinCPUReader struct{}

// GetCPUSample 获取CPU样本（macOS实现）
func (d *darwinCPUReader) GetCPUSample() (*CPUSample, error) {
	// 使用iostat命令获取CPU统计
	cmd := exec.Command("iostat", "-c", "2", "-n", "0")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute iostat: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 4 {
		return nil, fmt.Errorf("unexpected iostat output format")
	}

	// 解析第二组数据（跳过第一组，因为iostat需要两次采样）
	cpuLine := lines[3]
	fields := strings.Fields(cpuLine)
	if len(fields) < 6 {
		return nil, fmt.Errorf("insufficient fields in iostat output")
	}

	// 解析各个CPU时间百分比
	user, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user percentage: %w", err)
	}

	system, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse system percentage: %w", err)
	}

	idle, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse idle percentage: %w", err)
	}

	// 转换为ticks（简化处理，假设100 ticks per second）
	// 注意：iostat返回的是百分比，我们需要转换为ticks
	totalTicks := uint64(100)                 // 总ticks固定为100
	userTicks := uint64(user * 100 / 100)     // 用户态ticks
	systemTicks := uint64(system * 100 / 100) // 系统态ticks
	idleTicks := uint64(idle * 100 / 100)     // 空闲ticks
	ioWaitTicks := uint64(0)                  // macOS不直接提供iowait

	return &CPUSample{
		TotalTime:  totalTicks,
		UserTime:   userTicks,
		SystemTime: systemTicks,
		IdleTime:   idleTicks,
		IOWaitTime: ioWaitTicks,
		Timestamp:  time.Now(),
	}, nil
}

// GetLoadAverage 获取系统负载（macOS实现）
func (d *darwinCPUReader) GetLoadAverage() (float64, float64, float64, error) {
	cmd := exec.Command("uptime")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to execute uptime: %w", err)
	}

	// 解析uptime输出
	// 格式: 14:30  up 2 days, 3:45, 2 users, load averages: 1.23 2.34 3.45
	outputStr := string(output)
	if idx := strings.Index(outputStr, "load averages:"); idx != -1 {
		loadPart := outputStr[idx+14:] // 跳过"load averages: "
		fields := strings.Fields(loadPart)
		if len(fields) >= 3 {
			load1, err1 := strconv.ParseFloat(fields[0], 64)
			load5, err2 := strconv.ParseFloat(fields[1], 64)
			load15, err3 := strconv.ParseFloat(fields[2], 64)

			if err1 == nil && err2 == nil && err3 == nil {
				return load1, load5, load15, nil
			}
		}
	}

	return 0, 0, 0, fmt.Errorf("failed to parse load averages from uptime output")
}

// GetCPUTemperature 获取CPU温度（macOS实现）
func (d *darwinCPUReader) GetCPUTemperature() (float64, error) {
	// 尝试使用powermetrics（需要sudo）
	cmd := exec.Command("sudo", "-n", "powermetrics", "--samplers", "smc", "-n", "1")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to execute powermetrics: %w", err)
	}

	// 解析powermetrics输出查找CPU温度
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "CPU die temperature") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if strings.Contains(field, "temperature") && i+1 < len(fields) {
					// 提取温度值
					tempStr := strings.TrimSuffix(fields[i+1], "°C")
					temp, err := strconv.ParseFloat(tempStr, 64)
					if err != nil {
						return 0, fmt.Errorf("failed to parse temperature: %w", err)
					}
					return temp, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("CPU temperature not found in powermetrics output")
}

// GetCPUFrequency 获取CPU频率（macOS实现）
func (d *darwinCPUReader) GetCPUFrequency() (float64, error) {
	// 使用sysctl获取CPU频率信息
	cmd := exec.Command("sysctl", "-n", "hw.cpufrequency")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to execute sysctl: %w", err)
	}

	freqStr := strings.TrimSpace(string(output))
	freq, err := strconv.ParseFloat(freqStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CPU frequency: %w", err)
	}

	// 转换为MHz
	return freq / 1000000, nil // Hz to MHz
}

// GetPlatformName 获取平台名称
func (d *darwinCPUReader) GetPlatformName() string {
	return "darwin"
}

package cpu

import (
	"runtime"
	"time"
)

type PlatformCPUReader interface {
	GetCPUSample() (*CPUSample, error)
	GetLoadAverage() (float64, float64, float64, error)
	GetCPUTemperature() (float64, error)
	GetCPUFrequency() (float64, error)
	GetPlatformName() string
}

var platformReader PlatformCPUReader

func initPlatformReader() {
	switch getPlatform() {
	case "darwin":
		platformReader = &darwinCPUReader{}
		// case "linux":
		// 	platformReader = &linuxCPUReader{}
		// case "windows":
		// 	platformReader = &windowsCPUReader{}
		// default:
		// 	platformReader = &fallbackCPUReader{}
	}
}

// 获取当前平台
func getPlatform() string {
	return runtime.GOOS
}

func getPlatformCPUSample() *CPUSample {
	if platformReader == nil {
		initPlatformReader()
	}

	sample, err := platformReader.GetCPUSample()

	if err != nil {
		return getFallbackCPUSample()
	}

	return sample
}

func getPlatformLoadAverage() (float64, float64, float64) {
	if platformReader == nil {
		initPlatformReader()
	}

	load1, load5, load15, err := platformReader.GetLoadAverage()

	if err != nil {
		return 0.0, 0.0, 0.0
	}

	return load1, load5, load15
}

func getPlatformCPUTemperature() float64 {
	if platformReader == nil {
		initPlatformReader()
	}

	temp, err := platformReader.GetCPUTemperature()
	if err != nil {
		return 0.0
	}

	return temp
}

// getPlatformCPUFrequency 获取平台相关的CPU频率
func getPlatformCPUFrequency() float64 {
	if platformReader == nil {
		initPlatformReader()
	}

	freq, err := platformReader.GetCPUFrequency()
	if err != nil {
		return 0.0 // 如果无法获取频率，返回0
	}

	return freq
}

// 基于go runtime的fallback trails
func getFallbackCPUSample() *CPUSample {
	now := time.Now()

	var totalTicks uint64 = 1000

	numGoroutines := runtime.NumGoroutine()

	// 用gorountine的数量估算cpu是哦旅
	// 假设100个goroutines对应50%的CPU使用率
	estimatedUsage := float64(numGoroutines) / 200.0
	if estimatedUsage > 0.9 {
		estimatedUsage = 0.9 // 限制最大值
	}
	if estimatedUsage < 0.1 {
		estimatedUsage = 0.1
	}
	// 70%用户态
	userTicks := uint64(float64(totalTicks) * estimatedUsage * 0.7)
	// 20%为系统态
	systemTicks := uint64(float64(totalTicks) * estimatedUsage * 0.2)
	// 剩余为空闲
	idleTicks := totalTicks - userTicks - systemTicks

	return &CPUSample{
		TotalTime:  totalTicks,
		UserTime:   userTicks,
		SystemTime: systemTicks,
		IdleTime:   idleTicks,
		IOWaitTime: 0, // fallback不提供IO等待
		Timestamp:  now,
	}
}

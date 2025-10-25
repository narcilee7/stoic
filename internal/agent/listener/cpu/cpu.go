package cpu

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stoic/internal/agent/listener"
)

// CPUConfig CPU监听器配置
type CPUConfig struct {
	Enabled           bool          `toml:"enabled"`
	SampleInterval    time.Duration `toml:"sample_interval"`
	BufferSize        int           `toml:"buffer_size"`
	ThresholdWarning  float64       `toml:"threshold_warning"`  // 警告阈值 (0-1)
	ThresholdCritical float64       `toml:"threshold_critical"` // 严重阈值 (0-1)
	HistorySize       int           `toml:"history_size"`       // 历史数据保留数量
}

// DefaultCPUConfig 返回默认CPU配置
func DefaultCPUConfig() *CPUConfig {
	return &CPUConfig{
		Enabled:           true,
		SampleInterval:    1 * time.Second,
		BufferSize:        1000,
		ThresholdWarning:  0.7, // 70%
		ThresholdCritical: 0.9, // 90%
		HistorySize:       60,  // 保留60个样本（1分钟）
	}
}

// CPUStats CPU统计信息
type CPUStats struct {
	UsagePercent  float64   `json:"usage_percent"`   // CPU使用率百分比
	UserPercent   float64   `json:"user_percent"`    // 用户态使用率
	SystemPercent float64   `json:"system_percent"`  // 系统态使用率
	IdlePercent   float64   `json:"idle_percent"`    // 空闲率
	IoWaitPercent float64   `json:"io_wait_percent"` // IO等待率
	LoadAverage1  float64   `json:"load_average_1"`  // 1分钟负载
	LoadAverage5  float64   `json:"load_average_5"`  // 5分钟负载
	LoadAverage15 float64   `json:"load_average_15"` // 15分钟负载
	CoreCount     int       `json:"core_count"`      // CPU核心数
	ThreadCount   int       `json:"thread_count"`    // 线程数
	Temperature   float64   `json:"temperature"`     // CPU温度（如果可用）
	Frequency     float64   `json:"frequency"`       // CPU频率（MHz）
	Timestamp     time.Time `json:"timestamp"`
}

// CPUSample CPU样本数据
type CPUSample struct {
	TotalTime  uint64 // 总时间
	UserTime   uint64 // 用户态时间
	SystemTime uint64 // 系统态时间
	IdleTime   uint64 // 空闲时间
	IOWaitTime uint64 // IO等待时间（Linux）
	Timestamp  time.Time
}

// CPUListener CPU监听器接口
type CPUListener interface {
	Start() error
	Stop() error
	IsActive() bool
	GetName() string
	GetType() string
	GetCurrentStats() *CPUStats
	GetHistory() []CPUStats
	GetEventChannel() <-chan *listener.Event
	GetLoadAverage() (float64, float64, float64)
	GetStats() map[string]interface{}
}

// cpuListenerImpl CPU监听器实现
type cpuListenerImpl struct {
	config *CPUConfig
	ctx    context.Context
	cancel context.CancelFunc

	// 状态管理
	mu        sync.RWMutex
	running   atomic.Bool
	startTime time.Time

	// 数据收集
	lastSample   *CPUSample
	currentStats *CPUStats
	history      []CPUStats

	// 事件通道
	eventChan chan *listener.Event

	// 统计
	stats *CPUListenerStats
}

// CPUListenerStats 监听器统计
type CPUListenerStats struct {
	SamplesCollected  int64     `json:"samples_collected"`
	EventsGenerated   int64     `json:"events_generated"`
	WarningsIssued    int64     `json:"warnings_issued"`
	CriticalsIssued   int64     `json:"criticals_issued"`
	ErrorsEncountered int64     `json:"errors_encountered"`
	LastSampleTime    time.Time `json:"last_sample_time"`
	AverageUsage      float64   `json:"average_usage"`
	PeakUsage         float64   `json:"peak_usage"`
}

func NewCPUListener(config *CPUConfig) (CPUListener, error) {
	if config == nil {
		config = DefaultCPUConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	listener := &cpuListenerImpl{
		config:    config,
		ctx:       ctx,
		cancel:    cancel,
		history:   make([]CPUStats, 0, config.HistorySize),
		eventChan: make(chan *listener.Event, config.BufferSize),
		stats:     &CPUListenerStats{},
	}

	return listener, nil
}

func (c *cpuListenerImpl) Start() error {
	if c.running.Load() {
		return fmt.Errorf("CPU listener is already running")
	}

	c.running.Store(true)
	c.startTime = time.Now()

	// 启动数据收集的循环
	go c.collectLoop()

	return nil
}

func (c *cpuListenerImpl) Stop() error {
	if !c.running.Load() {
		return nil
	}

	c.cancel()
	c.running.Store(false)
	// 关闭
	close(c.eventChan)

	return nil
}

func (c *cpuListenerImpl) IsActive() bool {
	return c.running.Load()
}

func (c *cpuListenerImpl) GetName() string {
	return "cpu_listener"
}

func (c *cpuListenerImpl) GetType() string {
	return "system_monitor"
}

func (c *cpuListenerImpl) GetCurrentStats() *CPUStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.currentStats == nil {
		return &CPUStats{
			UsagePercent: 0,
			Timestamp:    time.Now(),
		}
	}

	stats := *c.currentStats
	return &stats
}

func (c *cpuListenerImpl) GetHistory() []CPUStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	history := make([]CPUStats, len(c.history))
	copy(history, c.history)

	return history
}

func (c *cpuListenerImpl) GetEventChannel() <-chan *listener.Event {
	return c.eventChan
}

func (c *cpuListenerImpl) GetLoadAverage() (float64, float64, float64) {
	v1, v5, v15, err := platformReader.GetLoadAverage()
	if err != nil {
		return 0.0, 0.0, 0.0
	}
	return v1, v5, v15
}

func (c *cpuListenerImpl) GetCPUSample() *CPUSample {
	sample, err := platformReader.GetCPUSample()
	if err != nil {
		return nil
	}
	return sample
}

func (c *cpuListenerImpl) GetCPUTemperature() float64 {
	temp, err := platformReader.GetCPUTemperature()
	if err != nil {
		return 0.0
	}
	return temp
}

func (c *cpuListenerImpl) collectLoop() {
	ticker := time.NewTicker(c.config.SampleInterval)
	defer ticker.Stop()

	c.lastSample = c.GetCPUSample()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.collectSample(); err != nil {
				c.stats.ErrorsEncountered++
				continue
			}
		}
	}
}

func (c *cpuListenerImpl) GetCPUFrequency() float64 {
	freq, err := platformReader.GetCPUFrequency()
	if err != nil {
		return 0.0
	}
	return freq
}

func (c *cpuListenerImpl) collectSample() error {
	// 获取新的cpu样本
	newSample := c.GetCPUSample()
	if newSample == nil {
		return fmt.Errorf("failed to get CPU sample")
	}

	stats, err := c.calculateCPUUsage(c.lastSample, newSample)
	if err == nil {
		return fmt.Errorf("failed to calculate CPU usage")
	}

	// 更新当前统计
	c.mu.Lock()
	c.currentStats = stats

	// 添加到历史记录
	c.addToHistory(stats)

	// 更新统计信息
	c.updateStats(stats)
	c.mu.Unlock()

	// 检查阈值并生成事件
	c.checkThresholds(stats)

	// 更新最后一个样本
	c.lastSample = newSample

	return nil
}

func (c *cpuListenerImpl) calculateCPUUsage(prev, curr *CPUSample) (*CPUStats, error) {
	if prev == nil || curr == nil {
		return nil, nil
	}

	totalDiff := curr.TotalTime - prev.TotalTime
	if totalDiff == 0 {
		return nil, nil
	}

	userDiff := curr.UserTime - prev.UserTime
	systemDiff := curr.SystemTime - prev.SystemTime
	idleDiff := curr.IdleTime - prev.IdleTime
	ioWaitDiff := curr.IOWaitTime - prev.IOWaitTime

	userPercent := float64(userDiff) / float64(totalDiff) * 100
	systemPercent := float64(systemDiff) / float64(totalDiff) * 100
	idlePercent := float64(idleDiff) / float64(totalDiff) * 100
	ioWaitPercent := float64(ioWaitDiff) / float64(totalDiff) * 100

	totalPercent := 100 - idlePercent

	coreCount := runtime.NumCPU()
	load1, load5, load15 := c.GetLoadAverage() // 使用已有的函数获取负载

	stats := &CPUStats{
		UsagePercent:  totalPercent,
		UserPercent:   userPercent,
		SystemPercent: systemPercent,
		IdlePercent:   idlePercent,
		IoWaitPercent: ioWaitPercent,
		LoadAverage1:  load1,
		LoadAverage5:  load5,
		LoadAverage15: load15,
		CoreCount:     coreCount,
		ThreadCount:   runtime.GOMAXPROCS(0),
		Timestamp:     curr.Timestamp,
	}

	stats.Temperature = c.GetCPUTemperature()
	stats.Frequency = c.GetCPUFrequency()

	return stats, nil
}

func (c *cpuListenerImpl) addToHistory(stats *CPUStats) error {
	c.history = append(c.history, *stats)

	if len(c.history) > c.config.HistorySize {
		c.history = c.history[1:]
	}

	return nil
}

func (c *cpuListenerImpl) updateStats(stats *CPUStats) {
	c.stats.SamplesCollected++
	c.stats.LastSampleTime = stats.Timestamp

	if c.stats.SamplesCollected == 1 {
		c.stats.AverageUsage = stats.UsagePercent
	} else {
		c.stats.AverageUsage = (c.stats.AverageUsage*float64(c.stats.SamplesCollected-1) + stats.UsagePercent) / float64(c.stats.SamplesCollected)
	}

	if stats.UsagePercent > c.stats.PeakUsage {
		c.stats.PeakUsage = stats.UsagePercent
	}
}

func (c *cpuListenerImpl) checkThresholds(stats *CPUStats) {
	// 检查警告阈值
	if stats.UsagePercent >= c.config.ThresholdCritical {
		c.generateEvent("cpu_critical", stats.UsagePercent, map[string]interface{}{
			"threshold": c.config.ThresholdCritical,
			"load_avg":  stats.LoadAverage1,
			"severity":  "critical",
		})
		c.stats.CriticalsIssued++
	} else if stats.UsagePercent >= c.config.ThresholdWarning {
		c.generateEvent("cpu_warning", stats.UsagePercent, map[string]interface{}{
			"threshold": c.config.ThresholdWarning,
			"load_avg":  stats.LoadAverage1,
			"severity":  "warning",
		})
		c.stats.WarningsIssued++
	}
}

// generateEvent 生成事件
func (c *cpuListenerImpl) generateEvent(eventType string, value float64, metadata map[string]interface{}) {
	// 转换事件类型字符串到EventType
	var evtType listener.EventType
	switch eventType {
	case "cpu_warning":
		evtType = listener.EventCPUWarning
	case "cpu_critical":
		evtType = listener.EventCPUCritical
	default:
		evtType = listener.EventType(eventType) // fallback
	}

	event := &listener.Event{
		ID:         fmt.Sprintf("cpu_%s_%d", eventType, time.Now().Unix()),
		Type:       evtType,
		Source:     "cpu_listener",
		Value:      value / 100.0, // 将百分比转换为0-1范围
		Metadata:   metadata,
		Timestamp:  time.Now(),
		Confidence: 0.9, // CPU数据置信度较高
		Processed:  false,
	}

	// 设置严重程度
	event.Severity = listener.DetermineSeverity(evtType, value)

	select {
	case c.eventChan <- event:
		c.stats.EventsGenerated++
	default:
		// 通道满了，丢弃事件
		// 在实际实现中可以记录日志
	}
}

func (c *cpuListenerImpl) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"running":            c.running.Load(),
		"start_time":         c.startTime,
		"samples_collected":  c.stats.SamplesCollected,
		"events_generated":   c.stats.EventsGenerated,
		"warning_issued":     c.stats.WarningsIssued,
		"criticals_issued":   c.stats.CriticalsIssued,
		"errors_encountered": c.stats.ErrorsEncountered,
		"last_sample_time":   c.stats.LastSampleTime,
		"average_usage":      c.stats.AverageUsage,
		"peak_usage":         c.stats.PeakUsage,
		"current_usage":      c.getCurrentUsage(),
	}
}

func (c *cpuListenerImpl) getCurrentUsage() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.currentStats != nil {
		return c.currentStats.UsagePercent
	}

	return 0
}

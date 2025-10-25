package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stoic/internal/agent/listener"
	"github.com/stoic/internal/agent/planner"
	"github.com/stoic/internal/agent/executor"
	"github.com/stoic/internal/infra/database"
)

// Engine Agent核心引擎
type Engine struct {
	ctx        context.Context
	cancel     context.CancelFunc
	db         *database.DB
	
	// 三大核心组件
	listener   listener.Manager    // 情境感知
	planner    planner.Manager     // 策略规划  
	executor   executor.Manager    // 行动执行
	
	// 状态管理
	mu         sync.RWMutex
	enabled    bool
	running    bool
	startTime  time.Time
	
	// 配置
	config     *Config
	
	// 事件通道
	eventChan  chan *Event
	resultChan chan *Result
}

// Config Agent配置
type Config struct {
	Enabled          bool          `toml:"enabled"`
	EventBufferSize  int           `toml:"event_buffer_size"`
	ProcessInterval  time.Duration `toml:"process_interval"`
	MaxEventsPerBatch int          `toml:"max_events_per_batch"`
	CooldownPeriod   time.Duration `toml:"cooldown_period"`
	PrivacyLevel     string        `toml:"privacy_level"` // standard, strict, minimal
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:           true,
		EventBufferSize:   1000,
		ProcessInterval:   5 * time.Second,
		MaxEventsPerBatch: 50,
		CooldownPeriod:    30 * time.Second,
		PrivacyLevel:      "standard",
	}
}


// Event 统一事件结构
type Event struct {
	ID          string                 `json:"id"`
	Type        string    `json:"type"`        // keyboard, cpu, git, idle, etc.
	Source      string                 `json:"source"`      // 事件来源组件
	Value       float64                `json:"value"`       // 事件数值
	Metadata    map[string]interface{} `json:"metadata"`    // 额外元数据
	Timestamp   time.Time              `json:"timestamp"`
	Processed   bool                   `json:"processed"`
}

// Result 处理结果
type Result struct {
	EventID     string                 `json:"event_id"`
	Intervention *Intervention          `json:"intervention,omitempty"`
	Decision    string                 `json:"decision"` // execute, ignore, defer
	Reason      string                 `json:"reason"`
	Confidence  float64                `json:"confidence"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Intervention 干预方案
type Intervention struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`        // breathe, scream, q, wcg, etc.
	Urgency     string                 `json:"urgency"`     // low, medium, high, critical
	Timing      string                 `json:"timing"`      // immediate, delayed, scheduled
	Parameters  map[string]interface{} `json:"parameters"`  // 干预参数
	PredictedEffectiveness float64     `json:"predicted_effectiveness"`
	Context     map[string]interface{} `json:"context"`     // 上下文信息
}

// NewEngine 创建新的Agent引擎
func NewEngine(db *database.Connection, config *Config) (*Engine, error) {
	if config == nil {
		config = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	engine := &Engine{
		ctx:        ctx,
		cancel:     cancel,
		db:         db,
		config:     config,
		eventChan:  make(chan *Event, config.EventBufferSize),
		resultChan: make(chan *Result, config.EventBufferSize),
		enabled:    config.Enabled,
	}

	// 初始化三大核心组件
	if err := engine.initComponents(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize components: %w", err)
	}

	return engine, nil
}

// initComponents 初始化核心组件
func (e *Engine) initComponents() error {
	// 初始化监听器管理器
	listenerConfig := listener.DefaultConfig()
	listenerConfig.PrivacyLevel = e.config.PrivacyLevel
	
	listenerManager, err := listener.NewManager(listenerConfig, e.db)
	if err != nil {
		return fmt.Errorf("failed to create listener manager: %w", err)
	}
	e.listener = listenerManager

	// 初始化规划器管理器
	plannerConfig := planner.DefaultConfig()
	plannerManager, err := planner.NewManager(plannerConfig, e.db)
	if err != nil {
		return fmt.Errorf("failed to create planner manager: %w", err)
	}
	e.planner = plannerManager

	// 初始化执行器管理器
	executorConfig := executor.DefaultConfig()
	executorManager, err := executor.NewManager(executorConfig, e.db)
	if err != nil {
		return fmt.Errorf("failed to create executor manager: %w", err)
	}
	e.executor = executorManager

	return nil
}

// Start 启动Agent引擎
func (e *Engine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return fmt.Errorf("engine is already running")
	}

	if !e.enabled {
		return fmt.Errorf("engine is disabled")
	}

	// 启动各个组件
	if err := e.listener.Start(); err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	if err := e.executor.Start(); err != nil {
		e.listener.Stop()
		return fmt.Errorf("failed to start executor: %w", err)
	}

	// 启动事件处理循环
	go e.eventLoop()

	e.running = true
	e.startTime = time.Now()

	return nil
}

// Stop 停止Agent引擎
func (e *Engine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running {
		return nil
	}

	// 停止事件处理循环
	e.cancel()

	// 停止各个组件
	e.listener.Stop()
	e.executor.Stop()
	e.planner.Stop()

	// 关闭通道
	close(e.eventChan)
	close(e.resultChan)

	e.running = false
	return nil
}

// eventLoop 事件处理主循环
func (e *Engine) eventLoop() {
	ticker := time.NewTicker(e.config.ProcessInterval)
	defer ticker.Stop()

	var eventBatch []*Event

	for {
		select {
		case <-e.ctx.Done():
			return
			
		case event, ok := <-e.eventChan:
			if !ok {
				return
			}
			eventBatch = append(eventBatch, event)
			
			// 如果批次满了，立即处理
			if len(eventBatch) >= e.config.MaxEventsPerBatch {
				e.processEventBatch(eventBatch)
				eventBatch = nil
			}
			
		case <-ticker.C:
			// 定期处理累积的事件
			if len(eventBatch) > 0 {
				e.processEventBatch(eventBatch)
				eventBatch = nil
			}
		}
	}
}

// processEventBatch 处理事件批次
func (e *Engine) processEventBatch(events []*Event) {
	for _, event := range events {
		result, err := e.processEvent(event)
		if err != nil {
			// 记录错误但不停止处理
			continue
		}
		
		if result != nil {
			select {
			case e.resultChan <- result:
			case <-e.ctx.Done():
				return
			}
		}
	}
}

// processEvent 处理单个事件
func (e *Engine) processEvent(event *Event) (*Result, error) {
	// 1. 情境分析 - 监听器处理
	context, err := e.listener.Analyze(event)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze event: %w", err)
	}

	// 2. 策略规划 - 规划器决策
	decision, err := e.planner.Plan(event, context)
	if err != nil {
		return nil, fmt.Errorf("failed to plan intervention: %w", err)
	}

	// 3. 如果决定执行干预，交给执行器
	if decision.Decision == "execute" && decision.Intervention != nil {
		if err := e.executor.Execute(decision.Intervention); err != nil {
			return nil, fmt.Errorf("failed to execute intervention: %w", err)
		}
	}

	return decision, nil
}

// SubmitEvent 提交事件到Agent
func (e *Engine) SubmitEvent(event *Event) error {
	if !e.running {
		return fmt.Errorf("engine is not running")
	}

	select {
	case e.eventChan <- event:
		return nil
	case <-e.ctx.Done():
		return fmt.Errorf("engine is shutting down")
	default:
		return fmt.Errorf("event buffer is full")
	}
}

// GetResults 获取处理结果（非阻塞）
func (e *Engine) GetResults() <-chan *Result {
	return e.resultChan
}

// IsRunning 检查引擎是否运行中
func (e *Engine) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

// IsEnabled 检查引擎是否启用
func (e *Engine) IsEnabled() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enabled
}

// SetEnabled 设置引擎启用状态
func (e *Engine) SetEnabled(enabled bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.enabled = enabled
}

// GetStats 获取Agent统计信息
func (e *Engine) GetStats() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	stats := map[string]interface{}{
		"running":     e.running,
		"enabled":     e.enabled,
		"start_time":  e.startTime,
		"uptime":      time.Since(e.startTime).Seconds(),
	}

	// 添加组件统计
	if e.listener != nil {
		stats["listener_stats"] = e.listener.GetStats()
	}
	if e.planner != nil {
		stats["planner_stats"] = e.planner.GetStats()
	}
	if e.executor != nil {
		stats["executor_stats"] = e.executor.GetStats()
	}

	return stats
}

import (
	"context"

	"github.com/stoic/internal/infra/database"
)

type Engine struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     *database.DB

	// listener *listener.LI
}

type Config struct {
	Enable          bool `toml:"enabled"`
	EventBufferSize int  `toml:"event_buffer_size"`
	Process
}

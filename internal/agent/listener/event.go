package listener

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type EventType string

const (
	// 系统监控
	EventCPUWarning  EventType = "cpu_warning"
	EventCPUCritical EventType = "cpu_critical"
	EventMemoryHigh  EventType = "memory_high"
	EventDiskFull    EventType = "disk_full"
	EventSystemLoad  EventType = "system_load"

	// 用户行为
	EventKeyboardBurst EventType = "keyboard_burst"
	EventMouseRapid    EventType = "mouse_rapid"
	EventIdleStart     EventType = "idle_start"
	EventIdleEnd       EventType = "idle_end"

	// 开发
	EventGitReset     EventType = "git_reset"
	EventGitCommit    EventType = "git_commit"
	EventBuildFailed  EventType = "build_failed"
	EventTestFailed   EventType = "test_failed"
	EventCompileError EventType = "compile_error"

	// 情绪相关事件
	EventMoodDrop    EventType = "mood_drop"
	EventStressHigh  EventType = "stress_high"
	EventFocusLost   EventType = "focus_lost"
	EventAnxietyHigh EventType = "anxiety_high"

	// 干预事件
	EventInterventionSuccess EventType = "intervention_success"
	EventInterventionFailed  EventType = "intervention_failed"
	EventInterventionIgnored EventType = "intervention_ignored"
)

type EventSeverity int

const (
	SeverityInfo EventSeverity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s EventSeverity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

type Event struct {
	ID         string                 `json:"id"`
	Type       EventType              `json:"type"`
	Source     string                 `json:"source"`
	Value      float64                `json:"value"`
	Metadata   map[string]interface{} `json:"metadata"`
	Timestamp  time.Time              `json:"timestamp"`
	Processed  bool                   `json:"processed"`
	Confidence float64                `json:"confidence"`
	Severity   EventSeverity          `json:"severity"`
}

type EventSource struct {
	Component  string                 `json:"component"`   // 组件名称
	Version    string                 `json:"version"`     // 组件版本
	InstanceID string                 `json:"instance_id"` // 实例ID
	Metadata   map[string]interface{} `json:"metadata"`    // 来源元数据
}

// EventPool 事件池（用于复用事件对象，减少GC压力）
type EventPool struct {
	mu      sync.Mutex
	events  []*Event
	maxSize int
}

func NewEventPool(maxSize int) *EventPool {
	return &EventPool{
		events:  make([]*Event, maxSize),
		maxSize: maxSize,
	}
}

func (p *EventPool) Get() *Event {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.events) > 0 {
		event := p.events[len(p.events)-1]
		p.events = p.events[:len(p.events)-1]
		return event
	}

	return &Event{
		Metadata: make(map[string]interface{}),
	}
}

// Put 将事件放回池中
func (p *EventPool) Put(event *Event) {
	// 重置事件
	event.ID = ""
	event.Type = ""
	event.Source = ""
	event.Value = 0
	event.Processed = false
	event.Confidence = 0
	event.Severity = SeverityInfo
	event.Timestamp = time.Time{}

	// 清空元数据但保留map
	for k := range event.Metadata {
		delete(event.Metadata, k)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.events) < p.maxSize {
		p.events = append(p.events, event)
	}
}

// 事件总线
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan *Event
	closed      bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: map[string][]chan *Event{},
	}
}

func (eb *EventBus) Subscribe(eventType EventType) <-chan *Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return nil
	}

	ch := make(chan *Event, 100)
	eb.subscribers[string(eventType)] = append(eb.subscribers[string(eventType)], ch)
	return ch
}

func (eb *EventBus) Publish(event *Event) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.closed {
		return fmt.Errorf("event bus is closed")
	}

	// 发布到具体类型订阅者
	if subscribers, exists := eb.subscribers[string(event.Type)]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// 通道满了，丢弃事件
			}
		}
	}

	// 发布到通配符订阅者(监听所有事件)
	if subscribers, exists := eb.subscribers["*"]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
				// 通道满了
			}
		}
	}

	return nil
}

func (eb *EventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return
	}

	eb.closed = true

	for _, subscribers := range eb.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
	}

	eb.subscribers = nil
}

type EventStore interface {
	Save(event *Event) error
	GetByID(id string) (*Event, error)
	GetByType(eventType EventType, limit int) ([]*Event, error)
	GetByTimeRange(start, end time.Time) ([]*Event, error)
	GetRecent(limit int) ([]*Event, error)
	DeleteBefore(before time.Time) (int64, error)
}

// EventValidator 事件验证器
type EventValidator interface {
	Validate(event *Event) error
}

// EventProcessor 事件处理器接口
type EventProcessor interface {
	Process(event *Event) error
	CanProcess(eventType EventType) bool
}

// EventFilter 事件过滤器
type EventFilter func(*Event) bool

func FilterEvents(events []*Event, filters ...EventFilter) []*Event {
	if len(events) == 0 {
		return events
	}

	filtered := make([]*Event, 0, len(events))

	for _, event := range events {
		match := true
		for _, filter := range filters {
			if !filter(event) {
				match = false
				break
			}
		}
		if match {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

func NewEvent(eventType EventType, source string, value float64) *Event {
	return &Event{
		ID:         generateEventID(),
		Type:       eventType,
		Source:     source,
		Value:      value,
		Metadata:   make(map[string]interface{}),
		Timestamp:  time.Now(),
		Processed:  false,
		Confidence: 0.9, // 默认置信度
		Severity:   DetermineSeverity(eventType, value),
	}
}

func generateEventID() string {
	timestamp := time.Now().UnixNano()

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Sprintf("evt_%d", timestamp)
	}

	return fmt.Sprintf("evt_%d_%s", timestamp, hex.EncodeToString(randomBytes))
}

// 根据事件类型和值确定严重程度
func DetermineSeverity(eventType EventType, value float64) EventSeverity {
	switch eventType {
	case EventCPUCritical, EventStressHigh:
		return SeverityCritical
	case EventCPUWarning, EventMoodDrop:
		if value >= 0.8 {
			return SeverityHigh
		}
		return SeverityMedium
	case EventKeyboardBurst, EventGitReset:
		if value >= 0.7 {
			return SeverityMedium
		}
		return SeverityLow
	default:
		return SeverityInfo
	}
}

func EventTypeFromString(s string) EventType {
	return EventType(s)
}

func IsValidEventType(eventType EventType) bool {
	switch eventType {
	case EventCPUWarning, EventCPUCritical, EventMemoryHigh, EventDiskFull,
		EventSystemLoad, EventKeyboardBurst, EventMouseRapid, EventIdleStart,
		EventIdleEnd, EventGitReset, EventGitCommit, EventBuildFailed,
		EventTestFailed, EventCompileError, EventMoodDrop, EventStressHigh,
		EventFocusLost, EventAnxietyHigh, EventInterventionSuccess,
		EventInterventionFailed, EventInterventionIgnored:
		return true
	default:
		return false
	}
}

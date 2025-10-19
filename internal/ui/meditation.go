package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MeditationMode struct {
	active       bool
	breathPhase  int // 0: 吸气, 1: 呼气
	cycleCount   int
	startTime    time.Time
	quotes       []string
	currentQuote int
	breathTimer  time.Duration
}

type meditationBreathMsg struct {
	phase int
	quote string
	cycle int
}

func NewMeditationMode() *MeditationMode {
	return &MeditationMode{
		quotes: []string{
			"深呼吸，让心灵归于平静...",
			"感受此刻的存在，放下过往的牵挂...",
			"呼吸如波浪，起伏皆有韵律...",
			"静心如水，映照万物而不动摇...",
			"每一次呼吸，都是与宇宙的对话...",
			"宁静不是无声，而是内心的和谐...",
			"在静默中，我们听见真正的自己...",
			"让思绪如云，来去自如...",
			"此刻即是永恒，当下就是全部...",
			"呼吸之间，生命在流动...",
		},
		breathTimer: 4 * time.Second,
	}
}

// Start 开始冥想
func (m *MeditationMode) Start() tea.Cmd {
	m.active = true
	m.startTime = time.Now()
	m.currentQuote = rand.Intn(len(m.quotes))
	m.breathPhase = 0
	m.cycleCount = 0
	return m.breathe()
}

// Stop 停止冥想
func (m *MeditationMode) Stop() {
	m.active = false
}

// IsActive 是否处于冥想模式
func (m *MeditationMode) IsActive() bool {
	return m.active
}

// breathe 呼吸节奏
func (m *MeditationMode) breathe() tea.Cmd {
	return tea.Tick(m.breathTimer, func(t time.Time) tea.Msg {
		if !m.active {
			return nil
		}
		m.breathPhase = (m.breathPhase + 1) % 2
		m.cycleCount++

		// 每3个循环更换一句名言
		if m.cycleCount%3 == 0 {
			m.currentQuote = rand.Intn(len(m.quotes))
		}

		return meditationBreathMsg{
			phase: m.breathPhase,
			quote: m.quotes[m.currentQuote],
			cycle: m.cycleCount,
		}
	})
}

// View 冥想视图
func (m *MeditationMode) View() string {
	if !m.active {
		return ""
	}

	var sb strings.Builder

	// 创建渐变背景效果
	backgroundStyle := lipgloss.NewStyle().
		Width(80).
		Height(25).
		Background(lipgloss.Color("17")).  // 深蓝色背景
		Foreground(lipgloss.Color("195")). // 浅蓝色文字
		Padding(2, 4)

	// 冥想标题
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("17")).
		Padding(1, 2).
		MarginBottom(2).
		Align(lipgloss.Center).
		Width(72)

	title := "🧘‍♀️ 冥想模式 🧘‍♀️"
	sb.WriteString(titleStyle.Render(title))
	sb.WriteString("\n\n")

	// 呼吸指示器区域
	breathSection := m.renderBreathSection()
	sb.WriteString(breathSection)
	sb.WriteString("\n\n")

	// 哲学名言区域
	quoteSection := m.renderQuoteSection()
	sb.WriteString(quoteSection)
	sb.WriteString("\n\n")

	// 进度指示器
	progressSection := m.renderProgressSection()
	sb.WriteString(progressSection)

	content := sb.String()
	return backgroundStyle.Render(content)
}

// renderBreathSection 渲染呼吸指示器
func (m *MeditationMode) renderBreathSection() string {
	var section strings.Builder

	// 呼吸状态标题
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("120")).
		Bold(true).
		Align(lipgloss.Center).
		Width(72)

	if m.breathPhase == 0 {
		section.WriteString(statusStyle.Render("🌬️  吸气 - 让清新的能量充满身心"))
	} else {
		section.WriteString(statusStyle.Render("🌊  呼气 - 释放所有的紧张与烦恼"))
	}
	section.WriteString("\n\n")

	// 视觉呼吸指示器
	breathIndicator := m.getBreathIndicator()
	indicatorStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(72).
		Foreground(lipgloss.Color("87"))

	section.WriteString(indicatorStyle.Render(breathIndicator))

	return section.String()
}

// renderQuoteSection 渲染名言区域
func (m *MeditationMode) renderQuoteSection() string {
	var section strings.Builder

	quoteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Italic(true).
		Align(lipgloss.Center).
		Width(60).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		MarginLeft(6)

	currentQuote := m.quotes[m.currentQuote]
	section.WriteString(quoteStyle.Render(currentQuote))

	return section.String()
}

// renderProgressSection 渲染进度区域
func (m *MeditationMode) renderProgressSection() string {
	var section strings.Builder

	// 冥想时长
	duration := time.Since(m.startTime)
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("251")).
		Align(lipgloss.Center).
		Width(72)

	timeText := fmt.Sprintf("冥想时长: %02d:%02d | 呼吸循环: %d", minutes, seconds, m.cycleCount)
	section.WriteString(timeStyle.Render(timeText))
	section.WriteString("\n")

	// 退出提示
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Width(72)

	helpText := "按 'q' 或 'esc' 退出冥想模式"
	section.WriteString(helpStyle.Render(helpText))

	return section.String()
}

// getBreathIndicator 获取呼吸指示器
func (m *MeditationMode) getBreathIndicator() string {
	phase := m.cycleCount % 8 // 8个阶段的呼吸动画

	if m.breathPhase == 0 {
		// 吸气动画 - 逐渐扩张
		return m.getInhaleAnimation(phase)
	} else {
		// 呼气动画 - 逐渐收缩
		return m.getExhaleAnimation(phase)
	}
}

// getInhaleAnimation 吸气动画
func (m *MeditationMode) getInhaleAnimation(phase int) string {
	size := phase + 1
	if size > 4 {
		size = 8 - phase
	}

	circle := "●"
	spaces := strings.Repeat(" ", size)

	// 创建渐变效果
	colors := []string{"120", "121", "122", "123", "124"}
	colorIndex := phase % len(colors)

	circleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[colorIndex]))

	return fmt.Sprintf("%s%s%s", spaces, circleStyle.Render(circle), spaces)
}

// getExhaleAnimation 呼气动画
func (m *MeditationMode) getExhaleAnimation(phase int) string {
	size := 4 - (phase % 4)
	if size < 0 {
		size = 0
	}

	circle := "○"
	spaces := strings.Repeat(" ", size)

	// 创建渐变效果
	colors := []string{"124", "123", "122", "121", "120"}
	colorIndex := phase % len(colors)

	circleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[colorIndex]))

	return fmt.Sprintf("%s%s%s", spaces, circleStyle.Render(circle), spaces)
}

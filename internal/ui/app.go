package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stoic/internal/config"
	"github.com/stoic/internal/philosopher"
)

type App struct {
	config             *config.Config
	philosopherManager *philosopher.Manager
	currentPhilosopher philosopher.Philosopher
	messages           []Message
	input              string
	typingInProgress   bool
	currentMessage     string
	typingMessage      string
	typingIndex        int
	lastTypeTime       int64
	commandHandler     *CommandHandler
	meditationMode     *MeditationMode
	inMeditationMode   bool
}

type Message struct {
	Content   string
	FromUser  bool
	Timestamp string
}

type philosopherResponseMsg struct {
	content string
}

type typingEffectMsg struct {
	content string
	index   int
}

func NewApp(cfg *config.Config, manager *philosopher.Manager) *App {
	app := &App{
		config:             cfg,
		philosopherManager: manager,
		currentPhilosopher: manager.GetCurrentPhilosopher(),
		messages:           []Message{},
	}
	app.commandHandler = NewCommandHandler(app)
	return app
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 冥想模式处理
	if a.inMeditationMode && a.meditationMode != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				a.inMeditationMode = false
				if a.meditationMode != nil {
					a.meditationMode.Stop()
				}
				a.messages = append(a.messages, Message{
					Content:   "🧘‍♀️ 已退出冥想模式，欢迎回来！",
					FromUser:  false,
					Timestamp: "",
				})
				return a, nil
			case tea.KeyRunes:
				if string(msg.Runes) == "q" || string(msg.Runes) == "Q" {
					a.inMeditationMode = false
					if a.meditationMode != nil {
						a.meditationMode.Stop()
					}
					a.messages = append(a.messages, Message{
						Content:   "🧘‍♀️ 已退出冥想模式，欢迎回来！",
						FromUser:  false,
						Timestamp: "",
					})
					return a, nil
				}
			}
		case meditationBreathMsg:
			// 冥想模式下的呼吸消息，继续计时
			return a, a.meditationMode.breathe()
		}
		// 冥想模式下不处理其他消息
		return a, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		case tea.KeyEnter:
			if a.input != "" && !a.typingInProgress {
				// 检查是否为命令
				if isCommand, cmd := a.commandHandler.ProcessCommand(a.input); isCommand {
					a.input = ""
					return a, cmd
				}
				// 普通对话输入
				return a.handleUserInput()
			}
		case tea.KeyBackspace:
			if len(a.input) > 0 {
				a.input = a.input[:len(a.input)-1]
			}
		default:
			if !a.typingInProgress {
				a.input += msg.String()
			}
		}
	case philosopherResponseMsg:
		a.handlePhilosopherResponse(msg.content)
		return a, a.startTypingEffect(msg.content)
	case typingEffectMsg:
		a.updateTypingEffect(msg)
		if a.typingIndex < len(a.currentMessage) {
			return a, a.continueTypingEffect()
		} else {
			a.finishTypingEffect()
		}
	case meditationBreathMsg:
		// 普通模式下的冥想消息忽略
		return a, nil
	}

	return a, nil

}

// View 渲染视图
func (a *App) View() string {
	// 冥想模式优先
	if a.inMeditationMode && a.meditationMode != nil {
		return a.meditationMode.View()
	}

	if a.config.UI.Theme == "calm" {
		return a.calmView()
	}
	return a.defaultView()
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}
	return nil
}

func (a *App) calmView() string {
	var sb strings.Builder

	// 渐变标题效果
	titleColors := []string{"86", "84", "82", "80"}
	titleText := "🧘 Stoic - Shell Philosopher"

	// 创建渐变效果
	for i, char := range titleText {
		colorIndex := i % len(titleColors)
		charStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(titleColors[colorIndex])).
			Background(lipgloss.Color("235"))
		sb.WriteString(charStyle.Render(string(char)))
	}
	sb.WriteString("\n\n")
	philStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Italic(true)

	// 确保currentPhilosopher不为nil
	if a.currentPhilosopher != nil {
		philInfo := fmt.Sprintf("Speaking with %s (%s)",
			a.currentPhilosopher.Name(),
			a.currentPhilosopher.School())
		sb.WriteString(philStyle.Render(philInfo))
	} else {
		sb.WriteString(philStyle.Render("No philosopher selected"))
	}
	sb.WriteString("\n\n")

	messageStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1).
		Width(80).
		Height(15)

	var messages strings.Builder
	for _, msg := range a.messages {
		if msg.FromUser {
			userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("87"))
			messages.WriteString(userStyle.Render("You: " + msg.Content))
		} else {
			if a.currentPhilosopher != nil {
				philStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
				messages.WriteString(philStyle.Render(a.currentPhilosopher.Name() + ": " + msg.Content))
			} else {
				// 如果没有哲学家，使用通用名称
				philStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
				messages.WriteString(philStyle.Render("Philosopher: " + msg.Content))
			}
		}
		messages.WriteString("\n")
	}

	if a.typingInProgress && a.typingMessage != "" && a.currentPhilosopher != nil {
		typingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
		messages.WriteString(typingStyle.Render(a.currentPhilosopher.Name() + ": " + a.typingMessage))
		messages.WriteString("█")
	}

	sb.WriteString(messageStyle.Render(messages.String()))
	sb.WriteString("\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1).
		Width(80)

	inputText := a.input
	if !a.typingInProgress {
		inputText += "█"
	}

	sb.WriteString(inputStyle.Render("Your thoughts: " + inputText))
	sb.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	sb.WriteString(helpStyle.Render("esc: quit | enter: send | backspace: delete"))

	return sb.String()
}

func (a *App) defaultView() string {
	return a.calmView()
}

func (a *App) handleUserInput() (tea.Model, tea.Cmd) {

	a.messages = append(a.messages, Message{
		Content:   a.input,
		FromUser:  true,
		Timestamp: "",
	})

	userInput := a.input
	a.input = ""
	a.typingInProgress = true

	return a, a.getPhilosopherResponse(userInput)
}

func (a *App) handlePhilosopherResponse(response string) {
	a.currentMessage = response
	a.typingIndex = 0
	a.typingMessage = ""
	a.typingInProgress = true
}

func (a *App) getPhilosopherResponse(userInput string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		response, err := a.currentPhilosopher.Respond(ctx, userInput)
		if err != nil {
			response = "I apologize, but I need a moment to gather my thoughts..."
		}
		return philosopherResponseMsg{content: response}
	}
}

func (a *App) startTypingEffect(content string) tea.Cmd {
	return func() tea.Msg {
		return typingEffectMsg{content: content, index: 0}
	}
}
func (a *App) continueTypingEffect() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return typingEffectMsg{content: a.currentMessage, index: a.typingIndex}
	})
}

func (a *App) updateTypingEffect(msg typingEffectMsg) {
	if msg.index < len(a.currentMessage) {
		a.typingMessage = a.currentMessage[:msg.index+1]
		a.typingIndex = msg.index + 1
	}
}

func (a *App) finishTypingEffect() {
	a.messages = append(a.messages, Message{
		Content:   a.currentMessage,
		FromUser:  false,
		Timestamp: "",
	})
	a.typingInProgress = false
	a.currentMessage = ""
	a.typingMessage = ""
	a.typingIndex = 0
}

package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/narcilee7/stoic/internal/config"
	"github.com/narcilee7/stoic/internal/philosopher"
)

type App struct {
	config             *config.Config
	philosopherManager *philosopher.Manager
	currentPhilosopher philosopher.Philosopher
	messages           []Message
	input              string
	typingInProgress   bool
	currentMessage     string
}

type Message struct {
	Content   string
	FromUser  bool
	Timestamp string
}

type philosopherResponseMsg struct {
	content string
}

func NewApp(cfg *config.Config, manager *philosopher.Manager) *App {
	return &App{
		config:             cfg,
		philosopherManager: manager,
		currentPhilosopher: manager.GetCurrentPhilosopher(),
		messages:           []Message{},
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		case tea.KeyEnter:
			if a.input != "" && !a.typingInProgress {
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
	}

	return a, nil
}

func (a *App) View() string {
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

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235")).
		Padding(1, 2).
		MarginBottom(1)

	sb.WriteString(titleStyle.Render("🧘 Stoic - Shell Philosopher"))
	sb.WriteString("\n\n")

	philStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Italic(true)

	philInfo := fmt.Sprintf("Speaking with %s (%s)",
		a.currentPhilosopher.Name(),
		a.currentPhilosopher.School())
	sb.WriteString(philStyle.Render(philInfo))
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
			philStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
			messages.WriteString(philStyle.Render(a.currentPhilosopher.Name() + ": " + msg.Content))
		}
		messages.WriteString("\n")
	}

	if a.typingInProgress && a.currentMessage != "" {
		typingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
		messages.WriteString(typingStyle.Render(a.currentPhilosopher.Name() + ": " + a.currentMessage))
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
	a.messages = append(a.messages, Message{
		Content:   response,
		FromUser:  false,
		Timestamp: "",
	})
	a.typingInProgress = false
	a.currentMessage = ""
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

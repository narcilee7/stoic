package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Command interface {
	Name() string
	Description() string
	Aliases() []string
	Execute(app *App, args []string) (bool, tea.Cmd)
	Help() string
}

type CommandResult struct {
	Success  bool
	Message  string
	Continue bool
	Command  tea.Cmd
}

type CommandHandler struct {
	app      *App
	commands map[string]Command
}

func NewCommandHandler(app *App) *CommandHandler {
	ch := &CommandHandler{app: app, commands: make(map[string]Command)}
	ch.registerBuiltinCommands()
	return ch
}

func (ch *CommandHandler) registerBuiltinCommands() {
	// 基础命令
	ch.RegisterCommand(&HelpCommand{})
	ch.RegisterCommand(&ClearCommand{})
	ch.RegisterCommand(&QuitCommand{})

	// 哲学家相关命令
	ch.RegisterCommand(&PhilosophersCommand{})
	ch.RegisterCommand(&SwitchCommand{})
	ch.RegisterCommand(&QuoteCommand{})
	ch.RegisterCommand(&MeditationCommand{})
}

func (ch *CommandHandler) RegisterCommand(cmd Command) {
	ch.commands[cmd.Name()] = cmd
	for _, alias := range cmd.Aliases() {
		ch.commands[alias] = cmd
	}
}

func (ch *CommandHandler) GetCommand(name string) (Command, bool) {
	cmd, exists := ch.commands[name]
	return cmd, exists
}

func (ch *CommandHandler) ProcessCommand(input string) (bool, tea.Cmd) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "/") {
		return false, nil
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false, nil
	}

	commandName := parts[0]
	args := parts[1:]

	// 查找命令
	cmd, exists := ch.GetCommand(commandName)
	if !exists {
		return ch.handleUnknownCommand(commandName)
	}

	// 执行命令
	return cmd.Execute(ch.app, args)
}

func (ch *CommandHandler) handleUnknownCommand(command string) (bool, tea.Cmd) {
	suggestion := ch.findSimilarCommand(command)

	var message string
	if suggestion != "" {
		message = fmt.Sprintf("❓ 未知命令: %s\n💡 你是想说 %s 吗?\n📖 输入 /help 查看所有可用命令", command, suggestion)
	} else {
		message = fmt.Sprintf("❓ 未知命令: %s\n📖 输入 /help 查看所有可用命令", command)
	}

	ch.app.messages = append(ch.app.messages, Message{
		Content:   message,
		FromUser:  false,
		Timestamp: "",
	})

	return true, nil
}

func (ch *CommandHandler) findSimilarCommand(input string) string {
	// 简单的相似度检查
	// TODO: 考虑更复杂的相似度算法
	commands := []string{"/help", "/clear", "/quit", "/meditate", "/philosophers", "/switch", "/quote"}

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, input) || strings.HasPrefix(input, cmd) {
			return cmd
		}
	}
	return ""
}

type HelpCommand struct{}

func (h *HelpCommand) Name() string { return "/help" }

func (h *HelpCommand) Description() string { return "显示帮助信息" }

func (h *HelpCommand) Aliases() []string { return []string{"/h", "/?"} }

func (h *HelpCommand) Execute(app *App, args []string) (bool, tea.Cmd) {
	// 如果用户请求特定命令的帮助
	if len(args) > 0 {
		return h.showCommandHelp(app, args[0])
	}

	// 显示一般帮助
	helpText := `🧘 Stoic - Shell Philosopher 命令帮助

📋 基础命令:
  /help, /h, /?     - 显示此帮助信息
  /clear, /c        - 清屏
  /quit, /q, /exit  - 退出程序
  
🧘‍♀️ 冥想相关:
  /meditate, /m, /zen, /calm, /breathe - 进入冥想模式
  
🧠 哲学家相关:
  /philosophers, /ph, /list, /philos   - 列出可用哲学家
  /switch <name>, /s, /change, /select - 切换哲学家
  /quote, /q, /wisdom, /inspire, /say  - 获取随机哲学名言
  
💡 提示:
  - 所有命令都以 '/' 开头
  - 在冥想模式下按 'q' 或 'esc' 退出
  - 使用 /help <命令> 查看特定命令的详细帮助
  - 使用 Tab 键可以自动补全命令`

	app.messages = append(app.messages, Message{
		Content:   helpText,
		FromUser:  false,
		Timestamp: "",
	})

	return true, nil
}

func (h *HelpCommand) Help() string {
	return "显示所有可用命令的详细帮助信息"
}

// showCommandHelp 显示特定命令的详细帮助
func (h *HelpCommand) showCommandHelp(app *App, commandName string) (bool, tea.Cmd) {
	// 移除可能的 '/' 前缀
	commandName = strings.TrimPrefix(commandName, "/")

	// 构建完整命令名
	fullCommandName := "/" + commandName

	// 查找命令处理器中的命令
	if cmd, exists := app.commandHandler.GetCommand(fullCommandName); exists {
		helpText := fmt.Sprintf(`📖 命令帮助: %s

📝 描述: %s

🔗 别名: %s

💡 用法: %s

📋 详细说明:
%s`,
			cmd.Name(),
			cmd.Description(),
			strings.Join(cmd.Aliases(), ", "),
			cmd.Name(),
			cmd.Help())

		app.messages = append(app.messages, Message{
			Content:   helpText,
			FromUser:  false,
			Timestamp: "",
		})
	} else {
		app.messages = append(app.messages, Message{
			Content:   fmt.Sprintf("❓ 未找到命令: /%s\n📖 使用 /help 查看所有可用命令", commandName),
			FromUser:  false,
			Timestamp: "",
		})
	}

	return true, nil
}

type ClearCommand struct{}

func (c *ClearCommand) Name() string { return "/clear" }

func (c *ClearCommand) Description() string { return "清空对话历史" }

func (c *ClearCommand) Aliases() []string { return []string{"/c", "/cls"} }

func (c *ClearCommand) Execute(app *App, args []string) (bool, tea.Cmd) {
	app.messages = []Message{}
	app.messages = append(app.messages, Message{
		Content:   "🧹 对话历史已清除，让我们开始新的思考之旅...",
		FromUser:  false,
		Timestamp: "",
	})
	return true, nil
}

func (c *ClearCommand) Help() string {
	return "清除所有对话历史"
}

type QuitCommand struct{}

func (q *QuitCommand) Name() string { return "/quit" }

func (q *QuitCommand) Description() string { return "退出程序" }

func (q *QuitCommand) Aliases() []string {
	return []string{"/q", "/exit", "/bye"}
}

func (q *QuitCommand) Execute(app *App, args []string) (bool, tea.Cmd) {
	app.messages = append(app.messages, Message{
		Content:   "👋 感谢使用 Stoic - Shell Philosopher，再见！",
		FromUser:  false,
		Timestamp: "",
	})

	return true, tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tea.Quit
	})
}

func (q *QuitCommand) Help() string {
	return "退出程序"
}

// MeditationCommand 冥想命令
type MeditationCommand struct{}

func (m *MeditationCommand) Name() string {
	return "/meditate"
}

func (m *MeditationCommand) Description() string {
	return "进入冥想模式"
}

func (m *MeditationCommand) Aliases() []string {
	return []string{"/m", "/zen", "/calm", "/breathe"}
}

func (m *MeditationCommand) Execute(app *App, args []string) (bool, tea.Cmd) {
	if app.meditationMode == nil {
		app.meditationMode = NewMeditationMode()
	}

	app.inMeditationMode = true
	app.messages = append(app.messages, Message{
		Content:   "🧘‍♀️ 进入冥想模式... 深呼吸，让心灵归于平静",
		FromUser:  false,
		Timestamp: "",
	})

	return true, app.meditationMode.Start()
}

func (m *MeditationCommand) Help() string {
	return "进入冥想模式，体验呼吸指导和哲学思考"
}

type PhilosophersCommand struct{}

func (p *PhilosophersCommand) Name() string {
	return "/philosophers"
}

func (p *PhilosophersCommand) Description() string {
	return "列出可用哲学家"
}

func (p *PhilosophersCommand) Aliases() []string {
	return []string{"/ph", "/list", "/philos"}
}

func (p *PhilosophersCommand) Execute(app *App, args []string) (bool, tea.Cmd) {
	philosophers := app.philosopherManager.ListPhilosophers()
	current := app.currentPhilosopher

	var list strings.Builder
	list.WriteString("🧠 可用哲学家:\n\n")

	for _, name := range philosophers {
		if name == "stoic" {
			marker := "  "
			if current.Name() == "Marcus Aurelius" {
				marker = "✓ "
			}
			list.WriteString(fmt.Sprintf("%s%s - Marcus Aurelius (斯多葛学派)\n", marker, name))
		} else if name == "taoist" {
			marker := "  "
			if current.Name() == "Laozi" {
				marker = "✓ "
			}
			list.WriteString(fmt.Sprintf("%s%s - Laozi (道家思想)\n", marker, name))
		}
	}

	list.WriteString("\n💡 使用 /switch <name> 切换哲学家")

	app.messages = append(app.messages, Message{
		Content:   list.String(),
		FromUser:  false,
		Timestamp: "",
	})
	return true, nil
}

func (p *PhilosophersCommand) Help() string {
	return "列出所有可用的哲学家及其学派信息"
}

// SwitchCommand 切换哲学家命令
type SwitchCommand struct{}

func (s *SwitchCommand) Name() string {
	return "/switch"
}

func (s *SwitchCommand) Description() string {
	return "切换哲学家"
}

func (s *SwitchCommand) Aliases() []string {
	return []string{"/s", "/change", "/select"}
}

func (s *SwitchCommand) Execute(app *App, args []string) (bool, tea.Cmd) {
	if len(args) == 0 {
		app.messages = append(app.messages, Message{
			Content:   "❓ 请指定哲学家名称。使用 /philosophers 查看可用选项。",
			FromUser:  false,
			Timestamp: "",
		})
		return true, nil
	}

	philosopherName := args[0]
	err := app.philosopherManager.SetCurrentPhilosopher(philosopherName)
	if err != nil {
		app.messages = append(app.messages, Message{
			Content:   fmt.Sprintf("❌ 哲学家 '%s' 未找到。使用 /philosophers 查看可用选项。", philosopherName),
			FromUser:  false,
			Timestamp: "",
		})
		return true, nil
	}

	app.currentPhilosopher = app.philosopherManager.GetCurrentPhilosopher()
	app.messages = append(app.messages, Message{
		Content:   fmt.Sprintf("✅ 已切换到 %s (%s)", app.currentPhilosopher.Name(), app.currentPhilosopher.School()),
		FromUser:  false,
		Timestamp: "",
	})
	return true, nil
}

func (s *SwitchCommand) Help() string {
	return "切换到指定的哲学家，如: /switch stoic"
}

// QuoteCommand 随机名言命令
type QuoteCommand struct{}

func (q *QuoteCommand) Name() string {
	return "/quote"
}

func (q *QuoteCommand) Description() string {
	return "获取随机哲学名言"
}

func (q *QuoteCommand) Aliases() []string {
	return []string{"/q", "/wisdom", "/inspire", "/say"}
}

func (q *QuoteCommand) Execute(app *App, args []string) (bool, tea.Cmd) {
	quotes := []string{
		"The happiness of your life depends upon the quality of your thoughts. - Marcus Aurelius",
		"You have power over your mind — not outside events. Realize this, and you will find strength. - Marcus Aurelius",
		"The best revenge is to be unlike him who performed the injury. - Marcus Aurelius",
		"Nature does not hurry, yet everything is accomplished. - Laozi",
		"The journey of a thousand miles begins with one step. - Laozi",
		"Knowing others is wisdom, knowing yourself is Enlightenment. - Laozi",
		"He who conquers others is strong; he who conquers himself is mighty. - Laozi",
		"The way of the sage is to act but not to compete. - Laozi",
		"Accept the things to which fate binds you, and love the people with whom fate brings you together. - Marcus Aurelius",
		"The soul becomes dyed with the colour of its thoughts. - Marcus Aurelius",
	}

	// 简单的随机选择（基于时间）
	quoteIndex := time.Now().Unix() % int64(len(quotes))
	quote := quotes[quoteIndex]

	app.messages = append(app.messages, Message{
		Content:   fmt.Sprintf("🌟 %s", quote),
		FromUser:  false,
		Timestamp: "",
	})
	return true, nil
}

func (q *QuoteCommand) Help() string {
	return "获取一条随机的哲学名言，带来智慧的启发"
}

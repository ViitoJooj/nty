package screens

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	configOptions = []string{"Register AI credentials", "Language"}
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9e0000")).Background(lipgloss.Color("#ffffff")).Bold(true)
)

type ConfigMenuScreen struct {
	cursor int
}

func NewConfigMenuScreen() *ConfigMenuScreen {
	return &ConfigMenuScreen{}
}

func (screen *ConfigMenuScreen) Init() tea.Cmd {
	return nil
}

func (screen *ConfigMenuScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return screen, nil
	}

	switch key.String() {
	case "up":
		if screen.cursor > 0 {
			screen.cursor--
		}
	case "down":
		if screen.cursor < len(configOptions)-1 {
			screen.cursor++
		}
	case "enter":
		switch screen.cursor {
		case 0:
			return NewClaudeLoginScreen(), nil
		case 1:
			return NewLanguageScreen(), nil
		}
	}

	return screen, nil
}

func (screen *ConfigMenuScreen) View() string {
	var builder strings.Builder

	builder.WriteString("Escolha o que configurar:\n\n")

	for i, opt := range configOptions {
		prefix := "  "

		if i == screen.cursor {
			prefix = "> "
			builder.WriteString(prefix)
			builder.WriteString(selectedStyle.Render(opt))
		} else {
			builder.WriteString(prefix)
			builder.WriteString(opt)
		}

		builder.WriteRune('\n')
	}

	return builder.String()
}

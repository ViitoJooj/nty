package screens

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/ViitoJooj/nty/internal/config"
)

var languages = []struct{ code, label string }{
	{"en", "English"},
	{"pt", "Português"},
}

type LanguageScreen struct {
	cursor int
	status string
	done   bool
}

func NewLanguageScreen() Screen {
	return &LanguageScreen{}
}

func (screen *LanguageScreen) Init() tea.Cmd {
	return nil
}

func (screen *LanguageScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
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
		if screen.cursor < len(languages)-1 {
			screen.cursor++
		}
	case "enter":
		if err := config.SaveLang(languages[screen.cursor].code); err != nil {
			screen.status = "falhou ao salvar: " + err.Error()
			return screen, nil
		}
		screen.status = "Idioma salvo: " + languages[screen.cursor].label
		screen.done = true
		return screen, tea.Quit
	}

	return screen, nil
}

func (screen *LanguageScreen) View() string {
	if screen.done {
		return screen.status + "\n"
	}

	var builder strings.Builder
	builder.WriteString("Escolha o idioma:\n\n")

	for i, lang := range languages {
		if i == screen.cursor {
			builder.WriteString("> ")
			builder.WriteString(selectedStyle.Render(lang.label))
		} else {
			builder.WriteString("  ")
			builder.WriteString(lang.label)
		}
		builder.WriteRune('\n')
	}

	return builder.String()
}

package screens

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/ViitoJooj/nty/internal/config"
)

type GithubEmailScreen struct {
	username string
	input    textinput.Model
}

func NewGithubEmailScreen(username string) *GithubEmailScreen {
	textInput := textinput.New()
	textInput.Placeholder = "Enter your GitHub email"
	textInput.Focus()

	return &GithubEmailScreen{username: username, input: textInput}
}

func (screen *GithubEmailScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (screen *GithubEmailScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if ok && key.String() == "enter" {
		if err := config.SaveGithub(screen.username, screen.input.Value()); err != nil {
			fmt.Println("erro ao salvar config:", err)
		}
		return screen, tea.Quit
	}

	var cmd tea.Cmd
	screen.input, cmd = screen.input.Update(msg)

	return screen, cmd
}

func (screen *GithubEmailScreen) View() string {
	return "GitHub Email:\n\n" + screen.input.View()
}

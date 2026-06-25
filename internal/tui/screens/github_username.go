package screens

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type GithubUsernameScreen struct {
	input textinput.Model
}

func NewGithubUsernameScreen() *GithubUsernameScreen {
	textInput := textinput.New()
	textInput.Placeholder = "Enter your GitHub username"
	textInput.Focus()

	return &GithubUsernameScreen{input: textInput}
}

func (screen *GithubUsernameScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (screen *GithubUsernameScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if ok && key.String() == "enter" {
		return NewGithubEmailScreen(screen.input.Value()), nil
	}

	var cmd tea.Cmd
	screen.input, cmd = screen.input.Update(msg)

	return screen, cmd
}

func (screen *GithubUsernameScreen) View() string {
	return "GitHub Username:\n\n" + screen.input.View()
}

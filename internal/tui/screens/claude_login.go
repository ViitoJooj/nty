package screens

import (
	"os/exec"
	"runtime"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/ViitoJooj/nty/internal/config"
	"github.com/ViitoJooj/nty/internal/services"
)

type ClaudeLoginScreen struct {
	loginURL string
	pkce     services.PKCE
	input    textinput.Model
	status   string
	done     bool
}

func NewClaudeLoginScreen() Screen {
	loginURL, pkce, err := services.Start()
	if err != nil {
		return &ClaudeLoginScreen{status: "erro ao iniciar login: " + err.Error(), done: true}
	}

	textInput := textinput.New()
	textInput.Placeholder = "Cole o code aqui"
	textInput.Focus()

	openBrowser(loginURL)

	return &ClaudeLoginScreen{loginURL: loginURL, pkce: pkce, input: textInput}
}

func (screen *ClaudeLoginScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (screen *ClaudeLoginScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		var cmd tea.Cmd
		screen.input, cmd = screen.input.Update(msg)
		return screen, cmd
	}

	if key.String() == "enter" && !screen.done {
		tok, err := services.Exchange(screen.input.Value(), screen.pkce)
		if err != nil {
			screen.status = "falhou: " + err.Error()
			return screen, nil
		}

		if err := config.SaveClaudeAuth(tok.AccessToken, tok.RefreshToken, tok.ExpiresAt()); err != nil {
			screen.status = "logou mas falhou ao salvar: " + err.Error()
			return screen, nil
		}

		screen.status = "Claude conectado!"
		screen.done = true
		return screen, tea.Quit
	}

	var cmd tea.Cmd
	screen.input, cmd = screen.input.Update(msg)
	return screen, cmd
}

func (screen *ClaudeLoginScreen) View() string {
	if screen.done {
		return screen.status + "\n"
	}

	out := "Login Claude (assinatura)\n"
	out += "1. Abra (deve ter aberto no browser):\n   " + screen.loginURL + "\n"
	out += "2. Autorize, copie o code da pagina e cole abaixo:\n"
	out += screen.input.View() + "\n"
	if screen.status != "" {
		out += "\n" + screen.status + "\n"
	}
	return out
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd, args = "rundll32", []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd, args = "open", []string{url}
	default:
		cmd, args = "xdg-open", []string{url}
	}
	_ = exec.Command(cmd, args...).Start()
}

package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ViitoJooj/nty/internal/tui/screens"
)

// Run launches the config TUI menu and blocks until it exits.
func Run() error {
	_, err := tea.NewProgram(NewApp()).Run()
	return err
}

// RunLang opens the language selection screen directly.
func RunLang() error {
	_, err := tea.NewProgram(NewAppWith(screens.NewLanguageScreen())).Run()
	return err
}

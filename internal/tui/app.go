package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ViitoJooj/nty/internal/tui/screens"
)

type App struct {
	screen screens.Screen
}

func NewApp() App {
	return App{screen: screens.NewConfigMenuScreen()}
}

func NewAppWith(screen screens.Screen) App {
	return App{screen: screen}
}

func (a App) Init() tea.Cmd {
	return a.screen.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return a, tea.Quit
	}

	screen, cmd := a.screen.Update(msg)
	a.screen = screen

	return a, cmd
}

func (a App) View() tea.View {
	return tea.NewView(a.screen.View())
}

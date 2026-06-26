package tui

import tea "charm.land/bubbletea/v2"

func Run() error {
	_, err := tea.NewProgram(NewApp()).Run()
	return err
}

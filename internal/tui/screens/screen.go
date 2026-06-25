package screens

import tea "charm.land/bubbletea/v2"

type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg) (Screen, tea.Cmd)
	View() string
}

package taskin

import "github.com/charmbracelet/lipgloss"

func Color(c lipgloss.TerminalColor, r string) string {
	/*
		if IsCI() {
			return r
		}
	*/
	return lipgloss.NewStyle().Foreground(c).Render(r)
}

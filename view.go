package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	appNameStyle    = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1).Border(lipgloss.RoundedBorder())
	faintStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
	enumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

func (m *model) View() string {
	var s strings.Builder

	s.WriteString(appNameStyle.Render("Notes App"))
	s.WriteString("\n\n")

	if m.state == errorView && m.err != nil {
		s.WriteString(m.err.Error())
		s.WriteString("\n\n")
		s.WriteString(faintStyle.Render("enter - continue\n\n"))
	}

	if m.state == titleView {
		s.WriteString("Note title: \n\n")
		s.WriteString(m.textinput.View())
		s.WriteString("\n\n")
		s.WriteString(faintStyle.Render("enter - continue, esc = discard\n\n"))
	}

	if m.state == bodyView {
		s.WriteString("Note body: \n\n")
		s.WriteString(m.textarea.View())
		s.WriteString("\n\n")
		s.WriteString(faintStyle.Render("ctrl-s - save note, esc = discard\n\n"))
	}

	if m.state == listView {
		maxLength := 100
		for i, note := range m.notes {
			prefix := " "
			if i == m.listIndex {
				prefix = ">"
			}

			body := strings.ReplaceAll(note.Body, "\n", "")
			if len(body) > maxLength {
				body = body[:maxLength]
			}

			s.WriteString(enumeratorStyle.Render(prefix))
			s.WriteString(note.Title)
			s.WriteString(" | ")
			s.WriteString(faintStyle.Render(body))
			s.WriteString("\n\n")
		}

		s.WriteString(faintStyle.Render("n - new note, q = quit\n\n"))
	}

	return s.String()
}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	titleView
	bodyView
	errorView
)

type model struct {
	state       uint
	store       *Store
	notes       []Note
	currentNote Note
	listIndex   int
	textinput   textinput.Model
	textarea    textarea.Model
	err         error
}

func NewModel(ctx context.Context, store *Store) (*model, error) {
	notes, err := store.GetNotes(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting notes: %w", err)
	}
	return &model{
		state:     listView,
		store:     store,
		notes:     notes,
		textarea:  textarea.New(),
		textinput: textinput.New(),
	}, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	ctx, done := context.WithTimeout(context.Background(), time.Second*10)
	defer done()

	if keyPress, ok := msg.(tea.KeyMsg); ok {
		// handle specific key bindings for different states
		switch m.state {
		case listView:
			{
				switch keyPress.String() {
				case "n":
					{
						m.textinput.SetValue("")
						cmd = m.textinput.Focus()
						cmds = append(cmds, cmd)
						m.currentNote = Note{}
						m.state = titleView
					}
				case "k", tea.KeyUp.String():
					if m.listIndex > 0 {
						m.listIndex -= 1
					}
				case "j", tea.KeyDown.String():
					if m.listIndex < (len(m.notes) - 1) {
						m.listIndex += 1
					}
				case "q", tea.KeyCtrlC.String():
					return m, tea.Quit

				case tea.KeyEnter.String():
					{
						m.currentNote = m.notes[m.listIndex]
						m.textarea.SetValue(m.currentNote.Body)
						m.textarea.Focus()
						m.textarea.CursorEnd()
						m.state = bodyView
					}
				}
			}

		case titleView:
			{
				switch keyPress.Type {
				case tea.KeyEsc:
					m.state = listView

				case tea.KeyEnter:
					{
						if title := m.textinput.Value(); title != "" {
							m.currentNote.Title = title
							m.textarea.SetValue("")
							m.textarea.Focus()
							m.textarea.CursorEnd()
							m.state = bodyView
						}

					}
				}
			}

		case bodyView:
			{
				switch keyPress.Type {
				case tea.KeyEsc:
					m.state = listView

				case tea.KeyCtrlS:
					{
						if body := m.textarea.Value(); body != "" {
							m.currentNote.Body = body
							err := m.store.Upsert(ctx, m.currentNote)
							if err != nil {
								m.err = err
								m.state = errorView
								return m, tea.Batch(cmds...)
							}
							notes, err := m.store.GetNotes(ctx)
							if err != nil {
								m.err = err
								m.state = errorView
								return m, tea.Batch(cmds...)
							}
							m.notes = notes
							m.state = listView
						}

					}
				}
			}

		case errorView:
			{
				switch keyPress.Type {
				case tea.KeyEnter:
					m.state = listView
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	s, err := NewStore(ctx)
	if err != nil {
		log.Fatalf("setting up store: %v", err)
	}

	m, err := NewModel(ctx, s)
	if err != nil {
		log.Fatalf("creating model: %v", err)
	}

	program := tea.NewProgram(m)
	go func() {
		<-ctx.Done()
		program.Quit()
	}()

	if _, err := program.Run(); err != nil {
		log.Fatalf("running tui: %v", err)
	}
}

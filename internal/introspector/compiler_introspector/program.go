package compiler_introspector

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func RunIntrospector() error {
	client, err := compiler_introspection.NewClient()
	if err != nil {
		return err
	}

	err = client.Connect()
	if err != nil {
		return fmt.Errorf("Failed to connect %w", err)
	}

	termenv.HideCursor()
	defer termenv.ShowCursor()
	err = tea.NewProgram(tui.InitialTUIModel(client), tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Printf("Failed to start TUI: %v\n", err)
		os.Exit(1)
	}

	return nil
}

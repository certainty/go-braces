package compiler_introspector

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/frame"
	tea "github.com/charmbracelet/bubbletea"
)

func RunIntrospector() error {
	client, err := compiler_introspection.NewClient()
	if err != nil {
		return err
	}

	// err = client.Connect()
	// if err != nil {
	// 	return fmt.Errorf("Failed to connect %w", err)
	// }

	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return fmt.Errorf("Failed to log to file %w", err)
	}
	defer logFile.Close()

	_, err = tea.NewProgram(frame.New(client), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Printf("Failed to start UI: %v\n", err)
		os.Exit(1)
	}

	return nil
}

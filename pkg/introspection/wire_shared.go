package introspection

import (
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
)

func SetupDirectories(scope string) error {
	if err := os.MkdirAll(filepath.Join(xdg.StateHome, "go-braces", scope), 0700); err != nil {
		return err
	}
	eventSockPath := EventSocketPath(scope)
	controlSockPath := ControlSocketPath(scope)
	os.Remove(eventSockPath)
	os.Remove(controlSockPath)
	return nil
}

func ControlSocketPath(scope string) string {
	return filepath.Join(xdg.StateHome, "go-braces", scope, "control.ipc")
}

func EventSocketPath(scope string) string {
	return filepath.Join(xdg.StateHome, "go-braces", scope, "events.ipc")
}

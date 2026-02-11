package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var initOnce sync.Once

// Log file permissions: owner read/write only (no terminal or other users).
const logFileMode = 0600

// Init initializes the default slog logger to write to the file at path.
// It is safe to call multiple times; only the first call opens the file.
// Creates the parent directory if needed (0755). The log file is created with 0600.
// If opening the file fails, the default logger is set to discard so no messages appear in the terminal.
func Init(path string) {
	initOnce.Do(func() {
		if path == "" {
			slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
			return
		}
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
			return
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, logFileMode)
		if err != nil {
			slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
			return
		}
		slog.SetDefault(slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelInfo})))
	})
}

// LogCLIAction logs a CLI invocation (command name and args).
func LogCLIAction(cmdName string, args []string) {
	slog.Info("cli action", "command", cmdName, "args", args)
}

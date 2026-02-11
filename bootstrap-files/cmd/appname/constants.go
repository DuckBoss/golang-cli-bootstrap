package main

// Version is set at build time via -ldflags "-X main.Version=..."
var Version = "0.0.1"

// DefaultLogPath is the default path for the application log file.
// Set at build time for dev builds via -ldflags "-X main.DefaultLogPath=..."
var DefaultLogPath = "/var/log/appname/appname.log"

const (
	ExitCodeSuccess     = 0
	ExitCodeError       = 1
	ExitCodeInvalidArgs = 2
)

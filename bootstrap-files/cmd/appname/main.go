package main

import (
	"fmt"
	"os"

	"github.com/ownername/appname/internal/config"
	"github.com/ownername/appname/internal/logging"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %v\n", err)
		os.Exit(GetExitCode(err))
	}
	os.Exit(ExitCodeSuccess)
}

func MainApp() *cli.App {
	return &cli.App{
		Name:    "appname",
		Usage:   "appname CLI",
		Version: Version,
		Authors: []*cli.Author{
			{
				Name:  "CLI-Bootstrap",
				Email: "me@example.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-file",
				Aliases: []string{"cfg"},
				Hidden:  true,
				Usage:   "Path to config file",
				EnvVars: []string{"APPNAME_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "environment-file",
				Hidden:  true,
				Aliases: []string{"env"},
				Usage:   "Path to environment file",
				EnvVars: []string{"APPNAME_ENV"},
			},
			&cli.StringFlag{
				Name:    "log",
				Aliases: []string{"l"},
				Hidden:  true,
				Usage:   "Path to log file",
				Value:   DefaultLogPath,
				EnvVars: []string{"APPNAME_LOG"},
			},
		},
		Commands: []*cli.Command{
			CmdVersion(),
		},
		Before: beforeApp,
	}
}

func beforeApp(c *cli.Context) error {
	cfg, err := config.Consolidate(
		c.String("config-file"),
		c.String("environment-file"),
		DefaultLogPath,
		c.String("log"),
	)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	logging.Init(cfg.LogPath)

	cmdName := "default"
	if c.Command != nil {
		cmdName = c.Command.FullName()
	}
	logging.LogCLIAction(cmdName, c.Args().Slice())

	return nil
}

// Run runs the CLI with the given args. Use os.Args from main.
func Run(args []string) error {
	return MainApp().Run(args)
}

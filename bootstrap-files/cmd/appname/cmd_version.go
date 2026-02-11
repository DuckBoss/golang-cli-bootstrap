package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// CmdVersion returns the version subcommand (exported for use from main).
func CmdVersion() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Show version",
		Action: versionAction,
	}
}

func versionAction(c *cli.Context) error {
	fmt.Printf("%s version %s\n", c.App.Name, c.App.Version)
	return nil
}

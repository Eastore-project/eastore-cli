package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// VersionCommand returns the CLI command for printing version information
func VersionCommand(version string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print the version",
		Action: func(cCtx *cli.Context) error {
			fmt.Printf("eastore version %s\n", version)
			return nil
		},
	}
}

package main

import (
	"log"
	"os"

	"github.com/eastore-project/eastore/cmd/eastore/commands"
	"github.com/urfave/cli/v2"
)

var (
	version = "0.1.0"
)

func main() {
	app := &cli.App{
		Name:    "eastore",
		Usage:   "Eastore CLI tool",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "private-key",
				EnvVars:  []string{"PRIVATE_KEY"},
				Required: true,
				Usage:    "Private key for signing transactions",
			},
			&cli.StringFlag{
				Name:    "rpc-url",
				EnvVars: []string{"RPC_URL"},
				Usage:   "RPC URL",
			},
			&cli.StringFlag{
				Name:    "contract",
				EnvVars: []string{"EASTORE_CONTRACT_ADDRESS"},
				Usage:   "Eastore contract address",
			},
		},
		Commands: []*cli.Command{
			commands.VersionCommand(version),
			commands.MakeDealCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

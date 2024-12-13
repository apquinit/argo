package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const version = "1.0.0"

func main() {
	app := &cli.App{
		Name:    "argo",
		Usage:   "a powerful, opinionated web framework for Go",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Display the CLI version",
				Action: func(c *cli.Context) error {
					fmt.Printf("Argo CLI version %s\n", version)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

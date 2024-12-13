package main

import (
	"argo/core/generator"
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
			{
				Name:  "new",
				Usage: "Create a new project",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "Enable verbose output",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("project name is required")
					}
					projectName := c.Args().Get(0)
					verbose := c.Bool("verbose")

					if verbose {
						fmt.Printf("Verbose mode enabled. Creating project: %s\n", projectName)
					}

					if _, err := os.Stat(projectName); !os.IsNotExist(err) {
						return fmt.Errorf("directory '%s' already exists", projectName)
					}

					if err := generator.CreateNewProject(projectName); err != nil {
						return fmt.Errorf("error creating project: %v", err)
					}

					fmt.Printf("Project '%s' created successfully.\n", projectName)
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

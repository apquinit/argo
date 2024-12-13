package main

import (
	"fmt"
	"os"
	"os/exec"

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
				Name:  "serve",
				Usage: "Run the server",
				Action: func(c *cli.Context) error {
					fmt.Println("Starting the server...")

					// Check if server.go exists
					if _, err := os.Stat("server.go"); os.IsNotExist(err) {
						return fmt.Errorf("server.go not found in the current directory")
					}

					// Run the server using `go run`
					cmd := exec.Command("go", "run", "server.go")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Run(); err != nil {
						return fmt.Errorf("failed to start the server: %v", err)
					}

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

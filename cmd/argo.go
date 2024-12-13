package main

import (
	"fmt"
	"os"

	"argo/core/generator"
	"argo/core/migration"
	"argo/core/seeder"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Argo: The Go Framework CLI")
		fmt.Println("Usage: argo <command> [options]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Usage: argo new <project-name>")
			os.Exit(1)
		}
		projectName := os.Args[2]
		generator.CreateNewProject(projectName)
	case "make:controller":
		if len(os.Args) < 3 {
			fmt.Println("Usage: argo make:controller <name>")
			os.Exit(1)
		}
		name := os.Args[2]
		generator.GenerateController(name)
	case "make:migration":
		if len(os.Args) < 3 {
			fmt.Println("Usage: argo make:migration <name>")
			os.Exit(1)
		}
		name := os.Args[2]
		generator.GenerateMigration(name)
	case "make:seeder":
		if len(os.Args) < 3 {
			fmt.Println("Usage: argo make:seeder <name>")
			os.Exit(1)
		}
		name := os.Args[2]
		generator.GenerateSeeder(name)
	case "migrate":
		migration.RunMigrations()
	case "seed":
		seeder.RunSeeders()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

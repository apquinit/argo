package main

import (
	"argo/core/router"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Welcome to Argo!")
		fmt.Println("Usage: argo <command> [options]")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "new":
		createNewProject()
	case "serve":
		serveApp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func createNewProject() {
	fmt.Println("Creating a new Argo project...")
	// Logic for initializing a new project goes here
}

func serveApp() {
	fmt.Println("Starting the Argo server...")
	r := router.InitRouter()
	r.Run(":8080") // Run the server on port 8080
}

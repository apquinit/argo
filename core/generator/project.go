package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateNewProject(projectName string) {
	directories := []string{
		"app/controllers",
		// "app/models",
		// "database/migrations",
		// "database/seeders",
		"public",
		"routes",
		"config",
	}

	// Create project directory
	fmt.Printf("Creating project: %s\n", projectName)
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		return
	}

	// Create subdirectories
	for _, dir := range directories {
		path := filepath.Join(projectName, dir)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
		}
	}

	fmt.Printf("Project %s created successfully!\n", projectName)
}

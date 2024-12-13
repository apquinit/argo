package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func CreateNewProject(projectName string) error {
	// Validate project name
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Define directories and files to create
	directories := []string{
		"app/controllers",
		"start",
		"config",
	}

	files := map[string]string{
		"app/controllers/health.go": "../templates/controllers/health.go.tmpl",
		"start/env.go":              "../templates/env.go.tmpl",
		"start/kernel.go":           "../templates/kernel.go.tmpl",
		"start/routes.go":           "../templates/routes.go.tmpl",
		".env.example":              "../templates/env.example.tmpl",
		".gitignore":                "../templates/gitignore.tmpl",
		"argo.go":                   "../templates/argo.go.tmpl",
	}

	// Check if project directory already exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		return fmt.Errorf("project directory %s already exists", projectName)
	}

	// Create project directory
	fmt.Printf("Creating project: %s\n", projectName)
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		return fmt.Errorf("error creating project directory: %v", err)
	}

	// Create subdirectories
	for _, dir := range directories {
		path := filepath.Join(projectName, dir)
		fmt.Printf("Creating directory: %s\n", path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}

	// Create files from templates
	for target, tmplPath := range files {
		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			return fmt.Errorf("template file %s does not exist", tmplPath)
		}

		targetPath := filepath.Join(projectName, target)
		fmt.Printf("Creating file: %s from template: %s\n", targetPath, tmplPath)
		err := createFileFromTemplate(targetPath, tmplPath, projectName)
		if err != nil {
			return fmt.Errorf("error creating file %s: %v", target, err)
		}
	}

	fmt.Printf("Project %s created successfully!\n", projectName)
	return nil
}

func createFileFromTemplate(targetPath, tmplPath, projectName string) error {
	// Parse template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create target file
	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Populate template data
	data := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}

	// Execute template
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

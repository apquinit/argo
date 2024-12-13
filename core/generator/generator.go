package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func CreateNewProject(projectName string) error {
	directories := []string{
		"app/controllers",
		"start",
		"config",
	}

	files := map[string]string{
		"start/env.go":    "../templates/env.go.tmpl",
		"start/kernel.go": "../templates/kernel.go.tmpl",
		"start/routes.go": "../templates/routes.go.tmpl",
		".env.example":    "../templates/env.example.tmpl",
		".gitignore":      "../templates/gitignore.tmpl",
		"argo.go":         "../templates/argo.go.tmpl",
	}

	// Create project directory
	fmt.Printf("Creating project: %s\n", projectName)
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		return fmt.Errorf("error creating project: %v", err)
	}

	// Create subdirectories
	for _, dir := range directories {
		path := filepath.Join(projectName, dir)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}

	// Create files from templates
	for target, tmplPath := range files {
		targetPath := filepath.Join(projectName, target)
		err := createFileFromTemplate(targetPath, tmplPath, projectName)
		if err != nil {
			return fmt.Errorf("error creating file %s: %v", target, err)
		}
	}

	fmt.Printf("Project %s created successfully!\n", projectName)
	return nil
}

func createFileFromTemplate(targetPath, tmplPath, projectName string) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	data := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

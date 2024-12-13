package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Dependencies []string `yaml:"dependencies"`
}

func CreateNewProject(projectName string) error {
	// Validate project name
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Define directories and files to create
	directories := []string{
		"app/controllers",
		"app/middleware",
		"config",
		"env",
		"internal/context",
		"start",
	}

	files := map[string]string{
		"app/controllers/health.go":       "../templates/app/controllers/health.go.tmpl",
		"app/middleware/logger.go":        "../templates/app/middleware/logger.go.tmpl",
		"config/app.go":                   "../templates/config/app.go.tmpl",
		"config/cors.go":                  "../templates/config/cors.go.tmpl",
		"env/env.go":                      "../templates/env/env.go.tmpl",
		"internal/context/context.go":     "../templates/internal/context/context.go.tmpl",
		"start/cors.go":                   "../templates/start/cors.go.tmpl",
		"start/kernel.go":                 "../templates/start/kernel.go.tmpl",
		"start/middleware.go":             "../templates/start/middleware.go.tmpl",
		"start/routes.go":                 "../templates/start/routes.go.tmpl",
		".env":                            "../templates/env.example.tmpl",
		".env.example":                    "../templates/env.example.tmpl",
		".gitignore":                      "../templates/gitignore.tmpl",
		fmt.Sprintf("%s.go", projectName): "../templates/server.go.tmpl",
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
		err := createFileFromTemplate(targetPath, tmplPath, projectName)
		if err != nil {
			return fmt.Errorf("error creating file %s: %v", target, err)
		}
	}

	// Initialize Go module
	fmt.Println("Initializing Go module...")
	if err := runCommand("go", projectName, "mod", "init", projectName); err != nil {
		return fmt.Errorf("error initializing Go module: %v", err)
	}

	// Run go mod tidy
	fmt.Println("Tidying up Go module...")
	if err := runCommand("go", projectName, "mod", "tidy"); err != nil {
		return fmt.Errorf("error running go mod tidy: %v", err)
	}

	// Load dependencies from YAML
	dependencies, err := loadDependenciesYAML("../dependencies.yaml")
	if err != nil {
		return fmt.Errorf("error loading dependencies: %v", err)
	}

	fmt.Println("Installing dependencies...")
	for _, dep := range dependencies {
		if err := runCommand("go", projectName, "get", dep); err != nil {
			return fmt.Errorf("error installing dependency %s: %v", dep, err)
		}
	}

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

func loadDependenciesYAML(filePath string) ([]string, error) {
	// Read YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	// Unmarshal YAML data
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %v", err)
	}

	return config.Dependencies, nil
}

func runCommand(command, dir string, args ...string) error {
	// Create the command with arguments
	cmd := exec.Command(command, args...)

	// Set the working directory for the command
	cmd.Dir = dir

	// Redirect command output to the standard output and error streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}

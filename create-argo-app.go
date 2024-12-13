package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const version = "1.0.0"

type Config struct {
	Dependencies []string `yaml:"dependencies"`
}

func main() {
	app := &cli.App{
		Name:    "create-argo-app",
		Usage:   "initialize and scaffold new projects for the Argo web framework",
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

					if err := createNewProject(projectName); err != nil {
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

func createNewProject(projectName string) error {
	// Define directories and files to create
	directories := []string{
		"app/controllers",
		"app/middleware",
		"config",
		"start",
	}

	files := map[string]string{
		"app/controllers/health.go": "templates/app/controllers/health.go.tmpl",
		"app/middleware/logger.go":  "templates/app/middleware/logger.go.tmpl",
		"config/app.go":             "templates/config/app.go.tmpl",
		"start/env.go":              "templates/start/env.go.tmpl",
		"start/kernel.go":           "templates/start/kernel.go.tmpl",
		"start/routes.go":           "templates/start/routes.go.tmpl",
		".env":                      "templates/env.example.tmpl",
		".env.example":              "templates/env.example.tmpl",
		".gitignore":                "templates/gitignore.tmpl",
		"server.go":                 "templates/server.go.tmpl",
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

	// Compile cli/argo.go into a binary and copy it to the project
	if err := buildAndCopyBinary(projectName); err != nil {
		return fmt.Errorf("error building and copying binary: %v", err)
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
	dependencies, err := loadDependenciesYAML("dependencies.yaml")
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

func buildAndCopyBinary(projectName string) error {
	// Build the binary from cli/argo.go
	fmt.Println("Building argo binary...")
	cmd := exec.Command("go", "build", "-o", "argo", "./cli/argo.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error building binary: %v", err)
	}

	// Move the binary to the new project's directory
	dest := filepath.Join(projectName, "argo")
	fmt.Printf("Copying binary to %s\n", dest)
	if err := os.Rename("argo", dest); err != nil {
		return fmt.Errorf("error copying binary: %v", err)
	}

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

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

func loadDependenciesYAML(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %v", err)
	}

	return config.Dependencies, nil
}

func runCommand(command, dir string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func GenerateController(name string) {
	filename := filepath.Join("app/controllers", name+".go")

	templateContent := `package controllers

import "github.com/gin-gonic/gin"

func {{.Name}}Controller(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello from {{.Name}}Controller"})
}
`

	tmpl, err := template.New("controller").Parse(templateContent)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, struct{ Name string }{Name: name})
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	fmt.Printf("Controller %s created successfully!\n", filename)
}

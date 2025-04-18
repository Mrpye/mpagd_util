package mpagd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/Mrpye/mpagd_util/project_template"
)

// ListTemplates lists all available templates in the embedded templates directory.
func ListTemplates() ([]Template, error) {
	var templates []Template

	// Define the directory to read templates from
	dir := "."
	entries, err := project_template.Templates.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Create a new APJFile instance for processing templates
	apj := NewAPJFile("")

	// Iterate through directory entries to find YAML templates
	for _, entry := range entries {
		ext := path.Ext(entry.Name())
		if ext == ".yaml" {
			// Read the template file
			file, err := project_template.Templates.ReadFile(entry.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to open template '%s': %w", entry.Name(), err)
			}

			// Load the template content into the APJFile instance
			err = apj.LoadYAMLFromString(file)
			if err != nil {
				return nil, fmt.Errorf("error reading template '%s': %w", entry.Name(), err)
			}

			// Append the template details to the list
			templates = append(templates, CreateTemplate(
				entry.Name(),
				apj.Description,
				"Project Template",
			))
		}
	}

	return templates, nil
}

// CreateProjectFromTemplate creates a new project file from a specified template.
func CreateProjectFromTemplate(projectFile string, templateName string) error {
	// Validate the project file name
	if projectFile == "" {
		return fmt.Errorf("project file name is empty")
	}

	// Ensure the project file has the correct extension
	if !strings.HasSuffix(projectFile, ".apj") {
		return fmt.Errorf("project file name must end with .apj")
	}

	// Extract the file name and directory path
	fileName := path.Base(projectFile)
	filePath := path.Dir(projectFile)

	// Create the directory if it doesn't exist
	if _, err := os.Stat(filePath); err != nil {
		err := os.MkdirAll(filePath, fs.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", filePath, err)
		}
	}

	// Ensure the template name has the correct extension
	if !strings.HasSuffix(templateName, ".yaml") {
		templateName = templateName + ".yaml"
	}

	// Read the template file
	file, err := project_template.Templates.ReadFile(templateName)
	if err != nil {
		return fmt.Errorf("failed to open template '%s': %w", templateName, err)
	}

	// Load the template content into a new APJFile instance
	apj := NewAPJFile(projectFile)
	err = apj.LoadYAMLFromString(file)
	if err != nil {
		return fmt.Errorf("error reading template '%s': %w", templateName, err)
	}

	// Set the file path and write the project file
	apj.FilePath = fileName
	err = apj.WriteAPJ(projectFile)
	if err != nil {
		return fmt.Errorf("failed to write project file '%s': %w", projectFile, err)
	}

	return nil
}

// GetStats returns statistics about the project.
func (apj *APJFile) DisplayStats() {

	fmt.Println("Project Statistics:")
	fmt.Printf("Blocks: %d\n", len(apj.Blocks))
	fmt.Printf("Sprites: %d\n", len(apj.Sprites))
	fmt.Printf("Screens: %d\n", len(apj.Screens))
	fmt.Printf("Objects: %d\n", len(apj.Objects))
	fmt.Printf("Maps: %d\n", len(apj.Map.Map))
	fmt.Printf("Fonts: %d\n", len(apj.Fonts))

}

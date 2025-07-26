package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project file.",
	Long:  `Manage project file.`,
}

// Cmd_Backup creates a command to back up a project file.
func Cmd_Backup() *cobra.Command {
	var code bool
	var cmd = &cobra.Command{
		Use:   "backup [project file]",
		Short: "Backup the project file.",
		Long:  `Create a backup of the project file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", len(args))
			}
			filePath := args[0]
			// Check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", filePath)
			}
			mpagd.LogMessage("Cmd_Backup", fmt.Sprintf("Creating backup for file: %s", filePath), "info", noColor)
			// Assuming APJFile is a struct with a method BackupProjectFile
			apj := mpagd.NewAPJFile(filePath)
			err := apj.BackupProjectFile(code)
			if err != nil {
				return fmt.Errorf("error creating backup: %v", err)
			}
			mpagd.LogMessage("Cmd_Backup", fmt.Sprintf("Backup created successfully for file: %s", filePath), "ok", noColor)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&code, "code", "c", false, "backup code files")
	return cmd
}

// Cmd_Restore creates a command to restore a project file from the last backup.
func Cmd_Restore() *cobra.Command {
	var code bool
	var cmd = &cobra.Command{
		Use:   "restore [project file]",
		Short: "Restore the project file from the last backup.",
		Args:  cobra.ExactArgs(1),
		Long:  `Restore the project file from the most recent backup.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", len(args))
			}
			filePath := args[0]
			// Check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				mpagd.LogMessage("Cmd_Restore", fmt.Sprintf("Project File %s does not exist", filePath), "warning", noColor)
				// Return an error if the file does not exist
				//return fmt.Errorf("file %s does not exist", filePath)
			}
			mpagd.LogMessage("Cmd_Restore", fmt.Sprintf("Restoring last backup for file: %s", filePath), "info", noColor)
			// Assuming APJFile is a struct with a method RestoreLastBackup
			//conver path to /

			filePath = strings.ReplaceAll(filePath, "\\", "/")
			backupDir := path.Dir(filePath)
			backupDir = path.Join(backupDir, "backups")
			apj := mpagd.NewAPJFile(filePath)

			lastBackup, err := apj.RestoreLastBackup(backupDir, code)
			if err != nil {
				return fmt.Errorf("error restoring the last backup: %v", err)

			}

			mpagd.LogMessage("Cmd_Restore", fmt.Sprintf("Last backup restored successfully: %s", lastBackup), "ok", noColor)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&code, "code", "c", false, "restore code files as well")
	return cmd
}

// Cmd_PurgeBackup creates a command to purge all backup files.
func Cmd_PurgeBackup() *cobra.Command {
	return &cobra.Command{
		Use:   "purge [project file]",
		Short: "Purge all backup files.",
		Args:  cobra.ExactArgs(1),
		Long:  `Delete all backup files for the project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", len(args))
			}
			filePath := args[0]
			// Check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", filePath)
			}
			mpagd.LogMessage("Cmd_PurgeBackup", fmt.Sprintf("Purging backup files for: %s", filePath), "info", noColor)
			// Assuming APJFile is a struct with a method RestoreLastBackup
			filePath = strings.ReplaceAll(filePath, "\\", "/")
			backupDir := path.Dir(filePath)
			backupDir = path.Join(backupDir, "backups")
			apj := mpagd.NewAPJFile(filePath)

			err := apj.PurgeBackupFiles(backupDir)
			if err != nil {
				return fmt.Errorf("error purging backup files: %v", err)

			}

			mpagd.LogMessage("Cmd_PurgeBackup", "All backup files purged successfully.", "ok", noColor)
			return nil
		},
	}
}

// Cmd_ListBackups creates a command to list all backup files.
func Cmd_ListBackups() *cobra.Command {
	return &cobra.Command{
		Use:   "backups [project file]",
		Short: "List all backup files.",
		Args:  cobra.ExactArgs(1),
		Long:  `Display a list of all backup files for the project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", len(args))
			}
			filePath := args[0]
			// Check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", filePath)
			}
			mpagd.LogMessage("Cmd_ListBackups", fmt.Sprintf("Listing backup files for: %s", filePath), "info", noColor)
			// Assuming APJFile is a struct with a method ListBackupProjectFiles
			filePath = strings.ReplaceAll(filePath, "\\", "/")
			backupDir := path.Dir(filePath)
			backupDir = path.Join(backupDir, "backups")
			apj := mpagd.NewAPJFile(filePath)

			backupFiles, err := apj.ListBackupProjectFiles(backupDir)
			if err != nil {
				return fmt.Errorf("error listing backup files: %v", err)
			}

			if len(backupFiles) == 0 {
				mpagd.LogMessage("Cmd_ListBackups", "No backup files found.", "warning", noColor)
				return nil
			}

			mpagd.LogMessage("Cmd_ListBackups", "Available backup files:", "ok", noColor)
			for _, file := range backupFiles {
				fmt.Println(file)
			}
			return nil
		},
	}
}

// Cmd_AutoBackup creates a command to enable automatic backups.
func Cmd_AutoBackup() *cobra.Command {
	var code bool
	var cmd = &cobra.Command{
		Use:   "auto-backup [project file]",
		Short: "Enable automatic backups for the project file.",
		Args:  cobra.ExactArgs(1),
		Long:  `Start monitoring the project file for changes and create backups automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			// Check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", filePath)
			}
			mpagd.LogMessage("Cmd_AutoBackup", fmt.Sprintf("Starting auto-backup for file: %s", filePath), "info", noColor)
			mpagd.LogMessage("Cmd_AutoBackup", "Press Ctrl+C to exit", "info", noColor)
			// Assuming APJFile is a struct with a method MonitorFileChanges
			apj := mpagd.NewAPJFile(filePath)
			apj.MonitorFileChanges(code)
			mpagd.LogMessage("Cmd_AutoBackup", "Auto-backup started successfully.", "ok", noColor)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&code, "code", "c", false, "backup code files")
	return cmd
}

// Cmd_SaveAsYAML creates a command to save the project file as a YAML file.
func Cmd_SaveAsYAML() *cobra.Command {
	return &cobra.Command{
		Use:   "save [project file] [output yaml file]",
		Short: "Save the project file as a YAML file.",
		Args:  cobra.ExactArgs(2),
		Long:  `Convert the project file to a YAML file and save it to the specified location.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			outputYAMLFile := args[1]
			mpagd.LogMessage("Cmd_SaveAsYAML", fmt.Sprintf("Saving project as YAML file: %s", outputYAMLFile), "info", noColor)
			// Check if the project file exists
			if _, err := os.Stat(projectFile); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", projectFile)
			}

			apj := mpagd.NewAPJFile(projectFile)
			if err := apj.SaveAsYAML(outputYAMLFile); err != nil {
				return fmt.Errorf("failed to save project as YAML: %v", err)
			}

			mpagd.LogMessage("Cmd_SaveAsYAML", fmt.Sprintf("Project saved successfully as YAML file: %s", outputYAMLFile), "ok", noColor)
			return nil
		},
	}
}

// Cmd_LoadYAML creates a command to load a project from a YAML file.
func Cmd_LoadYAML() *cobra.Command {
	return &cobra.Command{
		Use:   "load [yaml file] [output project file]",
		Short: "Load a project from a YAML file.",
		Args:  cobra.ExactArgs(2),
		Long:  `Load a project from a YAML file and save it as a project file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			yamlFile := args[0]
			outputProjectFile := args[1]
			mpagd.LogMessage("Cmd_LoadYAML", fmt.Sprintf("Loading project from YAML file: %s", yamlFile), "info", noColor)
			// Check if the YAML file exists
			if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", yamlFile)
			}

			apj := mpagd.NewAPJFile(outputProjectFile)
			if err := apj.LoadYAML(yamlFile); err != nil {
				return fmt.Errorf("failed to load project from YAML: %v", err)
			}
			mpagd.LogMessage("Cmd_LoadYAML", fmt.Sprintf("Project loaded successfully from YAML file: %s", yamlFile), "ok", noColor)
			return nil
		},
	}
}

// Cmd_ImportAGD creates a command to import all AGD elements into the project file.
func Cmd_ImportAGD() *cobra.Command {
	var replace bool
	var cmd = &cobra.Command{
		Use:   "import [project file] [agd file] [[output project file]]",
		Short: "Import all AGD elements into the project file.",
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		Long:  `Import all AGD elements (blocks, sprites, screens, etc.) into the project file and save the updated project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			agdFile := args[1]
			outputFile := projectFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outputFile = args[2]
			}
			mpagd.LogMessage("Cmd_ImportAGD", fmt.Sprintf("Importing AGD elements from %s to %s", agdFile, projectFile), "info", noColor)
			// Open and read the project file
			apj := mpagd.NewAPJFile(projectFile)

			if err := apj.ReadAPJ(); err != nil {
				mpagd.LogMessage("Cmd_ImportAGD", fmt.Sprintf("Project file %s does not exist", projectFile), "warning", noColor)
				replace = true // Set replace to false if the project file does not exist
			}

			// If the output file is the same as the input file, create a backup
			if outputFile == projectFile {
				// Check if the file exists
				if _, err := os.Stat(projectFile); os.IsNotExist(err) {
					mpagd.LogMessage("Cmd_ImportAGD", "New Project skipping backup", "warning", noColor)
				} else {
					err := apj.BackupProjectFile(false)
					if err != nil {
						return fmt.Errorf("failed to create backup: %w", err)
					}
					mpagd.LogMessage("Cmd_ImportAGD", fmt.Sprintf("Backup created: %s.bak", projectFile), "ok", noColor)
				}
			}

			// Import all AGD elements
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptionsFalse() // Import all elements
			options.SetOwOptions(replace, replace, replace, replace, replace, replace, replace, replace, replace)
			if err := apj.ImportAGD(agdFile, options); err != nil {
				return fmt.Errorf("failed to import AGD elements: %w", err)
			}

			// Write the updated project file
			if err := apj.WriteAPJ(outputFile); err != nil {
				return fmt.Errorf("failed to write updated project file: %w", err)
			}

			mpagd.LogMessage("Cmd_ImportAGD", fmt.Sprintf("AGD elements imported successfully. Updated project file saved to %s", outputFile), "ok", noColor)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace existing elements in the project file")
	return cmd
}

// Cmd_ImportAGDSelective creates a command to import selected AGD elements into the project file.
func Cmd_ImportAGDSelective() *cobra.Command {
	var replace bool
	var importBlocks, importSprites, importScreens, importObjects, importMaps, importFonts, importULAPalette, importWindows, importKeys bool

	var cmd = &cobra.Command{
		Use:   "import-selective [project file] [agd file] [[output project file]]",
		Short: "Import selected AGD elements into the project file.",
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		Long:  `Import selected AGD elements (blocks, sprites, screens, etc.) into the project file and save the updated project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			agdFile := args[1]
			outputFile := projectFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outputFile = args[2]
			}
			mpagd.LogMessage("Cmd_ImportAGDSelective", fmt.Sprintf("Importing AGD elements from %s to %s", agdFile, projectFile), "info", noColor)
			// Open and read the project file
			apj := mpagd.NewAPJFile(projectFile)
			if err := apj.ReadAPJ(); err != nil {
				mpagd.LogMessage("Cmd_ImportAGDSelective", fmt.Sprintf("Project file %s does not exist", projectFile), "warning", noColor)
				replace = true // Set replace to false if the project file does not exist
			}

			// If the output file is the same as the input file, create a backup
			if outputFile == projectFile {
				if _, err := os.Stat(projectFile); os.IsNotExist(err) {
					mpagd.LogMessage("Cmd_ImportAGDSelective", "New Project skipping backup", "warning", noColor)
				} else {
					err := apj.BackupProjectFile(false)
					if err != nil {
						return fmt.Errorf("failed to create backup: %w", err)
					}
					mpagd.LogMessage("Cmd_ImportAGDSelective", fmt.Sprintf("Backup created: %s.bak", projectFile), "ok", noColor)
				}
			}

			// Import selected AGD elements
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(!importWindows, !importKeys, !importBlocks, !importSprites, !importObjects, !importScreens, !importMaps, !importFonts, !importULAPalette)
			options.SetOwOptions(replace, replace, replace, replace, replace, replace, replace, replace, replace)

			if err := apj.ImportAGD(agdFile, options); err != nil {
				return fmt.Errorf("failed to import AGD elements: %w", err)
			}

			// Write the updated project file
			if err := apj.WriteAPJ(outputFile); err != nil {
				return fmt.Errorf("failed to write updated project file: %w", err)
			}

			mpagd.LogMessage("Cmd_ImportAGDSelective", fmt.Sprintf("AGD elements imported successfully. Updated project file saved to %s", outputFile), "ok", noColor)
			return nil
		},
	}

	// Define flags for selective import
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace existing elements in the project file")
	cmd.Flags().BoolVar(&importWindows, "window", false, "Import window")
	cmd.Flags().BoolVar(&importKeys, "keys", false, "Import keys")
	cmd.Flags().BoolVar(&importBlocks, "blocks", false, "Import blocks")
	cmd.Flags().BoolVar(&importSprites, "sprites", false, "Import sprites")
	cmd.Flags().BoolVar(&importScreens, "screens", false, "Import screens")
	cmd.Flags().BoolVar(&importObjects, "objects", false, "Import objects")
	cmd.Flags().BoolVar(&importMaps, "maps", false, "Import maps")
	cmd.Flags().BoolVar(&importFonts, "fonts", false, "Import fonts")
	cmd.Flags().BoolVar(&importULAPalette, "ula-palette", false, "Import ULA palette")

	return cmd
}

// Cmd_ListTemplates creates a command to list all available project templates.
func Cmd_ListTemplates() *cobra.Command {
	return &cobra.Command{
		Use:   "templates",
		Short: "List all available project templates.",
		Long:  `Display a list of all available project templates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mpagd.LogMessage("Cmd_ListTemplates", "Listing available templates", "info", noColor)
			templates, err := mpagd.ListTemplates()
			if err != nil {
				return fmt.Errorf("failed to list templates: %v", err)
			}
			if len(templates) == 0 {
				mpagd.LogMessage("Cmd_ListTemplates", "No Templates Available", "warning", noColor)
				return nil
			}
			mpagd.LogMessage("Cmd_ListTemplates", "Available Templates:", "ok", noColor)
			for _, template := range templates {
				fmt.Printf("Name: %s, Description:%s, Type:%s \n", template.Name, template.Description, template.Type)
			}
			mpagd.LogMessage("Cmd_ListTemplates", "Templates listed successfully.", "ok", noColor)
			return nil
		},
	}
}

// Cmd_CreateProjectFromTemplate creates a command to create a project from a template.
func Cmd_CreateProjectFromTemplate() *cobra.Command {
	return &cobra.Command{
		Use:   "create [project file] [template name]",
		Short: "Create a project from a template.",
		Args:  cobra.ExactArgs(2),
		Long:  `Create a new project file based on the specified template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			templateName := args[1]
			mpagd.LogMessage("Cmd_CreateProjectFromTemplate", fmt.Sprintf("Creating project from template '%s'", templateName), "info", noColor)
			err := mpagd.CreateProjectFromTemplate(projectFile, templateName)
			if err != nil {
				return fmt.Errorf("failed to create project from template: %v", err)
			}
			mpagd.LogMessage("Cmd_CreateProjectFromTemplate", fmt.Sprintf("Project created successfully from template '%s'", templateName), "ok", noColor)
			return nil
		},
	}
}

// Cmd_ProjectStats creates a command to display statistics about the project.
func Cmd_ProjectStats() *cobra.Command {
	return &cobra.Command{
		Use:   "stats [project file]",
		Short: "Display statistics about the project.",
		Args:  cobra.ExactArgs(1),
		Long:  `Show statistics about the project, such as the number of blocks, sprites, screens, etc.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			// Check if the project file exists
			if _, err := os.Stat(projectFile); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", projectFile)
			}

			mpagd.LogMessage("Cmd_ProjectStats", fmt.Sprintf("Gathering stats for project file: %s", projectFile), "info", noColor)
			apj := mpagd.NewAPJFile(projectFile)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read project file: %v", err)
			}

			mpagd.LogMessage("Cmd_ProjectStats", "Project Statistics:", "info", noColor)
			apj.DisplayStats() // Display the statistics in a user-friendly format
			mpagd.LogMessage("Cmd_ProjectStats", "Project stats displayed successfully.", "ok", noColor)
			return nil
		},
	}
}

// Cmd_CreateReadme creates a command to generate a README file for the project.
func Cmd_CreateReadme() *cobra.Command {
	return &cobra.Command{
		Use:   "create-readme [project file] [output readme file]",
		Short: "Generate a README file for the project.",
		Args:  cobra.ExactArgs(1),
		Long:  `Build a Markdown README file for the project using its data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			outputReadme := ""
			if len(args) > 1 {
				outputReadme = args[1]
			}
			if outputReadme == "" {
				outputReadme = filepath.Join(filepath.Dir(projectFile), "docs", "README.md")
			}
			// Check if the project file exists
			if _, err := os.Stat(projectFile); os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", projectFile)
			}
			mpagd.LogMessage("Cmd_CreateReadme", fmt.Sprintf("Generating README for project file: %s", projectFile), "info", noColor)
			apj := mpagd.NewAPJFile(projectFile)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read project file: %v", err)
			}
			// Build project info JSON

			// Use a version of BuildProjectInfoJson that returns the JSON string
			jsonStr, err := apj.BuildProjectInfoJson()
			if err != nil {
				return fmt.Errorf("failed to build project info JSON: %v", err)
			}
			err = mpagd.BuildProjectReadme(outputReadme, []byte(jsonStr))
			if err != nil {
				return fmt.Errorf("failed to build README: %v", err)
			}

			mpagd.LogMessage("Cmd_CreateReadme", fmt.Sprintf("README file created: %s", outputReadme), "ok", noColor)
			return nil
		},
	}
}

// Initialize all commands and add them to the root command.
func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(Cmd_Backup())
	projectCmd.AddCommand(Cmd_Restore())
	projectCmd.AddCommand(Cmd_PurgeBackup())
	projectCmd.AddCommand(Cmd_ListBackups())
	projectCmd.AddCommand(Cmd_AutoBackup())
	projectCmd.AddCommand(Cmd_SaveAsYAML())
	projectCmd.AddCommand(Cmd_LoadYAML())
	projectCmd.AddCommand(Cmd_ImportAGD())
	projectCmd.AddCommand(Cmd_ImportAGDSelective())
	projectCmd.AddCommand(Cmd_ListTemplates())
	projectCmd.AddCommand(Cmd_CreateProjectFromTemplate())
	projectCmd.AddCommand(Cmd_ProjectStats())
	projectCmd.AddCommand(Cmd_CreateReadme())
}

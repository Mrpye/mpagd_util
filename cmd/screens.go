package cmd

import (
	"fmt"
	"strconv"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// screensCmd is the root command for managing screens in the project.
var screensCmd = &cobra.Command{
	Use:   "screens",
	Short: "Manage screens in the project.",
	Long:  `Manage screens in the project. This command provides options to import screens from AGD files and render screens to bitmap images.`,
}

// Cmd_ImportScreens creates a command to import Screens elements from an AGD file into an APJ file.
func Cmd_ImportScreens() *cobra.Command {
	var replaceExistingScreens bool

	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Screens elements from an AGD file into an APJ file.",
		Long:  `Imports Screens elements from an AGD file into the corresponding section of an APJ file.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse input arguments
			apjFilePath := args[0]
			agdFilePath := args[1]
			outputFilePath := apjFilePath
			if len(args) == 3 {
				outputFilePath = args[2]
			}

			// Log the start of the import process
			mpagd.LogMessage("Cmd_ImportScreens", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info", noColor)

			// Read the APJ file
			apj := mpagd.NewAPJFile(apjFilePath)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// Create a backup if the output file is the same as the input file
			if outputFilePath == apjFilePath {
				if err := apj.BackupProjectFile(false); err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_ImportScreens", fmt.Sprintf("Backup created: %s", outputFilePath), "ok", noColor)
			}

			// Configure import options
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, true, true, true, true, true, false, true, true)
			options.SetOwOptions(false, false, false, false, false, false, replaceExistingScreens, false, false)

			// Perform the import operation
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Screens from AGD file: %w", err)
			}

			// Log the successful completion of the import
			mpagd.LogMessage("Cmd_ImportScreens", fmt.Sprintf("Screens imported successfully. Updated APJ file saved to %s", outputFilePath), "ok", noColor)
			return nil
		},
	}

	// Add a flag to control whether existing Screens should be replaced
	cmd.Flags().BoolVarP(&replaceExistingScreens, "replace", "r", false, "Replace the current Screens")
	return cmd
}

// Cmd_RenderScreensToBitmap creates a command to render a screen from an APJ file to a bitmap image.
func Cmd_RenderScreensToBitmap() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "render-bmp [apj file] [screen id] [bitmap file]",
		Short: "Render a screen from an APJ file to a bitmap image.",
		Long:  `Renders a specific screen from an APJ file and saves it as a bitmap image in PNG format.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse input arguments
			apjFilePath := args[0]
			screenID := args[1]
			outputFilePath := args[2]

			// Log the start of the render process
			mpagd.LogMessage("Cmd_RenderScreens", fmt.Sprintf("Starting render for file: %s, screen ID: %s", apjFilePath, screenID), "info", noColor)

			// Read the APJ file
			apj := mpagd.NewAPJFile(apjFilePath)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// Convert screen ID to an integer
			screenIndex, err := strconv.Atoi(screenID)
			if err != nil {
				return fmt.Errorf("invalid screen ID: %w", err)
			}

			// Render the screen to a bitmap image
			if err := apj.RenderScreenToBitmap(uint8(screenIndex), outputFilePath); err != nil {
				return fmt.Errorf("failed to render screen to bitmap: %w", err)
			}

			// Log the successful completion of the render
			mpagd.LogMessage("Cmd_RenderScreens", fmt.Sprintf("Screen rendered successfully. Output saved to %s", outputFilePath), "ok", noColor)
			return nil
		},
	}
	return cmd
}

// Cmd_ReorderScreens creates a command to reorder screens in an APJ file.
func Cmd_ReorderScreens() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "reorder [apj file] [order] [[output file]]",
		Short: "Reorder screens in an APJ file.",
		Long:  `Reorders the screens in an APJ file based on the provided order. The order should be a comma-separated list of screen indices.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			inFile := args[0]
			orderStr := args[1]
			outFile := inFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outFile = args[2]
			}

			mpagd.LogMessage("Cmd_ReorderScreens", fmt.Sprintf("Starting reorder of screens for file: %s", inFile), "info", noColor)

			apjFile := mpagd.NewAPJFile(inFile)
			err := apjFile.ReadAPJ()
			if err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// If the output file is the same as the input file, create a backup
			if outFile == inFile {
				err := apjFile.BackupProjectFile(false)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_ReorderScreens", fmt.Sprintf("Backup created: %s", inFile), "ok", noColor)
			}

			// Convert the order string to a slice of integers
			order := mpagd.CSVToIntSlice(orderStr)

			// Reorder sprites
			err = apjFile.ReorderScreens(order)
			if err != nil {
				return fmt.Errorf("failed to reorder sprites: %w", err)
			}

			// Write the updated APJ file
			err = apjFile.WriteAPJ(outFile)
			if err != nil {
				return fmt.Errorf("failed to write APJ file: %w", err)
			}

			mpagd.LogMessage("Cmd_ReorderScreens", "Reorder completed successfully", "ok", noColor)
			return nil
		},
	}
	return cmd
}

// Initialize the screens command and its subcommands
func init() {
	RootCmd.AddCommand(screensCmd)
	screensCmd.AddCommand(Cmd_ImportScreens())
	screensCmd.AddCommand(Cmd_RenderScreensToBitmap())
	screensCmd.AddCommand(Cmd_ReorderScreens())
}

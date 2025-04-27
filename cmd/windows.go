package cmd

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// windowsCmd is the root command for managing windows in the MPAGD project.
var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "Manage windows in the MPAGD project.",
	Long:  `Manage windows in the MPAGD project.`,
}

// Cmd_ImportWindows creates a command to import Windows elements from an AGD file into an APJ file.
func Cmd_ImportWindows() *cobra.Command {
	var replace bool

	// Define the import command.
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Windows elements from an AGD file into an APJ file.",
		Long:  `Imports Windows elements from an AGD file into the corresponding section of an APJ file.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse command arguments.
			apjFilePath := args[0]
			agdFilePath := args[1]
			outputFilePath := apjFilePath
			if len(args) == 3 {
				outputFilePath = args[2]
			}

			// Log the start of the import process.
			mpagd.LogMessage("Cmd_ImportWindows", fmt.Sprintf("Starting import for file: %s", apjFilePath), "ok", noColor)

			// Read the APJ file.
			apj := mpagd.NewAPJFile(apjFilePath)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// Create a backup if the output file is the same as the input file.
			if outputFilePath == apjFilePath {
				err := apj.BackupProjectFile(false)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_ImportWindows", fmt.Sprintf("Backup created: %s", outputFilePath), "ok", noColor)
			}

			// Configure import options.
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, false, true, true, true, true, true, true, true)
			options.SetOwOptions(false, replace, false, false, false, false, false, false, false)

			// Perform the import operation.
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Windows from AGD file: %w", err)
			}

			// Log the successful completion of the import process.
			mpagd.LogMessage("Cmd_ImportWindows", fmt.Sprintf("Windows imported successfully. Updated APJ file saved to %s", outputFilePath), "ok", noColor)
			return nil
		},
	}

	// Add a flag to control the replace behavior.
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current Windows")
	return cmd
}

// init initializes the windows command and adds it to the root command.
func init() {
	RootCmd.AddCommand(windowsCmd)
	windowsCmd.AddCommand(Cmd_ImportWindows())
}

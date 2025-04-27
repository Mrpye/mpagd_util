package cmd

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// objectsCmd is the root command for managing objects in the MPAGD project.
var objectsCmd = &cobra.Command{
	Use:   "objects ",
	Short: "Manage objects in the MPAGD project.",
	Long:  `Manage objects in the MPAGD project.`,
}

// Cmd_ImportObjects creates a command to import Objects elements from an AGD file into an APJ file.
func Cmd_ImportObjects() *cobra.Command {
	var replace bool

	// Define the import command.
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Objects elements from an AGD file into an APJ file.",
		Long:  `Imports Objects elements from an AGD file into the corresponding section of an APJ file.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse input arguments.
			apjFilePath := args[0]
			agdFilePath := args[1]
			outputFilePath := apjFilePath
			if len(args) == 3 {
				outputFilePath = args[2]
			}

			// Log the start of the import process.
			mpagd.LogMessage("Cmd_ImportObjects", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info", noColor)

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
				mpagd.LogMessage("Cmd_ImportObjects", fmt.Sprintf("Backup created: %s", outputFilePath), "ok", noColor)
			}

			// Configure import options.
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, true, true, true, true, false, true, true, true)
			options.SetOwOptions(false, false, false, false, false, replace, false, false, false)

			// Perform the import operation.
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Objects from AGD file: %w", err)
			}

			// Log the successful completion of the import process.
			mpagd.LogMessage("Cmd_ImportObjects", fmt.Sprintf("Objects imported successfully. Updated APJ file saved to %s", outputFilePath), "ok", noColor)
			return nil
		},
	}

	// Add a flag to control whether existing objects should be replaced.
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current Objects")
	return cmd
}

// init initializes the objects command and adds it to the root command.
func init() {
	RootCmd.AddCommand(objectsCmd)
	objectsCmd.AddCommand(Cmd_ImportObjects())
}

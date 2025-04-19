package cmd

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// ulaCmd is the root command for managing ULA in the MPAGD project.
var ulaCmd = &cobra.Command{
	Use:   "ula",
	Short: "Manage ula in the MPAGD project.",
	Long:  `Manage ula in the MPAGD project.`,
}

// Cmd_ImportULAPalette creates a command to import ULAPalette elements from an AGD file into an APJ file.
func Cmd_ImportULAPalette() *cobra.Command {
	var replace bool

	// Define the command and its arguments
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import ULAPalette elements from an AGD file into an APJ file.",
		Long:  `Imports ULAPalette elements from an AGD file into the corresponding section of an APJ file.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse arguments
			apjFilePath := args[0]
			agdFilePath := args[1]
			outputFilePath := apjFilePath
			if len(args) == 3 {
				outputFilePath = args[2]
			}

			// Log the start of the import process
			mpagd.LogMessage("Cmd_ImportULAPalette", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info")

			// Read the APJ file
			apj := mpagd.NewAPJFile(apjFilePath)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// Create a backup if the output file is the same as the input file
			if outputFilePath == apjFilePath {
				err := apj.BackupProjectFile(false)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_ImportULAPalette", fmt.Sprintf("Backup created: %s", outputFilePath), "ok")
			}

			// Configure import options
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, true, true, true, true, true, true, true, true)
			options.SetOwOptions(false, false, false, false, false, false, false, false, replace)

			// Perform the import
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import ULAPalette from AGD file: %w", err)
			}

			// Log the success of the import process
			mpagd.LogMessage("Cmd_ImportULAPalette", fmt.Sprintf("ULAPalette imported successfully. Updated APJ file saved to %s", outputFilePath), "ok")
			return nil
		},
	}

	// Add a flag to control whether to replace the current ULAPalette
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current ULAPalette")
	return cmd
}

// init adds the ula command and its subcommands to the root command.
func init() {
	RootCmd.AddCommand(ulaCmd)
	ulaCmd.AddCommand(Cmd_ImportULAPalette())
}

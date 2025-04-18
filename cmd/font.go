package cmd

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// fontsCmd represents the parent command for font-related operations.
var fontsCmd = &cobra.Command{
	Use:   "fonts",
	Short: "Manage fonts in the MPAGD project.",
	Long:  `Manage fonts in the MPAGD project.`,
}

// Cmd_ImportFont creates a command to import font elements from an AGD file into an APJ file.
func Cmd_ImportFont() *cobra.Command {
	var replace bool

	// Define the import command with its usage and description.
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Font elements from an AGD file into an APJ file.",
		Long:  `Imports Font elements from an AGD file into the corresponding section of an APJ file.`,
		Args:  cobra.MinimumNArgs(2), // Require at least two arguments: APJ file and AGD file.
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse input arguments.
			apjFilePath := args[0]
			agdFilePath := args[1]
			outputFilePath := apjFilePath
			if len(args) == 3 {
				outputFilePath = args[2]
			}

			// Log the start of the import process.
			mpagd.LogMessage("Cmd_ImportFont", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info")

			// Load the APJ file.
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
				mpagd.LogMessage("Cmd_ImportFont", fmt.Sprintf("Backup created: %s", outputFilePath), "ok")
			}

			// Configure import options.
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, true, true, true, true, true, true, true, false)
			options.SetOwOptions(false, false, false, false, false, false, false, false, replace)

			// Perform the import operation.
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Font from AGD file: %w", err)
			}

			// Log the successful completion of the import process.
			mpagd.LogMessage("Cmd_ImportFont", fmt.Sprintf("Font imported successfully. Updated APJ file saved to %s", outputFilePath), "ok")
			return nil
		},
	}

	// Add a flag to control whether existing fonts should be replaced.
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current Font")
	return cmd
}

// init initializes the fonts command and adds it to the root command.
func init() {
	RootCmd.AddCommand(fontsCmd)
	fontsCmd.AddCommand(Cmd_ImportFont())
}

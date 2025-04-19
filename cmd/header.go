package cmd

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// headerCmd defines the "header" command for managing headers in MPAGD projects.
var headerCmd = &cobra.Command{
	Use:   "header",
	Short: "Manage header in the MPAGD project.",
	Long:  `Manage header in the MPAGD project.`,
}

// Cmd_ImportHeader creates the "import" subcommand for importing header elements.
func Cmd_ImportHeader() *cobra.Command {
	var replace bool

	// Define the "import" command.
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Header elements from an AGD file into an APJ file.",
		Long:  `Imports Header elements from an AGD file into the corresponding section of an APJ file.`,
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
			mpagd.LogMessage("Cmd_ImportHeader", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info")

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
				mpagd.LogMessage("Cmd_ImportHeader", fmt.Sprintf("Backup created: %s", outputFilePath), "ok")
			}

			// Configure import options.
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(false, true, true, true, true, true, true, true, true)
			options.SetOwOptions(replace, false, false, false, false, false, false, false, false)

			// Perform the import operation.
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Header from AGD file: %w", err)
			}

			// Log the successful completion of the import process.
			mpagd.LogMessage("Cmd_ImportHeader", fmt.Sprintf("Header imported successfully. Updated APJ file saved to %s", outputFilePath), "ok")
			return nil
		},
	}

	// Add the "replace" flag to the command.
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current Header")
	return cmd
}

// init initializes the "header" command and its subcommands.
func init() {
	RootCmd.AddCommand(headerCmd)
	headerCmd.AddCommand(Cmd_ImportHeader())
}

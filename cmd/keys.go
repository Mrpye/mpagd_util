package cmd

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys in the MPAGD project.",
	Long:  `Manage keys in the MPAGD project.`,
}

func Cmd_ImportKeys() *cobra.Command {
	var replace bool
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Keys elements from an AGD file into an APJ file.",
		Long:  `Imports Keys elements from an AGD file into the corresponding section of an APJ file.`,
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
			mpagd.LogMessage("Cmd_ImportKeys", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info")

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
				mpagd.LogMessage("Cmd_ImportKeys", fmt.Sprintf("Backup created: %s", outputFilePath), "ok")
			}

			// Configure import options
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, true, true, false, true, true, true, true, true)
			options.SetOwOptions(false, false, false, replace, false, false, false, false, false)

			// Perform the import operation
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Keys from AGD file: %w", err)
			}

			// Log the success of the operation
			mpagd.LogMessage("Cmd_ImportKeys", fmt.Sprintf("Keys imported successfully. Updated APJ file saved to %s", outputFilePath), "ok")
			return nil
		},
	}

	// Add the replace flag to the command
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current Keys")
	return cmd
}

func init() {
	// Add the keys command and its subcommands to the root command
	RootCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(Cmd_ImportKeys())
}

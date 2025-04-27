package cmd

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Manage map in the MPAGD project.",
	Long:  `Manage map in the MPAGD project.`,
}

func Cmd_ImportMap() *cobra.Command {
	var replace bool
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Map elements from an AGD file into an APJ file.",
		Long:  `Imports Map elements from an AGD file into the corresponding section of an APJ file.`,
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
			mpagd.LogMessage("Cmd_ImportMap", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info", noColor)

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
				mpagd.LogMessage("Cmd_ImportMap", fmt.Sprintf("Backup created: %s", outputFilePath), "ok", noColor)
			}

			// Configure import options
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, true, true, true, true, true, true, false, true)
			options.SetOwOptions(false, false, false, false, false, false, false, replace, false)

			// Perform the import operation
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Map from AGD file: %w", err)
			}

			// Log the successful completion of the import
			mpagd.LogMessage("Cmd_ImportMap", fmt.Sprintf("Map imported successfully. Updated APJ file saved to %s", outputFilePath), "ok", noColor)
			return nil
		},
	}

	// Add the replace flag to the command
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current Map")
	return cmd
}

func init() {
	// Register the map command and its subcommands
	RootCmd.AddCommand(mapCmd)
	mapCmd.AddCommand(Cmd_ImportMap())
}

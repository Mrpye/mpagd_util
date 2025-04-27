/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var appVersion = "0.1.4"
var noColor = false

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mpagd_util",
	Short: `MPAGD Utility CLI version` + appVersion,
	Long:  `A command line interface for MPAGD utility functions.`,
}

func SetNoColor(no_color bool) {
	noColor = no_color
}

// GenerateDoc creates a new command to generate CLI documentation.
func GenerateDoc() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gen_docs",
		Short: "Generate CLI documentation",
		Long:  `This command generates markdown documentation for the CLI commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Log the start of the document generation process
			mpagd.LogMessage("GenerateDoc", "Generating documents...", "info", noColor)

			// Ensure the output directory exists
			os.MkdirAll("./documents", os.ModePerm)

			// Generate markdown documentation
			err := doc.GenMarkdownTree(RootCmd, "./documents")
			if err != nil {
				// Log an error message if document generation fails
				mpagd.LogMessage("GenerateDoc", "Failed to generate documents", "error", noColor)
				return err
			}

			// Log a success message upon successful document generation
			mpagd.LogMessage("GenerateDoc", "Documents generated successfully", "ok", noColor)
			return nil
		},
	}
	return cmd
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		mpagd.LogMessage("Execute", fmt.Sprintf("Error: %v", err), "error", noColor)
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")
	RootCmd.AddCommand(GenerateDoc())

}

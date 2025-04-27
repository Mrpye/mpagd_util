package cmd

import (
	"fmt"
	"strconv"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// blocksCmd represents the base command for managing blocks in the MPAGD project.
var blocksCmd = &cobra.Command{
	Use:   "blocks",
	Short: "Manage blocks in the MPAGD project.",
	Long:  `Manage blocks in the MPAGD project.`,
}

// blocksRotateCmd represents the subcommand for rotating blocks.
var blocksRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate blocks in the MPAGD project.",
	Long:  `Rotate blocks in the MPAGD project.`,
}

// Cmd_RotateBlockCCW90 creates a command to rotate a block 90 degrees counter-clockwise.
func Cmd_RotateBlockCCW90() *cobra.Command {
	var repeat int
	var add bool
	var startRange uint8
	var endRange uint8

	var cmd = &cobra.Command{
		Use:   "ccw [project file] [block number] [[output file]]",
		Short: "Rotate a block 90 degrees counter-clockwise.",
		Long:  `Rotate a block 90 degrees counter-clockwise.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			inFile := args[0]
			blockNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid block number: %w", err)
			}
			outFile := inFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outFile = args[2]
			}

			mpagd.LogMessage("Cmd_RotateBlockCCW90", fmt.Sprintf("Starting rotation for file: %s", inFile), "info")

			apjFile := mpagd.NewAPJFile(inFile)
			err = apjFile.ReadAPJ()
			if err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// If the output file is the same as the input file, create a backup
			if outFile == inFile {
				err := apjFile.BackupProjectFile(false)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_RotateBlockCCW90", fmt.Sprintf("Backup created: %s", inFile), "ok")
			}

			// Rotate blocks
			currentBlock := blockNumber
			if startRange == 0 && endRange == 0 {
				for i := 0; i < repeat; i++ {
					processedBlock, err := apjFile.RotateBlock(uint8(currentBlock), true, add)
					if err != nil {
						return fmt.Errorf("failed to rotate block: %w", err)
					}
					if add {
						currentBlock = int(processedBlock)
						mpagd.LogMessage("Cmd_RotateBlockCCW90", fmt.Sprintf("created new rotated block: %v", currentBlock), "ok")
					}
				}
			} else {
				for r := startRange; r < endRange; r++ {
					currentBlock = int(r)
					for i := 0; i < repeat; i++ {
						processedBlock, err := apjFile.RotateBlock(uint8(currentBlock), true, add)
						if err != nil {
							return fmt.Errorf("failed to rotate block: %w", err)
						}
						if add {
							currentBlock = int(processedBlock)
							mpagd.LogMessage("Cmd_RotateBlockCCW90", fmt.Sprintf("created new rotated block: %v", currentBlock), "ok")
						}
					}
				}
			}

			// Write the updated APJ file
			err = apjFile.WriteAPJ(outFile)
			if err != nil {
				return fmt.Errorf("failed to write APJ file: %w", err)
			}

			mpagd.LogMessage("Cmd_RotateBlockCCW90", fmt.Sprintf("Rotation completed for block number %d", blockNumber), "ok")
			return nil
		},
	}

	// Define flags for the command
	cmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "Repeat the operation this many times")
	cmd.Flags().BoolVarP(&add, "add", "a", false, "Add each rotated block to the project file")
	cmd.Flags().Uint8VarP(&startRange, "start", "s", 0, "Block start range for the operation")
	cmd.Flags().Uint8VarP(&endRange, "end", "e", 0, "Block end range for the operation")
	return cmd
}

// Cmd_RotateBlockCW90 creates a command to rotate a block 90 degrees clockwise.
func Cmd_RotateBlockCW90() *cobra.Command {
	var repeat int
	var add bool
	var startRange uint8
	var endRange uint8

	var cmd = &cobra.Command{
		Use:   "cw [project file] [block number] [[output file]]",
		Short: "Rotate a block 90 degrees clockwise.",
		Long:  `Rotate a block 90 degrees clockwise.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			inFile := args[0]
			blockNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid block number: %w", err)
			}
			outFile := inFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outFile = args[2]
			}

			mpagd.LogMessage("Cmd_RotateBlockCW90", fmt.Sprintf("Starting rotation for file: %s", inFile), "info")

			apjFile := mpagd.NewAPJFile(inFile)
			err = apjFile.ReadAPJ()
			if err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// If the output file is the same as the input file, create a backup
			if outFile == inFile {
				err := apjFile.BackupProjectFile(false)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_RotateBlockCW90", fmt.Sprintf("Backup created: %s", inFile), "ok")
			}

			// Rotate blocks
			currentBlock := blockNumber
			if startRange == 0 && endRange == 0 {
				for i := 0; i < repeat; i++ {
					processedBlock, err := apjFile.RotateBlock(uint8(currentBlock), false, add)
					if err != nil {
						return fmt.Errorf("failed to rotate block: %w", err)
					}
					if add {
						currentBlock = int(processedBlock)
						mpagd.LogMessage("Cmd_RotateBlockCW90", fmt.Sprintf("created new rotated block: %v", currentBlock), "ok")
					}
				}
			} else {
				for r := startRange; r < endRange; r++ {
					currentBlock = int(r)
					for i := 0; i < repeat; i++ {
						processedBlock, err := apjFile.RotateBlock(uint8(currentBlock), false, add)
						if err != nil {
							return fmt.Errorf("failed to rotate block: %w", err)
						}
						if add {
							currentBlock = int(processedBlock)
							mpagd.LogMessage("Cmd_RotateBlockCW90", fmt.Sprintf("created new rotated block: %v", currentBlock), "ok")
						}
					}
				}
			}

			// Write the updated APJ file
			err = apjFile.WriteAPJ(outFile)
			if err != nil {
				return fmt.Errorf("failed to write APJ file: %w", err)
			}

			mpagd.LogMessage("Cmd_RotateBlockCW90", fmt.Sprintf("Rotation completed for block number %d", blockNumber), "ok")
			return nil
		},
	}

	// Define flags for the command
	cmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "Repeat the operation this many times")
	cmd.Flags().BoolVarP(&add, "add", "a", false, "Add each rotated block to the project file")
	cmd.Flags().Uint8VarP(&startRange, "start", "s", 0, "Block start range for the operation")
	cmd.Flags().Uint8VarP(&endRange, "end", "e", 0, "Block end range for the operation")
	return cmd
}

// Cmd_ImportBlocks creates a command to import DEFINEBLOCK elements from an AGD file into an APJ file.
func Cmd_ImportBlocks() *cobra.Command {
	var replace bool

	var cmd = &cobra.Command{
		Use:   "import [project file] [agd file] [[output file]]",
		Short: "Import DEFINEBLOCK elements from an AGD file into an APJ file.",
		Long:  `Imports DEFINEBLOCK elements from an AGD file into the Blocks section of an APJ file and updates the Nr_of_Blocks field.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			apjFilePath := args[0]
			agdFilePath := args[1]
			outputFilePath := apjFilePath // Default to the input file if no output file is provided
			if len(args) == 3 {
				outputFilePath = args[2]
			}

			mpagd.LogMessage("Cmd_ImportBlocks", fmt.Sprintf("Starting import for file: %s", apjFilePath), "info")

			// Open and read the APJ file
			apj := mpagd.NewAPJFile(apjFilePath)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			// If the output file is the same as the input file, create a backup
			if outputFilePath == apjFilePath {
				err := apj.BackupProjectFile(false)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_ImportBlocks", fmt.Sprintf("Backup created: %s", outputFilePath), "ok")
			}

			// Import blocks from the AGD file
			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptionsFalse()
			options.SetIgnoreOptions(true, true, false, true, true, true, true, true, false)
			options.SetOwOptions(false, false, replace, false, false, false, false, false, true)
			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import blocks from AGD file: %w", err)
			}

			// Write the updated APJ file
			if err := apj.WriteAPJ(outputFilePath); err != nil {
				return fmt.Errorf("failed to write updated APJ file: %w", err)
			}

			mpagd.LogMessage("Cmd_ImportBlocks", fmt.Sprintf("Blocks imported successfully. Updated APJ file saved to %s", outputFilePath), "ok")
			return nil
		},
	}

	// Define flags for the command
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current blocks")
	return cmd
}

// Cmd_RenderBlock creates a command to render blocks to the terminal.
func Cmd_RenderBlock() *cobra.Command {
	var reorderStr string

	var cmd = &cobra.Command{
		Use:   "render [project file] [[start block]] [[end block]]",
		Short: "Render blocks to the terminal.",
		Long:  `Render specific blocks from the project file to the terminal for visualization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			startBlock := 0
			endBlock := 0

			apjFile := mpagd.NewAPJFile(projectFile)
			if err := apjFile.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read project file: %w", err)
			}
			if len(args) == 1 {
				startBlock = 0
				endBlock = len(apjFile.Blocks) - 1
			} else if len(args) == 2 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid blocks number: %w", err)
				}
				startBlock = blockNumber
				endBlock = blockNumber + 1
			} else if len(args) == 3 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid blocks number: %w", err)
				}
				blockNumber2, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("invalid blocks number: %w", err)
				}
				startBlock = blockNumber
				endBlock = blockNumber2
			} else {
				return fmt.Errorf("invalid number of arguments")
			}

			if startBlock < 0 || startBlock >= len(apjFile.Blocks) {
				return fmt.Errorf("start block out of range: %d", startBlock)
			}
			if endBlock < 0 || endBlock > len(apjFile.Blocks) {
				return fmt.Errorf("end block out of range: %d", endBlock)
			}
			mpagd.LogMessage("Cmd_RenderBlock", fmt.Sprintf("Rendering blocks %d to %d", startBlock, endBlock-1), "info")
			var reorder []int
			if reorderStr != "" {
				reorder = mpagd.CSVToIntSlice(reorderStr)
				mpagd.LogMessage("Cmd_RenderBlock", fmt.Sprintf("Reordering blocks: %v", reorder), "ok")
			}
			if err := apjFile.RenderBlockToTerminal(startBlock, endBlock, reorder); err != nil {
				return fmt.Errorf("failed to render blocks: %w", err)
			}
			mpagd.LogMessage("Cmd_RenderBlock", fmt.Sprintf("Blocks %d to %d successfully rendered", startBlock, endBlock-1), "ok")
			return nil
		},
	}

	// Define flags for the command
	cmd.Flags().StringVarP(&reorderStr, "reorder", "r", "", "Reorder the blocks in the output file")
	return cmd
}

// Cmd_RenderBlocksToBitmap creates a command to render blocks to a bitmap file.
func Cmd_RenderBlocksToBitmap() *cobra.Command {
	var reorderStr string
	var offset int
	var cmd = &cobra.Command{
		Use:   "render-bmp [project file] [[start block]] [[end block]] [output file]",
		Short: "Render blocks to a bitmap file.",
		Long:  `Render specific blocks from the project file to a bitmap file for visualization.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			outputFile := args[len(args)-1]
			startIndex := 0
			endIndex := 0

			apjFile := mpagd.NewAPJFile(projectFile)
			if err := apjFile.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read project file: %w", err)
			}

			if len(args) == 2 {
				startIndex = 0
				endIndex = len(apjFile.Blocks) - 1
			} else if len(args) == 3 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid blocks number: %w", err)
				}
				startIndex = blockNumber
				endIndex = blockNumber + 1
			} else if len(args) == 4 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid blocks number: %w", err)
				}
				blockNumber2, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("invalid blocks number: %w", err)
				}
				startIndex = blockNumber
				endIndex = blockNumber2
			} else {
				return fmt.Errorf("invalid number of arguments")
			}

			if startIndex < 0 || startIndex >= len(apjFile.Blocks) {
				return fmt.Errorf("start blocks out of range: %d", startIndex)
			}
			if endIndex < 0 || endIndex > len(apjFile.Blocks) {
				return fmt.Errorf("end blocks out of range: %d", endIndex)
			}
			mpagd.LogMessage("Cmd_RenderBlocksToBitmap", fmt.Sprintf("Rendering blocks %d to %d to bitmap file: %s", startIndex, endIndex-1, outputFile), "info")
			//for i := startIndex; i < endIndex; i++ {
			outputFileWithIndex := outputFile
			//split reorder string by comma
			var reorder []int
			if reorderStr != "" {
				reorder = mpagd.CSVToIntSlice(reorderStr)
				mpagd.LogMessage("Cmd_RenderBlocksToBitmap", fmt.Sprintf("Reordering blocks: %v", reorder), "ok")
			}
			if err := apjFile.RenderBlockToBitmap(uint8(startIndex), uint8(endIndex), outputFileWithIndex, reorder, offset); err != nil {
				return fmt.Errorf("failed to render blocks to bitmap: %w", err)
			}
			//}

			mpagd.LogMessage("Cmd_RenderBlocksToBitmap", fmt.Sprintf("blocks %d to %d successfully rendered to %s", startIndex, endIndex-1, outputFile), "ok")
			return nil
		},
	}

	// Define flags for the command
	cmd.Flags().StringVarP(&reorderStr, "reorder", "r", "", "Reorder the blocks in the output file")
	cmd.Flags().IntVarP(&offset, "offset", "o", 0, "Offset for the start of the reordering blocks")
	return cmd
}

// Cmd_ReorderBlocks creates a command to reorder blocks in the MPAGD project.
func Cmd_ReorderBlocks() *cobra.Command {
	var offset int
	var cmd = &cobra.Command{
		Use:   "reorder [project file] [order] [[output file]]",
		Short: "Reorder blocks in the MPAGD project.",
		Long:  `Reorder blocks in the MPAGD project based on the specified order.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			inFile := args[0]
			orderStr := args[1]
			outFile := inFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outFile = args[2]
			}

			mpagd.LogMessage("Cmd_ReorderBlocks", fmt.Sprintf("Starting reorder for file: %s", inFile), "info")

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
				mpagd.LogMessage("Cmd_ReorderBlocks", fmt.Sprintf("Backup created: %s", inFile), "ok")
			}

			// Convert the order string to a slice of integers
			order := mpagd.CSVToIntSlice(orderStr)

			// Reorder blocks
			err = apjFile.ReorderBlocks(order, offset)
			if err != nil {
				return fmt.Errorf("failed to reorder blocks: %w", err)
			}

			// Write the updated APJ file
			err = apjFile.WriteAPJ(outFile)
			if err != nil {
				return fmt.Errorf("failed to write APJ file: %w", err)
			}

			mpagd.LogMessage("Cmd_ReorderBlocks", "Reorder completed successfully", "ok")
			return nil
		},
	}
	cmd.Flags().IntVarP(&offset, "offset", "o", 0, "Offset for the start of the reordering blocks")
	return cmd
}

// Initialize the commands and add them to the root command.
func init() {
	RootCmd.AddCommand(blocksCmd)
	blocksCmd.AddCommand(Cmd_ImportBlocks())
	blocksCmd.AddCommand(blocksRotateCmd)
	blocksRotateCmd.AddCommand(Cmd_RotateBlockCCW90())
	blocksRotateCmd.AddCommand(Cmd_RotateBlockCW90())
	blocksCmd.AddCommand(Cmd_RenderBlock())
	blocksCmd.AddCommand(Cmd_RenderBlocksToBitmap())
	blocksCmd.AddCommand(Cmd_ReorderBlocks())
}

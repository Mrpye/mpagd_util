package cmd

import (
	"fmt"
	"strconv"

	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

var spriteCmd = &cobra.Command{
	Use:   "sprites",
	Short: "Manage sprite in the MPAGD project.",
	Long:  `Manage sprite in the MPAGD project.`,
}

var spriteRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate sprite in the MPAGD project.",
	Long:  `Rotate sprite in the MPAGD project.`,
}

func Cmd_ImportSprites() *cobra.Command {
	var replace bool
	var cmd = &cobra.Command{
		Use:   "import [apj file] [agd file] [[output file]]",
		Short: "Import Sprites elements from an AGD file into an APJ file.",
		Long:  `Imports Sprites elements from an AGD file into the corresponding section of an APJ file.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			apjFilePath := args[0]
			agdFilePath := args[1]
			outputFilePath := apjFilePath
			if len(args) == 3 {
				outputFilePath = args[2]
			}

			mpagd.LogMessage("Cmd_ImportSprites", fmt.Sprintf("Starting import for file: %s", apjFilePath), "ok")

			apj := mpagd.NewAPJFile(apjFilePath)
			if err := apj.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read APJ file: %w", err)
			}

			if outputFilePath == apjFilePath {
				err := apj.BackupProjectFile(false)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				mpagd.LogMessage("Cmd_ImportSprites", fmt.Sprintf("Backup created: %s", outputFilePath), "ok")
			}

			options := mpagd.CreateImportOptions()
			options.SetIgnoreOptions(true, true, true, true, false, true, true, true, true)
			options.SetOwOptions(false, false, false, false, replace, false, false, false, false)

			if err := apj.ImportAGD(agdFilePath, options); err != nil {
				return fmt.Errorf("failed to import Sprites from AGD file: %w", err)
			}

			mpagd.LogMessage("Cmd_ImportSprites", fmt.Sprintf("Sprites imported successfully. Updated APJ file saved to %s", outputFilePath), "ok")
			return nil
		},
	}
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "Replace the current Sprites")
	return cmd
}

func Cmd_RotateSpritesCCW90() *cobra.Command {
	var repeat int
	var startRange uint8
	var endRange uint8
	var add bool
	var cmd = &cobra.Command{
		Use:   "ccw [project file] [sprite number] [[output file]]",
		Short: "Rotate a sprite 90 degrees counterclockwise.",
		Long:  `Rotate a sprite 90 degrees counterclockwise.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			inFile := args[0]
			spriteNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid sprite number: %w", err)
			}
			outFile := inFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outFile = args[2]
			}

			mpagd.LogMessage("Cmd_RotateSpritesCCW90", fmt.Sprintf("Starting rotation for file: %s", inFile), "ok")

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
				mpagd.LogMessage("Cmd_RotateSpritesCW90", fmt.Sprintf("Backup created: %s", inFile), "ok")
			}

			// Rotate sprites
			currentSprite := spriteNumber
			if startRange == 0 && endRange == 0 {
				for i := 0; i < repeat; i++ {
					processedSprite, err := apjFile.RotateSprite(uint8(currentSprite), true, add)
					if err != nil {
						return fmt.Errorf("failed to rotate sprite: %w", err)
					}
					if add {
						currentSprite = int(processedSprite)
					}
				}
			} else {
				for r := startRange; r < endRange; r++ {
					currentSprite = int(r)
					for i := 0; i < repeat; i++ {
						processedSprite, err := apjFile.RotateSprite(uint8(currentSprite), true, add)
						if err != nil {
							return fmt.Errorf("failed to rotate sprite: %w", err)
						}
						if add {
							currentSprite = int(processedSprite)
						}
					}
				}
			}

			// Write the updated APJ file
			err = apjFile.WriteAPJ(outFile)
			if err != nil {
				return fmt.Errorf("failed to write APJ file: %w", err)
			}

			mpagd.LogMessage("Cmd_RotateSpritesCCW90", fmt.Sprintf("Rotation completed for sprite number %d", spriteNumber), "ok")
			return nil
		},
	}
	cmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "Repeat the operation this many times")
	cmd.Flags().BoolVarP(&add, "add", "a", false, "Add each rotated sprite to the project file")
	cmd.Flags().Uint8VarP(&startRange, "start", "s", 0, "Sprite Start range for the operation")
	cmd.Flags().Uint8VarP(&endRange, "end", "e", 0, "Sprite End range for the operation")
	return cmd
}

func Cmd_RotateSpritesCW90() *cobra.Command {
	var repeat int
	var startRange uint8
	var endRange uint8
	var add bool
	var cmd = &cobra.Command{
		Use:   "cw [project file] [sprite number] [[output file]]",
		Short: "Rotate a sprite 90 degrees clockwise.",
		Long:  `Rotate a sprite 90 degrees clockwise.`,
		Args:  cobra.MinimumNArgs(2), // Ensure at least 2 arguments are provided
		RunE: func(cmd *cobra.Command, args []string) error {
			inFile := args[0]
			spriteNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid sprite number: %w", err)
			}
			outFile := inFile // Default to the input file if no output file is provided
			if len(args) == 3 {
				outFile = args[2]
			}

			mpagd.LogMessage("Cmd_RotateSpritesCW90", fmt.Sprintf("Starting rotation for file: %s", inFile), "ok")

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
				mpagd.LogMessage("Cmd_RotateSpritesCW90", fmt.Sprintf("Backup created: %s", inFile), "ok")
			}

			// Rotate sprites
			currentSprite := spriteNumber
			if startRange == 0 && endRange == 0 {
				for i := 0; i < repeat; i++ {
					processedSprite, err := apjFile.RotateSprite(uint8(currentSprite), false, add)
					if add {
						currentSprite = int(processedSprite)
					}
					if err != nil {
						return fmt.Errorf("failed to rotate sprite: %w", err)
					}
				}
			} else {
				for r := startRange; r < endRange; r++ {
					currentSprite = int(r)
					for i := 0; i < repeat; i++ {
						processedSprite, err := apjFile.RotateSprite(uint8(currentSprite), false, add)
						if add {
							currentSprite = int(processedSprite)
						}
						if err != nil {
							return fmt.Errorf("failed to rotate sprite: %w", err)
						}
					}
				}
			}

			// Write the updated APJ file
			err = apjFile.WriteAPJ(outFile)
			if err != nil {
				return fmt.Errorf("failed to write APJ file: %w", err)
			}

			mpagd.LogMessage("Cmd_RotateSpritesCW90", fmt.Sprintf("Rotation completed for sprite number %d", spriteNumber), "ok")
			return nil
		},
	}
	cmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "Repeat the operation this many times")
	cmd.Flags().BoolVarP(&add, "add", "a", false, "Add each rotated sprite to the project file")
	cmd.Flags().Uint8VarP(&startRange, "start", "s", 0, "Sprite Start range for the operation")
	cmd.Flags().Uint8VarP(&endRange, "end", "e", 0, "Sprite End range for the operation")
	return cmd
}

func Cmd_RenderSprite() *cobra.Command {
	var frame uint8
	var reorderStr string
	var offset int
	var cmd = &cobra.Command{
		Use:   "render [project file]  [[start block]] [[end block]]",
		Short: "Render a sprite to the terminal.",
		Long:  `Render a specific sprite from the project file to the terminal for visualization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			startIndex := 0
			endIndex := 0

			apjFile := mpagd.NewAPJFile(projectFile)
			if err := apjFile.ReadAPJ(); err != nil {
				return fmt.Errorf("failed to read project file: %w", err)
			}
			if len(args) == 1 {
				startIndex = 0
				endIndex = len(apjFile.Sprites) - 1
			} else if len(args) == 2 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid sprite number: %w", err)
				}
				startIndex = blockNumber
				endIndex = blockNumber + 1
			} else if len(args) == 3 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid sprite number: %w", err)
				}
				blockNumber2, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("invalid sprite number: %w", err)
				}
				startIndex = blockNumber
				endIndex = blockNumber2
			} else {
				return fmt.Errorf("invalid number of arguments")
			}

			if startIndex < 0 || startIndex >= len(apjFile.Blocks) {
				return fmt.Errorf("start Sprite out of range: %d", startIndex)
			}
			if endIndex < 0 || endIndex > len(apjFile.Blocks) {
				return fmt.Errorf("end Sprite out of range: %d", endIndex)
			}
			var reorder []int
			if reorderStr != "" {
				reorder = mpagd.CSVToIntSlice(reorderStr)
			}
			if err := apjFile.RenderSpriteToTerminal(uint8(startIndex), uint8(endIndex), reorder, offset); err != nil {
				return fmt.Errorf("failed to render sprite: %w", err)
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&reorderStr, "reorder", "r", "", "Reorder the blocks in the output file")
	cmd.Flags().Uint8VarP(&frame, "frame", "f", 0, "Sprite frame to render")
	cmd.Flags().IntVarP(&offset, "offset", "o", 0, "Offset for the start of the reordering blocks")
	return cmd
}

func Cmd_RenderSpriteToBitmap() *cobra.Command {
	var frame uint8
	var reorderStr string
	var offset int
	var cmd = &cobra.Command{
		Use:   "render-bmp [project file] [[start block]] [[end block]] [output file]",
		Short: "Render sprites to a bitmap file.",
		Long:  `Render specific sprites from the project file to a bitmap file for visualization.`,
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
				endIndex = len(apjFile.Sprites) - 1
			} else if len(args) == 3 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid sprite number: %w", err)
				}
				startIndex = blockNumber
				endIndex = blockNumber + 1
			} else if len(args) == 4 {
				blockNumber, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid sprite number: %w", err)
				}
				blockNumber2, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("invalid sprite number: %w", err)
				}
				startIndex = blockNumber
				endIndex = blockNumber2
			} else {
				return fmt.Errorf("invalid number of arguments")
			}

			if startIndex < 0 || startIndex >= len(apjFile.Sprites) {
				return fmt.Errorf("start sprite out of range: %d", startIndex)
			}
			if endIndex < 0 || endIndex > len(apjFile.Sprites) {
				return fmt.Errorf("end sprite out of range: %d", endIndex)
			}

			//for i := startIndex; i < endIndex; i++ {

			var reorder []int
			if reorderStr != "" {
				reorder = mpagd.CSVToIntSlice(reorderStr)
			}
			if err := apjFile.RenderSpriteToBitmap(uint8(startIndex), uint8(endIndex), outputFile, reorder, offset); err != nil {
				return fmt.Errorf("failed to render sprite to bitmap: %w", err)
			}
			//}

			fmt.Printf("Sprites %d to %d successfully rendered to %s\n", startIndex, endIndex-1, outputFile)
			return nil
		},
	}
	cmd.Flags().StringVarP(&reorderStr, "reorder", "r", "", "Reorder the blocks in the output file")
	cmd.Flags().Uint8VarP(&frame, "frame", "f", 0, "Sprite frame to render")
	cmd.Flags().IntVarP(&offset, "offset", "o", 0, "Offset for the start of the reordering blocks")
	return cmd
}

func init() {
	RootCmd.AddCommand(spriteCmd)
	spriteCmd.AddCommand(Cmd_ImportSprites())
	spriteCmd.AddCommand(Cmd_RenderSprite())
	spriteCmd.AddCommand(spriteRotateCmd)
	spriteRotateCmd.AddCommand(Cmd_RotateSpritesCCW90())
	spriteRotateCmd.AddCommand(Cmd_RotateSpritesCW90())
	spriteCmd.AddCommand(Cmd_RenderSpriteToBitmap())
}

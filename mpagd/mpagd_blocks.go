package mpagd

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/msoap/tcg"
)

// Block represents a game block with various platform-specific data.
type Block struct {
	ID         uint8   // Unique identifier for the block
	Spectrum   []uint8 // Spectrum platform data
	Timex      []uint8 // Timex platform data
	CPC        []uint8 // CPC platform data
	Atom       []uint8 // Atom platform data
	MSX        []uint8 // MSX platform data
	AtomColour []uint8 // Atom color data
}

// createBlock initializes a new Block with the given ID and allocates memory for its platform-specific data.
func (apj *APJFile) createBlock(blockID uint8) Block {
	return Block{
		ID:         blockID,
		Spectrum:   make([]uint8, 9),
		Timex:      make([]uint8, 16),
		CPC:        make([]uint8, 24),
		Atom:       make([]uint8, 8),
		MSX:        make([]uint8, 16),
		AtomColour: make([]uint8, 8),
	}
}

// BlockInit initializes the Blocks slice. If overwrite is true, it resets the Blocks and NrOfBlocks.
func (apj *APJFile) BlockInit(overwrite bool) {
	if apj.Blocks == nil || overwrite {
		apj.Blocks = make([]Block, 0)
		apj.NrOfBlocks = 0
	}
}

// BlockDefault sets up default blocks if none exist. It initializes one block with default values.
func (apj *APJFile) BlockDefault() {
	if len(apj.Blocks) == 0 {
		apj.Blocks = make([]Block, 0)
		for J := 0; J < 1; J++ {
			for i := 0; i < 1; i++ {
				apj.Blocks = append(apj.Blocks, apj.createBlock(uint8(i)))
			}
			// Set default values for each block type
			for i := range apj.Blocks {
				apj.Blocks[i].Spectrum = make([]uint8, 9)
				apj.Blocks[i].Spectrum[len(apj.Blocks[i].Spectrum)-1] = uint8(71)

				apj.Blocks[i].Timex = make([]uint8, 16)
				apj.Blocks[i].CPC = make([]uint8, 24)
				apj.Blocks[i].Atom = make([]uint8, 8)
				apj.Blocks[i].MSX = make([]uint8, 16)
				apj.Blocks[i].AtomColour = make([]uint8, 8)
			}
		}
		apj.NrOfBlocks = uint8(len(apj.Blocks))
	}
}

// readBlocks reads block data from the provided reader and populates the Blocks slice.
func (apj *APJFile) readBlocks(f io.Reader) error {
	var nrOfBlocks uint8
	if err := binary.Read(f, binary.LittleEndian, &nrOfBlocks); err != nil {
		return err
	}
	apj.NrOfBlocks = nrOfBlocks

	blocks := make([]Block, nrOfBlocks)
	for i := 0; i < int(nrOfBlocks); i++ {
		var blockType uint8
		if err := binary.Read(f, binary.LittleEndian, &blockType); err != nil {
			return err
		}
		blocks[i] = apj.createBlock(blockType)
	}
	// Read block image data for each platform
	for i := range blocks {
		blocks[i].Spectrum = apj.readChunk(f, 9)
	}
	for i := range blocks {
		blocks[i].Timex = apj.readChunk(f, 16)
	}
	for i := range blocks {
		blocks[i].CPC = apj.readChunk(f, 24)
	}
	for i := range blocks {
		blocks[i].Atom = apj.readChunk(f, 8)
	}
	for i := range blocks {
		blocks[i].MSX = apj.readChunk(f, 16)
	}
	for i := range blocks {
		blocks[i].AtomColour = apj.readChunk(f, 8)
	}
	for i := range blocks {
		blocks[i].ID = uint8(i)
	}
	apj.State.Blocks = true
	apj.Blocks = blocks
	return nil
}

// writeBlocks writes the block data to the provided writer.
func (apj *APJFile) writeBlocks(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, apj.NrOfBlocks); err != nil {
		return err
	}
	for _, block := range apj.Blocks {
		if err := binary.Write(f, binary.LittleEndian, block.ID); err != nil {
			return err
		}
	}
	blockFields := []func(Block) []uint8{
		func(b Block) []uint8 { return b.Spectrum },
		func(b Block) []uint8 { return b.Timex },
		func(b Block) []uint8 { return b.CPC },
		func(b Block) []uint8 { return b.Atom },
		func(b Block) []uint8 { return b.MSX },
		func(b Block) []uint8 { return b.AtomColour },
	}
	for _, field := range blockFields {
		for _, block := range apj.Blocks {
			if _, err := f.Write(field(block)); err != nil {
				return err
			}
		}
	}
	return nil
}

// ImportBlocks imports block data from a slice of strings.
func (apj *APJFile) ImportBlocks(lines []string) error {
	if len(apj.Blocks) == 0 {
		apj.BlockDefault()
	}
	blockPattern := regexp.MustCompile(`^DEFINEBLOCK\s+(\w+)`)
	if match := blockPattern.FindStringSubmatch(lines[0]); match != nil {
		blockID := GetBlockIDByType(match[1])
		currentBlock := apj.createBlock(blockID)
		blockValues := strings.Fields(lines[1])
		imageData := make([]uint8, len(blockValues))
		for i, value := range blockValues {
			imageData[i] = strUint8(value)
		}
		currentBlock.Spectrum = imageData[:9] // Example: Spectrum uses 9 bytes
		apj.Blocks = append(apj.Blocks, currentBlock)
		apj.NrOfBlocks = uint8(len(apj.Blocks))
		apj.State.Blocks = true
	}
	return nil
}

// RotateBlock rotates a block's Spectrum data 90 degrees counter-clockwise or clockwise.
func (apj *APJFile) RotateBlock(blockIndex uint8, ccw bool, retain bool) (uint8, error) {
	if blockIndex >= uint8(len(apj.Blocks)) {
		return blockIndex, fmt.Errorf("block index out of range: %d", blockIndex)
	}
	var block Block
	if retain {
		block = apj.createBlock(uint8(len(apj.Blocks)))
	} else {
		block = apj.Blocks[blockIndex]
	}

	// Rotate the block's Spectrum data 90 degrees counter-clockwise
	block.Spectrum = apj.blockRotate(apj.Blocks[blockIndex].Spectrum, ccw)
	//apj.Blocks[blockIndex].Timex = rotate90CCW(apj.Blocks[blockIndex].Timex, 8, 8)
	//apj.Blocks[blockIndex].CPC = rotate90CCW(apj.Blocks[blockIndex].CPC, 8, 8)
	//apj.Blocks[blockIndex].Atom = rotate90CCW(apj.Blocks[blockIndex].Atom, 8, 8)
	//apj.Blocks[blockIndex].MSX = rotate90CCW(apj.Blocks[blockIndex].MSX, 8, 8)
	//apj.Blocks[blockIndex].AtomColour = rotate90CCW(apj.Blocks[blockIndex].AtomColour, 8, 8)
	if retain {
		apj.Blocks = append(apj.Blocks, block)
		apj.NrOfBlocks++
	}
	return uint8(len(apj.Blocks) - 1), nil
}

// blockRotate rotates the given block data 90 degrees counter-clockwise or clockwise.
func (apj *APJFile) blockRotate(data []uint8, ccw bool) []uint8 {
	rotated := make([]uint8, len(data))
	// need to convert the number into binary 10101010 into a 2d array 8x8
	//data[0] is first row, data[1] is second row, data[2] is third row, data[3] is fourth row, data[4] is fifth row, data[5] is sixth row, data[6] is seventh row, data[7] is eighth row
	var blockData [8][8]uint8
	var rotatedData [8][8]uint8
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			blockData[i][j] = (data[i] >> (7 - j)) & 1
		}
	}

	// then rotate it 90 degrees counter-clockwis
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if ccw {
				rotatedData[7-x][y] = blockData[y][x]
			} else {
				rotatedData[x][y] = blockData[y][7-x]
			}
		}
	}
	// convert the 2d array back into a 1d array
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			rotated[i] |= (rotatedData[i][j] << (7 - j))
		}
	}

	//put the last 8 bits into the last byte of the array
	// rotated[7] = 0
	rotated[len(rotated)-1] = data[len(data)-1]
	return rotated
}

// RenderBlockToTerminal renders blocks to the terminal.
/*func (apj *APJFile) RenderBlockToTerminal(startBlock, endBlock uint8, reorder []int) error {
	if startBlock >= apj.NrOfBlocks || endBlock > apj.NrOfBlocks {
		return fmt.Errorf("block index out of range: %d", startBlock)
	}
	apj.renderBlockToTerminal(int(startBlock), int(endBlock), reorder)
	//apj.Blocks[blockIndex].Timex = rotate90CW(apj.Blocks[blockIndex].Timex, 8, 8)
	//apj.Blocks[blockIndex].CPC = rotate90CW(apj.Blocks[blockIndex].CPC, 8, 8)
	//apj.Blocks[blockIndex].Atom = rotate90CW(apj.Blocks[blockIndex].Atom, 8, 8)
	//apj.Blocks[blockIndex].MSX = rotate90CW(apj.Blocks[blockIndex].MSX, 8, 8)
	//apj.Blocks[blockIndex].AtomColour = rotate90CW(apj.Blocks[blockIndex].AtomColour, 8, 8)
	return nil
}*/

// BlockTo2DArray converts block data to a 2D array.
func (apj *APJFile) BlockTo2DArray(data []uint8) ([][]uint8, error) {
	var screen = make([][]uint8, 8)
	for i := range screen {
		screen[i] = make([]uint8, 8)
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			screen[y][x] = (data[y] >> (7 - x)) & 1
		}
	}
	return screen, nil
}

// renderBlockToTerminal renders blocks to the terminal.
func (apj *APJFile) RenderBlockToTerminal(start, end int, reorder []int) error {
	tg, err := tcg.New(tcg.Mode1x2) // each terminal symbol contains a 2x3 pixels grid, also you can use 1x1, 1x2, and 2x2 modes
	if err != nil {
		return fmt.Errorf("create tcg: %s", err)
	}

	var data []Block
	if len(reorder) > 0 {
		// Reorder the blocks based on the provided order
		data, err = apj.GetReorderedBlocks(reorder)
		if err != nil {
			return fmt.Errorf("failed to reorder blocks: %s", err)
		}
	} else {
		data = apj.Blocks
	}

	// Calculate the number of sprites and layout
	size := 8
	columns := int(16) // Fixed number of columns
	for Index := start; Index < end; Index++ {
		blockData := data[Index].Spectrum
		sp, err := apj.BlockTo2DArray(blockData)
		if err != nil {
			return fmt.Errorf("failed to convert block to array: %s", err)
		}
		xOffset, yOffset := CalcImageOffSet(uint8(Index), uint8(start), columns, size)
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				pixel := sp[y][x]
				rowIndex := (yOffset + y) // BMP is bottom-up
				colIndex := (xOffset + x)
				if pixel == 1 {
					tg.Buf.Set(colIndex, rowIndex, 1) // draw one pixel with color from 10,10
				} else {
					tg.Buf.Set(colIndex, rowIndex, 0) // draw one pixel with color from 10,10
				}
			}
		}

	}

	tg.Show() // synchronize buffer with screen
	//tg.PrintStr(0, 0, "aa") // draw one pixel with color from 10,10
	return nil
}

// RenderBlockToBitmap renders blocks to a bitmap file.
func (apj *APJFile) RenderBlockToBitmap(startIndex, endIndex uint8, filePath string, reorder []int) error {
	// Validate sprite indexes
	if startIndex >= uint8(len(apj.Blocks)) || endIndex > uint8(len(apj.Blocks)) || startIndex >= endIndex {
		return fmt.Errorf("blocks index range out of bounds: %d-%d", startIndex, endIndex)
	}

	// Calculate the number of sprites and layout
	size := 8
	columns := int(8)
	imageWidth, imageHeight := CalcImageSize(startIndex, endIndex, columns, size)
	var err error
	img := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))

	// Reorder the blocks if needed
	var data []Block
	if len(reorder) > 0 {
		// Reorder the blocks based on the provided order
		data, err = apj.GetReorderedBlocks(reorder)
		if err != nil {
			return fmt.Errorf("failed to reorder blocks: %s", err)
		}
	} else {
		data = apj.Blocks
	}

	for Index := startIndex; Index < endIndex; Index++ {
		sp, err := apj.BlockTo2DArray(data[Index].Spectrum)
		if err != nil {
			return fmt.Errorf("failed to convert sprite to array: %s", err)
		}

		// Calculate sprite position in the image
		xOffset, yOffset := CalcImageOffSet(Index, startIndex, columns, size)
		// sp is a 2D array of uint8 representing the block 1 is black, 0 is white
		attr := data[Index].Spectrum[len(data[Index].Spectrum)-1]
		// get the color from the attribute byte
		fg, bg := SpectrumAttrToColors(attr)
		// Write sprite pixels into the pixel data
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				rowIndex := (yOffset + y) // BMP is bottom-up
				colIndex := (xOffset + x)
				if sp[y][x] == 1 {
					// Set the pixel color to white (or any other color you want)
					img.Set(colIndex, rowIndex, fg)
				} else {
					// Set the pixel color to black (or any other color you want)
					img.Set(colIndex, rowIndex, bg)
				}
			}
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}

// ReorderBlocks reorders the blocks based on the provided order.
func (apj *APJFile) ReorderBlocks(order []int) error {
	newBlocks, err := apj.GetReorderedBlocks(order)
	if err != nil {
		return fmt.Errorf("failed to reorder blocks: %s", err)
	}

	// Make a copy of the screen data
	updatedScreens := make([]Screen, len(apj.Screens))
	for i, screen := range apj.Screens {
		updatedScreens[i] = Screen{
			ScreenData: make([][]uint8, len(screen.ScreenData)),
		}
		for j, row := range screen.ScreenData {
			updatedScreens[i].ScreenData[j] = make([]uint8, len(row))
			copy(updatedScreens[i].ScreenData[j], row)
		}
	}

	// Update the copied screen data
	for i, block := range newBlocks {
		for j := range updatedScreens {
			for k, row := range updatedScreens[j].ScreenData {
				for l := range row {
					if apj.Screens[j].ScreenData[k][l] == block.ID {
						updatedScreens[j].ScreenData[k][l] = uint8(i)
					}
				}
			}
		}
	}

	// Update block IDs
	for i := range newBlocks {
		newBlocks[i].ID = uint8(i)
	}

	// Overwrite the original blocks and screen data
	apj.Blocks = newBlocks
	apj.NrOfBlocks = uint8(len(apj.Blocks))
	apj.Screens = updatedScreens

	return nil
}

// GetReorderedBlocks returns a new slice of blocks reordered based on the provided order.
func (apj *APJFile) GetReorderedBlocks(order []int) ([]Block, error) {
	if len(order) > len(apj.Blocks) {
		return nil, fmt.Errorf("order list exceeds the number of blocks")
	}
	//update the id of the blocks in the new order
	for i := range apj.Blocks {
		apj.Blocks[i].ID = uint8(i)
	}
	// Create a map to track which blocks are already reordered
	reordered := make(map[int]bool)
	newBlocks := make([]Block, 0, len(apj.Blocks))

	// Add blocks in the specified order
	for _, index := range order {
		if index < 0 || index >= len(apj.Blocks) {
			return nil, fmt.Errorf("block index out of range: %d", index)
		}
		newBlocks = append(newBlocks, apj.Blocks[index])
		reordered[index] = true
	}

	// Append remaining blocks that were not specified in the order
	for i, block := range apj.Blocks {
		if !reordered[i] {
			newBlocks = append(newBlocks, block)
		}
	}

	//apj.Blocks = newBlocks
	//apj.NrOfBlocks = uint8(len(apj.Blocks))
	return newBlocks, nil
}

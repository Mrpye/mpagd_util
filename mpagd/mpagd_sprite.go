package mpagd

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/msoap/tcg"
)

// SpriteFrame represents a single frame of sprite data
type SpriteFrame struct {
	Frame     int
	ImageData []uint8
}

// Sprite represents the structure of a sprite
type Sprite struct {
	SpriteID   uint8
	OffSet     uint8
	Frames     uint8
	Spectrum   []SpriteFrame
	Timex      []SpriteFrame
	CPC        []SpriteFrame
	Atom       []SpriteFrame
	AtomColour []SpriteFrame
	VZColour   []SpriteFrame
}

// CreateSprite creates a new sprite with the given ID, offset, and number of frames
// It initializes the sprite's frame data for different platforms.
func (apj *APJFile) createSprite(id, offset, frames uint8) Sprite {
	return Sprite{
		SpriteID:   id,
		OffSet:     offset,
		Frames:     frames,
		Spectrum:   make([]SpriteFrame, 0),
		Timex:      make([]SpriteFrame, 0),
		CPC:        make([]SpriteFrame, 0),
		Atom:       make([]SpriteFrame, 0),
		AtomColour: make([]SpriteFrame, 0),
		VZColour:   make([]SpriteFrame, 0),
	}
}

// readSpriteData reads sprite data for a given number of frames and size
func (apj *APJFile) readSpriteData(f io.Reader, frames int, size int) []SpriteFrame {
	items := make([]SpriteFrame, frames)
	for fr := 0; fr < frames; fr++ {
		imageData := make([]uint8, size)
		if err := binary.Read(f, binary.LittleEndian, &imageData); err != nil {
			return nil
		}
		items[fr] = SpriteFrame{
			Frame:     fr,
			ImageData: imageData,
		}
	}
	return items
}

// SpriteInit initializes the sprite data structure
func (apj *APJFile) SpriteInit(overwrite bool) {
	if apj.Sprites == nil || overwrite {
		apj.Sprites = make([]Sprite, 0)
		apj.NrOfSprites = 0
	}
}

// SetFrameDefaults initializes the frame data for a sprite
func (apj *APJFile) SetFrameDefaults(sprite *Sprite) {
	sprite.Spectrum = make([]SpriteFrame, 0)
	sprite.Timex = make([]SpriteFrame, 0)
	sprite.CPC = make([]SpriteFrame, 0)
	sprite.Atom = make([]SpriteFrame, 0)
	sprite.AtomColour = make([]SpriteFrame, 0)
	sprite.VZColour = make([]SpriteFrame, 0)

	for i := 0; i < int(sprite.Frames); i++ {
		sprite.Spectrum = append(sprite.Spectrum, SpriteFrame{
			Frame:     i,
			ImageData: make([]uint8, 32),
		})
		sprite.Timex = append(sprite.Timex, SpriteFrame{
			Frame:     i,
			ImageData: make([]uint8, 32),
		})
		sprite.CPC = append(sprite.CPC, SpriteFrame{
			Frame:     i,
			ImageData: make([]uint8, 80),
		})
		sprite.Atom = append(sprite.Atom, SpriteFrame{
			Frame:     i,
			ImageData: make([]uint8, 32),
		})
		sprite.AtomColour = append(sprite.AtomColour, SpriteFrame{
			Frame:     i,
			ImageData: make([]uint8, 32),
		})
		sprite.VZColour = append(sprite.VZColour, SpriteFrame{
			Frame:     i,
			ImageData: make([]uint8, 16),
		})
	}
}

// SpriteDefault initializes the default sprite data structure
func (apj *APJFile) SpriteDefault() {
	if len(apj.Sprites) == 0 {
		apj.Sprites = make([]Sprite, 0)
		apj.NrOfSprites = 0
		apj.Sprites = append(apj.Sprites, Sprite{
			SpriteID:   0,
			OffSet:     0,
			Frames:     1,
			Spectrum:   make([]SpriteFrame, 0),
			Timex:      make([]SpriteFrame, 0),
			CPC:        make([]SpriteFrame, 0),
			Atom:       make([]SpriteFrame, 0),
			AtomColour: make([]SpriteFrame, 0),
			VZColour:   make([]SpriteFrame, 0),
		})
		apj.NrOfSprites = 1
		for i := 0; i < 1; i++ {
			apj.SetFrameDefaults(&apj.Sprites[0])
		}
	}
}

// ImportSprites imports sprite data from a list of strings
func (apj *APJFile) ImportSprites(lines []string) error {

	// Setup vars and regex
	blockPattern := regexp.MustCompile(`^DEFINESPRITE\s+(\w+)`)
	blocks := apj.Sprites

	var currentBlock Sprite

	var currentImageData []uint8
	// Setup some flags
	isData := false
	isSecondData := false
	frameCount := 0
	//loop though the lines
	for _, v := range lines {

		// Clean the line
		line := strings.TrimSpace(v)
		if line == "" {
			continue
		}

		// Extract the header data
		if match := blockPattern.FindStringSubmatch(line); match != nil && !isData {
			frames := strUint8(match[1])
			id := uint8(len(blocks))
			currentBlock = apj.createSprite(id, 0, frames)
			isData = true
			continue
		}

		//extract the line data
		if isData {
			blockValues := strings.Fields(line)
			for _, value := range blockValues {
				currentImageData = append(currentImageData, strUint8(value))
			}
			if isSecondData {
				currentBlock.Spectrum = append(currentBlock.Spectrum, SpriteFrame{
					Frame:     frameCount,
					ImageData: currentImageData,
				})

				currentBlock.Timex = append(currentBlock.Timex, SpriteFrame{
					Frame:     frameCount,
					ImageData: make([]uint8, 32),
				})
				currentBlock.CPC = append(currentBlock.CPC, SpriteFrame{
					Frame:     frameCount,
					ImageData: make([]uint8, 80),
				})
				currentBlock.Atom = append(currentBlock.Atom, SpriteFrame{
					Frame:     frameCount,
					ImageData: make([]uint8, 32),
				})
				currentBlock.AtomColour = append(currentBlock.AtomColour, SpriteFrame{
					Frame:     frameCount,
					ImageData: make([]uint8, 32),
				})
				currentBlock.VZColour = append(currentBlock.VZColour, SpriteFrame{
					Frame:     frameCount,
					ImageData: make([]uint8, 16),
				})
				frameCount++
				currentImageData = make([]uint8, 0)
				apj.State.Sprites = true
			}
			isSecondData = !isSecondData
		}
	}

	apj.Sprites = append(apj.Sprites, currentBlock)
	apj.NrOfSprites = uint8(len(apj.Sprites))

	return nil
}

// CalcImageSize calculates the image size based on the number of sprites and layout
func (apj *APJFile) CalcOffset() {
	offset := uint8(0)
	for i, sprite := range apj.Sprites {
		apj.Sprites[i].OffSet = offset
		offset += sprite.Frames
	}
}

// readSprite reads sprite data from a binary file and populates the APJFile structure
func (apj *APJFile) readSprite(f io.Reader) error {
	var nrOfSprites uint8
	if err := binary.Read(f, binary.LittleEndian, &nrOfSprites); err != nil {
		return err
	}
	apj.NrOfSprites = nrOfSprites

	sprites := make([]Sprite, nrOfSprites)
	for i := 0; i < int(nrOfSprites); i++ {
		var offset, frames uint8
		if err := binary.Read(f, binary.LittleEndian, &offset); err != nil {
			return err
		}
		if err := binary.Read(f, binary.LittleEndian, &frames); err != nil {
			return err
		}
		sprites[i] = apj.createSprite(uint8(i), offset, frames)
	}

	// Read sprite image data
	for i := range sprites {
		sprites[i].Spectrum = apj.readSpriteData(f, int(sprites[i].Frames), 32)
	}
	for i := range sprites {
		sprites[i].Timex = apj.readSpriteData(f, int(sprites[i].Frames), 32)
	}
	for i := range sprites {
		sprites[i].CPC = apj.readSpriteData(f, int(sprites[i].Frames), 80)
	}
	for i := range sprites {
		sprites[i].Atom = apj.readSpriteData(f, int(sprites[i].Frames), 32)
	}
	for i := range sprites {
		sprites[i].AtomColour = apj.readSpriteData(f, int(sprites[i].Frames), 32)
	}
	for i := range sprites {
		sprites[i].VZColour = apj.readSpriteData(f, int(sprites[i].Frames), 16)
	}
	apj.State.Sprites = true
	apj.Sprites = sprites
	return nil
}

// writeSprite writes sprite data to a binary file
func (apj *APJFile) writeSprite(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, apj.NrOfSprites); err != nil {
		return err
	}
	offset := uint8(0)
	for _, sprite := range apj.Sprites {
		if err := binary.Write(f, binary.LittleEndian, offset); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, sprite.Frames); err != nil {
			return err
		}
		offset += sprite.Frames
	}

	for _, sprite := range apj.Sprites {
		for _, frame := range sprite.Spectrum { // Replace with the correct field dynamically if needed
			if _, err := f.Write(frame.ImageData); err != nil {
				return err
			}
		}
	}
	for _, sprite := range apj.Sprites {
		for _, frame := range sprite.Timex { // Replace with the correct field dynamically if needed
			if _, err := f.Write(frame.ImageData); err != nil {
				return err
			}
		}
	}
	for _, sprite := range apj.Sprites {
		for _, frame := range sprite.CPC { // Replace with the correct field dynamically if needed
			if _, err := f.Write(frame.ImageData); err != nil {
				return err
			}
		}
	}
	for _, sprite := range apj.Sprites {
		for _, frame := range sprite.Atom { // Replace with the correct field dynamically if needed
			if _, err := f.Write(frame.ImageData); err != nil {
				return err
			}
		}
	}
	for _, sprite := range apj.Sprites {
		for _, frame := range sprite.AtomColour { // Replace with the correct field dynamically if needed
			if _, err := f.Write(frame.ImageData); err != nil {
				return err
			}
		}
	}
	for _, sprite := range apj.Sprites {
		for _, frame := range sprite.VZColour { // Replace with the correct field dynamically if needed
			if _, err := f.Write(frame.ImageData); err != nil {
				return err
			}
		}
	}
	return nil
}

// RotateSprite rotates a sprite by 90 degrees clockwise or counter-clockwise
func (apj *APJFile) RotateSprite(spriteIndex uint8, ccw bool, retain bool) (uint8, error) {
	// Check if the sprite index is valid
	if spriteIndex >= uint8(len(apj.Sprites)) {
		return spriteIndex, fmt.Errorf("block index out of range: %d", spriteIndex)
	}
	// Create a new sprite if retain is true
	// Otherwise, use the existing sprite
	var sprite Sprite
	if retain {
		sprite = apj.createSprite(uint8(len(apj.Sprites)), apj.Sprites[spriteIndex].OffSet, apj.Sprites[spriteIndex].Frames)
	} else {
		sprite = apj.Sprites[spriteIndex]
	}

	// Rotate the sprite's Spectrum data 90 degrees clockwise
	for i := 0; i < len(apj.Sprites[spriteIndex].Spectrum); i++ {
		if retain && len(sprite.Spectrum) < len(apj.Sprites[spriteIndex].Spectrum) {
			sprite.Spectrum = append(sprite.Spectrum, SpriteFrame{
				Frame:     i,
				ImageData: make([]uint8, 32),
			})
			apj.SetFrameDefaults(&sprite)
		}
		sprite.Spectrum[i].ImageData = spriteRotate90(apj.Sprites[spriteIndex].Spectrum[i].ImageData, ccw)
	}

	if retain {
		apj.Sprites = append(apj.Sprites, sprite)
		apj.NrOfSprites++
	}
	//apj.Blocks[spriteIndex].Timex = rotate90CW(apj.Blocks[spriteIndex].Timex, 8, 8)
	//apj.Blocks[spriteIndex].CPC = rotate90CW(apj.Blocks[spriteIndex].CPC, 8, 8)
	//apj.Blocks[spriteIndex].Atom = rotate90CW(apj.Blocks[spriteIndex].Atom, 8, 8)
	//apj.Blocks[spriteIndex].MSX = rotate90CW(apj.Blocks[spriteIndex].MSX, 8, 8)
	//apj.Blocks[spriteIndex].AtomColour = rotate90CW(apj.Blocks[spriteIndex].AtomColour, 8, 8)
	return uint8(len(apj.Sprites) - 1), nil
}

// RotateSprite rotates a sprite by 90 degrees clockwise or counter-clockwise
func spriteRotate90(data []uint8, ccw bool) []uint8 {
	rotated := make([]uint8, len(data))
	// need to convert the number into binary 10101010 into a 2d array 8x8
	//data[0] is first row, data[1] is second row, data[2] is third row, data[3] is fourth row, data[4] is fifth row, data[5] is sixth row, data[6] is seventh row, data[7] is eighth row
	var blockData [16][16]uint8
	var rotatedData [16][16]uint8
	// convert the number into binary 10101010 into a 2d array 16x16
	dataCursor := 0
	subId := 0
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			if j == 8 {
				dataCursor++
				subId = -8
			} else if j == 0 && i > 0 {
				dataCursor++
				subId = 0
			}

			v := (j + subId)

			blockData[i][j] = (data[dataCursor] >> (7 - v)) & 1
		}
	}
	if ccw {
		// then rotate it 90 degrees counter-clockwis
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				rotatedData[7-x][y] = blockData[y][x]

			}
		}
	} else {
		// then rotate it 90 degrees clockwis
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				rotatedData[x][y] = blockData[y][15-x]

			}
		}
	}
	// convert the 2d array back into a 1d array
	dataCursor = 0
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			if j == 8 {
				dataCursor++
				subId = -8
			} else if j == 0 && i > 0 {
				dataCursor++
				subId = 0
			}

			v := (j + subId)

			rotated[dataCursor] |= (rotatedData[i][j] << (7 - v))
		}
	}

	return rotated
}

// SpriteTo2DArray converts a sprite's image data into a 2D array representation
// It assumes the sprite is 16x16 pixels and each pixel is represented by a bit in the byte array.
func (apj *APJFile) SpriteTo2DArray(data []uint8) ([][]uint8, error) {

	var screen = make([][]uint8, 16)
	for i := range screen {
		screen[i] = make([]uint8, 16)
	}

	// convert the number into binary 10101010 into a 2d array 16x16
	dataCursor := 0
	subId := 0

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			if j == 8 {
				dataCursor++
				subId = -8
			} else if j == 0 && i > 0 {
				dataCursor++
				subId = 0
			}
			v := (j + subId)
			if (data[dataCursor]>>(7-v))&1 == 1 {
				screen[i][j] = 1 // draw one pixel with color from 10,10
			} else {
				screen[i][j] = 0 // draw one pixel with color from 10,10
			}
		}

	}

	return screen, nil

}

// GetReorderedSprites returns a reordered list of sprites based on the provided order and offset
func (apj *APJFile) GetReorderedSprites(order []int, offset int) ([]Sprite, error) {
	if len(order)+offset > len(apj.Sprites) {
		return nil, fmt.Errorf("order list exceeds the number of sprites")
	}

	// Create a map to track which sprites are already reordered
	reordered := make(map[int]bool)
	newSprites := make([]Sprite, 0, len(apj.Sprites))

	for i := range apj.Sprites {
		if offset <= i {
			break // Skip blocks before the offset
		}
		newSprites = append(newSprites, apj.Sprites[i])
		reordered[i] = true
	}

	// Add sprites in the specified order
	for _, index := range order {
		index += offset
		if index < 0 || index >= len(apj.Sprites) {
			return nil, fmt.Errorf("sprite index out of range: %d", index)
		}
		newSprites = append(newSprites, apj.Sprites[index])
		reordered[index] = true
	}

	// Append remaining sprites that were not specified in the order
	for i, sprite := range apj.Sprites {
		if !reordered[i] {
			newSprites = append(newSprites, sprite)
		}
	}

	return newSprites, nil
}
func (apj *APJFile) RenderSpriteToSeperateBitmap(startIndex, endIndex uint8, ImagePath string, reorder []int, offset int) error {
	// Validate sprite indexes
	if startIndex >= uint8(len(apj.Sprites)) || endIndex > uint8(len(apj.Sprites)) || startIndex >= endIndex {
		return fmt.Errorf("sprite index range out of bounds: %d-%d", startIndex, endIndex)
	}

	if _, err := os.Stat(ImagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(ImagePath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory for block images: %s", err)
		}
	}

	// Create separate bitmap files for each sprite
	for i := startIndex; i < endIndex; i++ {
		spriteFilePath := fmt.Sprintf("%s/sprite_%d.png", ImagePath, i)
		if err := apj.RenderSpriteToBitmap(i, i+1, spriteFilePath, reorder, offset); err != nil {
			return fmt.Errorf("failed to render sprite %d to bitmap: %s", i, err)
		}
	}

	return nil
}

// RenderSpriteToBitmap renders a range of sprites to a bitmap file
func (apj *APJFile) RenderSpriteToBitmap(startIndex, endIndex uint8, filePath string, reorder []int, offset int) error {
	// Validate sprite indexes
	if startIndex >= uint8(len(apj.Sprites)) || endIndex > uint8(len(apj.Sprites)) || startIndex >= endIndex {
		return fmt.Errorf("sprite index range out of bounds: %d-%d", startIndex, endIndex)
	}

	// Calculate layout and image dimensions
	size := 16
	columns := int(8)
	imageWidth, imageHeight := CalcImageSize(startIndex, endIndex, columns, size)

	var err error
	img := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))

	var data []Sprite
	if len(reorder) > 0 {
		// Reorder the blocks based on the provided order
		data, err = apj.GetReorderedSprites(reorder, offset)
		if err != nil {
			return fmt.Errorf("failed to reorder blocks: %s", err)
		}
	} else {
		data = apj.Sprites
	}
	// Render sprites into pixel data
	for spriteIndex := startIndex; spriteIndex < endIndex; spriteIndex++ {
		sp, err := apj.SpriteTo2DArray(data[spriteIndex].Spectrum[0].ImageData)
		if err != nil {
			return fmt.Errorf("failed to convert sprite to array: %s", err)
		}
		// Calculate sprite position in the image
		xOffset, yOffset := CalcImageOffSet(spriteIndex, startIndex, columns, size)
		fg := color.RGBA{192, 192, 0, 255} // Yellow
		bg := color.RGBA{0, 0, 0, 255}     // Black
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

// RenderSpriteToTerminal renders a range of sprites to the terminal using tcg
// This function assumes that the terminal supports 2x3 pixel characters
func (apj *APJFile) RenderSpriteToTerminal(startIndex, endIndex uint8, reorder []int, offset int) error {
	if startIndex >= uint8(len(apj.Sprites)) || endIndex > uint8(len(apj.Sprites)) {
		return fmt.Errorf("sprite index out of range: %d", startIndex)
	}
	tg, err := tcg.New(tcg.Mode1x2) // each terminal symbol contains a 2x3 pixels grid, also you can use 1x1, 1x2, and 2x2 modes
	if err != nil {
		return fmt.Errorf("create tcg: %s", err)
	}

	// Calculate the number of sprites and layout
	//spriteCount := int(endSprite - startSprite)
	size := 16
	columns := int(8)
	var data []Sprite
	if len(reorder) > 0 {
		// Reorder the blocks based on the provided order
		data, err = apj.GetReorderedSprites(reorder, offset)
		if err != nil {
			return fmt.Errorf("failed to reorder blocks: %s", err)
		}
	} else {
		data = apj.Sprites
	}
	for spriteIndex := startIndex; spriteIndex < endIndex; spriteIndex++ {
		sp, err := apj.SpriteTo2DArray(data[spriteIndex].Spectrum[0].ImageData)
		if err != nil {
			return fmt.Errorf("failed to convert sprite to array: %s", err)
		}
		// Calculate sprite position in the image
		xOffset, yOffset := CalcImageOffSet(spriteIndex, startIndex, columns, size)

		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
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

// ReorderSprites reorders the sprites based on the provided order and offset
// It updates the sprite IDs in the SpriteInfo array accordingly
func (apj *APJFile) ReorderSprites(order []int, offset int) error {
	newSprites, err := apj.GetReorderedSprites(order, offset)
	if err != nil {
		return fmt.Errorf("failed to reorder sprites: %s", err)
	}

	//Update the sprite IDs in the SpriteInfo array
	reordered := make(map[int]bool, len(apj.SpriteInfo))
	for i, sprite := range newSprites {
		for j, spritePos := range apj.SpriteInfo {
			if spritePos.Image == sprite.SpriteID && !reordered[j] {
				reordered[j] = true
				apj.SpriteInfo[j].Image = uint8(i)
			}
		}
		newSprites[i].SpriteID = uint8(i)
	}

	apj.Sprites = newSprites
	apj.NrOfSprites = uint8(len(apj.Sprites))
	apj.CalcOffset() // Recalculate offsets after reordering
	return nil
}

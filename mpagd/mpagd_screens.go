package mpagd

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"
)

// Define a strong type for APJFile data
type Screen struct {
	ScreenID   uint8
	ScreenData [][]uint8
}

func (apj *APJFile) ScreensInit(force bool) {
	if apj.Screens == nil || force {
		apj.Screens = make([]Screen, 0)
		apj.NrOfScreens = 0
	}
}
func (apj *APJFile) ScreensDefault() {
	if len(apj.Screens) == 0 {
		apj.Screens = make([]Screen, 0)
		// make a blank screen
		screenData := make([][]uint8, apj.Windows.Height)
		for y := 0; y < int(apj.Windows.Height); y++ {
			row := make([]uint8, apj.Windows.Width)
			screenData[y] = row
		}
		apj.NrOfScreens = uint8(1)
		apj.Screens = append(apj.Screens, Screen{
			ScreenID:   0,
			ScreenData: screenData,
		})
	}
}
func (apj *APJFile) ImportScreens(lines []string) error {
	screenData := make([][]uint8, 0)
	var currentData []uint8
	var spritePosLines []string
	// Setup some flags
	isPos := false
	line := ""
	// Loop through the lines
	for _, v := range lines {
		line = v
		// See what type of data we have
		if strings.Contains(line, "SPRITEPOSITION") {
			isPos = true
			spritePosLines = append(spritePosLines, line)
			line = strings.ReplaceAll(line, "SPRITEPOSITION", "")
		} else {
			line = strings.ReplaceAll(line, "DEFINESCREEN", "")
		}

		// Clean the line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !isPos {
			blockValues := strings.Fields(line)
			for _, value := range blockValues {
				currentData = append(currentData, strUint8(value))
			}

			screenData = append(screenData, currentData)
			currentData = make([]uint8, 0)
			apj.State.Screens = true
		}
	}

	id := uint8(len(apj.Screens))
	currentBlock := Screen{
		ScreenID:   id,
		ScreenData: screenData,
	}
	apj.Screens = append(apj.Screens, currentBlock)
	apj.NrOfScreens = uint8(len(apj.Screens))
	// Process the sprite pos
	apj.ImportSpritePos(id, spritePosLines)
	return nil
}

// make a copy of the screen data then remap the blocks
func (apj *APJFile) RemapScreens(screenIndex uint8, BlockOffSet uint8) {
	screen := apj.Screens[screenIndex]
	newScreenData := make([][]uint8, len(screen.ScreenData))
	for y, row := range screen.ScreenData {
		newRow := make([]uint8, len(row))
		for x, _ := range row {
			// remap the block ID in the spectrum data
			if row[x] == 0 {
				continue // skip empty blocks
			}
			newRow[x] = row[x] + BlockOffSet // remap the block ID in the spectrum data
		}
		newScreenData[y] = newRow
	}
	apj.Screens[screenIndex].ScreenData = newScreenData

}

// readScreens reads screen data
func (apj *APJFile) readScreens(f io.Reader) error {
	var nrOfScreens uint8
	if err := binary.Read(f, binary.LittleEndian, &nrOfScreens); err != nil {
		return err
	}
	apj.NrOfScreens = nrOfScreens

	screens := make([]Screen, nrOfScreens)
	for i := 0; i < int(nrOfScreens); i++ {
		screenData := make([][]uint8, apj.Windows.Height)
		for y := 0; y < int(apj.Windows.Height); y++ {
			row := make([]uint8, apj.Windows.Width)
			if err := binary.Read(f, binary.LittleEndian, &row); err != nil {
				return err
			}
			screenData[y] = row
		}
		screen := Screen{
			ScreenID:   uint8(i),
			ScreenData: screenData,
		}
		screens[i] = screen
	}
	apj.State.Screens = true
	apj.Screens = screens
	return nil
}

func (apj *APJFile) writeScreens(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, apj.NrOfScreens); err != nil {
		return err
	}
	for _, screen := range apj.Screens {
		for _, row := range screen.ScreenData {
			if _, err := f.Write(row); err != nil {
				return err
			}
		}
	}
	return nil
}

// RenderScreenToBitmap renders the specified screen to a bitmap and saves it as a PNG file.
func (apj *APJFile) RenderScreenToBitmap(screenIndex uint8, filePath string) error {
	if int(screenIndex) >= len(apj.Screens) {
		return fmt.Errorf("screen index out of range")
	}

	screen := apj.Screens[screenIndex]
	img := image.NewRGBA(image.Rect(0, 0, int(apj.Windows.Width*8), int(apj.Windows.Height*8)))

	for y, row := range screen.ScreenData {
		for x, blockID := range row {
			// Map blockID to a color (example: grayscale based on blockID)
			// get the block
			data := apj.Blocks[blockID].Spectrum
			sp, err := apj.BlockTo2DArray(data)
			if err != nil {
				return fmt.Errorf("failed to convert sprite to array: %s", err)
			}
			// sp is a 2D array of uint8 representing the block 1 is black, 0 is white
			attr := apj.Blocks[blockID].Spectrum[len(apj.Blocks[blockID].Spectrum)-1]
			// get the color from the attribute byte
			fg, bg := SpectrumAttrToColors(attr)

			// draw the block in the image
			for blockY := 0; blockY < int(8); blockY++ {
				for blockX := 0; blockX < int(8); blockX++ {
					//calc the pixle colour from spectrum attributew byte
					// clear the upper bits
					if sp[blockY][blockX] == 1 {
						// Set the pixel color to white (or any other color you want)
						img.Set(x*int(8)+blockX, y*int(8)+blockY, fg)
					} else {
						// Set the pixel color to black (or any other color you want)
						img.Set(x*int(8)+blockX, y*int(8)+blockY, bg)

					}
				}
			}

			// draw the block
			// Here we just use a simple grayscale mapping for demonstration
			// You can replace this with your own color mapping logic
			// For example, you can use a color palette or a more complex mapping

			//grayValue := uint8(blockID * 16) // Example mapping
			//col := color.RGBA{grayValue, grayValue, grayValue, 255}
			//img.Set(x, y, col)
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

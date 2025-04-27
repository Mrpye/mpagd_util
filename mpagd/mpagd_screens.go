package mpagd

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"sort"
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

// ReorderScreens reorders the screens based on the provided order slice.
// The order slice should contain the new indices for the screens.
func (apj *APJFile) ReorderScreens(order []int) error {
	if len(order) != len(apj.Screens) {
		return fmt.Errorf("order length does not match the number of screens")
	}

	newScreens := make([]Screen, len(apj.Screens))
	for newIndex, oldIndex := range order {
		if int(oldIndex) >= len(apj.Screens) {
			return fmt.Errorf("invalid screen index in order: %d", oldIndex)
		}
		newScreens[newIndex] = apj.Screens[oldIndex]
	}

	//update the apj.Map
	for y := 0; y < len(apj.Map.Map); y++ {
		for x := 0; x < len(apj.Map.Map[y]); x++ {
			screenID := apj.Map.Map[y][x]
			if int(screenID) < len(order) {
				apj.Map.Map[y][x] = uint8(order[screenID])
			}
		}
	}
	apj.Map.StartScreen = uint8(order[apj.Map.StartScreen]) // Update the start screen ID in the map

	//update spriteinfo
	for i := 0; i < len(apj.SpriteInfo); i++ {
		screenID := apj.SpriteInfo[i].Screen
		if int(screenID) < len(order) {
			apj.SpriteInfo[i].Screen = uint8(order[screenID])
		}
	}
	//sorting the apj.SpriteInfo by screen ID
	sort.Slice(apj.SpriteInfo, func(i, j int) bool {
		return apj.SpriteInfo[i].Screen < apj.SpriteInfo[j].Screen
	})

	// Update the ScreenID for each screen in the new order
	for i := 0; i < len(newScreens); i++ {
		newScreens[i].ScreenID = uint8(i) // Update the ScreenID to match the new order
	}
	apj.NrOfScreens = uint8(len(newScreens))
	apj.Screens = newScreens
	return nil
}

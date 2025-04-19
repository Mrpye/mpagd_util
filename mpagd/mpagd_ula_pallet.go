package mpagd

import (
	"encoding/binary"
	"io"
	"strings"
)

// ULAPalette represents the ULAPlus palette data.
// The ULAPlus palette is a 16-color palette used in the ULAPlus graphics system.
// The palette is represented as a slice of uint8 values, where each value corresponds to a color index.
type ULAPalette struct {
	Colors []uint8
}

// ULAPaletteInit initializes the ULAPlus palette with default values.
// If the overwrite flag is true, it will reset the palette even if it already contains data.
func (apj *APJFile) ULAPaletteInit(overwrite bool) {
	// Initialize the palette with default colors if empty or overwrite is true.
	if len(apj.ULAPalette.Colors) == 0 || overwrite {
		apj.ULAPalette = ULAPalette{
			Colors: []uint8{0, 66, 24, 146, 195, 152, 252, 109, 0, 44, 156, 15, 195, 131, 190, 253},
		}
	}
}

// ImportULAPalette imports ULAPlus palette data from a slice of strings.
// Each string represents a line of palette data, and the function parses and updates the palette colors.
func (apj *APJFile) ImportULAPalette(lines []string) error {
	palette := &apj.ULAPalette
	colorIndex := 0

	// Process each line to extract color values.
	for _, v := range lines {
		line := strings.ReplaceAll(v, "DEFINEPALETTE", "")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse color values from the line.
		blockValues := strings.Fields(line)
		for _, value := range blockValues {
			if colorIndex >= len(palette.Colors) {
				break
			}
			palette.Colors[colorIndex] = strUint8(value)
			colorIndex++
		}
		apj.State.ULAPalette = true
	}

	return nil
}

// readULAPalette reads ULAPlus palette data from an io.Reader.
// It populates the palette with 16 colors read in little-endian format.
func (apj *APJFile) readULAPalette(f io.Reader) error {
	apj.ULAPalette.Colors = make([]uint8, 16)
	if err := binary.Read(f, binary.LittleEndian, &apj.ULAPalette.Colors); err != nil {
		return err
	}
	apj.State.ULAPalette = true
	return nil
}

// writeULAPalette writes the ULAPlus palette data to an io.Writer.
// The palette is written in little-endian format.
func (apj *APJFile) writeULAPalette(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, apj.ULAPalette.Colors); err != nil {
		return err
	}
	return nil
}

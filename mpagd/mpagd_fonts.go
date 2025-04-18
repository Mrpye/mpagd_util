package mpagd

import (
	"encoding/binary"
	"io"
	"strings"
)

type Font struct {
	ID   int
	Data []uint8
}

// FontInit initializes the Font slice in apj.Data
func (apj *APJFile) FontInit(overwrite bool) {
	if apj.Fonts == nil || len(apj.Fonts) == 0 || overwrite {
		font := make([]Font, 96)
		for i := 0; i < 96; i++ {
			data := make([]uint8, 8)
			for j := 0; j < 8; j++ {
				data[j] = uint8(i + j)
			}
			font[i] = Font{
				ID:   i,
				Data: data,
			}
		}
		apj.Fonts = font
	}
}

// create fonts with default values
func (apj *APJFile) FontDefault() {
	if apj.Fonts == nil {
		apj.Fonts = make([]Font, 96)
		for i := 0; i < 96; i++ {
			data := make([]uint8, 8)
			apj.Fonts[i] = Font{
				ID:   i,
				Data: data,
			}
		}
	}

}

// ImportFont imports font data from lines and populates apj.Data["Font"]
func (apj *APJFile) ImportFont(lines []string) error {
	var font []Font
	var currentData []uint8

	// Loop through the lines
	for i, v := range lines {
		line := strings.ReplaceAll(v, "DEFINEFONT", "")
		// Clean the line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		blockValues := strings.Fields(line)
		for _, value := range blockValues {
			currentData = append(currentData, strUint8(value))
		}
		font = append(font, Font{
			ID:   i,
			Data: currentData,
		})
		currentData = make([]uint8, 0)
		apj.State.Fonts = true
	}

	apj.Fonts = font

	return nil
}

// readFont reads font data from an io.Reader
func (apj *APJFile) readFont(f io.Reader) error {
	font := make([]Font, 96)
	for i := 0; i < 96; i++ {
		data := make([]uint8, 8)
		if err := binary.Read(f, binary.LittleEndian, &data); err != nil {
			return err
		}
		font[i] = Font{
			ID:   i,
			Data: data,
		}
	}
	apj.State.Fonts = true
	apj.Fonts = font
	return nil
}

// writeFont writes font data to an io.Writer
func (apj *APJFile) writeFont(f io.Writer) error {
	for _, font := range apj.Fonts {
		if _, err := f.Write(font.Data); err != nil {
			return err
		}
	}
	return nil
}

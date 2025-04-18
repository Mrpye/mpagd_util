package mpagd

import (
	"encoding/binary"
	"io"
	"strings"
)

// MapData holds map-related data
type Map struct {
	Height      uint8
	Width       uint8
	StartRow    uint8
	StartColumn uint8
	StartScreen uint8
	Map         [][]uint8
}

func (apj *APJFile) MapInit(overwrite bool) {
	if len(apj.Map.Map) == 0 || overwrite {
		apj.Map.Height = 0
		apj.Map.Width = 0
		apj.Map.StartRow = 0
		apj.Map.StartColumn = 0
		apj.Map.Map = make([][]uint8, 0)
		apj.Map.StartScreen = 0
	}
}

func (apj *APJFile) MapDefault() {
	if len(apj.Map.Map) == 0 {
		apj.Map.Height = 10
		apj.Map.Width = 16
		apj.Map.StartRow = 4
		apj.Map.StartColumn = 7
		apj.Map.StartScreen = 0

		//populate the map with 255
		for i := 0; i < int(apj.Map.Height); i++ {
			row := make([]uint8, apj.Map.Width)
			for j := 0; j < int(apj.Map.Width); j++ {
				row[j] = 255
			}
			apj.Map.Map = append(apj.Map.Map, row)
		}
	}
}

func (apj *APJFile) ImportMap(lines []string) error {
	// Initialize variables for map properties
	var mapWidth uint8
	var startScreen uint8
	var startRow uint8
	var startColumn uint8
	var mapLayout [][]uint8
	isMapSection := false

	// Loop through the provided lines
	for _, line := range lines {
		// Trim whitespace from the line
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for the start of the MAP section
		if strings.HasPrefix(line, "MAP") {
			isMapSection = true
			line = strings.ReplaceAll(line, "MAP", "")
			line = strings.TrimSpace(line)
		}

		// Process lines within the MAP section
		if isMapSection {
			// Parse the map width
			if strings.HasPrefix(line, "WIDTH") {
				parts := strings.Fields(line)
				if len(parts) == 2 {
					mapWidth = uint8(strUint8(parts[1]))
				}
				continue
			}

			// Parse the start screen
			if strings.HasPrefix(line, "STARTSCREEN") {
				parts := strings.Fields(line)
				if len(parts) == 2 {
					startScreen = uint8(strUint8(parts[1]))
				}
				continue
			}

			// End the MAP section
			if strings.HasPrefix(line, "ENDMAP") {
				isMapSection = false
				continue
			}

			// Parse map layout rows
			rowValues := strings.Fields(line)
			row := make([]uint8, len(rowValues))
			for i, value := range rowValues {
				row[i] = uint8(strUint8(value))
			}
			mapLayout = append(mapLayout, row)
			continue
		}
	}

	// Find the position of the start screen in the map layout
	for i, row := range mapLayout {
		for j, value := range row {
			if value == startScreen {
				startRow = uint8(i)
				startColumn = uint8(j)
				break
			}
		}
	}

	// Update the APJFile map properties
	apj.Map.Width = mapWidth
	apj.Map.Height = uint8(len(mapLayout))
	apj.Map.StartRow = startRow
	apj.Map.StartColumn = startColumn
	apj.Map.StartScreen = startScreen
	apj.Map.Map = mapLayout

	// Mark the map state as initialized
	apj.State.Map = true

	// Return nil to indicate success
	return nil
}

// readMap reads map data
func (apj *APJFile) readMap(f io.Reader) error {
	if err := binary.Read(f, binary.LittleEndian, &apj.Map.Height); err != nil {
		return err
	}
	if err := binary.Read(f, binary.LittleEndian, &apj.Map.Width); err != nil {
		return err
	}
	if err := binary.Read(f, binary.LittleEndian, &apj.Map.StartRow); err != nil {
		return err
	}
	if err := binary.Read(f, binary.LittleEndian, &apj.Map.StartColumn); err != nil {
		return err
	}

	mapData := make([][]uint8, apj.Map.Height)
	for y := 0; y < int(apj.Map.Height); y++ {
		row := make([]uint8, apj.Map.Width)
		if err := binary.Read(f, binary.LittleEndian, &row); err != nil {
			return err
		}
		mapData[y] = row
	}
	apj.State.Map = true
	apj.Map.Map = mapData
	return nil
}

func (apj *APJFile) writeMap(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, apj.Map.Height); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, apj.Map.Width); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, apj.Map.StartRow); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, apj.Map.StartColumn); err != nil {
		return err
	}
	for _, row := range apj.Map.Map {
		if _, err := f.Write(row); err != nil {
			return err
		}
	}
	return nil
}

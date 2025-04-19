package mpagd

import (
	"encoding/binary"
	"io"
	"regexp"
)

// Windows represents the size and position of a window.
type Windows struct {
	Top    uint8
	Left   uint8
	Height uint8
	Width  uint8
}

// CreateWindow initializes a Windows struct with the given parameters.
func CreateWindow(top, left, height, width uint8) Windows {
	return Windows{Top: top, Left: left, Height: height, Width: width}
}

// WindowsInit initializes the Windows struct in the APJFile.
// If overwrite is true, it resets the window to default values.
func (apj *APJFile) WindowsInit(overwrite bool) {
	// Check if the window is uninitialized or overwrite is requested.
	if apj.Windows.Top == 0 && apj.Windows.Left == 0 && apj.Windows.Height == 0 && apj.Windows.Width == 0 || overwrite {
		apj.Windows = CreateWindow(1, 2, 22, 22) // Default window values.
	}
}

// readWindows reads the window size and position from the given reader.
func (apj *APJFile) readWindows(f io.Reader) error {
	var winTop, winLeft, winHeight, winWidth uint8

	// Read window properties in LittleEndian format.
	if err := binary.Read(f, binary.LittleEndian, &winTop); err != nil {
		return err
	}
	if err := binary.Read(f, binary.LittleEndian, &winLeft); err != nil {
		return err
	}
	if err := binary.Read(f, binary.LittleEndian, &winHeight); err != nil {
		return err
	}
	if err := binary.Read(f, binary.LittleEndian, &winWidth); err != nil {
		return err
	}

	// Update the APJFile's window properties.
	apj.Windows = CreateWindow(winTop, winLeft, winHeight, winWidth)
	apj.State.Windows = true
	return nil
}

// updateWindows updates the window properties in the APJFile.
func (apj *APJFile) updateWindows(winTop, winLeft, winHeight, winWidth uint8) error {
	apj.Windows = CreateWindow(winTop, winLeft, winHeight, winWidth)
	return nil
}

// importWindows parses a DEFINEWINDOW line and updates the window properties.
// It also resets screen data if it exists.
func (apj *APJFile) importWindows(line string) error {
	// Regular expression to match the DEFINEWINDOW format.
	pattern := regexp.MustCompile(`^DEFINEWINDOW\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)$`)
	if match := pattern.FindStringSubmatch(line); match != nil {
		// Update window properties using parsed values.
		apj.updateWindows(strUint8(match[1]), strUint8(match[2]), strUint8(match[3]), strUint8(match[4]))

		// Clear existing screen data.
		if len(apj.Screens) > 0 {
			apj.Screens = make([]Screen, 0)
		}
		apj.NrOfScreens = 0
		apj.ScreensInit(true)
		apj.State.Windows = true
	}
	return nil
}

// writeWindows writes the window size and position to the given writer.
func (apj *APJFile) writeWindows(f io.Writer) error {
	// Write window properties in LittleEndian format.
	if err := binary.Write(f, binary.LittleEndian, apj.Windows.Top); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, apj.Windows.Left); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, apj.Windows.Height); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, apj.Windows.Width); err != nil {
		return err
	}
	return nil
}

package mpagd

import (
	"encoding/binary"
	"io"
	"regexp"
	"strconv"
)

// readKeys reads the key data from the provided reader and updates the APJFile's Keys field.
func (apj *APJFile) readKeys(f io.Reader) error {
	keys := make([]uint8, 11)
	if err := binary.Read(f, binary.LittleEndian, &keys); err != nil {
		return err
	}
	apj.Keys = keys
	apj.State.Keys = true
	return nil
}

// KeysInit initializes the Keys field with default values if it is empty or if overwrite is true.
func (apj *APJFile) KeysInit(overwrite bool) {
	if len(apj.Keys) == 0 || apj.Keys == nil || overwrite {
		// Default key values
		apj.Keys = []uint8{87, 83, 65, 68, 32, 74, 72, 49, 50, 51, 52}
	}
}

// updateKeys updates the Keys field with the provided slice of uint8 values.
func (apj *APJFile) updateKeys(Keys []uint8) error {
	apj.Keys = Keys
	return nil
}

// importKeys parses a string to extract key values and updates the Keys field.
func (apj *APJFile) importKeys(line string) error {
	// Regular expression to match quoted characters or standalone numbers
	pattern := regexp.MustCompile(`'([^']*)'|\b(\d+)\b`)
	matches := pattern.FindAllStringSubmatch(line, -1)

	var key []uint8
	for _, match := range matches {
		if match[1] != "" { // Content inside quotes
			key = append(key, uint8([]rune(match[1])[0]))
		} else if match[2] != "" { // Numbers
			num, _ := strconv.Atoi(match[2])
			key = append(key, uint8(num))
		}
	}

	apj.State.Keys = true
	return apj.updateKeys(key)
}

// writeKeys writes the Keys field data to the provided writer.
func (apj *APJFile) writeKeys(f io.Writer) error {
	if _, err := f.Write(apj.Keys); err != nil {
		return err
	}
	return nil
}

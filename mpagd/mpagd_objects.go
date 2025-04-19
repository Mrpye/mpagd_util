package mpagd

import (
	"encoding/binary"
	"io"
	"regexp"
	"strings"
)

// Object represents the structure of an object with various platform-specific data.
type Object struct {
	ID         uint8   // Unique identifier for the object
	Spectrum   []uint8 // Data specific to the Spectrum platform
	Timex      []uint8 // Data specific to the Timex platform
	CPC        []uint8 // Data specific to the CPC platform
	Atom       []uint8 // Data specific to the Atom platform
	MSX        []uint8 // Data specific to the MSX platform
	AtomColour []uint8 // Colour data for the Atom platform
	VZColour   []uint8 // Colour data for the VZ platform
}

// createObject initializes a new Object with default padded slices for each platform.
func (apj *APJFile) createObject(id uint8) Object {
	return Object{
		ID:         id,
		Spectrum:   padSlice(36),
		Timex:      padSlice(35),
		CPC:        padSlice(67),
		Atom:       padSlice(35),
		MSX:        padSlice(67),
		AtomColour: padSlice(35),
		VZColour:   padSlice(19),
	}
}

// ObjectDefault initializes the Objects slice with a single default object if it is empty.
func (apj *APJFile) ObjectDefault() {
	if len(apj.Objects) == 0 {
		apj.Objects = make([]Object, 0)
		apj.Objects = append(apj.Objects, apj.createObject(0))
		apj.NrOfObjects = uint8(len(apj.Objects))
	}
}

// ImportObjects parses object definitions and data from a list of strings.
func (apj *APJFile) ImportObjects(lines []string) error {
	pattern := regexp.MustCompile(`DEFINEOBJECT\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	var objects []Object
	var currentData []uint8
	isData := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Match object definition header
		if match := pattern.FindStringSubmatch(line); match != nil && !isData {
			currentData = append(currentData, strUint8(match[1]), strUint8(match[2]), strUint8(match[3]), strUint8(match[4]))
			isData = true
			continue
		}

		// Collect object data
		if isData {
			for _, value := range strings.Fields(line) {
				currentData = append(currentData, strUint8(value))
			}
			apj.State.Objects = true
		}
	}

	// Create and append the object
	id := uint8(len(objects))
	currentObject := apj.createObject(id)
	currentObject.Spectrum = currentData
	objects = append(objects, currentObject)

	apj.Objects = objects
	apj.NrOfObjects = uint8(len(objects))
	return nil
}

// ObjectInit initializes or resets the Objects slice based on the overwrite flag.
func (apj *APJFile) ObjectInit(overwrite bool) {
	if len(apj.Objects) == 0 || overwrite {
		apj.Objects = make([]Object, 0)
		apj.NrOfObjects = 0
	}
}

// DeleteObjects clears all objects from the APJFile.
func (apj *APJFile) DeleteObjects() {
	apj.Objects = make([]Object, 0)
	apj.NrOfObjects = 0
}

// DeleteObject removes a specific object by its ID.
func (apj *APJFile) DeleteObject(objectID int) {
	if objectID >= 0 && objectID < len(apj.Objects) {
		apj.Objects = append(apj.Objects[:objectID], apj.Objects[objectID+1:]...)
		apj.NrOfObjects = uint8(len(apj.Objects))
	}
}

// readObjects reads object data from a binary file.
func (apj *APJFile) readObjects(f io.Reader) error {
	var nrOfObjects uint8
	if err := binary.Read(f, binary.LittleEndian, &nrOfObjects); err != nil {
		return err
	}
	apj.NrOfObjects = nrOfObjects

	objects := make([]Object, nrOfObjects)
	for i := 0; i < int(nrOfObjects); i++ {
		objects[i] = apj.createObject(uint8(i))
	}

	// Read platform-specific data for each object
	for i := range objects {
		objects[i].Spectrum = apj.readChunk(f, 36)
	}
	for i := range objects {
		objects[i].Timex = apj.readChunk(f, 35)
	}
	for i := range objects {
		objects[i].CPC = apj.readChunk(f, 67)
	}
	for i := range objects {
		objects[i].Atom = apj.readChunk(f, 35)
	}
	for i := range objects {
		objects[i].MSX = apj.readChunk(f, 67)
	}
	for i := range objects {
		objects[i].AtomColour = apj.readChunk(f, 35)
	}
	for i := range objects {
		objects[i].VZColour = apj.readChunk(f, 19)
	}
	apj.State.Objects = true
	apj.Objects = objects
	return nil
}

// writeObjects writes object data to a binary file.
func (apj *APJFile) writeObjects(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, apj.NrOfObjects); err != nil {
		return err
	}

	// Write platform-specific data for each object
	for _, obj := range apj.Objects {
		if _, err := f.Write(obj.Spectrum); err != nil {
			return err
		}
	}
	for _, obj := range apj.Objects {
		if _, err := f.Write(obj.Timex); err != nil {
			return err
		}
	}
	// Write object image data
	for _, obj := range apj.Objects {
		if _, err := f.Write(obj.CPC); err != nil {
			return err
		}
	}
	for _, obj := range apj.Objects {
		if _, err := f.Write(obj.Atom); err != nil {
			return err
		}
	}
	for _, obj := range apj.Objects {
		if _, err := f.Write(obj.MSX); err != nil {
			return err
		}
	}
	for _, obj := range apj.Objects {
		if _, err := f.Write(obj.AtomColour); err != nil {
			return err
		}
	}
	for _, obj := range apj.Objects {
		if _, err := f.Write(obj.VZColour); err != nil {
			return err
		}
	}
	return nil
}

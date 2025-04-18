package mpagd

import (
	"encoding/binary"
	"io"
)

// readHeader reads the main header and version of the APJ file from the provided reader.
// f: io.Reader - the input source to read the header and version from.
func (apj *APJFile) readHeader(f io.Reader) error {
	// Read the 4-byte header
	header := make([]byte, 4)
	if _, err := f.Read(header); err != nil {
		return err
	}
	apj.Header = header

	// Read the 4-byte version
	var version uint32
	if err := binary.Read(f, binary.LittleEndian, &version); err != nil {
		return err
	}
	apj.Version = version
	return nil
}

// HeaderInit initializes the APJFile with default header and version values.
func (apj *APJFile) HeaderInit() {
	apj.Header = []uint8{65, 71, 68, 42} // "APJ\0"
	apj.Version = uint32(10)
}

// writeHeader writes the main header and version of the APJ file to the provided writer.
// f: io.Writer - the output destination to write the header and version to.
func (apj *APJFile) writeHeader(f io.Writer) error {
	// Write the 4-byte header
	if _, err := f.Write(apj.Header); err != nil {
		return err
	}

	// Write the 4-byte version
	if err := binary.Write(f, binary.LittleEndian, apj.Version); err != nil {
		return err
	}
	return nil
}

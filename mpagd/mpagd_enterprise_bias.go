package mpagd

import (
	"encoding/binary"
	"io"
)

// APJData represents the data structure for APJ files, including the EnterpriseBias setting.
type APJData struct {
	EnterpriseBias uint8 // EnterpriseBias is a single-byte value representing the enterprise bias setting.
}

// readEnterpriseBiasSetting reads the Enterprise bias setting from the provided reader.
// It updates the EnterpriseBias field of the APJFile instance.
func (apj *APJFile) readEnterpriseBiasSetting(f io.Reader) error {
	var bias uint8
	if err := binary.Read(f, binary.LittleEndian, &bias); err != nil {
		return err
	}
	apj.EnterpriseBias = bias
	return nil
}

// writeEnterpriseBiasSetting writes the Enterprise bias setting to the provided writer.
// It serializes the EnterpriseBias field of the APJFile instance using little-endian encoding.
func (apj *APJFile) writeEnterpriseBiasSetting(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, apj.EnterpriseBias); err != nil {
		return err
	}
	return nil
}

package mpagd

import (
	"io"
	"strings"
)

// ASMPathInit initializes the ASM path as an empty slice.
func (apj *APJFile) ASMPathInit() {
	apj.AsmPath = make([]uint8, 0)
}

// readASMPath reads the ASM path from the provided reader.
// It reads up to 256 bytes or until a null byte is encountered.
// Trims any trailing null bytes and assigns the result to AsmPath.
func (apj *APJFile) readASMPath(f io.Reader) error {
	path := make([]byte, 256)
	if _, err := f.Read(path); err != nil {
		return err // Return the error if reading fails.
	}
	apj.AsmPath = []uint8(strings.TrimRight(string(path), "\x00"))
	return nil
}

// writeASMPath writes the ASM path to the provided writer.
// Returns an error if the write operation fails.
func (apj *APJFile) writeASMPath(f io.Writer) error {
	// if apj.AsmPath =="" then padding with 0x00
	if len(apj.AsmPath) == 0 {
		apj.AsmPath = make([]uint8, 256)
	} else if len(apj.AsmPath) > 256 {
		// if the length of AsmPath is greater than 256, truncate it to 256 bytes
		apj.AsmPath = apj.AsmPath[:256]
	} else {
		// if the length of AsmPath is less than 256, pad it with 0x00
		padding := make([]uint8, 256-len(apj.AsmPath))
		apj.AsmPath = append(apj.AsmPath, padding...)
	}

	if _, err := f.Write(apj.AsmPath); err != nil {
		return err
	}
	return nil
}

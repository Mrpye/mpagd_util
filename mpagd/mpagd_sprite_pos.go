package mpagd

import (
	"encoding/binary"
	"io"
	"strings"
)

// SpriteInfo represents the structure of a sprite's position data.
type SpriteInfo struct {
	Type    uint8 // Type of the sprite
	Image   uint8 // Image ID associated with the sprite
	Unknown uint8 // Placeholder for unknown data
	Screen  uint8 // Screen ID where the sprite is located
	X       uint8 // X-coordinate of the sprite
	Y       uint8 // Y-coordinate of the sprite
}

// SpriteInfoInit initializes the SpriteInfo slice. If `overwrite` is true, it clears the existing data.
func (apj *APJFile) SpriteInfoInit(overwrite bool) {
	if len(apj.SpriteInfo) == 0 || overwrite {
		apj.SpriteInfo = make([]SpriteInfo, 0)
	}
}

// ImportSpritePos parses sprite position data from a list of strings and associates it with a specific screen ID.
func (apj *APJFile) ImportSpritePos(screenId uint8, lines []string) error {
	// Initialize SpriteInfo if necessary
	// apj.SpriteInfoInit()

	// Loop through the lines
	for _, v := range lines {
		line := strings.ReplaceAll(v, "SPRITEPOSITION", "")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		blockValues := strings.Fields(line)

		// Parse and append sprite data
		spriteData := SpriteInfo{
			Type:    strUint8(blockValues[0]),
			Image:   strUint8(blockValues[1]),
			Unknown: 15, // Default value for unknown field
			Screen:  screenId,
			X:       strUint8(blockValues[2]),
			Y:       strUint8(blockValues[3]),
		}

		apj.SpriteInfo = append(apj.SpriteInfo, spriteData)
	}

	return nil
}

// readSpritePos reads sprite position data from a binary file and populates the SpriteInfo slice.
func (apj *APJFile) readSpritePos(f io.Reader) error {
	spriteInfo := []SpriteInfo{}
	for i := 0; i < int(apj.NrOfScreens); i++ {
		for {
			var spriteType uint8
			if err := binary.Read(f, binary.LittleEndian, &spriteType); err != nil {
				return err
			}
			if spriteType == 0xFF { // End of screen marker
				break
			}
			var image, unknown, x, y uint8
			if err := binary.Read(f, binary.LittleEndian, &image); err != nil {
				return err
			}
			if err := binary.Read(f, binary.LittleEndian, &unknown); err != nil {
				return err
			}
			if err := binary.Read(f, binary.LittleEndian, &x); err != nil {
				return err
			}
			if err := binary.Read(f, binary.LittleEndian, &y); err != nil {
				return err
			}
			spriteData := SpriteInfo{
				Type:    spriteType,
				Image:   image,
				Unknown: unknown,
				Screen:  uint8(i),
				X:       x,
				Y:       y,
			}
			spriteInfo = append(spriteInfo, spriteData)
		}
	}
	apj.SpriteInfo = spriteInfo
	return nil
}

// writeSpritePos writes sprite position data to a binary file.
func (apj *APJFile) writeSpritePos(f io.Writer) error {
	currScreen := uint8(0)
	for _, sprite := range apj.SpriteInfo {
		// Write screen separator if the screen changes
		if sprite.Screen != currScreen {
			currScreen = sprite.Screen
			if err := binary.Write(f, binary.LittleEndian, uint8(0xFF)); err != nil {
				return err
			}
		}
		// Write sprite data
		if err := binary.Write(f, binary.LittleEndian, sprite.Type); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, sprite.Image); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, sprite.Unknown); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, sprite.X); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, sprite.Y); err != nil {
			return err
		}
	}
	// Write final screen separator if there are sprites
	//if len(apj.SpriteInfo) > 0 {
	if err := binary.Write(f, binary.LittleEndian, uint8(0xFF)); err != nil {
		return err
	}
	//}
	return nil
}

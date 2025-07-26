package mpagd

import (
	"archive/tar"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	fc "github.com/fatih/color"
	"golang.org/x/sys/windows"
)

var user32_dll = windows.NewLazyDLL("user32.dll")
var GetKeyState = user32_dll.NewProc("GetKeyState")

// padSlice ensures the slice has the specified size by appending zeros if necessary.
// size: int - the desired size of the slice.
// Returns a slice of uint8 with the specified size.
func padSlice(size int) []uint8 {
	slice := make([]uint8, size)
	return slice
}

// strUint8 converts a string to a uint8 value.
// val: string - the input string to convert.
// Returns the uint8 representation of the string.
func strUint8(val string) uint8 {
	spriteTypeInt, _ := strconv.Atoi(val)
	return uint8(spriteTypeInt)
}

// iCount counts the number of elements in a slice of maps.
// v: interface{} - the input slice of maps.
// Returns the count as a uint8.
func iCount(v interface{}) uint8 {
	if blocks, ok := v.([]map[string]interface{}); ok {
		return uint8(len(blocks))
	} else {
		return 0
	}
}

// CopyFile copies a file from the source path to the destination filepath.
// src: string - the source file filepath.
// dst: string - the destination file filepath.
// Returns an error if the operation fails.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// LogMessage displays a log message with a specific format and color-coded output.
// element: string - the element being logged.
// message: string - the log message.
// level: string - the log level ("ok", "warning", "error").
func LogMessage(element, message, level string, noColor bool) {
	// Define color functions for different log levels
	yellow := fc.New(fc.FgYellow).SprintFunc()
	red := fc.New(fc.FgRed).SprintFunc()
	green := fc.New(fc.FgGreen).SprintFunc()
	orange := fc.New(fc.FgHiYellow).SprintFunc()
	blue := fc.New(fc.FgBlue).SprintFunc()
	white := fc.New(fc.FgWhite).SprintFunc()
	reset := fc.New(fc.Reset).SprintFunc()

	// Determine color based on log level
	var elementColor, messageColor string
	switch level {
	case "ok":
		elementColor = green(element)
		messageColor = white(message)
	case "warning":
		elementColor = orange(element)
		messageColor = white(message)
	case "error":
		elementColor = red(element)
		messageColor = white(message)
	case "info": // Added yellow level
		elementColor = yellow(element)
		messageColor = white(message)
	default:
		elementColor = white(element)
		messageColor = white(message)
	}

	// Format and display the log message
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if noColor {
		fmt.Printf("[%s] [%s] %s\n", timestamp, element, message)
	} else {
		fmt.Fprintf(fc.Output, "[%s] [%s] %s %s\n", blue(timestamp), elementColor, messageColor, reset())
	}

}

// SpectrumAttrToColors converts a Spectrum attribute byte to foreground and background colors.
// attr: uint8 - the Spectrum attribute byte.
// Returns the foreground and background colors as color.Color.
func SpectrumAttrToColors(attr uint8) (fg color.Color, bg color.Color) {
	// Extract the brightness bit (bit 6)
	bright := (attr & 0x40) >> 6

	// Extract the ink (foreground) and paper (background) colors
	ink := attr & 0x07          // Bits 0-2
	paper := (attr & 0x38) >> 3 // Bits 3-5

	// Define the base colors (0 = black, 1 = blue, 2 = red, etc.)
	baseColors := []color.Color{
		color.RGBA{0, 0, 0, 255},       // Black
		color.RGBA{0, 0, 192, 255},     // Blue
		color.RGBA{192, 0, 0, 255},     // Red
		color.RGBA{192, 0, 192, 255},   // Magenta
		color.RGBA{0, 192, 0, 255},     // Green
		color.RGBA{0, 192, 192, 255},   // Cyan
		color.RGBA{192, 192, 0, 255},   // Yellow
		color.RGBA{192, 192, 192, 255}, // White
	}

	// Define the bright colors
	brightColors := []color.Color{
		color.RGBA{0, 0, 0, 255},       // Black
		color.RGBA{0, 0, 255, 255},     // Bright Blue
		color.RGBA{255, 0, 0, 255},     // Bright Red
		color.RGBA{255, 0, 255, 255},   // Bright Magenta
		color.RGBA{0, 255, 0, 255},     // Bright Green
		color.RGBA{0, 255, 255, 255},   // Bright Cyan
		color.RGBA{255, 255, 0, 255},   // Bright Yellow
		color.RGBA{255, 255, 255, 255}, // Bright White
	}

	// Select the appropriate color sets based on the brightness bit
	if bright == 1 {
		fg = brightColors[ink]
		bg = brightColors[paper]
	} else {
		fg = baseColors[ink]
		bg = baseColors[paper]
	}

	return fg, bg
}

// CalcImageSize calculates the dimensions of an image based on sprite layout.
// startIndex: uint8 - the starting index of the sprites.
// endIndex: uint8 - the ending index of the sprites.
// columns: int - the number of columns in the layout.
// size: int - the size of each sprite.
// Returns the width and height of the image.
func CalcImageSize(startIndex, endIndex uint8, columns, size int) (int, int) {
	// Calculate the number of sprites and layout
	spriteCount := int(endIndex - startIndex)
	spriteWidth, spriteHeight := size, size
	rows := (spriteCount + columns - 1) / columns
	imageWidth := 0
	imageHeight := 0
	if columns == 1 {
		imageWidth = columns * spriteWidth
		imageHeight = rows * spriteHeight
	} else {
		imageWidth = columns*spriteWidth + (size * 1)
		imageHeight = rows*spriteHeight + (rows * 1)
	}

	return imageWidth, imageHeight
}

// CalcImageOffSet calculates the offset of a sprite in the image layout.
// Index: uint8 - the index of the sprite.
// startIndex: uint8 - the starting index of the sprites.
// columns: int - the number of columns in the layout.
// size: int - the size of each sprite.
// Returns the x and y offsets of the sprite.
func CalcImageOffSet(Index, startIndex uint8, columns, size int) (int, int) {
	// Calculate sprite position in the image
	spriteOffset := int(Index - startIndex)
	col := spriteOffset % columns
	row := spriteOffset / columns
	xOffset := col * (size + 1)
	yOffset := row * (size + 1)
	return xOffset, yOffset
}

// CSVToIntSlice converts a CSV string to a slice of integers.
// csv: string - the input CSV string.
// Returns a slice of integers.
func CSVToIntSlice(csv string) []int {
	parts := strings.Split(csv, ",")
	result := make([]int, len(parts))
	for i, part := range parts {
		value, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil {
			result[i] = int(value)
		}
	}
	return result
}

// ensureDirExists ensures that the directory for the given file path exists
func ensureDirExists(filePath string) error {
	// convert path to /
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// addFileToTar adds a file to the tar archive
func AddFileToTar(tarWriter *tar.Writer, filePath string) error {
	// Open the file to be added to the tar archive
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file info
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar header for the file
	header := &tar.Header{
		Name:    fileInfo.Name(),
		Size:    fileInfo.Size(),
		Mode:    int64(fileInfo.Mode()),
		ModTime: fileInfo.ModTime(),
	}

	// Write the header to the tar archive
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// Write the file content to the tar archive
	if _, err := io.Copy(tarWriter, file); err != nil {

		return err
	}
	return nil
}

// IsKeyPressed checks if a specific key is pressed.
func IsESCKeyPressed() bool {
	r1, _, _ := GetKeyState.Call(27) // Call API to get ESC key state.
	return r1 == 65409               // Code for KEY_UP event of ESC key.
}

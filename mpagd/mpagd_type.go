package mpagd

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Template represents a template with metadata such as name, file name, type, and description.
type Template struct {
	Name        string
	FileName    string
	Type        string
	Description string
}

// CreateTemplate creates a new Template instance by extracting the name from the file name.
func CreateTemplate(fileName, templateType, description string) Template {
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	return Template{
		Name:        name,
		FileName:    fileName,
		Type:        templateType,
		Description: description,
	}
}

// State represents the state of an APJ file with various flags and settings.
type State struct {
	FilePath       bool
	Windows        bool
	Header         bool
	Version        bool
	AsmPath        bool
	Blocks         bool
	Screens        bool
	EnterpriseBias bool
	LivesScore     bool
	Map            bool
	Fonts          bool
	Keys           bool
	Objects        bool
	SpriteInfo     bool
	Sprites        bool
	ULAPalette     bool
	BlocksOffSet   uint8
}

// APJFile represents the structure of an APJ file with all its components.
type APJFile struct {
	noColor        bool
	FilePath       string
	Description    string
	Windows        Windows
	Header         []uint8
	Version        uint32
	AsmPath        []uint8
	Blocks         []Block
	NrOfBlocks     uint8
	Screens        []Screen
	NrOfScreens    uint8
	EnterpriseBias uint8
	LivesScore     LivesScore
	Map            Map
	State          State
	Fonts          []Font
	Keys           []uint8
	Objects        []Object
	NrOfObjects    uint8
	SpriteInfo     []SpriteInfo
	Sprites        []Sprite
	NrOfSprites    uint8
	ULAPalette     ULAPalette
}

func (apj *APJFile) SetNoColorOutput(noColor bool) {
	apj.noColor = noColor
}

// NewAPJFile creates a new APJFile instance with default values.
func NewAPJFile(filePath string) *APJFile {
	o := &APJFile{
		FilePath: filePath,
		noColor:  false,
		Header:   make([]uint8, 4),
		Version:  10,
		Keys:     []uint8{87, 83, 65, 68, 32, 74, 72, 49, 50, 51, 52}, // Default key mappings
	}
	o.initStruct(true)
	o.ResetState()
	return o
}

// Display prints the details of the APJFile instance in a structured format.
func (apj *APJFile) Display() {
	fmt.Printf("FilePath: %s\n", apj.FilePath)
	fmt.Printf("Version: %d\n", apj.Version)
	fmt.Printf("NrOfBlocks: %d\n", apj.NrOfBlocks)
	fmt.Printf("NrOfScreens: %d\n", apj.NrOfScreens)
	fmt.Printf("NrOfObjects: %d\n", apj.NrOfObjects)
	fmt.Printf("NrOfSprites: %d\n", apj.NrOfSprites)
	fmt.Printf("EnterpriseBias: %d\n", apj.EnterpriseBias)

	apj.printSection("Windows", apj.Windows)
	apj.printSection("Header", apj.Header)
	apj.printSection("AsmPath", apj.AsmPath)
	apj.printBlocks()
	apj.printScreens()
	apj.printSection("LivesScore", apj.LivesScore)
	apj.printSection("Map", apj.Map)
	apj.printFonts()
	apj.printSection("Keys", apj.Keys)
	apj.printObjects()
	apj.printSpriteInfo()
	apj.printSprites()
	apj.printSection("ULAPalette", apj.ULAPalette)
}

// Helper function to print a generic section.
func (apj *APJFile) printSection(name string, data interface{}) {
	fmt.Printf("%s:\n", name)
	fmt.Printf("%+v\n", data)
}

// Helper function to print blocks.
func (apj *APJFile) printBlocks() {
	fmt.Println("Blocks:")
	for _, block := range apj.Blocks {
		fmt.Printf("%+v\n", block)
	}
}

// Helper function to print screens.
func (apj *APJFile) printScreens() {
	fmt.Println("Screens:")
	for _, screen := range apj.Screens {
		fmt.Printf("%+v\n", screen)
	}
}

// Helper function to print fonts.
func (apj *APJFile) printFonts() {
	fmt.Println("Fonts:")
	for _, font := range apj.Fonts {
		fmt.Printf("%+v\n", font)
	}
}

// Helper function to print objects.
func (apj *APJFile) printObjects() {
	fmt.Println("Objects:")
	for _, object := range apj.Objects {
		fmt.Printf("%+v\n", object)
	}
}

// Helper function to print sprite information.
func (apj *APJFile) printSpriteInfo() {
	fmt.Println("SpriteInfo:")
	for _, spriteInfo := range apj.SpriteInfo {
		fmt.Printf("%+v\n", spriteInfo)
	}
}

// Helper function to print sprites.
func (apj *APJFile) printSprites() {
	fmt.Println("Sprites:")
	for _, sprite := range apj.Sprites {
		fmt.Printf("%+v\n", sprite)
	}
}

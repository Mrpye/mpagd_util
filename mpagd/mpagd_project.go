package mpagd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Mrpye/mpagd_util/project_template"
)

// ListTemplates lists all available templates in the embedded templates directory.
func ListTemplates() ([]Template, error) {
	var templates []Template

	// Define the directory to read templates from
	dir := "."
	entries, err := project_template.Templates.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Create a new APJFile instance for processing templates
	apj := NewAPJFile("")

	// Iterate through directory entries to find YAML templates
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if ext == ".yaml" {
			// Read the template file
			file, err := project_template.Templates.ReadFile(entry.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to open template '%s': %w", entry.Name(), err)
			}

			// Load the template content into the APJFile instance
			err = apj.LoadYAMLFromString(file)
			if err != nil {
				return nil, fmt.Errorf("error reading template '%s': %w", entry.Name(), err)
			}

			// Append the template details to the list
			templates = append(templates, CreateTemplate(
				entry.Name(),
				apj.Description,
				"Project Template",
			))
		}
	}

	return templates, nil
}

// CreateProjectFromTemplate creates a new project file from a specified template.
func CreateProjectFromTemplate(projectFile string, templateName string) error {
	// Validate the project file name
	if projectFile == "" {
		return fmt.Errorf("project file name is empty")
	}

	// Ensure the project file has the correct extension
	if !strings.HasSuffix(projectFile, ".apj") {
		return fmt.Errorf("project file name must end with .apj")
	}

	// Extract the file name and directory path
	fileName := filepath.Base(projectFile)
	filePath := filepath.Dir(projectFile)

	// Create the directory if it doesn't exist
	if _, err := os.Stat(filePath); err != nil {
		err := os.MkdirAll(filePath, fs.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", filePath, err)
		}
	}

	// Ensure the template name has the correct extension
	if !strings.HasSuffix(templateName, ".yaml") {
		templateName = templateName + ".yaml"
	}

	// Read the template file
	file, err := project_template.Templates.ReadFile(templateName)
	if err != nil {
		return fmt.Errorf("failed to open template '%s': %w", templateName, err)
	}

	// Load the template content into a new APJFile instance
	apj := NewAPJFile(projectFile)
	err = apj.LoadYAMLFromString(file)
	if err != nil {
		return fmt.Errorf("error reading template '%s': %w", templateName, err)
	}

	// Set the file path and write the project file
	apj.FilePath = fileName
	err = apj.WriteAPJ(projectFile)
	if err != nil {
		return fmt.Errorf("failed to write project file '%s': %w", projectFile, err)
	}

	return nil
}

// GetStats returns statistics about the project.
func (apj *APJFile) DisplayStats() {

	fmt.Println("Project Statistics:")
	fmt.Printf("Blocks: %d\n", len(apj.Blocks))
	fmt.Printf("Sprites: %d\n", len(apj.Sprites))
	fmt.Printf("Screens: %d\n", len(apj.Screens))
	fmt.Printf("Objects: %d\n", len(apj.Objects))
	fmt.Printf("Maps: %d\n", len(apj.Map.Map))
	fmt.Printf("Fonts: %d\n", len(apj.Fonts))

}

type BlocksDoc struct {
	Count int             `json:"blocks"`
	Info  []BlocksDocInfo `json:"info"` // List of block documentation info
}
type BlocksDocInfo struct {
	ID        uint8  `json:"id"`
	Type      string `json:"type"`
	ForeColor string `json:"foreground_color"`
	BackColor string `json:"background_color"`
}
type SpriteDoc struct {
	Count int             `json:"sprites"`
	Info  []SpriteDocInfo `json:"info"` // List of sprite documentation info
}
type SpriteDocInfo struct {
	ID          uint8  `json:"id"`
	Description string `json:"description"`
	Frames      uint8  `json:"frames"`
}

type ScreenDoc struct {
	Count int             `json:"screens"`
	Info  []ScreenDocInfo `json:"info"` // List of screen documentation info
}

type BlockTypesDocInfo struct {
	Type  string `json:"type"`
	ID    uint8  `json:"id"`
	Count int    `json:"count"` // Count of blocks of this type
}

type ScreenDocInfo struct {
	ID          uint8                         `json:"id"`
	Description string                        `json:"description"`
	BlockTypes  map[uint8]BlockTypesDocInfo   `json:"block_types"`
	Sprites     map[uint8]ScreenSpriteDocInfo `json:"sprite_info"`
}

type ScreenSpriteDocInfo struct {
	Count  int   `json:"count"`
	Type   uint8 `json:"type"`   // Type of the sprite
	Image  uint8 `json:"image"`  // Image ID associated with the sprite
	Screen uint8 `json:"screen"` // Screen ID where the sprite is located
	X      uint8 `json:"x"`      // X-coordinate of the sprite
	Y      uint8 `json:"y"`      // Y-coordinate of the sprite
}

type Variables struct {
	Variable    string   `json:"variable"`
	Description string   `json:"description"` // Description of the variable
	Locations   []string `json:"locations"`   // Type of the variable
	Scope       string   `json:"scope"`       // Scope of the variable
}
type SpriteTypeImage struct {
	Image       string   `json:"image"`
	Description string   `json:"description"`
	FrameInfo   []string `json:"frame_info"`
}

type SpriteFrameDesc struct {
	FrameRange  string `json:"frame_range"` // Frame range in the format "start-end"
	Description string `json:"description"`
}

type SpriteImageDesc struct {
	ImageID    int               `json:"image_id"`    // Image ID of the sprite
	ImageName  string            `json:"image_name"`  // Image name of the sprite
	FrameDescs []SpriteFrameDesc `json:"frame_descs"` // Frame descriptions for the sprite
}

type SpriteType struct {
	EventType         string            `json:"event_type"`         // Type of the sprite event
	EventDescription  string            `json:"event_description"`  // Description of the sprite event
	ImageDescriptions []SpriteImageDesc `json:"image_descriptions"` // List of images and their frames
}
type ProjectInfo struct {
	FilePath   string       `json:"file_path"` // Path to the project fil
	Name       string       `json:"name"`      // Name of the project
	Blocks     BlocksDoc    `json:"blocks"`
	Sprites    SpriteDoc    `json:"sprites"`
	Screens    ScreenDoc    `json:"screens"`
	Objects    int          `json:"objects"`
	Maps       int          `json:"maps"`
	Fonts      int          `json:"fonts"`
	Variables  []Variables  `json:"variables"`
	SpriteType []SpriteType `json:"sprite_type"`
}

func (apj *APJFile) BuildProjectInfoJson() (string, error) {
	//Extract the path from the APJFile
	if apj.FilePath == "" {
		return "", errors.New("project file path is empty")
	}
	directoryPath := filepath.Dir(apj.FilePath)
	if directoryPath == "" {
		return "", errors.New("project file path is empty, cannot build project info")
	}

	name := filepath.Base(apj.FilePath)
	if name == "" {
		return "", errors.New("project file name is empty, cannot build project info")
	}
	// remove the extension from the name
	name = strings.TrimSuffix(name, filepath.Ext(name))

	projectInfo := ProjectInfo{
		Blocks: BlocksDoc{
			Count: len(apj.Blocks),
			Info:  make([]BlocksDocInfo, 0, len(apj.Blocks)),
		},
		Sprites: SpriteDoc{
			Count: len(apj.Sprites),
			Info:  make([]SpriteDocInfo, 0, len(apj.Sprites)),
		},
		Screens: ScreenDoc{
			Count: len(apj.Screens),
			Info:  make([]ScreenDocInfo, 0, len(apj.Screens)),
		},
		Objects:    len(apj.Objects),
		Maps:       len(apj.Map.Map),
		Fonts:      len(apj.Fonts),
		Variables:  ExtractVariablesFromCodeFiles(directoryPath),
		SpriteType: make([]SpriteType, 0),
		Name:       name,
		FilePath:   apj.FilePath,
	}

	// Blocks info
	for _, block := range apj.Blocks {
		attr := block.Spectrum[len(block.Spectrum)-1]
		fg, bg := SpectrumAttrToColors(attr)
		projectInfo.Blocks.Info = append(projectInfo.Blocks.Info, BlocksDocInfo{
			ID:   block.ID,
			Type: GetBlockTypeByTypeID(block.Type),
			ForeColor: func() string {
				r, g, b, a := fg.RGBA()
				return fmt.Sprintf("rgba(%d,%d,%d,%d)", r, g, b, a)
			}(),
			BackColor: func() string {
				r, g, b, a := bg.RGBA()
				return fmt.Sprintf("rgba(%d,%d,%d,%d)", r, g, b, a)
			}(),
		})
	}

	spriteTypes, err := ParseSpriteTypeFiles(directoryPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse sprite type files: %w", err)
	}
	projectInfo.SpriteType = spriteTypes
	// Sprites info
	for _, sprite := range apj.Sprites {
		projectInfo.Sprites.Info = append(projectInfo.Sprites.Info, SpriteDocInfo{
			ID:          sprite.SpriteID,
			Description: "",
			Frames:      sprite.Frames,
		})
	}

	// Screens info
	for _, screen := range apj.Screens {
		blockTypes := make(map[uint8]BlockTypesDocInfo)
		spriteTypes := make(map[uint8]ScreenSpriteDocInfo)

		// Collect block types
		for _, row := range screen.ScreenData {
			for _, blockID := range row {
				blockType := GetBlockTypeByTypeID(apj.Blocks[blockID].Type)
				if _, exists := blockTypes[blockID]; !exists {
					blockTypesDo := BlockTypesDocInfo{
						ID:   blockID, // I
						Type: blockType,
						// D is not available here, set to 0 or update if needed
						Count: 1,
					}
					blockTypes[blockID] = blockTypesDo
				} else {
					bt := blockTypes[blockID]
					bt.Count++
					bt.Type = blockType
					blockTypes[blockID] = bt
				}
			}
		}

		// Collect sprite info for this screen
		for _, spriteInfo := range apj.SpriteInfo {
			if spriteInfo.Screen != screen.ScreenID {
				continue
			}
			spriteID := spriteInfo.Image
			if int(spriteID) >= len(apj.Sprites) {
				continue
			}
			if s, exists := spriteTypes[spriteInfo.Type]; exists {
				s.Count++
				spriteTypes[spriteInfo.Type] = s
			} else {
				spriteTypes[spriteInfo.Type] = ScreenSpriteDocInfo{
					Count:  1,
					Type:   spriteInfo.Type,
					Image:  spriteID,
					Screen: screen.ScreenID,
					X:      spriteInfo.X,
					Y:      spriteInfo.Y,
				}
			}
		}

		projectInfo.Screens.Info = append(projectInfo.Screens.Info, ScreenDocInfo{
			ID:          screen.ScreenID,
			Description: "",
			BlockTypes:  blockTypes,
			Sprites:     spriteTypes,
		})
	}

	//marshal the projectInfo into JSON
	jsonData, err := json.MarshalIndent(projectInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal project info: %w", err)
	}

	// save the images for sprites and blocks
	if len(apj.Blocks) > 0 {
		blockDir := filepath.Join(directoryPath, "docs", "images", "blocks")
		err := apj.RenderBlockToSeperateBitmap(0, uint8(len(apj.Blocks)), blockDir, nil, 0)
		if err != nil {
			return "", fmt.Errorf("failed to render blocks: %w", err)
		}

	}
	if len(apj.Sprites) > 0 {
		spriteDir := filepath.Join(directoryPath, "docs", "images", "sprites")
		err := apj.RenderSpriteToSeperateBitmap(0, uint8(len(apj.Sprites)), spriteDir, nil, 0)
		if err != nil {
			return "", fmt.Errorf("failed to render sprites: %w", err)
		}

	}
	if len(apj.Screens) > 0 {
		screenDir := filepath.Join(directoryPath, "docs", "images", "screens")
		err := os.MkdirAll(screenDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create screen directory: %w", err)
		}
		for _, st := range apj.Screens {
			screenBMP := filepath.Join(screenDir, fmt.Sprintf("screen_%d.png", st.ScreenID))
			err := apj.RenderScreenToBitmap(st.ScreenID, screenBMP)
			if err != nil {
				return "", fmt.Errorf("failed to render screens: %w", err)
			}
		}

	}
	return string(jsonData), nil
}

// ExtractVariablesFromCodeFiles scans code files for variables named A-Z and documents their file locations.
func ExtractVariablesFromCodeFiles(codeDir string) []Variables {
	variableMap := make(map[string]Variables)
	varNames := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	varRegex := regexp.MustCompile(`\b([A-Z])\b`)

	files, err := ioutil.ReadDir(codeDir)
	if err != nil {
		fmt.Printf("Error reading code directory: %v\n", err)
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Only process .aXX files
		if !strings.HasSuffix(file.Name(), ".a00") &&
			!strings.HasSuffix(file.Name(), ".a01") &&
			!strings.HasSuffix(file.Name(), ".a02") &&
			!strings.HasSuffix(file.Name(), ".a03") &&
			!strings.HasSuffix(file.Name(), ".a04") &&
			!strings.HasSuffix(file.Name(), ".a05") &&
			!strings.HasSuffix(file.Name(), ".a06") &&
			!strings.HasSuffix(file.Name(), ".a07") &&
			!strings.HasSuffix(file.Name(), ".a08") &&
			!strings.HasSuffix(file.Name(), ".a09") &&
			!strings.HasSuffix(file.Name(), ".a10") &&
			!strings.HasSuffix(file.Name(), ".a11") &&
			!strings.HasSuffix(file.Name(), ".a12") &&
			!strings.HasSuffix(file.Name(), ".a13") &&
			!strings.HasSuffix(file.Name(), ".a14") &&
			!strings.HasSuffix(file.Name(), ".a15") &&
			!strings.HasSuffix(file.Name(), ".a16") &&
			!strings.HasSuffix(file.Name(), ".a17") &&
			!strings.HasSuffix(file.Name(), ".a18") &&
			!strings.HasSuffix(file.Name(), ".a19") &&
			!strings.HasSuffix(file.Name(), ".a20") {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(codeDir, file.Name()))
		if err != nil {
			continue
		}
		lines := strings.Split(string(content), "\n")
		eventType := ""
		if len(lines) > 0 {
			eventType = strings.TrimSpace(lines[0])
		}
		seen := make(map[string]bool)
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			// Only check for global variable assignments (LET X = ...) at the top level
			if strings.HasPrefix(trimmed, "LET ") {
				// Example: LET S = 0      ; seconds
				parts := strings.SplitN(trimmed, "=", 2)
				if len(parts) == 2 {
					left := strings.TrimSpace(parts[0][4:]) // after "LET "
					if len(left) == 1 && strings.Contains(varNames, left) {
						varName := left
						comment := ""
						if idx := strings.Index(parts[1], ";"); idx != -1 {
							comment = strings.TrimSpace(parts[1][idx+1:])
						}
						location := fmt.Sprintf("%s (%s)", file.Name(), eventType)
						if v, ok := variableMap[varName]; ok {
							found := false
							for _, loc := range v.Locations {
								if loc == location {
									found = true
									break // Skip if location already exists
								}
							}
							if !found {
								v.Locations = append(v.Locations, location)
							}
							scope := "local"
							if eventType == "EVENT RESTARTSCREEN" || eventType == "EVENT GAMEINIT" {
								scope = "global"
								v.Scope = scope
								if comment != "" {
									v.Description = comment // store comment in Scope for documentation
								}
							}
							variableMap[varName] = v
						} else {
							scope := "local"
							if eventType == "EVENT RESTARTSCREEN" || eventType == "EVENT GAMEINIT" {
								scope = "global"
							}
							variableMap[varName] = Variables{
								Variable:  varName,
								Locations: []string{location},
								Scope:     scope,
							}
						}
						seen[varName] = true
					}
				}
			}
		}
		// Also scan for usage elsewhere in the file
		matches := varRegex.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			varName := match[1]
			if strings.Contains(varNames, varName) && !seen[varName] {
				seen[varName] = true
				location := fmt.Sprintf("%s (%s)", file.Name(), eventType)
				// If the variable already exists, append the location
				// Otherwise, create a new entry with the variable name and location
				// and set the scope based on the event type
				// If the variable already exists, append the location

				if v, ok := variableMap[varName]; ok {
					//check if location already 	exist in v.Locations
					found := false
					for _, loc := range v.Locations {
						if loc == location {
							found = true
							break // Skip if location already exists
						}
					}
					if !found {
						v.Locations = append(v.Locations, location)
						variableMap[varName] = v
					}
				} else {
					scope := "local"
					if eventType == "EVENT RESTARTSCREEN" || eventType == "EVENT GAMEINIT" {
						scope = "global"
					}
					variableMap[varName] = Variables{
						Variable:  varName,
						Locations: []string{location},
						Scope:     scope,
					}
				}
			}
		}
	}
	// Convert the map to a slice for easier processing
	var variableMapSlice []Variables
	for _, v := range variableMap {
		if v.Variable == "" {
			continue // Skip empty variable entries
		}
		variableMapSlice = append(variableMapSlice, v)
	}
	// Sort the slice by variable name for consistency
	if len(variableMapSlice) > 0 {
		sort.Slice(variableMapSlice, func(i, j int) bool {
			return variableMapSlice[i].Variable < variableMapSlice[j].Variable
		})
	}

	return variableMapSlice
}

// ParseSpriteTypeFiles reads EVENT SPRITETYPE files a00 to a08 and parses header comments into SpriteType structs.
func ParseSpriteTypeFiles(codeDir string) ([]SpriteType, error) {
	var spriteTypes []SpriteType
	for i := 0; i <= 8; i++ {
		filename := filepath.Join(codeDir, fmt.Sprintf("splat.a%02d", i))
		contentBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			// skip missing files
			continue
		}
		content := string(contentBytes)
		lines := strings.Split(content, "\n")
		var eventType, eventDesc string
		var imageDescs []SpriteImageDesc

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "EVENT ") {
				eventType = strings.TrimPrefix(line, "EVENT ")
			}
			if strings.HasPrefix(line, ";Event Description:") {
				eventDesc = strings.TrimSpace(strings.TrimPrefix(line, ";Event Description:"))
			}
			if strings.HasPrefix(line, ";Image Description:") {
				imgDesc := strings.TrimSpace(strings.TrimPrefix(line, ";Image Description:"))
				// Example: IMAGE 0,Player Left, Frame: (0-1,Player Idle),(1-2,Start Jump),(3,Fly)
				parts := strings.SplitN(imgDesc, ",", 3)
				if len(parts) >= 3 {
					imageIDStr := strings.TrimPrefix(parts[0], "IMAGE ")
					imageID := 0
					fmt.Sscanf(imageIDStr, "%d", &imageID)
					imageName := strings.TrimSpace(parts[1])
					framePart := strings.TrimSpace(parts[2])
					framePart = strings.TrimPrefix(framePart, "Frame:")
					frameDescs := []SpriteFrameDesc{}
					for _, frameDesc := range strings.Split(framePart, "),") {
						frameDesc = strings.TrimSpace(frameDesc)
						frameDesc = strings.TrimPrefix(frameDesc, "(")
						frameDesc = strings.TrimSuffix(frameDesc, ")")
						frameFields := strings.SplitN(frameDesc, ",", 2)
						if len(frameFields) == 2 {
							frameDescs = append(frameDescs, SpriteFrameDesc{
								FrameRange:  strings.TrimSpace(frameFields[0]),
								Description: strings.TrimSpace(frameFields[1]),
							})
						}
					}
					imageDescs = append(imageDescs, SpriteImageDesc{
						ImageID:    imageID,
						ImageName:  imageName,
						FrameDescs: frameDescs,
					})
				}
			}
		}
		spriteTypes = append(spriteTypes, SpriteType{
			EventType:         eventType,
			EventDescription:  eventDesc,
			ImageDescriptions: imageDescs,
		})
	}
	return spriteTypes, nil
}

// BuildProjectReadme generates a Markdown README file from project info JSON data.
func BuildProjectReadme(fileName string, projectInfoJson []byte) error {
	var data ProjectInfo
	if err := json.Unmarshal(projectInfoJson, &data); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("# Project %s README\n\n", data.Name))
	buf.WriteString("## Overview\n\n")
	buf.WriteString(fmt.Sprintf("Project File Path: %s\n", data.FilePath))
	buf.WriteString("\n---\n\n")
	buf.WriteString("Number of elements in the project:\n")
	buf.WriteString(fmt.Sprintf("- **Blocks:** %d\n", data.Blocks.Count))
	buf.WriteString(fmt.Sprintf("- **Sprites:** %d\n", data.Sprites.Count))
	buf.WriteString(fmt.Sprintf("- **Screens:** %d\n", data.Screens.Count))
	buf.WriteString(fmt.Sprintf("- **Objects:** %d\n", data.Objects))
	buf.WriteString(fmt.Sprintf("- **Maps:** %d\n", data.Maps))
	buf.WriteString(fmt.Sprintf("- **Fonts:** %d\n", data.Fonts))
	buf.WriteString("\n---\n\n")

	buf.WriteString("## Block Types\n\n")
	buf.WriteString("### Block Type Count\n\n")
	blockTypeCount := make(map[string]int)
	for _, b := range data.Blocks.Info {
		blockTypeCount[b.Type]++
	}
	for t, c := range blockTypeCount {
		buf.WriteString(fmt.Sprintf("- %s: %d\n", t, c))
	}
	buf.WriteString("\n---\n\n")

	buf.WriteString("## Sprite Types\n\n")
	for _, st := range data.SpriteType {
		buf.WriteString(fmt.Sprintf("### %s\n", st.EventType))
		buf.WriteString(fmt.Sprintf("- Description: %s\n", st.EventDescription))
		if len(st.ImageDescriptions) > 0 {
			buf.WriteString("#### Images:\n")
			for _, img := range st.ImageDescriptions {
				buf.WriteString(fmt.Sprintf("- **Image %d:** %s ![](images/sprites/sprite_%d.png)\n", img.ImageID, img.ImageName, img.ImageID))
				for _, frame := range img.FrameDescs {
					buf.WriteString(fmt.Sprintf("  - Frame %s: %s\n", frame.FrameRange, frame.Description))
				}
			}
		}
		buf.WriteString("\n")
	}
	buf.WriteString("---\n\n")

	buf.WriteString("## Variables\n\n")
	for _, v := range data.Variables {
		buf.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", v.Variable, v.Scope, v.Description))
		if len(v.Locations) > 0 {
			buf.WriteString("  - Used in:\n")
			for _, loc := range v.Locations {
				buf.WriteString(fmt.Sprintf("    - %s\n", loc))
			}
		}
	}
	buf.WriteString("\n---\n\n")

	buf.WriteString("## Screens\n\n")
	for _, s := range data.Screens.Info {

		buf.WriteString(fmt.Sprintf("### Screen %d\n", s.ID))
		if s.Description != "" {
			buf.WriteString(fmt.Sprintf("- Description: %s\n", s.Description))
		}
		buf.WriteString(fmt.Sprintf("![](images/screens/screen_%d.png)\n", s.ID))
		buf.WriteString("- Block Types:\n")
		for _, id := range s.BlockTypes {

			buf.WriteString(fmt.Sprintf("  - %s (ID: %d) (Count: %d) ![](images/blocks/block_%d.png)\n", id.Type, id.ID, id.Count, id.ID))
		}
		if len(s.Sprites) > 0 {
			buf.WriteString("- Sprites:\n")
			for _, si := range s.Sprites {
				typeDesc := data.SpriteType[si.Type].EventDescription
				//add image description if available
				buf.WriteString(fmt.Sprintf("  - Type %d (%s), Image %d, Count %d, Pos (%d,%d) ![](images/sprites/sprite_%d.png)\n", si.Type, typeDesc, si.Image, si.Count, si.X, si.Y, si.Image))
			}
		}
		buf.WriteString("\n")
	}
	buf.WriteString("---\n")
	//writew to file
	os.WriteFile(fileName, buf.Bytes(), 0644)
	return nil
}

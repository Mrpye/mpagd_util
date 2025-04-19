package mpagd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadAPJ reads the entire APJ file and processes its components sequentially.
func (apj *APJFile) ReadAPJ() error {
	file, err := os.Open(apj.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Sequentially read and process each component of the APJ file.
	readers := []func(io.Reader) error{
		apj.readHeader,
		apj.readWindows,
		apj.readLivesScore,
		apj.readKeys,
		apj.readBlocks,
		apj.readSprite,
		apj.readObjects,
		apj.readScreens,
		apj.readMap,
		apj.readSpritePos,
		apj.readFont,
		apj.readULAPalette,
		apj.readEnterpriseBiasSetting,
		apj.readASMPath,
	}

	for _, reader := range readers {
		if err := reader(file); err != nil {
			return err
		}
	}
	return nil
}

// CreateBlank initializes a blank APJ file with default values.
func (apj *APJFile) CreateBlank() {
	apj.initStruct(true)
	apj.initDefaults()
	apj.CalcOffset()
}

// CreateState initializes and returns a default State object.
func (apj *APJFile) CreateState() State {
	return State{
		FilePath:       false,
		Windows:        false,
		Header:         false,
		Version:        false,
		AsmPath:        false,
		Blocks:         false,
		Screens:        false,
		EnterpriseBias: false,
		LivesScore:     false,
		Map:            false,
		Fonts:          false,
		Keys:           false,
		Objects:        false,
		SpriteInfo:     false,
		Sprites:        false,
		ULAPalette:     false,
		BlocksOffSet:   uint8(0),
	}
}

// ResetState resets the APJ file's state to its default values.
func (apj *APJFile) ResetState() {
	apj.State = apj.CreateState()
}

// initDefaults initializes default values for uninitialized components.
func (apj *APJFile) initDefaults() {
	if !apj.State.Screens {
		apj.ScreensDefault()
	}
	if !apj.State.Map {
		apj.MapDefault()
	}
	if !apj.State.Fonts {
		apj.FontDefault()
	}
	if !apj.State.Objects {
		apj.ObjectDefault()
	}
	if !apj.State.Sprites {
		apj.SpriteDefault()
	}
	if !apj.State.Blocks {
		apj.BlockDefault()
	}

	apj.CalcOffset()
}

// initStruct initializes the structure of the APJ file, optionally overwriting existing data.
func (apj *APJFile) initStruct(overwrite bool) {
	apj.HeaderInit()
	apj.ASMPathInit()
	apj.BlockInit(overwrite)
	apj.SpriteInit(overwrite)
	apj.ObjectInit(overwrite)
	apj.ScreensInit(overwrite)
	apj.SpriteInfoInit(overwrite)
	apj.FontInit(overwrite)
	apj.WindowsInit(overwrite)
	apj.ULAPaletteInit(overwrite)
	apj.LivesScoreInit(overwrite)
	apj.KeysInit(overwrite)
}

// ImportAGD imports data from an AGD file into the APJ file.
func (apj *APJFile) ImportAGD(agdFilePath string, options ImportOptions) error {
	file, err := os.Open(agdFilePath)
	if err != nil {
		return fmt.Errorf("failed to open AGD file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineBuf []string
	lineBufFunction := ""
	state := apj.CreateState()

	// Process each line in the AGD file.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Match and process specific AGD directives.
		switch {
		case strings.Contains(line, "DEFINEWINDOW") && !options.ignoreWindow:
			lineBuf = append(lineBuf, line)
			lineBufFunction = "DEFINEWINDOW"
			apj.handleDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINEWINDOW", apj.WindowsInit, apj.importWindows, options.owWindow, &state.Windows)
		case strings.Contains(line, "DEFINECONTROLS") && !options.ignoreKeys:
			lineBuf = append(lineBuf, line)
			lineBufFunction = "DEFINECONTROLS"
			apj.handleDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINECONTROLS", apj.KeysInit, apj.importKeys, options.owKeys, &state.Keys)
		case strings.Contains(line, "DEFINEBLOCK") && !options.ignoreBlocks:
			apj.handleBufferedDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINEBLOCK", apj.BlockInit, apj.ImportBlocks, options.owBlocks, &state.Blocks, &state.BlocksOffSet)
		case strings.Contains(line, "DEFINESPRITE") && !options.ignoreSprites:
			apj.handleBufferedDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINESPRITE", apj.SpriteInit, apj.ImportSprites, options.owSprites, &state.Sprites, nil)
		case strings.Contains(line, "DEFINEOBJECT") && !options.ignoreObjects:
			apj.handleBufferedDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINEOBJECT", apj.ObjectInit, apj.ImportObjects, options.owObjects, &state.Objects, nil)
		case strings.Contains(line, "DEFINESCREEN") && !options.ignoreScreens:
			apj.handleBufferedDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINESCREEN", apj.ScreensInit, apj.ImportScreens, options.owScreens, &state.Screens, nil)
		case strings.HasPrefix(line, "MAP") && !options.ignoreMaps:
			apj.handleBufferedDirective(&lineBufFunction, line, &lineBuf, state, options, "MAP", apj.MapInit, nil, options.owMaps, &state.Map, nil)
		case strings.HasPrefix(line, "DEFINEFONT") && !options.ignoreFonts:
			apj.handleBufferedDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINEFONT", apj.FontInit, apj.ImportFont, options.owFonts, &state.Fonts, nil)
		case strings.HasPrefix(line, "DEFINEPALETTE") && !options.ignoreULAPalette:
			apj.handleBufferedDirective(&lineBufFunction, line, &lineBuf, state, options, "DEFINEPALETTE", apj.ULAPaletteInit, apj.ImportULAPalette, options.owULAPalette, &state.ULAPalette, nil)
		case strings.HasPrefix(line, "DEFINEMESSAGES"), strings.HasPrefix(line, "EVENT"):
			apj.processData(&lineBufFunction, &lineBuf, state, options)
			lineBufFunction = ""
		default:
			if line != "" && !apj.checkIgnore(lineBufFunction, options) {
				lineBuf = append(lineBuf, line)
			} else {
				apj.processData(&lineBufFunction, &lineBuf, state, options)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading AGD file: %w", err)
	}

	apj.initDefaults()
	return nil
}

// handleDirective processes a single-line directive.
func (apj *APJFile) handleDirective(lineBufFunction *string, line string, lineBuf *[]string, state State, options ImportOptions, directive string, initFunc func(bool), importFunc func(string) error, overwrite bool, stateField *bool) {
	apj.processData(lineBufFunction, lineBuf, state, options)
	if !*stateField {
		*stateField = true
		initFunc(overwrite)
	}
	importFunc(line)
}

// handleBufferedDirective processes a buffered directive.
func (apj *APJFile) handleBufferedDirective(lineBufFunction *string, line string, lineBuf *[]string, state State, options ImportOptions, directive string, initFunc func(bool), importFunc func([]string) error, overwrite bool, stateField *bool, offsetField *uint8) {
	apj.processData(lineBufFunction, lineBuf, state, options)
	*lineBufFunction = directive
	if !*stateField {
		*stateField = true
		initFunc(overwrite)
		if offsetField != nil && !overwrite {
			*offsetField = uint8(len(apj.Blocks))
		}
	}
	*lineBuf = append(*lineBuf, line)
}

// function to check if the line is been ignored or not
func (apj *APJFile) checkIgnore(lineBufFunction string, options ImportOptions) bool {
	if strings.Contains(lineBufFunction, "DEFINEWINDOW") && options.ignoreWindow {
		return true
	} else if strings.Contains(lineBufFunction, "DEFINECONTROLS") && options.ignoreKeys {
		return true
	} else if strings.Contains(lineBufFunction, "DEFINEBLOCK") && options.ignoreBlocks {
		return true
	} else if strings.Contains(lineBufFunction, "DEFINESPRITE") && options.ignoreSprites {
		return true
	} else if strings.Contains(lineBufFunction, "DEFINEOBJECT") && options.ignoreObjects {
		return true
	} else if strings.Contains(lineBufFunction, "DEFINESCREEN") && options.ignoreScreens {
		return true
	} else if strings.HasPrefix(lineBufFunction, "MAP") && options.ignoreMaps {
		return true
	} else if strings.HasPrefix(lineBufFunction, "DEFINEFONT") && options.ignoreFonts {
		return true
	} else if strings.HasPrefix(lineBufFunction, "DEFINEPALETTE") && options.ignoreULAPalette {
		return true
	} else if lineBufFunction == "" {
		return true
	}
	return false
}
func (apj *APJFile) processData(lineBufFunction *string, lineBuf *[]string, state State, options ImportOptions) error {
	switch *lineBufFunction {
	case "DEFINEBLOCK":
		err := apj.ImportBlocks(*lineBuf)
		if err != nil {
			return fmt.Errorf("error importing blocks: %w", err)
		}
		*lineBufFunction = ""
		*lineBuf = make([]string, 0)
	case "DEFINESPRITE":
		err := apj.ImportSprites(*lineBuf)
		if err != nil {
			return fmt.Errorf("error importing sprites: %w", err)
		}
		*lineBufFunction = ""
		*lineBuf = make([]string, 0)
	case "DEFINEOBJECT":
		err := apj.ImportObjects(*lineBuf)
		if err != nil {
			return fmt.Errorf("error importing objects: %w", err)
		}
		*lineBufFunction = ""
		*lineBuf = make([]string, 0)
	case "DEFINESCREEN":
		err := apj.ImportScreens(*lineBuf)
		if err != nil {
			return fmt.Errorf("error importing screens: %w", err)
		}
		screen := uint8(len(apj.Screens) - 1)
		if state.BlocksOffSet > 0 && !options.owScreens {
			apj.RemapScreens(screen, state.BlocksOffSet)
		}
		*lineBufFunction = ""
		*lineBuf = make([]string, 0)
	case "MAP":
		//apj.ImportMap(*lineBuf)
		*lineBufFunction = ""
		*lineBuf = make([]string, 0)
	case "DEFINEFONT":
		err := apj.ImportFont(*lineBuf)
		if err != nil {
			return fmt.Errorf("error importing font: %w", err)
		}
		*lineBufFunction = ""
		*lineBuf = make([]string, 0)
	case "DEFINEPALETTE":
		err := apj.ImportULAPalette(*lineBuf)
		if err != nil {
			return fmt.Errorf("error importing ULA palette: %w", err)
		}
		*lineBufFunction = ""
		*lineBuf = make([]string, 0)
	}

	return nil
}

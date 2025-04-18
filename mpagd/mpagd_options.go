package mpagd

// ImportOptions defines options for importing various components.
// Each field represents whether to overwrite or ignore a specific component during import.
type ImportOptions struct {
	owWindow         bool // Overwrite window
	owKeys           bool // Overwrite keys
	owBlocks         bool // Overwrite blocks
	owSprites        bool // Overwrite sprites
	owObjects        bool // Overwrite objects
	owScreens        bool // Overwrite screens
	owMaps           bool // Overwrite maps
	owFonts          bool // Overwrite fonts
	owULAPalette     bool // Overwrite ULA palette
	ignoreWindow     bool // Ignore window
	ignoreKeys       bool // Ignore keys
	ignoreBlocks     bool // Ignore blocks
	ignoreSprites    bool // Ignore sprites
	ignoreObjects    bool // Ignore objects
	ignoreScreens    bool // Ignore screens
	ignoreMaps       bool // Ignore maps
	ignoreFonts      bool // Ignore fonts
	ignoreULAPalette bool // Ignore ULA palette
}

// CreateImportOptions initializes an ImportOptions instance with default values (all false).
func CreateImportOptions() ImportOptions {
	return ImportOptions{
		owWindow:         false,
		owKeys:           false,
		owBlocks:         false,
		owSprites:        false,
		owObjects:        false,
		owScreens:        false,
		owMaps:           false,
		owFonts:          false,
		owULAPalette:     false,
		ignoreWindow:     false,
		ignoreKeys:       false,
		ignoreBlocks:     false,
		ignoreSprites:    false,
		ignoreObjects:    false,
		ignoreScreens:    false,
		ignoreMaps:       false,
		ignoreULAPalette: false,
	}
}

// SetOwOptions sets the overwrite options for all components.
func (o *ImportOptions) SetOwOptions(owWindow, owKeys, owBlocks, owSprites, owObjects, owScreens, owMaps, owFonts, owULAPalette bool) {
	o.owBlocks = owBlocks
	o.owFonts = owFonts
	o.owMaps = owMaps
	o.owObjects = owObjects
	o.owScreens = owScreens
	o.owSprites = owSprites
	o.owULAPalette = owULAPalette
	o.owKeys = owKeys
	o.owWindow = owWindow
}

// SetOwOptionsTrue enables overwrite for all components.
func (o *ImportOptions) SetOwOptionsTrue() {
	o.SetOwOptions(true, true, true, true, true, true, true, true, true)
}

// SetOwOptionsFalse disables overwrite for all components.
func (o *ImportOptions) SetOwOptionsFalse() {
	o.SetOwOptions(false, false, false, false, false, false, false, false, false)
}

// SetIgnoreOptionsTrue enables ignore for all components.
func (o *ImportOptions) SetIgnoreOptionsTrue() {
	o.SetIgnoreOptions(true, true, true, true, true, true, true, true, true)
}

// SetIgnoreOptionsFalse disables ignore for all components.
func (o *ImportOptions) SetIgnoreOptionsFalse() {
	o.SetIgnoreOptions(false, false, false, false, false, false, false, false, false)
}

// SetIgnoreOptions sets the ignore options for all components.
func (o *ImportOptions) SetIgnoreOptions(ignoreWindow, ignoreKeys, ignoreBlocks, ignoreSprites, ignoreObjects, ignoreScreens, ignoreMaps, ignoreFonts, ignoreULAPalette bool) {
	o.ignoreBlocks = ignoreBlocks
	o.ignoreFonts = ignoreFonts
	o.ignoreMaps = ignoreMaps
	o.ignoreObjects = ignoreObjects
	o.ignoreScreens = ignoreScreens
	o.ignoreSprites = ignoreSprites
	o.ignoreULAPalette = ignoreULAPalette
	o.ignoreKeys = ignoreKeys
	o.ignoreWindow = ignoreWindow
}

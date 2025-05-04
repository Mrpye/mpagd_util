package mpagd

// IdToBlockType maps block IDs to their corresponding block type names.
var IdToBlockType = map[uint8]string{
	0: "EMPTYBLOCK",
	1: "PLATFORMBLOCK",
	2: "WALLBLOCK",
	3: "LADDERBLOCK",
	4: "FODDERBLOCK",
	5: "DEADLYBLOCK",
	// Add other block types as needed
}

// BlockTypeToId maps block type names to their corresponding block IDs.
var BlockTypeToId = map[string]uint8{
	"EMPTYBLOCK":    0,
	"PLATFORMBLOCK": 1,
	"WALLBLOCK":     2,
	"LADDERBLOCK":   3,
	"FODDERBLOCK":   4,
	"DEADLYBLOCK":   5,
	// Add other block types as needed
}

// GetBlockTypeByID returns the block type name for a given block ID.
// If the ID does not exist, it returns "UNKNOWN".
func GetBlockTypeByTypeID(id uint8) string {
	if blockType, exists := IdToBlockType[id]; exists {
		return blockType
	}
	return "UNKNOWN"
}

// GetBlockIDByType returns the block ID for a given block type name.
// If the block type does not exist, it returns 0 as the default value.
func GetBlockTypeIDByType(blockType string) uint8 {
	if id, exists := BlockTypeToId[blockType]; exists {
		return id
	}
	return 0 // Default value for unknown block types.
}

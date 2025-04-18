package mpagd

import (
	"bytes"
	"reflect"
)

// CompareData compares the data in the current APJFile with another APJFile
// and returns a map of differences.
func (apj *APJFile) CompareData(other *APJFile) map[string]interface{} {
	differences := make(map[string]interface{})

	// Compare Objects
	if !deepEqual(apj.Objects, other.Objects) {
		differences["Objects"] = map[string]interface{}{
			"self":  apj.Objects,
			"other": other.Objects,
		}
	}

	// Compare Sprites
	compareSprites(apj, other, differences)

	// Compare SpriteInfo
	if !deepEqual(apj.SpriteInfo, other.SpriteInfo) {
		differences["SpriteInfo"] = map[string]interface{}{
			"self":  apj.SpriteInfo,
			"other": other.SpriteInfo,
		}
	}

	// Compare ULAPalette
	if !deepEqual(apj.ULAPalette, other.ULAPalette) {
		differences["ULAPalette"] = map[string]interface{}{
			"self":  apj.ULAPalette,
			"other": other.ULAPalette,
		}
	}

	// Compare other fields if necessary
	// ...existing code for other fields...

	return differences
}

// compareSprites compares the sprites of two APJFile objects and updates the differences map.
func compareSprites(apj, other *APJFile, differences map[string]interface{}) {
	// Iterate over the sprites and compare each one
	for i := 0; i < int(apj.NrOfSprites); i++ {
		if i >= int(other.NrOfSprites) {
			differences["Sprites"] = map[string]interface{}{
				"self":  apj.Sprites[i],
				"other": nil,
			}
			continue
		}
		for j := 0; j < len(apj.Sprites[i].Spectrum); j++ {
			if !deepEqual(apj.Sprites[i].Spectrum[j], other.Sprites[i].Spectrum[j]) {
				differences["Sprites"] = map[string]interface{}{
					"self":  apj.Sprites[i],
					"other": other.Sprites[i],
				}
			}
		}
	}

	// Compare the overall Sprites field
	if !deepEqual(apj.Sprites, other.Sprites) {
		differences["Sprites"] = map[string]interface{}{
			"self":  apj.Sprites,
			"other": other.Sprites,
		}
	}
}

// deepEqual performs a deep comparison of two values, including nested types and arrays.
func deepEqual(value1, value2 interface{}) bool {
	// ...existing code...
	if reflect.DeepEqual(value1, value2) {
		return true
	}

	// Handle byte slices separately
	if v1, ok := value1.([]byte); ok {
		if v2, ok := value2.([]byte); ok {
			return bytes.Equal(v1, v2)
		}
	}

	// Handle slices of other types
	if v1, ok := value1.([]interface{}); ok {
		if v2, ok := value2.([]interface{}); ok {
			if len(v1) != len(v2) {
				return false
			}
			for i := range v1 {
				if !deepEqual(v1[i], v2[i]) {
					return false
				}
			}
			return true
		}
	}

	// Handle nested slices (e.g., [][]interface{})
	if v1, ok := value1.([][]interface{}); ok {
		if v2, ok := value2.([][]interface{}); ok {
			if len(v1) != len(v2) {
				return false
			}
			for i := range v1 {
				if !deepEqual(v1[i], v2[i]) {
					return false
				}
			}
			return true
		}
	}

	// Handle maps
	if m1, ok := value1.(map[string]interface{}); ok {
		if m2, ok := value2.(map[string]interface{}); ok {
			if len(m1) != len(m2) {
				return false
			}
			for key, val1 := range m1 {
				if val2, exists := m2[key]; !exists || !deepEqual(val1, val2) {
					return false
				}
			}
			return true
		}
	}

	return false
}

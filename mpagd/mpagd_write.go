package mpagd

import (
	"os"
)

// writeAPJ writes the APJ file to the specified output file path
func (apj *APJFile) WriteAPJ(outputFilePath string) error {
	// Ensure the output directory exists
	if err := ensureDirExists(outputFilePath); err != nil {
		return err
	}

	// Create the output file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Sequentially write all components of the APJ file
	if err := apj.writeHeader(file); err != nil {
		return err
	}

	if err := apj.writeWindows(file); err != nil {
		return err
	}

	if err := apj.writeLivesScore(file); err != nil {
		return err
	}

	if err := apj.writeKeys(file); err != nil {
		return err
	}

	if err := apj.writeBlocks(file); err != nil {
		return err
	}

	if err := apj.writeSprite(file); err != nil {
		return err
	}

	if err := apj.writeObjects(file); err != nil {
		return err
	}

	if err := apj.writeScreens(file); err != nil {
		return err
	}

	if err := apj.writeMap(file); err != nil {
		return err
	}

	if err := apj.writeSpritePos(file); err != nil {
		return err
	}

	if err := apj.writeFont(file); err != nil {
		return err
	}

	if err := apj.writeULAPalette(file); err != nil {
		return err
	}

	if err := apj.writeEnterpriseBiasSetting(file); err != nil {
		return err
	}

	if err := apj.writeASMPath(file); err != nil {
		return err
	}

	return nil
}

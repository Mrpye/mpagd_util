package tests

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Mrpye/mpagd_util/cmd"
	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

// CleanOutputFolder cleans the output folder by removing all files and directories within it
func CleanOutputFolder() {
	// empty the output folder
	if err := os.RemoveAll("output"); err != nil {
		panic(err)
	}
	if err := os.Mkdir("output", 0755); err != nil {
		panic(err)
	}
}

// Helper function to execute a command and capture its output
func executeCommand(t *testing.T, rootCmd *cobra.Command, args []string, output string) (string, error) {
	cmd.SetNoColor(true)
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//var outputBuffer bytes.Buffer
	//rootCmd.SetOut(os.Stdout)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Error executing import command: %v", err)
	}

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	// Check the output
	outputBuf := string(out)

	if !containsExpectedOutput(outputBuf, output) {
		t.Errorf("Unexpected output for command: %s", output)
	}
	return outputBuf, err
}

// Helper function to check if the output contains the expected text
func containsExpectedOutput(output, expected string) bool {
	return strings.Contains(output, expected)
}

// TestProjectBackup tests the project backup and restore functionality
func TestProjectBackup(t *testing.T) {
	// empty the output folder
	CleanOutputFolder()
	// Create a new project file for testing
	args := []string{"project", "import", "output/output.apj", "testproject.agd"}
	executeCommand(t, cmd.RootCmd, args, "AGD elements imported successfully")

	// Check if the file exists
	if _, err := os.Stat("output/output.apj"); os.IsNotExist(err) {
		t.Fatalf("File does not exist: %v", err)
	}

	// Capture the output of the commands
	args = []string{"project", "backup", "output/output.apj", "-c"}
	executeCommand(t, cmd.RootCmd, args, "Backup created successfully")

	args = []string{"project", "backups", "output/output.apj"}
	executeCommand(t, cmd.RootCmd, args, "Listing backup files for")

	// delete the project file on output folder
	os.Remove("output/output.apj")

	// now restore the project file
	args = []string{"project", "restore", "output/output.apj", "-c"}
	executeCommand(t, cmd.RootCmd, args, "Last backup restored successfully")

	//check if the file exists
	if _, err := os.Stat("output/output.apj"); os.IsNotExist(err) {
		t.Fatalf("Restored file does not exist: %v", err)
	}

	// now restore the project file
	args = []string{"project", "purge", "output/output.apj"}
	executeCommand(t, cmd.RootCmd, args, "All backup files purged successfully")

	CleanOutputFolder()
}

// TestSelectiveImport tests the project import functionality
func TestSelectiveImport(t *testing.T) {
	// empty the output folder
	CleanOutputFolder()
	// Create a new project file for testing
	args := []string{"project", "import-selective", "output/output.apj", "testproject.agd", "--blocks"}
	executeCommand(t, cmd.RootCmd, args, "AGD elements imported successfully")

	args = []string{"project", "stats", "output/output.apj"}
	executeCommand(t, cmd.RootCmd, args, "Blocks: 68")
	CleanOutputFolder()
}

// TestImportWriteCompare tests the import and write functionality of the APJ file
func TestImportWriteCompare(t *testing.T) {
	CleanOutputFolder()
	filePath := "blank.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)

	agdFilePath := "testproject.agd"
	opt := mpagd.CreateImportOptions()
	opt.SetOwOptionsTrue()
	opt.SetIgnoreOptions(false, false, false, false, false, false, false, false, false)
	apjFile.ImportAGD(agdFilePath, opt)

	// // Example comparison
	outputFilePath := "output/output.apj" // Replace with the other file path
	apjFile.WriteAPJ(outputFilePath)

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	otherAPJFile.ReadAPJ()

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		t.Errorf("Differences found: %+v", differences)
	}
	CleanOutputFolder()
}

// TestReadProjectWriteCompare tests the read and write functionality of the APJ file
func TestReadProjectWriteCompare(t *testing.T) {
	CleanOutputFolder()
	filePath := "splat.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	apjFile.ReadAPJ()

	outputFilePath := "output/output.apj" // Replace with the other file path
	err := apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	if err := otherAPJFile.ReadAPJ(); err != nil {
		t.Fatalf("Error reading other APJ file: %v", err)
	}

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		t.Errorf("Differences found: %+v", differences)
	}
	CleanOutputFolder()
}

// TestReadBlankProjectWriteCompare tests the read and write functionality of a blank APJ file
func TestReadBlankProjectWriteCompare(t *testing.T) {
	CleanOutputFolder()
	filePath := "../project_template/blank.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	apjFile.ReadAPJ()
	//apjFile.CreateBlank()

	outputFilePath := "output/output.apj" // Replace with the other file path
	err := apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	if err := otherAPJFile.ReadAPJ(); err != nil {
		t.Fatalf("Error reading other APJ file: %v", err)
	}
	// ***************************************************
	// use this to do a comparison of the yaml files using
	// something like https://www.text-comparer.com/yaml
	// ***************************************************
	//apjFile.SaveAsYAML("output/this.yaml")
	//otherAPJFile.SaveAsYAML("output/other.yaml")

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		t.Errorf("Differences found: %+v", differences)
	}
	CleanOutputFolder()

}

// TestCreateBlankProjectWriteCompare tests the creation of a blank project and its write functionality
func TestCreateBlankProjectWriteCompare(t *testing.T) {
	CleanOutputFolder()
	filePath := "../project_template/blank.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	//apjFile.ReadAPJ()
	apjFile.CreateBlank()

	outputFilePath := "output/output.apj" // Replace with the other file path
	err := apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	if err := otherAPJFile.ReadAPJ(); err != nil {
		t.Fatalf("Error reading other APJ file: %v", err)
	}
	// ***************************************************
	// use this to do a comparison of the yaml files using
	// something like https://www.text-comparer.com/yaml
	// ***************************************************
	//apjFile.SaveAsYAML("output/this.yaml")
	//otherAPJFile.SaveAsYAML("output/other.yaml")

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		t.Errorf("Differences found: %+v", differences)
	}
	CleanOutputFolder()
}

// TestReadProjectImportWriteCompare tests the read, import, and write functionality of the APJ file
func TestReadProjectImportWriteCompare(t *testing.T) {
	CleanOutputFolder()
	filePath := "splat.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	apjFile.ReadAPJ()

	agdFilePath := "testproject.agd"
	opt := mpagd.CreateImportOptions()
	opt.SetOwOptionsFalse()
	opt.SetIgnoreOptions(true, true, false, true, true, false, true, true, true)
	apjFile.ImportAGD(agdFilePath, opt)

	outputFilePath := "output/output.apj" // Replace with the other file path
	err := apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	if err := otherAPJFile.ReadAPJ(); err != nil {
		t.Fatalf("Error reading other APJ file: %v", err)
	}

	apjFile.SaveAsYAML("output/apj.yaml")
	otherAPJFile.SaveAsYAML("output/other_apj.yaml")

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		t.Errorf("Differences found: %+v", differences)
	}
	CleanOutputFolder()
}

// Test Sprite reordering
func TestSpriteReorder(t *testing.T) {
	CleanOutputFolder()
	filePath := "debug_test.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	err := apjFile.ReadAPJ()
	if err != nil {
		t.Fatalf("Error reading APJ file: %v", err)
	}

	// Store old spectrum data for comparison
	oldData := [][]byte{
		apjFile.Sprites[0].Spectrum[0].ImageData,
		apjFile.Sprites[1].Spectrum[0].ImageData,
	}

	// Reorder sprites
	err = apjFile.ReorderSprites([]int{1, 0}, 0)
	if err != nil {
		t.Fatalf("Error reordering sprites: %v", err)
	}

	outputFilePath := "output/output.apj" // Replace with the other file path
	err = apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}

	// Compare spectrum data after reordering
	compareSpriteSpectrumData(t, oldData[0], apjFile.Sprites[1].Spectrum[0].ImageData, 1)
	compareSpriteSpectrumData(t, oldData[1], apjFile.Sprites[0].Spectrum[0].ImageData, 0)

	// Validate sprite order
	checkOrder := []int{1, 0, 2, 1, 0, 2}
	for i := 0; i < len(checkOrder); i++ {
		if int(apjFile.SpriteInfo[i].Image) != checkOrder[i] {
			t.Errorf("SpriteInfo[%d].Image should be %d, got %d", i, checkOrder[i], apjFile.SpriteInfo[i].Image)
		}
	}
}

// Helper function to compare sprite spectrum data
func compareSpriteSpectrumData(t *testing.T, expected, actual []byte, spriteIndex int) {
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Errorf("Data mismatch at sprite %d, index %d: expected %d, got %d", spriteIndex, i, expected[i], actual[i])
		}
	}
}

func TestSpriteReorderOffSet(t *testing.T) {
	CleanOutputFolder()
	filePath := "debug_test.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	err := apjFile.ReadAPJ()
	if err != nil {
		t.Fatalf("Error reading APJ file: %v", err)
	}

	// Store old spectrum data for comparison
	oldData := [][]byte{
		apjFile.Sprites[1].Spectrum[0].ImageData,
		apjFile.Sprites[2].Spectrum[0].ImageData,
	}

	// Reorder sprites
	err = apjFile.ReorderSprites([]int{1, 0}, 1)
	if err != nil {
		t.Fatalf("Error reordering sprites: %v", err)
	}

	outputFilePath := "output/output.apj" // Replace with the other file path
	err = apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}

	// Compare spectrum data after reordering
	compareSpriteSpectrumData(t, oldData[0], apjFile.Sprites[2].Spectrum[0].ImageData, 2)
	compareSpriteSpectrumData(t, oldData[1], apjFile.Sprites[1].Spectrum[0].ImageData, 1)

	// Validate sprite order
	checkOrder := []int{0, 2, 1, 0, 2, 1}
	for i := 0; i < len(checkOrder); i++ {
		if int(apjFile.SpriteInfo[i].Image) != checkOrder[i] {
			t.Errorf("SpriteInfo[%d].Image should be %d, got %d", i, checkOrder[i], apjFile.SpriteInfo[i].Image)
		}
	}
}

// Helper function to compare spectrum data
func compareSpectrumData(t *testing.T, expected, actual []byte, blockIndex int) {
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Errorf("Data mismatch at block %d, index %d: expected %d, got %d", blockIndex, i, expected[i], actual[i])
		}
	}
}

func TestBlockReorder(t *testing.T) {
	CleanOutputFolder()
	filePath := "debug_test.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	err := apjFile.ReadAPJ()
	if err != nil {
		t.Fatalf("Error reading APJ file: %v", err)
	}

	// Store old spectrum data for comparison
	oldData := make([][]byte, len(apjFile.Blocks))
	for i := range apjFile.Blocks {
		oldData[i] = apjFile.Blocks[i].Spectrum
	}

	// Reorder blocks
	err = apjFile.ReorderBlocks([]int{5, 4, 3, 2, 1, 0}, 0)
	if err != nil {
		t.Fatalf("Error reordering blocks: %v", err)
	}

	outputFilePath := "output/output.apj" // Replace with the other file path
	err = apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}

	// Compare spectrum data after reordering
	for i, newIndex := range []int{5, 4, 3, 2, 1, 0, 6, 7, 8, 9, 10} {
		compareSpectrumData(t, oldData[i], apjFile.Blocks[newIndex].Spectrum, newIndex)
	}
}

// Test Block reordering
func TestBlockReorderOffset(t *testing.T) {
	CleanOutputFolder()
	filePath := "debug_test.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	err := apjFile.ReadAPJ()
	if err != nil {
		t.Fatalf("Error reading APJ file: %v", err)
	}

	// Store old spectrum data for comparison
	oldData := make([][]byte, len(apjFile.Blocks))
	for i := range apjFile.Blocks {
		oldData[i] = apjFile.Blocks[i].Spectrum
	}
	// Reorder sprites
	err = apjFile.ReorderBlocks([]int{5, 4, 3, 2, 1, 0}, 5)
	if err != nil {
		t.Fatalf("Error reordering sprites: %v", err)
	}

	outputFilePath := "output/output.apj" // Replace with the other file path
	err = apjFile.WriteAPJ(outputFilePath)
	if err != nil {
		t.Fatalf("Error writing APJ file: %v", err)
	}
	for i, newIndex := range []int{0, 1, 2, 3, 4, 10, 9, 8, 7, 6, 5} {
		compareSpectrumData(t, oldData[i], apjFile.Blocks[newIndex].Spectrum, newIndex)
	}
}

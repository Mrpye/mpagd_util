package tests

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/Mrpye/mpagd_util/cmd"
	"github.com/Mrpye/mpagd_util/mpagd"
	"github.com/spf13/cobra"
)

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
	rootCmd.SetArgs(args)
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//var outputBuffer bytes.Buffer
	//rootCmd.SetOut(&outputBuffer)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Error executing import command: %v", err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old

	// Check the output
	outputBuf := buf.String()
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

func TestImportWriteCompare(t *testing.T) {
	CleanOutputFolder()
	filePath := "blank.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	//apjFile.ReadAPJ()
	//apjFile.Display()

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
	//apjFile.SaveAsYAML("output/this.yaml")
	//otherAPJFile.SaveAsYAML("output/other.yaml")
	//apjFile.SaveAsYAML("output/apj.yaml")
	//otherAPJFile.SaveAsYAML("output/other_apj.yaml")

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		t.Errorf("Differences found: %+v", differences)
	}
	CleanOutputFolder()
}
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
	//apjFile.SaveAsYAML("output/this.yaml")
	//otherAPJFile.SaveAsYAML("output/other.yaml")
	//apjFile.SaveAsYAML("output/apj.yaml")
	//otherAPJFile.SaveAsYAML("output/other_apj.yaml")

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		t.Errorf("Differences found: %+v", differences)
	}
	CleanOutputFolder()
}

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

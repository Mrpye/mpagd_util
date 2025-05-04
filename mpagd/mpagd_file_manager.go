package mpagd

import (
	"archive/tar"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

// add a function to save the APJ file as a yaml file
func (apj *APJFile) SaveAsYAML(filePath string) error {
	// Implement the logic to save the APJ file as a YAML file
	obj, err := yaml.Marshal(apj)
	if err != nil {
		return err
	}
	os.WriteFile(filePath, obj, 0644)
	return nil
}

func (apj *APJFile) LoadYAMLFromString(yamlString []byte) error {
	err := yaml.Unmarshal(yamlString, apj)
	if err != nil {
		return err
	}
	return nil
}

// LoadYAML loads the APJ file from a YAML file
func (apj *APJFile) LoadYAML(filePath string) error {
	// Implement the logic to load the APJ file from a YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return apj.LoadYAMLFromString(data)
}

// List the backup files in the backup directory
func (apj *APJFile) ListBackupProjectFiles(backupDir string) ([]string, error) {
	// Open the backup directory
	dir, err := os.Open(backupDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory entries
	entries, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// Filter and collect backup files
	var backupFiles []string
	for _, entry := range entries {
		if strings.HasPrefix(entry, "backup_") && strings.HasSuffix(entry, ".tar") {
			backupFiles = append(backupFiles, entry)
		}
	}
	//sort the backup files by date
	sort.Slice(backupFiles, func(i, j int) bool {
		return backupFiles[i] < backupFiles[j]
	})

	return backupFiles, nil
}

// PurgeBackupFiles purges the backup files in the backup directory
func (apj *APJFile) PurgeBackupFiles(backupDir string) error {
	// Log start message with parameters
	LogMessage("PurgeBackupFiles", fmt.Sprintf("Starting purge for backup directory: %s", backupDir), "ok", apj.noColor)
	// List the backup files
	backupFiles, err := apj.ListBackupProjectFiles(backupDir)
	if err != nil {
		return err
	}
	for _, file := range backupFiles {
		// Construct the full path to the backup file
		backupFilePath := filepath.Join(backupDir, file)
		// Delete the backup file
		err := os.Remove(backupFilePath)
		if err != nil {
			return fmt.Errorf("failed to delete backup file: %w", err)
		}
		// Log success message
		LogMessage("PurgeBackupFiles", fmt.Sprintf("Successfully deleted backup file: %s", backupFilePath), "ok", apj.noColor)
		// List the backup files)

	}
	return nil
}

// Restore the last backup file
func (apj *APJFile) RestoreLastBackup(backupDir string, code bool) (string, error) {
	// Log start message with parameters
	LogMessage("RestoreLastBackup", fmt.Sprintf("Starting restore from backup directory: %s", backupDir), "ok", apj.noColor)
	// List the backup files)

	// List the backup files
	backupFiles, err := apj.ListBackupProjectFiles(backupDir)
	if err != nil {
		return "", err
	}

	// Check if there are no backup files
	if len(backupFiles) == 0 {
		return "", fmt.Errorf("no backup files found in %s", backupDir)
	}

	// Get the last backup file (the most recent one)
	lastBackupFile := backupFiles[len(backupFiles)-1]
	lastBackupFilePath := filepath.Join(backupDir, lastBackupFile)

	// Open the tar file
	tarFile, err := os.Open(lastBackupFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open backup file: %w", err)
	}
	defer tarFile.Close()

	// Create a tar reader
	tarReader := tar.NewReader(tarFile)

	// Extract files from the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar archive: %w", err)
		}

		// Construct the full path for the extracted file
		extractedFilePath := filepath.Join(filepath.Dir(apj.FilePath), header.Name)

		// If code is true, add the code to the file name
		if filepath.Ext(extractedFilePath) != ".apj" && !code {
			// Skip files that are not APJ files
			continue
		}

		// Get the file name without the extension
		// Create the extracted file
		extractedFile, err := os.Create(extractedFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to create extracted file: %w", err)
		}

		// Copy the file content from the tar archive
		if _, err := io.Copy(extractedFile, tarReader); err != nil {
			extractedFile.Close()
			return "", fmt.Errorf("failed to extract file: %w", err)
		}
		extractedFile.Close()
	}

	// Log success message
	LogMessage("RestoreLastBackup", fmt.Sprintf("Successfully restored backup file: %s", lastBackupFile), "ok", apj.noColor)
	// List the backup files)
	return lastBackupFile, nil
}

// Backup creates a backup of the APJ file in the same directory as the original file
func (apj *APJFile) BackupProjectFile(code bool) error {
	// Log start message with parameters
	LogMessage("BackupProjectFile", fmt.Sprintf("Starting backup for file: %s", apj.FilePath), "info", apj.noColor)

	// extract the path from apj.FilePath

	backupDir := filepath.Dir(apj.FilePath)
	backupDir = filepath.Join(backupDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	//extract the file name from the path
	fileName := filepath.Base(apj.FilePath)
	// get the file name without the extension
	ext := filepath.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, ext)

	// Check if the file name is empty
	if fileName == "" {
		return os.ErrInvalid
	}

	//if code is true, add the code to the file name
	// tar project file and files that begin with the project name excluding extension xxxx.a00   xxxx.a01 ...
	var projectFiles []string
	//read the project folder
	projectFolder := filepath.Dir(apj.FilePath)
	if code {
		// get the file name without the extension
		fileName = strings.TrimSuffix(fileName, ".apj")
		files, err := os.ReadDir(projectFolder)
		if err != nil {
			return err
		}
		// create a list of files that begin with the project name

		for _, file := range files {
			regex := regexp.MustCompile(fmt.Sprintf("^%s\\.a\\d{2}$", fileName))
			if strings.HasPrefix(file.Name(), fileName) && regex.MatchString(file.Name()) {
				projectFiles = append(projectFiles, file.Name())
			}
		}

	}

	currentTime := strings.ReplaceAll(time.Now().Format("2006-01-02_15-04-05"), ":", "-")
	currentTime = strings.ReplaceAll(currentTime, " ", "_")
	backupFileName := fmt.Sprintf("%s/backup_%s_%s%s", backupDir, fileName, currentTime, ".tar")

	//tar the project file and files that begin with the project name excluding extension xxxx.a00   xxxx.a01 ...
	// create a tar file
	tarFile, err := os.Create(backupFileName)
	if err != nil {
		return err
	}
	defer tarFile.Close()
	// create a tar writer
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()
	// add the project file to the tar file
	err = AddFileToTar(tarWriter, apj.FilePath)
	if err != nil {
		return err
	}
	// add the project files to the tar file
	for _, file := range projectFiles {
		filePath := filepath.Join(projectFolder, file)
		err = AddFileToTar(tarWriter, filePath)
		if err != nil {
			return err
		}
	}
	// close the tar writer
	err = tarWriter.Close()
	if err != nil {
		return err
	}

	//CopyFile(apj.FilePath, backupFileName)

	// Log success message after creating the backup
	LogMessage("BackupProjectFile", fmt.Sprintf("Backup created: %s", backupFileName), "ok", apj.noColor)
	return nil
}

// readChunk reads a chunk of bytes from the file
func (apj *APJFile) readChunk(f io.Reader, size int) []uint8 {
	items := make([]uint8, size)
	err := binary.Read(f, binary.LittleEndian, &items)
	if err != nil {
		fmt.Println("Error reading chunk:", err)
		return nil
	}
	return items
}

// MonitorFileChanges monitors the specified file for changes and creates a backup when changes are detected.
// It uses a polling mechanism to check the file's last modified time every 5 seconds.
func (apj *APJFile) MonitorFileChanges(code bool) {

	// Check if the file exists
	if _, err := os.Stat(apj.FilePath); os.IsNotExist(err) {
		LogMessage("MonitorFileChanges", fmt.Sprintf("File does not exist: %s", apj.FilePath), "error", apj.noColor)
		return
	}

	// Get the initial file info and last modified time
	fileInfo, err := os.Stat(apj.FilePath)
	if err != nil {
		LogMessage("MonitorFileChanges", "Error getting file info", "error", apj.noColor)
		return
	}
	lastModified := fileInfo.ModTime()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	//Check for ESC key pressed
	go func() {
		<-done
		LogMessage("MonitorFileChanges", "Exiting file monitoring loop", "info", apj.noColor)
		os.Exit(0)
	}()
	// go func() {
	// 	for {
	// 		if IsESCKeyPressed() {
	// 			LogMessage("MonitorFileChanges", "Exiting file monitoring loop", "info", apj.noColor)
	// 			os.Exit(0)
	// 		}
	// 	}
	// }()

	// Monitor for changes in a loop
	for {
		//sleep for 5 seconds
		time.Sleep(3 * time.Second)
		// Get the current file info
		fileInfo, err = os.Stat(apj.FilePath)
		if err != nil {
			LogMessage("MonitorFileChanges", "Error getting file info", "error", apj.noColor)
			return
		}

		// Check if the file's last modified time has changed
		if fileInfo.ModTime() != lastModified {
			LogMessage("MonitorFileChanges", "File has changed, creating backup...", "warning", apj.noColor)
			apj.createBackupOnChange(code) // Extracted backup logic into a helper function
			lastModified = fileInfo.ModTime()
		}

	}
}

// createBackupOnChange handles the backup creation when a file change is detected.
func (apj *APJFile) createBackupOnChange(code bool) {
	if err := apj.BackupProjectFile(code); err != nil {
		LogMessage("MonitorFileChanges", fmt.Sprintf("Failed to create backup: %v", err), "error", apj.noColor)
	}
}

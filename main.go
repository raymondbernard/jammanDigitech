package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"gopkg.in/ini.v1"
)

// Config structure to hold INI values
type Config struct {
	JamManType    string
	SSDVolumeName string
	WavFileLoc    string
}

// Function to load the INI configuration file
func loadConfig(file string) (*Config, error) {
	cfg, err := ini.Load(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read INI file: %v", err)
	}

	// Load values from the INI file
	config := &Config{
		JamManType:    cfg.Section("JamManSettings").Key("JamManType").String(),
		SSDVolumeName: cfg.Section("JamManSettings").Key("SSDLocation").String(),
		WavFileLoc:    cfg.Section("JamManSettings").Key("wavFileLoc").String(),
	}

	return config, nil
}

// Function to find the drive letter for the given volume name
func findDriveByVolume(volumeName string) (string, error) {
	// Use `wmic` to list logical drives with their volume names
	cmd := exec.Command("wmic", "logicaldisk", "get", "VolumeName,DeviceID")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute wmic command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the volume name matches
		if strings.Contains(line, volumeName) {
			// Extract the drive letter from the line
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[0], nil
			}
		}
	}

	return "", fmt.Errorf("drive with volume name '%s' not found", volumeName)
}

// Function to check if the directory exists and ask for user input
func checkDirectoryAndPrompt(rootDir string) (string, error) {
	// Check if the root directory exists
	if _, err := os.Stat(rootDir); !os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Directory '%s' already exists. Do you want to (A)ppend or (O)verwrite? ", rootDir)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToUpper(response))

		if response == "O" {
			// Overwrite: Remove the existing directory and its contents
			err := os.RemoveAll(rootDir)
			if err != nil {
				return "", fmt.Errorf("failed to remove existing directory: %v", err)
			}
			fmt.Println("Directory and contents removed. Proceeding with overwrite...")
		} else if response != "A" {
			return "", fmt.Errorf("invalid option selected")
		}
	}

	return rootDir, nil
}

// Function to create directory structure based on JamManType
func createDirectoryStructure(driveLetter, jamManType string) {
	var rootDir string

	// Determine the root directory based on JamMan type
	switch jamManType {
	case "JamManStereo":
		rootDir = filepath.Join(driveLetter, "JamManStereo")
	case "JamManSoloXT":
		rootDir = filepath.Join(driveLetter, "JamManSoloXT")
	default:
		log.Fatalf("Unsupported JamManType: %s", jamManType)
	}

	// Check if the directory exists and ask user whether to append or overwrite
	_, err := checkDirectoryAndPrompt(rootDir)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create directories for patches (Patch01 to Patch99) and PhraseA subfolders
	for i := 1; i <= 99; i++ {
		patchFolder := fmt.Sprintf("Patch%02d", i)
		patchPath := filepath.Join(rootDir, patchFolder, "PhraseA")

		// Create the full directory path
		err := os.MkdirAll(patchPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", patchPath, err)
		}

		fmt.Printf("Created directory: %s\n", patchPath)
	}

	fmt.Println("All directories created successfully.")
}

func main() {
	// Path to your INI file
	iniPath := "jamman.ini"

	// Load the configuration from the INI file
	config, err := loadConfig(iniPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Find the drive letter by the volume name (e.g., 'JAMMAN')
	driveLetter, err := findDriveByVolume(config.SSDVolumeName)
	if err != nil {
		log.Fatalf("Error finding drive by volume name: %v", err)
	}

	// Confirm the drive letter with the user
	fmt.Printf("Found SSD at drive '%s'. Is this correct? (y/n): ", driveLetter)
	var response string
	fmt.Scanln(&response)

	if response != "y" {
		fmt.Print("Please enter the correct drive letter: ")
		fmt.Scanln(&driveLetter)
	}

	// Create the directory structure based on the JamManType
	createDirectoryStructure(driveLetter, config.JamManType)
}

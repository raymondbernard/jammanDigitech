package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"gopkg.in/ini.v1"
)

// Config structure to hold INI values
type Config struct {
	JamManType    string
	SSDVolumeName string
	WavFileLoc    string
}

// PatchXML template structure
type PatchXMLData struct {
	PatchName   string
	RhythmType  string
	StopMode    string
	ID          string
}

// PhraseXML template structure
type PhraseXMLData struct {
	BeatsPerMinute   string
	BeatsPerMeasure  string
	IsLoop           string
	ID               string
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

// Function to generate the patch.xml file for each patch in PatchXX
func generatePatchXML(patchDir string, patchData PatchXMLData) error {
	patchTemplate := `<?xml version="1.0" encoding="UTF-8" ?>
<JamManPatch xmlns="http://schemas.digitech.com/JamMan/Patch" device="JamManSoloXT" version="1">
    <PatchName>{{.PatchName}}</PatchName>
    <RhythmType>{{.RhythmType}}</RhythmType>
    <StopMode>{{.StopMode}}</StopMode>
    <SettingsVersion>1</SettingsVersion>
    <ID>{{.ID}}</ID>
    <Metadata />
</JamManPatch>`

	// Save patch.xml in PatchXX directory
	patchFilePath := filepath.Join(patchDir, "patch.xml")
	f, err := os.Create(patchFilePath)
	if err != nil {
		return fmt.Errorf("failed to create patch.xml: %v", err)
	}
	defer f.Close()

	// Execute the template and write to file
	tmpl, err := template.New("patchXML").Parse(patchTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse patch.xml template: %v", err)
	}
	err = tmpl.Execute(f, patchData)
	if err != nil {
		return fmt.Errorf("failed to write patch.xml: %v", err)
	}

	return nil
}

// Function to generate the phrase.xml file for each patch inside PhraseA subdirectory
func generatePhraseXML(phraseDir string, phraseData PhraseXMLData) error {
	phraseTemplate := `<?xml version="1.0" encoding="UTF-8" ?>
<JamManPhrase xmlns="http://schemas.digitech.com/JamMan/Phrase" version="1">
    <BeatsPerMinute>{{.BeatsPerMinute}}</BeatsPerMinute>
    <BeatsPerMeasure>{{.BeatsPerMeasure}}</BeatsPerMeasure>
    <IsLoop>{{.IsLoop}}</IsLoop>
    <ID>{{.ID}}</ID>
    <Metadata />
</JamManPhrase>`

	// Save phrase.xml in PhraseA subdirectory
	phraseFilePath := filepath.Join(phraseDir, "phrase.xml")
	f, err := os.Create(phraseFilePath)
	if err != nil {
		return fmt.Errorf("failed to create phrase.xml: %v", err)
	}
	defer f.Close()

	// Execute the template and write to file
	tmpl, err := template.New("phraseXML").Parse(phraseTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse phrase.xml template: %v", err)
	}
	err = tmpl.Execute(f, phraseData)
	if err != nil {
		return fmt.Errorf("failed to write phrase.xml: %v", err)
	}

	return nil
}

// Function to create or open and write to songs.csv
func createCSV(driveLetter string) (*csv.Writer, *os.File, error) {
	csvFilePath := filepath.Join(driveLetter, "songs.csv")
	_, err := os.Stat(csvFilePath)
	var csvFile *os.File

	// If file doesn't exist, create it and write headers
	if os.IsNotExist(err) {
		csvFile, err = os.Create(csvFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create songs.csv: %v", err)
		}

		// Create CSV writer and write headers
		writer := csv.NewWriter(csvFile)
		headers := []string{"songName", "bpMeasure", "bpMinute", "StopMode", "RhythmType", "patchNumber", "wavFileLoc"}
		if err := writer.Write(headers); err != nil {
			return nil, nil, fmt.Errorf("failed to write CSV headers: %v", err)
		}
		writer.Flush()
	} else {
		// If file exists, open in append mode
		csvFile, err = os.OpenFile(csvFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open songs.csv: %v", err)
		}
	}

	writer := csv.NewWriter(csvFile)
	return writer, csvFile, nil
}

// Function to handle .wav files, generate XML files, and update CSV
func handleWavFilesAndCreateCSV(driveLetter, rootDir, wavFileLoc string) error {
	// Create or open the songs.csv file
	writer, csvFile, err := createCSV(driveLetter)
	if err != nil {
		return fmt.Errorf("failed to create or open CSV: %v", err)
	}
	defer csvFile.Close()

	// List the WAV files from the wavFileLoc directory
	files, err := os.ReadDir(wavFileLoc)
	if err != nil {
		return fmt.Errorf("failed to read wav files directory: %v", err)
	}

	// Process each WAV file
	for i, file := range files {
		if filepath.Ext(file.Name()) != ".wav" {
			continue // Skip non-wav files
		}

		patchNumber := fmt.Sprintf("Patch%02d", i+1)
		patchDir := filepath.Join(rootDir, patchNumber)
		phraseDir := filepath.Join(patchDir, "PhraseA")
		destPath := filepath.Join(phraseDir, "phrase.wav")

		// Copy the WAV file and rename it to phrase.wav
		srcFile := filepath.Join(wavFileLoc, file.Name())
		err := copyFile(srcFile, destPath)
		if err != nil {
			return fmt.Errorf("failed to copy wav file: %v", err)
		}

		// Write song details to the CSV
		record := []string{
			file.Name(),                   // songName
			"4",                           // bpMeasure (default)
			"124.9213180542",              // bpMinute (default)
			"StopInstantly",               // StopMode (default)
			"StudioKickAndHighHat",        // RhythmType (default)
			patchNumber,                   // Patch number
			filepath.Join(wavFileLoc, file.Name()), // wavFileLoc (absolute path)
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write to CSV: %v", err)
		}
		writer.Flush()

		// Generate patch.xml and phrase.xml for this patch
		patchData := PatchXMLData{
			PatchName:   file.Name(),
			RhythmType:  "StudioKickAndHighHat",
			StopMode:    "StopInstantly",
			ID:          fmt.Sprintf("patch-%d", i+1), // Unique ID
		}
		phraseData := PhraseXMLData{
			BeatsPerMinute:   "124.9213180542", // Default BPM
			BeatsPerMeasure:  "4",              // Default beats per measure
			IsLoop:           "1",              // Default loop setting
			ID:               fmt.Sprintf("phrase-%d", i+1), // Unique ID
		}

		// Generate patch.xml in PatchXX directory
		err = generatePatchXML(patchDir, patchData)
		if err != nil {
			return fmt.Errorf("failed to generate patch.xml: %v", err)
		}

		// Generate phrase.xml in PhraseA subdirectory
		err = generatePhraseXML(phraseDir, phraseData)
		if err != nil {
			return fmt.Errorf("failed to generate phrase.xml: %v", err)
		}

		fmt.Printf("Processed file: %s -> %s\n", file.Name(), destPath)
	}

	fmt.Println("CSV file, XML files, and directories updated successfully.")
	return nil
}

// Function to copy a file
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// Function to create directory structure based on the number of WAV files
func createDirectoryStructure(driveLetter, jamManType, wavFileLoc string) {
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

	// Count the number of WAV files in wavFileLoc directory
	files, err := os.ReadDir(wavFileLoc)
	if err != nil {
		log.Fatalf("Error reading wav file directory: %v", err)
	}

	// Create only the necessary number of directories based on the number of WAV files
	for i, file := range files {
		if filepath.Ext(file.Name()) != ".wav" {
			continue // Skip non-wav files
		}

		patchFolder := fmt.Sprintf("Patch%02d", i+1)
		patchPath := filepath.Join(rootDir, patchFolder, "PhraseA")

		// Create the full directory path
		err := os.MkdirAll(patchPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", patchPath, err)
		}

		fmt.Printf("Created directory: %s\n", patchPath)
	}

	fmt.Println("Necessary directories created successfully.")

	// Handle copying WAV files and creating CSV and XML
	err = handleWavFilesAndCreateCSV(driveLetter, rootDir, wavFileLoc)
	if err != nil {
		log.Fatalf("Error handling WAV files and creating CSV/XML: %v", err)
	}
}

func main() {
	// Path to your INI file
	fmt.Println("Welcome to JamMan utility created by Ray Bernard")
	fmt.Println("contact : ray.bernard@outlook.com")

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
	fmt.Printf("Found MicroSD card at drive '%s'. Is this correct? (y/n): ", driveLetter)
	var response string
	fmt.Scanln(&response) 

	if response != "y" {
		fmt.Print("Please enter the correct drive letter: ")
		fmt.Scanln(&driveLetter)
	}

	// Create the directory structure and handle WAV files and XML generation
	createDirectoryStructure(driveLetter, config.JamManType, config.WavFileLoc)
}

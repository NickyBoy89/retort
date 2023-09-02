package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// `listMetadataFiles` looks in the device and returns a list of all the
// file names for the metadata files
func ListMetadataFiles() ([]string, error) {
	entries, err := os.ReadDir(MetadataDir)
	if err != nil {
		return nil, err
	}

	fileNames := []string{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".metadata") {
			fileNames = append(fileNames, entry.Name())
		}
	}

	return fileNames, nil
}

func ReloadFiles() error {
	return exec.Command("systemctl", "restart", "xochitl").Run()
}

func UploadFile(fileName string, writeFiles bool) error {
	switch filepath.Ext(fileName) {
	case ".pdf":
		// To upload a file, we need to complete several steps:

		fileId := uuid.New()
		destFolder := filepath.Join(MetadataDir, fileId.String())
		thumbnailFolder := filepath.Join(MetadataDir, fileId.String()+".thumbnails")
		contentFile := filepath.Join(MetadataDir, fileId.String()+".content")
		metadataFile := filepath.Join(MetadataDir, fileId.String()+".metadata")
		outputFile := filepath.Join(MetadataDir, fileId.String()+filepath.Ext(fileName))

		fmt.Printf("Making new directory %s\n", destFolder)
		fmt.Printf("Making new thumbnail directory %s\n", thumbnailFolder)
		fmt.Printf("Making new content file %s\n", contentFile)
		fmt.Printf("Making new metadata file %s\n", metadataFile)
		fmt.Printf("Copying file from %s to %s\n", fileName, outputFile)

		if writeFiles {
			log.Warn("Writing files")
			// Create a folder
			if err := os.MkdirAll(destFolder, 0755); err != nil {
				return err
			}
			// Make a folder for the thumbnails
			if err := os.MkdirAll(thumbnailFolder, 0755); err != nil {
				return err
			}
			// Upload a blank content file
			if err := os.WriteFile(contentFile, []byte("{}"), 0644); err != nil {
				return err
			}

			// Upload the relevant metadata file
			metadata := NewMetadataForFile(fileName)

			metadataFile, err := os.OpenFile(metadataFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer metadataFile.Close()

			encoder := json.NewEncoder(metadataFile)
			encoder.SetIndent("", "    ")
			if err := encoder.Encode(metadata); err != nil {
				return err
			}

			// Upload the file to the correct location
			contents, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}

			if err := os.WriteFile(outputFile, contents, 0600); err != nil {
				return err
			}

			if err := ReloadFiles(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Unknown file type to upload: %s", fileName)
	}

	return nil
}

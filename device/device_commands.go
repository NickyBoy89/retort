package device

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/NickyBoy89/retort/metadata"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// createSupplementalFiles is an internal function that creates all of the
// necessary files for uploading a file to the reMarkable device
//
// `fileId` should be the UUID of the file being uploaded
//
// The final upload of the actual data file needs to be handled separately, as
// well as the metadate for that file
func createSupplementalFiles(fileUUID uuid.UUID) error {
	fileId := fileUUID.String()
	destFolder := filepath.Join(metadata.MetadataDir, fileId)
	thumbnailFolder := filepath.Join(metadata.MetadataDir, fileId+".thumbnails")
	contentFile := filepath.Join(metadata.MetadataDir, fileId+".content")

	// Create a folder
	log.Debugf("Creating new file directory at %s", destFolder)
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return err
	}

	// Make a folder for the thumbnails
	log.Debugf("Creating new thumbnail folder at %s", thumbnailFolder)
	if err := os.MkdirAll(thumbnailFolder, 0755); err != nil {
		return err
	}

	// Upload a blank content file
	log.Debugf("Creating new content file at %s", contentFile)
	if err := os.WriteFile(contentFile, []byte("{}"), 0644); err != nil {
		return err
	}

	return nil
}

// ReloadFiles restarts the reMarkable's user interface, which allows for newly
// uploaded files to be shown immediately
func ReloadFiles() error {
	return exec.Command("systemctl", "restart", "xochitl").Run()
}

// ReceiveFileBuffer takes in a file, given by its buffer, to the reMarkable device
//
// # Detection of the file's type is done through specifying the file's name
//
// This function also takes an optional refresh argument that determines if the
// device should be refreshed once the file is uploaded. This needs to be done
// before the file can be used, but refreshing too fast can cause the device
// to reset
func ReceiveFileBuffer(fileName string, fileBuffer io.Reader, shouldRefresh bool) error {
	switch filepath.Ext(fileName) {
	case ".pdf", ".epub":
		// To upload a file, we need to complete several steps:

		fileId := uuid.New()

		// Create the supplementary files
		createSupplementalFiles(fileId)

		// Now, we only need to handle the upload of the final file, and the metadata
		metadataFile := filepath.Join(metadata.MetadataDir, fileId.String()+".metadata")
		outputFile := filepath.Join(metadata.MetadataDir, fileId.String()+filepath.Ext(fileName))

		// Upload the relevant metadata file
		metadata := metadata.NewMetadataForFile(fileName)

		outputMetadata, err := os.OpenFile(metadataFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer outputMetadata.Close()

		encoder := json.NewEncoder(outputMetadata)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(metadata); err != nil {
			return err
		}

		// Upload the document file to the correct location
		destFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, fileBuffer); err != nil {
			return err
		}

		if shouldRefresh {
			if err := ReloadFiles(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Unknown file type to upload: %s", fileName)
	}

	return nil
}

package metadata

import (
	"os"
	"strings"
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

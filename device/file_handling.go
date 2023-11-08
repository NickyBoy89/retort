package device

import "github.com/NickyBoy89/retort/metadata"

// FileWithNameExists tests for the existence of a file with the same name on
// the device
//
// This is intended to allow preventing duplicate files being uploaded to the
// device
func FileWithNameExists(name string) (bool, error) {
	metadataFiles, err := metadata.ListMetadataFiles()
	if err != nil {
		return false, err
	}

	for _, metadataFile := range metadataFiles {
		meta, err := metadata.FromFilename(metadataFile)
		if err != nil {
			return false, err
		}

		if meta.VisibleName == name {
			return true, nil
		}
	}

	return false, nil
}

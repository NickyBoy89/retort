package metadata

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

const (
	MetadataDir = ".local/share/remarkable/xochitl/"
)

type DocumentType byte

const (
	DocumentTypeCollection DocumentType = iota
	DocumentTypeDocument
)

func (dt DocumentType) String() string {
	switch dt {
	case DocumentTypeCollection:
		return "CollectionType"
	case DocumentTypeDocument:
		return "DocumentType"
	}

	panic(fmt.Sprintf("Unknown document type: %v", int(dt)))
}

type FileMetadata struct {
	// VisibleName is the name of the file in the user interface
	VisibleName string `json:"visibleName"`
	// Parent is the document's parent folder in the document hierarchy
	// If the document has been deleted, this will be `"trash"`
	Parent string `json:"parent"`
	// A Unix timestamp, converted into a string, that represents
	// when the file was last modified
	LastModified string `json:"lastModified"`
	// A Unix timestamp formatted as a string
	LastOpened       string `json:"lastOpened"`
	LastOpenedPage   int    `json:"lastOpenedPage"`
	MetadataModified bool   `json:"metadatamodified"`
	Modified         bool   `json:"modified"`
	Pinned           bool   `json:"pinned"`
	Synced           bool   `json:"synced"`
	// The type of the file. This has been observed to be one of `MetadataType`
	Type    DocumentType `json:"type"`
	Version int          `json:"version"`
	Deleted bool         `json:"deleted"`
}

func FromFilename(fileName string) (FileMetadata, error) {
	var metadata FileMetadata

	metadataFile, err := os.Open(fileName)
	if err != nil {
		return metadata, err
	}
	defer metadataFile.Close()

	if err := json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
		return metadata, err
	}

	if metadata.Deleted || metadata.Parent == "trash" {
		return FileMetadata{}, nil
	}

	return metadata, nil
}

// `FromUUID` retrives the metadate for a given document identified by its UUID
func FromUUID(id *uuid.UUID) (FileMetadata, error) {
	return FromFilename(MetadataDir + id.String() + ".metadata")
}

// `FromName` returns a list of all the documents that have the same name
// as `visibleName`
func FromName(visibleName string) ([]FileMetadata, error) {
	metaFiles, uuids, err := listMetadataFiles()
	if err != nil {
		return nil, err
	}

	matchingMetadata := []FileMetadata{}

	for _, metadataPath := range metaFiles {
		meta, err := FromFilename(metadataPath)
		if err != nil {
			return matchingMetadata, err
		}

		if meta.VisibleName == visibleName {
			matchingMetadata = append(matchingMetadata, meta)
		}
	}

	return matchingMetadata, nil
}

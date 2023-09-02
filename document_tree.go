package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"nicholasnovak.io/retort/metadata"
)

var (
	ErrUnknownDocumentType = errors.New("error: Unknown document type")
)

type NodeFiletype byte

const (
	FiletypePDF NodeFiletype = iota
	FiletypeEpub
	FiletypeFolder
)

func ParseFiletype(raw string) NodeFiletype {
	switch raw {
	case "pdf":
		return FiletypePDF
	case "epub":
		return FiletypeEpub
	default:
		panic("Unknown filetype: " + raw)
	}
}

// `Node` represents part of a document tree on the reMarkable
type Node struct {
	Name     string
	Filetype NodeFiletype
	Doctype  metadata.DocumentType
	Doc      string
	Parent   *Node
	Children []*Node

	Id           uuid.UUID
	Exists       bool
	GetsModified bool
}

// `NewNode` creates a blank node that represents any type of file node
//
// Users creating a Folder or Document would want to use the specific constructor
// functions `NewFolder` or `NewDocument`
func NewNode(name string, parent *Node, filetype NodeFiletype, document string) *Node {
	var node Node

	node.Name = name
	node.Filetype = filetype

	if filetype == FiletypeFolder {
		node.Doctype = metadata.DocumentTypeCollection
	} else {
		node.Doctype = metadata.DocumentTypeDocument
	}

	node.Parent = parent
	node.Children = []*Node{}

	switch filetype {
	case FiletypePDF, FiletypeEpub:
		if document != "" {
			node.Doc = document
		} else {
			panic(fmt.Sprintf("No document provided for file node %s", name))
		}
	}

	node.Id = uuid.Nil
	node.Exists = false
	node.GetsModified = false

	metadataFiles, err := metadata.ListMetadataFiles()
	if err != nil {
		panic(err)
	}

	// TODO: Remove this array
	filteredMetadata := []metadata.FileMetadata{}
	filteredMetadataIds := []uuid.UUID{}

	// Get all the matching files that have the same name
	for _, fileName := range metadataFiles {
		meta, err := metadata.FromFilename(fileName)
		if err != nil {
			panic(err)
		}

		fileId := uuid.MustParse(strings.TrimSuffix(filepath.Base(fileName), ".metadata"))

		isRootNode := parent == nil && meta.Parent == ""
		hasMatchingParent := parent != nil && parent.Id == uuid.MustParse(meta.Parent)

		locationMatch := isRootNode || hasMatchingParent
		typeMatch := node.Doctype == meta.Type

		if locationMatch && typeMatch {
			filteredMetadata = append(filteredMetadata, meta)
			filteredMetadataIds = append(filteredMetadataIds, fileId)
		}
	}

	// There is a document already here, unpack its data
	switch len(filteredMetadata) {
	case 1:
		node.Id = filteredMetadataIds[0]
		node.Exists = true
	default:
		log.Fatalf("Unknown state: Selected file occurs many times at destination")
	}

	return &node
}

func NewDocument(filePath string, parent *Node) *Node {
	fileType := filepath.Ext(filePath)
	fileType = strings.TrimPrefix(fileType, ".")

	fileName := filepath.Base(filePath)

	return NewNode(fileName, parent, ParseFiletype(fileType), filePath)
}

func NewFolder(dir, parent string) *Node {
	panic("Unimplemented")
}

func (n *Node) AddChild(node *Node) {
	n.Children = append(n.Children, node)
}

func (n *Node) String() string {
	return n.GetFullPath()
}

func (n *Node) GetFullPath() string {
	if n.Parent == nil {
		return n.Name
	}
	return n.Parent.GetFullPath() + "/" + n.Name
}

func (n *Node) RenderCommon(outputDir string) {
	panic("Unimplemented")
}

// Renders a `DocumentType` tree node
func (n *Node) RenderDocument(outputDir string) error {
	if !n.Exists {
		n.RenderCommon(outputDir)

		if err := os.MkdirAll(outputDir+"/"+n.Id.String(), 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(outputDir+"/"+n.Id.String()+".thumbnails", 0755); err != nil {
			return err
		}

		src := n.Doc
		dst := fmt.Sprintf("%s/%v.%v", outputDir, n.Id, n.Filetype)

		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0666); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) RenderFolder(outputDir string) error {
	if !n.Exists {
		n.RenderCommon(outputDir)
	}

	for _, child := range n.Children {
		child.Render(outputDir)
	}

	return nil
}

func (n *Node) Render(outputDir string) error {
	panic("Unimplemented")
}

// This recursively constructs the document tree based on the top-level
// document/folder data structure on disk that we put in initially
func constructNodeTreeFromDisk(basepath, parent string) (*Node, error) {
	stat, err := os.Stat(basepath)
	if err != nil {
		return nil, err
	}

	var node *Node
	if stat.IsDir() {
		node = NewFolder(basepath, parent)
	} else {
		node = NewDocument(basepath, parent)

		// If the document already exists
		if node.Exists {
			// NOTE: Handle conflict resolution protocol here
			log.Warn("Document already exists, skipping")
		}
	}

	return node, nil
}

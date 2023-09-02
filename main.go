package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"nicholasnovak.io/retort/metadata"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "retort",
	}

	rootCmd.AddCommand(ListFilesCommand)
	rootCmd.AddCommand(UploadFileCommand)
	rootCmd.AddCommand(SearchByHashCommand)
	rootCmd.AddCommand(ReceiveFilesCommand)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

var ListFilesCommand = &cobra.Command{
	Use: "status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fileNames, err := metadata.ListMetadataFiles()
		if err != nil {
			return err
		}

		fmt.Println("FILE\tNAME\tTYPE")

		for _, filename := range fileNames {
			meta, err := metadata.FromFilename(filename)
			if err != nil {
				return err
			}

			fmt.Printf("%s\t%s\t%s\n", filename, meta.VisibleName, meta.Type)
		}
		return nil
	},
}

var UploadFileCommand = &cobra.Command{
	Use: "upload",
	RunE: func(cmd *cobra.Command, args []string) error {
		return metadata.UploadFile(args[0])
	},
}

var SearchByHashCommand = &cobra.Command{
	Use:  "hash-search",
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		originalFile, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		testHash := sha256.Sum256(originalFile)

		log.Infof("Hash of original file is %v", testHash)

		pdfFiles, err := filepath.Glob(filepath.Join(metadata.MetadataDir, "*.pdf"))
		if err != nil {
			return err
		}

		for _, existingFile := range pdfFiles {
			fileData, err := os.ReadFile(existingFile)
			if err != nil {
				return err
			}

			fileHash := sha256.Sum256(fileData)
			log.Infof("Hash of file %s is %v", existingFile, fileHash)
			if fileHash == testHash {
				log.Warnf("Hashes matched for pdf file: %s", existingFile)
			}
		}

		return nil
	},
}

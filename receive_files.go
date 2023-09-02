package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"nicholasnovak.io/retort/metadata"
	"tailscale.com/client/tailscale"
)

var ReceiveFilesCommand = &cobra.Command{
	Use: "receive-files",
	RunE: func(cmd *cobra.Command, args []string) error {

		var client tailscale.LocalClient

		log.Info("Waiting to receive any files")

		files, err := client.WaitingFiles(context.Background())
		if err != nil {
			return err
		}

		log.Infof("Received file list: %v", files)

		for _, file := range files {
			log.Infof("Fetching waiting file %s", file.Name)
			fileBuf, _, err := client.GetWaitingFile(context.Background(), file.Name)
			if err != nil {
				return err
			}

			log.Info("Uploading to device")
			if err := metadata.UploadFileBuffer(file.Name, fileBuf); err != nil {
				return err
			}

			log.Info("Done")

			fileBuf.Close()
		}

		return nil
	},
}
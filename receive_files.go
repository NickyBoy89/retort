package main

import (
	"context"

	"github.com/NickyBoy89/retort/device"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"tailscale.com/client/tailscale"
)

var ReceiveFilesCommand = &cobra.Command{
	Use:   "receive-files",
	Short: "Listens for taildrop files, and uploads them to the device",
	RunE: func(cmd *cobra.Command, args []string) error {

		var client tailscale.LocalClient

		log.Info("Waiting to receive any files")

		for {
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
				if err := device.UploadFileBuffer(file.Name, fileBuf); err != nil {
					log.Errorf("Uploading file to device failed with error: %v", err)
				}

				log.Info("Done")
				if err := client.DeleteWaitingFile(context.Background(), file.Name); err != nil {
					return err
				}

				fileBuf.Close()
			}
		}
	},
}

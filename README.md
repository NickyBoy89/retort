# retort

`retort` is a command line application for the reMarkable 2 tablet, with the goal of being able to connect the device to a Tailscale tailnet and upload files to the device using [Taildrop](https://tailscale.com/taildrop/)

## Prerequisites

* A Go toolchain to compile the project
* A local Tailscale daemon running on device

## Building

A `Makefile` has been included in the repo for help with cross-compiling the Go binary for the reMarkable tablet

```bash
make build-arm
```

Then, you'll need to copy the binary to the device somehow

## Running

The daemon can be started on device with `retort receive-files`.

This will connect to the local Tailscale API and wait for files to be uploaded to the device. Once one or more files have been uploaded, `retort` creates all the necessary metadata for the reMarkable format, and uploads the file to the local device.

## Random Commands

This is the command that I used to create a custom sleep screen background for the device:
```bash
convert <input-file> -resize 20%x20% -gravity center -background white -extent 1404x1872 <output-file>
```
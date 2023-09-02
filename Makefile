all: build-arm

build-arm:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build .

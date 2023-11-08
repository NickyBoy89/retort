BINARY := retort

.PHONY: all clean remote ship-binary build-arm
all: build-arm

remote: build-arm ship-binary

ship-binary:
	scp $(BINARY) remarkable:/home/root

build-arm:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build .

clean:
	rm $(BINARY)

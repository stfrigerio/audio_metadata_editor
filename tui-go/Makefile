.PHONY: build clean

BINARY_NAME=audio-metadata-edi

build:
	go build -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

.PHONY: build clean run

BINARY_NAME=audio-meta-edit

build:
	go build -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

.PHONY: run build clean

run: build
	go run .

build:
	go build -o bsky-schwartz *.go

clean:
	rm -f bsky-schwartz

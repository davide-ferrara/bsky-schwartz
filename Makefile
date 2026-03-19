.PHONY: run build clean

run: build
	./bsky-schwartz

build:
	go build -o bsky-schwartz main.go

clean:
	rm -f bsky-schwartz

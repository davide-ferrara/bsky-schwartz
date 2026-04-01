.PHONY: build run clean dev

build:
	templ generate
	go build -o bin/web .

run: build
	./bin/web

dev:
	air

clean:
	rm -rf bin
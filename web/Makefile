.PHONY: build run clean

build:
	templ generate
	go build -o bin/web

run: build
	./bin/web

clean:
	rm -rf bin
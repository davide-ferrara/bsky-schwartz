.PHONY: run build clean test

run:
	go run .

build:
	go build -o feedgen .

clean:
	rm -f feedgen
	rm -f feed_*.txt

test:
	go test -v ./...
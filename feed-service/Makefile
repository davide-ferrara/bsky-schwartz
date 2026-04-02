.PHONY: help build run dev clean test

build: templ ## Build del binario
	go build -o bin/feed-service

run: build ## Build ed esegui
	./bin/feed-service

dev: ## Sviluppo con hot reload (Air)
	air

templ: ## Genera template Templ
	templ generate

test: ## Esegui test
	go test -v ./...

clean: ## Rimuovi binari e file temporanei
	rm -rf bin/ tmp/

deps: ## Installa dipendenze
	go mod download
	go mod tidy

setup: deps ## Setup iniziale
	go mod tidy
	templ generate
	go build -o bin/feed-service

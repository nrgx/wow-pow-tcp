.PHONY: help server client up down deps lint

help: ## Shows helpful message
	@echo "Makefile commands:"
	@echo "help            Show available commands"
	@echo "server          Build and run the server using Docker"
	@echo "client          Build and run the client using Docker"
	@echo "up              Run both client and server using docker-compose"
	@echo "down            Stop both client and server using docker-compose"
	@echo "deps            Download Go dependencies"
	@echo "lint            Run the linter on the project"
	@echo ""
	@echo "Usage examples:"
	@echo "make server      - to run the server only"
	@echo "make client      - to run the client only"
	@echo "make up          - to run both server and client with docker-compose"

server: ## Build and run the server via Docker
	docker build -t wow-server -f ./dockerfiles/server.Dockerfile .
	docker run -d -p 9090:9090 --name wow-server wow-server

client: ## Build and run the client via Docker
	docker build -t wow-client -f ./dockerfiles/client.Dockerfile .
	docker run --name wow-client wow-client

up: ## Run both client and server via docker-compose
	docker-compose up --build -d

down: ## Stop both client and server via docker-compose
	docker-compose down

deps: ## Download Go dependencies
	go mod download

lint: ## Run Go linter
	golangci-lint run ./...

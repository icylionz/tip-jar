.PHONY: build run dev clean test migration-up migration-down migration-create docker-build docker-run

# Build the application
build:
	go build -o bin/tipjar ./cmd/server

# Run the application
run: build
	./bin/tipjar

# Development mode with live reload
dev:
	air

# Clean build artifacts
clean:
	rm -rf bin/ tmp/ 

# Run tests
test:
	go test -v ./...

# Database migrations
migration-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migration-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migration-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Docker commands
docker-build:
	docker build -f docker/Dockerfile -t tipjar .

docker-run:
	docker-compose up

docker-dev:
	docker-compose up postgres

# Generate templ files
templ:
	templ generate

# Generate sqlc files  
sqlc:
	sqlc generate

# Install development dependencies
deps:
	go install github.com/cosmtrek/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Setup development environment
setup: deps
	cp .env.example .env
	echo "Please update .env with your configuration"
	echo "Run 'make docker-dev' to start PostgreSQL"
	echo "Run 'make migration-up' to run migrations"
	echo "Run 'make templ' to generate templates"
	echo "Run 'make dev' to start development server"

# Production build with optimizations
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/tipjar ./cmd/server

# Generate all code
generate: templ sqlc

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	templ fmt .
.PHONY: all run build test clean tidy docker-up docker-down

APP_NAME = server
BUILD_DIR = bin
ENTRY_POINT = ./cmd/server/main.go

all: run

tidy:
	@echo "==> Tidying dependencies..."
	go mod tidy

run:
	@echo "==> Running application..."
	go run $(ENTRY_POINT)

build: tidy
	@echo "==> Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(ENTRY_POINT)
	@echo "==> Binary built at $(BUILD_DIR)/$(APP_NAME)"

test:
	@echo "==> Running unit tests..."
	go test -v ./...

clean:
	@echo "==> Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f app.db

docker-up:
	@echo "==> Starting Docker containers..."
	docker compose up --build -d

docker-down:
	@echo "==> Stopping Docker containers..."
	docker compose down
